package gencode

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/packages"
)

func GenCodeFromFile(filePath, typeName, tapeName string) (string, error) {
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return "", fmt.Errorf("resolve file path: %w", err)
	}

	cfg := &packages.Config{
		Mode:  packages.NeedName | packages.NeedFiles | packages.NeedCompiledGoFiles | packages.NeedTypes | packages.NeedSyntax,
		Dir:   filepath.Dir(absPath),
		Tests: true,
	}

	pkgs, err := packages.Load(cfg, "file="+absPath)
	if err != nil {
		return "", fmt.Errorf("load package: %w", err)
	}
	if loadErr := packagesLoadError(pkgs); loadErr != nil {
		return "", loadErr
	}

	if len(pkgs) == 0 || pkgs[0].Types == nil {
		return "", fmt.Errorf("no package found for %s", absPath)
	}
	pkg := pkgs[0]

	typeStrOverride := ""
	var targetType types.Type
	obj := pkg.Types.Scope().Lookup(typeName)
	if obj != nil {
		typeNameObj, ok := obj.(*types.TypeName)
		if !ok {
			return "", fmt.Errorf("%s is not a type in package %s", typeName, pkg.Name)
		}
		targetType = typeNameObj.Type()
	} else {
		pos := filePos(pkg, absPath)
		tv, err := types.Eval(pkg.Fset, pkg.Types, pos, typeName)
		if err != nil {
			return "", fmt.Errorf("type %s not found in package %s", typeName, pkg.Name)
		}
		if !tv.IsType() {
			return "", fmt.Errorf("%s is not a type expression in package %s", typeName, pkg.Name)
		}
		targetType = tv.Type
		typeStrOverride = typeName
	}

	gen := newTypeGenerator(pkg.Types, pkg.Name)
	return gen.Generate(targetType, tapeName, typeStrOverride)
}

func GenCodeFromDir(dir, tags string) (map[string]string, error) {
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return nil, fmt.Errorf("resolve dir: %w", err)
	}

	cfg := &packages.Config{
		Mode:  packages.NeedName | packages.NeedFiles | packages.NeedCompiledGoFiles | packages.NeedSyntax | packages.NeedTypes | packages.NeedTypesInfo,
		Dir:   absDir,
		Tests: true,
	}
	if tags != "" {
		cfg.BuildFlags = []string{"-tags=" + tags}
	}

	pkgs, err := packages.Load(cfg, ".")
	if err != nil {
		return nil, fmt.Errorf("load package: %w", err)
	}
	if loadErr := packagesLoadError(pkgs); loadErr != nil {
		return nil, loadErr
	}

	basePkg := selectBasePackage(pkgs)
	if basePkg == nil {
		return nil, fmt.Errorf("no package found in %s", absDir)
	}

	outputs := make(map[string]string)
	baseCode, err := generatePackageCode(basePkg, false)
	if err != nil {
		return nil, err
	}
	if baseCode != "" {
		outputs["fastape_gen.go"] = baseCode
	}

	if testPkg := selectTestPackage(pkgs, basePkg); testPkg != nil {
		testCode, err := generatePackageCode(testPkg, true)
		if err != nil {
			return nil, err
		}
		if testCode != "" {
			outputs["fastape_gen_test.go"] = testCode
		}
	}

	if len(outputs) == 0 {
		return nil, fmt.Errorf("no fastape:generate annotations found in %s", basePkg.PkgPath)
	}
	return outputs, nil
}

func packagesLoadError(pkgs []*packages.Package) error {
	var messages []string
	for _, pkg := range pkgs {
		for _, err := range pkg.Errors {
			messages = append(messages, err.Error())
		}
	}
	if len(messages) == 0 {
		return nil
	}
	return fmt.Errorf("load package: %s", strings.Join(messages, "; "))
}

func selectBasePackage(pkgs []*packages.Package) *packages.Package {
	for _, pkg := range pkgs {
		if pkg.ForTest != "" {
			continue
		}
		if strings.HasSuffix(pkg.Name, "_test") {
			continue
		}
		return pkg
	}
	return nil
}

func selectTestPackage(pkgs []*packages.Package, basePkg *packages.Package) *packages.Package {
	if basePkg == nil {
		return nil
	}
	for _, pkg := range pkgs {
		if pkg.ForTest == "" {
			continue
		}
		if pkg.Name != basePkg.Name {
			continue
		}
		return pkg
	}
	return nil
}

type annotatedType struct {
	name string
	typ  types.Type
}

