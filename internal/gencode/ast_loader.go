package gencode

import (
	"fmt"
	"go/token"
	"go/types"
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
