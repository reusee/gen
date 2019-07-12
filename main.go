package main

import (
	"go/ast"
	"go/format"
	"go/types"
	"os"
	"regexp"
	"strings"

	"golang.org/x/tools/go/packages"
)

var templates = map[string]Template{
	"map": {
		Path: "github.com/reusee/gen/map",
		Args: []string{"Name", "Key", "Value"},
	},

	"heap": {
		Path: "github.com/reusee/gen/heap",
		Args: []string{"T"},
	},
}

type Template struct {
	Path string
	Args []string
}

func main() {
	src, err := Gen(
		os.Args[1],
		os.Args[2:]...,
	)
	if err != nil {
		panic(err)
	}
	println(src)
}

func Gen(
	name string,
	args ...string,
) (
	src string,
	err error,
) {

	template, ok := templates[name]
	if !ok {
		panic("bad template name")
	}

	argMap := make(map[string]string)
	for i, arg := range args {
		argMap[template.Args[i]] = arg
	}

	// load template package
	pkgs, err := packages.Load(&packages.Config{
		Mode: 0 |
			packages.NeedImports |
			packages.NeedDeps |
			packages.NeedTypes |
			packages.NeedSyntax |
			packages.NeedTypesInfo,
	}, template.Path)
	if err != nil {
		panic(err)
	}
	if packages.PrintErrors(pkgs) > 0 {
		return
	}
	if len(pkgs) != 1 {
		panic("bad import path")
	}
	pkg := pkgs[0]
	mappings := make(map[types.Object]string)

	// type mappings
	typeDecls := make(map[ast.Decl]bool)
	for _, file := range pkg.Syntax {
		for _, decl := range file.Decls {
			genDecl, ok := decl.(*ast.GenDecl)
			if !ok {
				continue
			}
			for _, spec := range genDecl.Specs {
				typeSpec, ok := spec.(*ast.TypeSpec)
				if !ok {
					continue
				}
				for i, tplTypeName := range template.Args {
					if tplTypeName != typeSpec.Name.Name {
						continue
					}
					obj := pkg.TypesInfo.Defs[typeSpec.Name]
					mappings[obj] = args[i]
					typeDecls[decl] = true
				}
			}
		}
	}

	// renames
	for _, file := range pkg.Syntax {
		for _, decl := range file.Decls {
			var doc *ast.CommentGroup
			switch decl := decl.(type) {
			case *ast.GenDecl:
				doc = decl.Doc
			case *ast.FuncDecl:
				doc = decl.Doc
			}
			if doc == nil {
				continue
			}
			matches := namePattern.FindStringSubmatch(doc.Text())
			if len(matches) == 0 {
				continue
			}

			name := strings.TrimSpace(matches[1])
			for k, v := range argMap {
				name = strings.ReplaceAll(name, "{"+k+"}", v)
			}
			doc.List = nil

			var ident *ast.Ident
			switch decl := decl.(type) {
			case *ast.GenDecl:
				spec := decl.Specs[0]
				switch spec := spec.(type) {
				case *ast.TypeSpec:
					ident = spec.Name
				}
			case *ast.FuncDecl:
				ident = decl.Name
			}
			if ident == nil {
				continue
			}
			obj := pkg.TypesInfo.Defs[ident]
			mappings[obj] = name

		}
	}

	// bundle sources
	buf := new(strings.Builder)
	for _, file := range pkg.Syntax {

		// apply mappings
		ast.Inspect(file, func(node ast.Node) bool {
			ident, ok := node.(*ast.Ident)
			if !ok {
				return true
			}
			obj := pkg.TypesInfo.Uses[ident]
			if obj == nil {
				obj = pkg.TypesInfo.Defs[ident]
				if obj == nil {
					return true
				}
			}
			mapTo, ok := mappings[obj]
			if !ok {
				return true
			}
			ident.Name = mapTo
			return true
		})

		// write
		for _, decl := range file.Decls {
			if _, ok := typeDecls[decl]; ok {
				continue
			}
			err = format.Node(buf, pkg.Fset, decl)
			if err != nil {
				return
			}
			buf.WriteString("\n")
		}
	}
	src = buf.String()

	return
}

var namePattern = regexp.MustCompile(`name:(.*)`)