func generatePackageCode(pkg *packages.Package, wantTestFiles bool) (string, error) {
	if pkg == nil || pkg.Types == nil {
		return "", fmt.Errorf("package type info missing")
	}
	targets, err := collectAnnotatedTypes(pkg, wantTestFiles)
	if err != nil {
		return "", err
	}
	if len(targets) == 0 {
		return "", nil
	}

	gen := newTypeGenerator(pkg.Types, pkg.Name)
	for _, target := range targets {
		if err := gen.AddTarget(target.typ, target.name+"Tape", target.name); err != nil {
			return "", err
		}
	}
	return gen.Render(), nil
}

func collectAnnotatedTypes(pkg *packages.Package, wantTestFiles bool) ([]annotatedType, error) {
	if pkg.TypesInfo == nil {
		return nil, fmt.Errorf("package types info missing")
	}

	seen := make(map[string]struct{})
	var out []annotatedType
	lineCache := make(map[string][]string)
	for i, file := range pkg.Syntax {
		filename := syntaxFilename(pkg, i, file)
		if wantTestFiles != strings.HasSuffix(filename, "_test.go") {
			continue
		}
		lines, err := fileLines(filename, lineCache)
		if err != nil {
			return nil, err
		}

		for _, decl := range file.Decls {
			genDecl, ok := decl.(*ast.GenDecl)
			if !ok || genDecl.Tok != token.TYPE {
				continue
			}
			for _, spec := range genDecl.Specs {
				typeSpec, ok := spec.(*ast.TypeSpec)
				if !ok {
					continue
				}
				if !hasFastapeAnnotation(pkg.Fset, file, genDecl, typeSpec, lines) {
					continue
				}
				obj := pkg.TypesInfo.Defs[typeSpec.Name]
				typeName, ok := obj.(*types.TypeName)
				if !ok {
					return nil, fmt.Errorf("failed to resolve type for %s", typeSpec.Name.Name)
				}
				if _, exists := seen[typeName.Name()]; exists {
					continue
				}
				seen[typeName.Name()] = struct{}{}
				out = append(out, annotatedType{name: typeName.Name(), typ: typeName.Type()})
			}
		}
	}
	return out, nil
}

const annotationMarker = "fastape:generate"

func hasFastapeAnnotation(fset *token.FileSet, file *ast.File, decl *ast.GenDecl, spec *ast.TypeSpec, lines []string) bool {
	if spec == nil || fset == nil || file == nil {
		return false
	}
	doc := spec.Doc
	if doc == nil {
		doc = decl.Doc
	}
	if commentGroupHasMarker(doc) {
		return false
	}

	anchorLine := fset.Position(spec.Pos()).Line
	if doc != nil {
		anchorLine = fset.Position(doc.Pos()).Line
	}
	if anchorLine == 0 {
		return false
	}

	candidateEnd := 0
	for _, cg := range file.Comments {
		if !commentGroupHasMarker(cg) {
			continue
		}
		endLine := fset.Position(cg.End()).Line
		if endLine >= anchorLine {
			continue
		}
		if endLine > candidateEnd {
			candidateEnd = endLine
		}
	}
	if candidateEnd == 0 || anchorLine-candidateEnd < 2 {
		return false
	}
	if !blankLinesOnly(lines, candidateEnd+1, anchorLine-1) {
		return false
	}
	return true
}

func commentGroupHasMarker(group *ast.CommentGroup) bool {
	if group == nil {
		return false
	}
	for _, comment := range group.List {
		if strings.Contains(comment.Text, annotationMarker) {
			return true
		}
	}
	return false
}

func blankLinesOnly(lines []string, startLine, endLine int) bool {
	if startLine > endLine {
		return true
	}
	if startLine < 1 {
		startLine = 1
	}
	if endLine > len(lines) {
		endLine = len(lines)
	}
	for i := startLine; i <= endLine; i++ {
		if strings.TrimSpace(lines[i-1]) != "" {
			return false
		}
	}
	return true
}

func fileLines(filename string, cache map[string][]string) ([]string, error) {
	if cached, ok := cache[filename]; ok {
		return cached, nil
	}
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("read file %s: %w", filename, err)
	}
	lines := strings.Split(string(data), "\n")
	cache[filename] = lines
	return lines, nil
}

func syntaxFilename(pkg *packages.Package, index int, file *ast.File) string {
	if pkg != nil && index >= 0 && index < len(pkg.CompiledGoFiles) {
		return pkg.CompiledGoFiles[index]
	}
	if pkg != nil && pkg.Fset != nil && file != nil {
		return pkg.Fset.Position(file.Pos()).Filename
	}
	return ""
}

func filePos(pkg *packages.Package, filePath string) token.Pos {
	if pkg == nil || pkg.Fset == nil {
		return token.NoPos
	}
	clean := filepath.Clean(filePath)
	for _, file := range pkg.Syntax {
		if pkg.Fset.Position(file.Pos()).Filename == clean {
			return file.Pos()
		}
	}
	return token.NoPos
}
