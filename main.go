package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	fset := token.NewFileSet()
	files, err := parseFiles(fset, []string{"."})
	if err != nil {
		fmt.Printf("Failed to parse files: %v\n", err)
		return
	}

	for _, file := range files {
		outputFile, err := generateToPbMethods(file, fset)
		if err != nil {
			fmt.Printf("Failed to generate ToPb methods for %s: %v\n", file.Name.Name, err)
			continue
		}

		if outputFile != "" {
			fmt.Printf("ToPb methods generated in %s\n", outputFile)
		}
	}
}

func parseFiles(fset *token.FileSet, paths []string) ([]*ast.File, error) {
	var files []*ast.File
	for _, path := range paths {
		packages, err := parser.ParseDir(fset, path, nil, parser.ParseComments)
		if err != nil {
			return nil, err
		}

		for _, pkg := range packages {
			for _, file := range pkg.Files {
				files = append(files, file)
			}
		}
	}
	return files, nil
}

func generateToPbMethods(file *ast.File, fset *token.FileSet) (string, error) {
	inFile, ok := getInputFileFromComments(file.Comments)
	if !ok {
		return "", nil // Skip if no input file is specified
	}

	outputDir, _ := filepath.Split(inFile)
	outputFileName := fmt.Sprintf("autogen_topb_%s.go", strings.TrimSuffix(filepath.Base(inFile), filepath.Ext(inFile)))
	outputFile := filepath.Join(outputDir, outputFileName)

	f, err := os.Create(outputFile)
	if err != nil {
		return "", err
	}
	defer f.Close()

	fmt.Fprintf(f, "// Code generated by topb; DO NOT EDIT.\n\n")
	fmt.Fprintf(f, "package %s\n\n", file.Name.Name)

	fmt.Fprintf(f, "import \"github.com/tikbox/topb/pb\"\n\n")

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

			structType, ok := typeSpec.Type.(*ast.StructType)
			if !ok {
				continue
			}

			if isGenToPb(typeSpec.Doc.List) {
				generateToPbMethod(f, typeSpec.Name.Name, structType)
			}
		}
	}

	return outputFile, nil
}

func isGenToPb(comments []*ast.Comment) bool {
	for _, c := range comments {
		if c.Text == "// gen:topb" {
			return true
		}
	}
	return false
}

func getInputFileFromComments(comments []*ast.CommentGroup) (string, bool) {
	pattern := `(?m)^//go:generate\s+topb\s+-in\s+(\S+)$`
	re := regexp.MustCompile(pattern)

	for _, group := range comments {
		for _, c := range group.List {
			match := re.FindStringSubmatch(c.Text)
			if len(match) > 1 {
				return match[1], true
			}
		}
	}

	return "", false
}

func generateToPbMethod(f *os.File, structName string, structType *ast.StructType) {
	fmt.Fprintf(f, "func (m *%s) ToPb() *pb.%s {\n", structName, structName)
	fmt.Fprintf(f, "    return &pb.%s{\n", structName)

	for _, field := range structType.Fields.List {
		if field.Names == nil {
			continue
		}
		fieldName := field.Names[0].Name

		fmt.Fprintf(f, "        %s: m.%s,\n", fieldName, fieldName)
	}

	fmt.Fprintf(f, "    }\n")
	fmt.Fprintf(f, "}\n")
}
