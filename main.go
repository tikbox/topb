// main.go
package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	inFile := flag.String("in", "", "")
	pbPkg := flag.String("pb", "", "")
	flag.Parse()

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, *inFile, nil, parser.ParseComments)
	if err != nil {
		fmt.Println(err)
		return
	}

	processFile(*inFile, file, *pbPkg)
}

func processFile(filename string, file *ast.File, pbPkg string) {
	for _, decl := range file.Decls {
		genDecl, ok := processDecl(decl, pbPkg, file.Name.Name)
		if ok {
			writeToFile(filename, genDecl)
		}
	}
}

func processDecl(decl ast.Decl, pbPkg string, pkgName string) ([]byte, bool) {
	genDecl, ok := decl.(*ast.GenDecl)
	if !ok || genDecl.Tok != token.TYPE {
		return nil, false
	}

	var buf bytes.Buffer
	for _, spec := range genDecl.Specs {
		typeSec, ok := spec.(*ast.TypeSpec)
		if !ok {
			continue
		}

		structDecl, ok := typeSec.Type.(*ast.StructType)
		if !ok {
			continue
		}

		// 获取结构体声明的文档注释
		doc := typeSec.Doc
		if doc == nil {
			doc = genDecl.Doc
		}

		if hasTopbAnnotation(doc) {
			// 生成 ToPb 方法
			buf.WriteString(fmt.Sprintf("package %s\n\n", pkgName))
			buf.WriteString(fmt.Sprintf("import \"%s\"\n\n", pbPkg))
			buf.WriteString(fmt.Sprintf("func (u *%s) ToPb() *pb.%s {\n", typeSec.Name, typeSec.Name))
			buf.WriteString("    return &pb" + "." + typeSec.Name.Name + "{\n")

			// 遍历结构体字段，生成对应的赋值语句
			for _, field := range structDecl.Fields.List {
				names := field.Names
				if len(names) == 0 {
					// 如果字段没有名称，则使用默认名称 f1, f2, ...
					names = []*ast.Ident{{Name: fmt.Sprintf("f%d", len(field.Names))}}
				}
				for _, name := range names {
					buf.WriteString(fmt.Sprintf("        %s: u.%s,\n", name.Name, name.Name))
				}
			}

			buf.WriteString("    }\n")
			buf.WriteString("}\n")
			return buf.Bytes(), true
		}
	}

	return nil, false
}

func hasTopbAnnotation(doc *ast.CommentGroup) bool {
	if doc == nil {
		return false
	}

	for _, comment := range doc.List {
		if strings.Contains(comment.Text, "gen:topb") {
			return true
		}
	}

	return false
}

func writeToFile(filename string, content []byte) {
	outFilename := generateOutputFilename(filename)
	f, err := os.Create(outFilename)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	f.Write(content)
}

func generateOutputFilename(filename string) string {
	dir, file := filepath.Split(filename)
	ext := filepath.Ext(file)
	base := strings.TrimSuffix(file, ext)
	return filepath.Join(dir, "autogen_topb_"+base+".go")
}
