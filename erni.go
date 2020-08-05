package erni

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"

	"golang.org/x/tools/go/ast/astutil"
)

func init() {
	Traverse()
}

func Traverse() {
	fset := token.NewFileSet()
	pkg, err := parser.ParseDir(fset, ".", nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	for _, p := range pkg {
		for nam, fil := range p.Files {
			// inspect for func declaration
			astutil.Apply(fil, func(c *astutil.Cursor) bool {
				fd, ok := c.Node().(*ast.FuncDecl)
				if !ok {
					return true
				}
				nfd := astutil.Apply(fd, func(c *astutil.Cursor) bool {
					bs, ok := c.Node().(*ast.BlockStmt)
					if !ok {
						return true
					}
					nbs := astutil.Apply(bs, func(c *astutil.Cursor) bool {
						_, ok := c.Node().(*ast.BlockStmt)
						if ok {
							return true
						}
						as, ok := c.Node().(*ast.AssignStmt)
						if !ok {
							_, ok := c.Parent().(*ast.BlockStmt)
							if !ok {
								return true
							}
							c.InsertBefore(c.Node())
							c.Delete()
							return true
						}
						under := make(map[int]struct{})
						for i, e := range as.Lhs {
							ide, ok := e.(*ast.Ident)
							if !ok {
								continue
							}
							if ide.String() != "_" {
								continue
							}
							under[i] = struct{}{}
						}
						for _, e := range as.Rhs {
							cal, ok := e.(*ast.CallExpr)
							if !ok {
								continue
							}
							// check if imported
							_, ok = cal.Fun.(*ast.SelectorExpr)
							if ok {
								// traverse original package to check for function signature
								// fmt.Printf("imported: %v\n", cal.Fun)
								continue
							}
							ide, ok := cal.Fun.(*ast.Ident)
							if !ok {
								continue
							}
							fd, ok := ide.Obj.Decl.(*ast.FuncDecl)
							if !ok {
								continue
							}
							if fd.Type.Results != nil {
								for i, f := range fd.Type.Results.List {
									ide, ok := f.Type.(*ast.Ident)
									if !ok {
										continue
									}
									_, ok = under[i]
									if ide.String() == "error" && ok {
										// add fmt just in case
										astutil.AddImport(fset, fil, "fmt")
									}
									if ide.String() == "error" && ok {
										lhs, _ := as.Lhs[i].(*ast.Ident)
										lhs.Name = fmt.Sprintf("err%d", i)
										c.InsertBefore(as)
										c.InsertBefore(iferni(fd, lhs.Name, 1))
										c.Delete()
									}
								}
							}
						}
						return true
					}, nil)
					c.Replace(nbs)
					return true
				}, nil)
				c.Replace(nfd)
				return true
			}, nil)
			var buf bytes.Buffer
			format.Node(&buf, fset, fil)
			f, _ := os.Open(nam)
			s, _ := f.Stat()
			mod := s.Mode()
			f.Close()
			err := ioutil.WriteFile(nam, buf.Bytes(), mod)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

func iferni(fd *ast.FuncDecl, err string, t token.Pos) *ast.IfStmt {
	var bod *ast.BlockStmt
	// if funcdecl with return then return
	// else fmt.Println
	bod = &ast.BlockStmt{
		List: []ast.Stmt{
			&ast.ExprStmt{
				X: &ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   &ast.Ident{Name: "fmt"},
						Sel: &ast.Ident{Name: "Println"},
					},
					Args: []ast.Expr{
						&ast.Ident{
							Name: err,
							Obj: &ast.Object{
								Kind: ast.Var,
								Name: err,
								Decl: fd,
								Data: nil,
								Type: nil,
							},
						},
					},
					Ellipsis: token.NoPos,
				},
			},
		},
	}
	return &ast.IfStmt{
		If:   t,
		Init: nil,
		Cond: &ast.BinaryExpr{
			X: &ast.Ident{
				// NamePos: t,
				Name: err,
				Obj: &ast.Object{
					Kind: ast.Var,
					Name: err,
					Decl: fd,
					Data: nil,
					Type: nil,
				},
			},
			// OpPos: t,
			Op: token.NEQ,
			Y: &ast.Ident{
				// NamePos: t,
				Name: "nil",
				Obj:  nil,
			},
		},
		Body: bod,
		Else: nil,
	}
}
