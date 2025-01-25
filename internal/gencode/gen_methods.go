package gencode

import (
	"fmt"
	"strings"
)

// func generateStructMethods(tapeTypeStr string, tapeName string) string {
// 	var methods strings.Builder
// 	var fieldStack []string
// 	var methodStack []string
// 	// Parse field names and tape types from the tapeTypeStr
// 	lines := strings.Split(tapeTypeStr, "\n")
// 	fields := lines[1 : len(lines)-1] // Ignore the first and last lines (struct { ... })
// 	// Sizeof method
// 	methods.WriteString(fmt.Sprintf("func (t *%s) Sizeof(v %s) int {\n", tapeName, strings.TrimSuffix(tapeName, "Tape")))
// 	methods.WriteString("\tsize := 0\n")
// 	for _, field := range fields {
// 		parts := strings.Fields(field)
// 		if len(parts) == 1 && parts[0] == "}" {
// 			fieldStack = fieldStack[:len(fieldStack)-1]
// 			methodStack = methodStack[:len(methodStack)-1]
// 			continue
// 		}
// 		fieldName := strings.TrimSuffix(parts[0], "Tape")
// 		if parts[1] == "struct" {
// 			fieldStack = append(fieldStack, fieldName)
// 			methodStack = append(methodStack, fieldName+"Tape")
// 		} else {
// 			fullFieldName := strings.Join(append(fieldStack, fieldName), ".")
// 			methodPath := strings.Join(methodStack, ".")
// 			if len(fieldStack) > 0 {
// 				methods.WriteString(fmt.Sprintf("\tsize += t.%s.%s.Sizeof(v.%s)\n", methodPath, parts[0], fullFieldName))
// 			} else {
// 				methods.WriteString(fmt.Sprintf("\tsize += t.%s.Sizeof(v.%s)\n", parts[0], fieldName))
// 			}
// 		}
// 	}
// 	methods.WriteString("\treturn size\n")
// 	methods.WriteString("}\n\n")
// 	// Roll method
// 	methods.WriteString(fmt.Sprintf("func (t *%s) Roll(v %s, bs []byte) (n int, err error) {\n", tapeName, strings.TrimSuffix(tapeName, "Tape")))
// 	methods.WriteString("\tvar m int\n")
// 	for _, field := range fields {
// 		parts := strings.Fields(field)
// 		if len(parts) == 1 && parts[0] == "}" {
// 			fieldStack = fieldStack[:len(fieldStack)-1]
// 			methodStack = methodStack[:len(methodStack)-1]
// 			continue
// 		}
// 		fieldName := strings.TrimSuffix(parts[0], "Tape")
// 		if parts[1] == "struct" {
// 			fieldStack = append(fieldStack, fieldName)
// 			methodStack = append(methodStack, fieldName+"Tape")
// 		} else {
// 			fullFieldName := strings.Join(append(fieldStack, fieldName), ".")
// 			methodPath := strings.Join(methodStack, ".")
// 			if len(fieldStack) > 0 {
// 				methods.WriteString(fmt.Sprintf("\tn, err = t.%s.%s.Roll(v.%s, bs[m:])\n", methodPath, parts[0], fullFieldName))
// 			} else {
// 				methods.WriteString(fmt.Sprintf("\tn, err = t.%s.Roll(v.%s, bs[m:])\n", parts[0], fieldName))
// 			}
// 			methods.WriteString("\tif err != nil { return }\n")
// 			methods.WriteString("\tm += n\n")
// 		}
// 	}
// 	methods.WriteString("\treturn m, nil\n")
// 	methods.WriteString("}\n\n")
// 	// Unroll method
// 	methods.WriteString(fmt.Sprintf("func (t *%s) Unroll(bs []byte, v *%s) (n int, err error) {\n", tapeName, strings.TrimSuffix(tapeName, "Tape")))
// 	methods.WriteString("\tvar m int\n")
// 	for _, field := range fields {
// 		parts := strings.Fields(field)
// 		if len(parts) == 1 && parts[0] == "}" {
// 			fieldStack = fieldStack[:len(fieldStack)-1]
// 			methodStack = methodStack[:len(methodStack)-1]
// 			continue
// 		}
// 		fieldName := strings.TrimSuffix(parts[0], "Tape")
// 		if parts[1] == "struct" {
// 			fieldStack = append(fieldStack, fieldName)
// 			methodStack = append(methodStack, fieldName+"Tape")
// 		} else {
// 			fullFieldName := strings.Join(append(fieldStack, fieldName), ".")
// 			methodPath := strings.Join(methodStack, ".")
// 			if len(fieldStack) > 0 {
// 				methods.WriteString(fmt.Sprintf("\tn, err = t.%s.%s.Unroll(bs[m:], &v.%s)\n", methodPath, parts[0], fullFieldName))
// 			} else {
// 				methods.WriteString(fmt.Sprintf("\tn, err = t.%s.Unroll(bs[m:], &v.%s)\n", parts[0], fieldName))
// 			}
// 			methods.WriteString("\tif err != nil { return }\n")
// 			methods.WriteString("\tm += n\n")
// 		}
// 	}
// 	methods.WriteString("\treturn m, nil\n")
// 	methods.WriteString("}\n")

// 	return methods.String()
// }

func generateNonStructMethods(typeName string, tapeType string, elemType string) string {
	var methods strings.Builder

	// Sizeof method
	methods.WriteString(fmt.Sprintf("func (t *%s) Sizeof(v %s) int {\n", typeName, elemType))
	methods.WriteString(fmt.Sprintf("\treturn ((*%s)(t)).Sizeof(v)\n", tapeType))
	methods.WriteString("}\n\n")

	// Roll method
	methods.WriteString(fmt.Sprintf("func (t *%s) Roll(v %s, bs []byte) (n int, err error) {\n", typeName, elemType))
	methods.WriteString(fmt.Sprintf("\treturn ((*%s)(t)).Roll(v, bs)\n", tapeType))
	methods.WriteString("}\n\n")

	// Unroll method
	methods.WriteString(fmt.Sprintf("func (t *%s) Unroll(bs []byte, v *%s) (n int, err error) {\n", typeName, elemType))
	methods.WriteString(fmt.Sprintf("\treturn ((*%s)(t)).Unroll(bs, v)\n", tapeType))
	methods.WriteString("}\n")

	return methods.String()
}

func generateStructMethods(tapeTypeStr string, tapeName string) string {
	var methods strings.Builder
	var fieldStack []string
	var methodStack []string

	// Parse field names and tape types from the tapeTypeStr
	lines := strings.Split(tapeTypeStr, "\n")
	fields := lines[1 : len(lines)-1] // Ignore the first and last lines (struct { ... })

	// Sizeof method
	methods.WriteString(fmt.Sprintf("func (t *%s) Sizeof(v %s) int {\n", tapeName, strings.TrimSuffix(tapeName, "Tape")))
	methods.WriteString("\tsize := 0\n")
	for _, field := range fields {
		parts := strings.Fields(field)
		if len(parts) == 1 && parts[0] == "}" {
			fieldStack = fieldStack[:len(fieldStack)-1]
			methodStack = methodStack[:len(methodStack)-1]
			continue
		}
		fieldName := strings.TrimPrefix(strings.TrimSuffix(parts[0], "Tape"), "a")
		tapeFieldName := "a" + fieldName
		if parts[1] == "struct" {
			fieldStack = append(fieldStack, fieldName)
			methodStack = append(methodStack, tapeFieldName+"Tape")
		} else {
			fullFieldName := strings.Join(append(fieldStack, fieldName), ".")
			methodPath := strings.Join(methodStack, ".")
			if len(fieldStack) > 0 {
				methods.WriteString(fmt.Sprintf("\tsize += t.%s.%s.Sizeof(v.%s)\n", methodPath, tapeFieldName, fullFieldName))
			} else {
				methods.WriteString(fmt.Sprintf("\tsize += t.%s.Sizeof(v.%s)\n", tapeFieldName, fieldName))
			}
		}
	}
	methods.WriteString("\treturn size\n")
	methods.WriteString("}\n\n")

	// Roll method
	methods.WriteString(fmt.Sprintf("func (t *%s) Roll(v %s, bs []byte) (n int, err error) {\n", tapeName, strings.TrimSuffix(tapeName, "Tape")))
	methods.WriteString("\tvar m int\n")
	for _, field := range fields {
		parts := strings.Fields(field)
		if len(parts) == 1 && parts[0] == "}" {
			fieldStack = fieldStack[:len(fieldStack)-1]
			methodStack = methodStack[:len(methodStack)-1]
			continue
		}
		fieldName := strings.TrimPrefix(strings.TrimSuffix(parts[0], "Tape"), "a")
		tapeFieldName := "a" + fieldName
		if parts[1] == "struct" {
			fieldStack = append(fieldStack, fieldName)
			methodStack = append(methodStack, tapeFieldName+"Tape")
		} else {
			fullFieldName := strings.Join(append(fieldStack, fieldName), ".")
			methodPath := strings.Join(methodStack, ".")
			if len(fieldStack) > 0 {
				methods.WriteString(fmt.Sprintf("\tn, err = t.%s.%s.Roll(v.%s, bs[m:])\n", methodPath, tapeFieldName, fullFieldName))
			} else {
				methods.WriteString(fmt.Sprintf("\tn, err = t.%s.Roll(v.%s, bs[m:])\n", tapeFieldName, fieldName))
			}
			methods.WriteString("\tif err != nil { return }\n")
			methods.WriteString("\tm += n\n")
		}
	}
	methods.WriteString("\treturn m, nil\n")
	methods.WriteString("}\n\n")

	// Unroll method
	methods.WriteString(fmt.Sprintf("func (t *%s) Unroll(bs []byte, v *%s) (n int, err error) {\n", tapeName, strings.TrimSuffix(tapeName, "Tape")))
	methods.WriteString("\tvar m int\n")
	for _, field := range fields {
		parts := strings.Fields(field)
		if len(parts) == 1 && parts[0] == "}" {
			fieldStack = fieldStack[:len(fieldStack)-1]
			methodStack = methodStack[:len(methodStack)-1]
			continue
		}
		fieldName := strings.TrimPrefix(strings.TrimSuffix(parts[0], "Tape"), "a")
		tapeFieldName := "a" + fieldName
		if parts[1] == "struct" {
			fieldStack = append(fieldStack, fieldName)
			methodStack = append(methodStack, tapeFieldName+"Tape")
		} else {
			fullFieldName := strings.Join(append(fieldStack, fieldName), ".")
			methodPath := strings.Join(methodStack, ".")
			if len(fieldStack) > 0 {
				methods.WriteString(fmt.Sprintf("\tn, err = t.%s.%s.Unroll(bs[m:], &v.%s)\n", methodPath, tapeFieldName, fullFieldName))
			} else {
				methods.WriteString(fmt.Sprintf("\tn, err = t.%s.Unroll(bs[m:], &v.%s)\n", tapeFieldName, fieldName))
			}
			methods.WriteString("\tif err != nil { return }\n")
			methods.WriteString("\tm += n\n")
		}
	}
	methods.WriteString("\treturn m, nil\n")
	methods.WriteString("}\n")

	return methods.String()
}
