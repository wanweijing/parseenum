package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"regexp"
	"strings"
)

type pkg struct {
	structInfos map[string]*ast.StructType
	mod         string
	fullPath    string
	mgr         *global

	// 包名、包路径
	imports map[string]string

	// type 定义
	types map[string]string
}

func NewPkg(mgr *global, mod string, fullPath string) *pkg {
	p := &pkg{
		structInfos: make(map[string]*ast.StructType),
		mod:         mod,
		fullPath:    fullPath,
		mgr:         mgr,
		imports:     make(map[string]string),
		types:       make(map[string]string),
	}

	return p
}

func (p *pkg) printType(t ast.Expr) string {
	fn := func(prefix string, typeExpr ast.Expr) string {
		switch typeExpr.(type) {
		case *ast.StarExpr:
			dd := typeExpr.(*ast.StarExpr)
			// fmt.Println(dd)
			return prefix + "*" + p.printType(dd.X)
		case *ast.ArrayType:
			dd := typeExpr.(*ast.ArrayType)
			return prefix + "[]" + p.printType(dd.Elt)
		case *ast.InterfaceType:
			return "interface{}"
		case *ast.StructType:
			return "fuck" // 没啥用
		case *ast.MapType:
			dd := typeExpr.(*ast.MapType)
			keyType := p.printType(dd.Key)
			valueType := p.printType(dd.Value)

			return fmt.Sprintf("map<%v,%v>", keyType, valueType)

		default:
			flag := fmt.Sprintf("%v", typeExpr)
			if !strings.Contains(flag, " ") {
				return prefix + flag
			}

			a := regexp.MustCompile(`^&{(\w+) (\w+)}`)
			f := a.FindAllStringSubmatch(flag, -1)
			// fmt.Println(f)
			return f[0][1] + "." + f[0][2]
		}
	}

	if ty, ok := p.types[fmt.Sprintf("%v", t)]; ok {
		return ty
	}

	return fn("", t)
}

func (p *pkg) GetStructInfoByName(structName string) []fieldInfo {
	info, ok := p.structInfos[structName]
	if !ok {
		return nil
	}

	return p.GetStructInfo(info, structName)
}

func (p *pkg) GetAllStructName() []string {
	var names []string
	for k, _ := range p.structInfos {
		names = append(names, k)
	}

	return names
}

func (p *pkg) GetStructInfo(structInfo *ast.StructType, structName string) []fieldInfo {

	var fieldInfos []fieldInfo
	for _, v := range structInfo.Fields.List {
		if v.Tag == nil || p.parseTag(v.Tag.Value) == "" {
			continue
		}

		if strTag := p.parseTag(v.Tag.Value); strTag != "" && strTag != "-" {
			if strTag == ",inline" || strTag == "inline" {
				isOutter, pkgName, stcutName := p.parseType(v.Type)
				if isOutter {
					if otherPkg := p.imports[pkgName]; otherPkg != "" {
						if other := p.mgr.GetPkg(otherPkg, p.isOutterProject(otherPkg)); other != nil {
							// info.Type = p.printType(v.Type)
							// info.SubFieldInfo = other.GetStructInfoByName(strName)
							temps := other.GetStructInfoByName(stcutName)
							fieldInfos = append(fieldInfos, temps...)
						}
					}
				} else {
					fieldInfos = append(fieldInfos, p.GetStructInfoByName(stcutName)...)
				}
				continue
			}
			// if v.Names[0].Name == "StayTime" {
			// 	fmt.Println(1)
			// }
			info := fieldInfo{
				Name: v.Names[0].Name,
				// Type:    p.printType(v.Type),
				Tag:     strTag,
				Comment: p.parseComment(v),
			}

			if p.isGolangPrimitiveType(v.Type) {
				info.Type = p.printType(v.Type)
			} else {
				dependOther, pkgName, strName := p.parseType(v.Type)
				if !dependOther {
					switch v.Type.(type) {
					case *ast.StarExpr:
						dd := v.Type.(*ast.StarExpr)
						info.Type = p.printType(dd.X)
					case *ast.StructType:
						dd := v.Type.(*ast.StructType)
						info.SubFieldInfo = p.GetStructInfo(dd, v.Names[0].Name)
						info.Type = "object"
					case *ast.MapType:
						// dd := v.Type.(*ast.MapType)
						info.Type = p.printType(v.Type)
						// info.SubFieldInfo = p.parseNoNameStruct(v)
					default:
						if !p.isGolangPrimitiveType(v.Type) {
							if strName != structName {
								if temp, ok := p.structInfos[ /*v.Names[0].Name*/ strName]; ok {
									info.Type = p.printType(v.Type)
									info.SubFieldInfo = p.GetStructInfo(temp, strName)
								}
							} else {
								info.Type = "外层结构体类型(递归定义)"
							}
						}
						// info.SubFieldInfo = p.GetStructInfo(v.Names[0].Name)
						// info.Type = p.printType(v.Type)
					}

				} else {
					if otherPkg := p.imports[pkgName]; otherPkg != "" {
						if other := p.mgr.GetPkg(otherPkg, p.isOutterProject(otherPkg)); other != nil {
							info.Type = p.printType(v.Type)
							info.SubFieldInfo = other.GetStructInfoByName(strName)
						}
					}

				}
			}

			fieldInfos = append(fieldInfos, info)
		}
	}

	return fieldInfos
}

func (p *pkg) parseTag(tag string) string {
	tag = strings.Trim(tag, "`")
	flag := "bson"
	temps := strings.Split(tag, " ")
	for _, v := range temps {
		if strings.HasPrefix(v, flag+":") {
			temps2 := strings.Split(v, `:"`)
			temp := strings.Split(temps2[1], `"`)[0]
			return strings.TrimSuffix(temp, ",omitempty")
		}
	}

	return ""
}

func (p *pkg) parseComment(f *ast.Field) string {
	if f.Comment.Text() != "" {
		return strings.TrimSuffix(f.Comment.Text(), "\n")
	} else {
		return strings.TrimSuffix(f.Doc.Text(), "\n")
	}
}

func (p *pkg) isGolangPrimitiveType(typ ast.Expr) bool {

	switch typ.(type) {
	case *ast.ArrayType:

		dd := typ.(*ast.ArrayType)
		return p.isGolangPrimitiveType(dd.Elt)
		// ...

	}

	typeName := fmt.Sprintf("%v", typ)
	if _, ok := p.types[typeName]; ok {
		return true // 有问题
	}
	switch typeName {
	case "uint",
		"int",
		"uint8",
		"int8",
		"uint16",
		"int16",
		"byte",
		"uint32",
		"int32",
		"rune",
		"uint64",
		"int64",
		"float32",
		"float64",
		"bool",
		"string":
		return true
	default:
		return false
	}
}

func (p *pkg) Parse() {
	files := p.getAllGoFile(p.fullPath)
	for _, v := range files {
		p.parseFile(p.fullPath, p.fullPath+"/"+v, nil)
	}
}

func (p *pkg) getAllGoFile(path string) []string {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil
	}

	var fileList []string
	for _, v := range files {
		if !v.IsDir() && strings.HasSuffix(v.Name(), ".go") {
			fileList = append(fileList, v.Name())
		}
	}

	return fileList
}

func (p *pkg) parseFile(pkgPath string, filename string, src []byte) {
	var err error
	if src == nil {
		src, err = ioutil.ReadFile(filename)
		if err != nil {
			panic(err)
		}
	}
	fset := token.NewFileSet()

	file, err := parser.ParseFile(fset, filename, src, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	m := file.Scope.Objects
	for k, v := range m {
		if temp, ok := v.Decl.(*ast.ValueSpec); ok {
			for _, subv := range temp.Values {
				if temp2, ok := subv.(*ast.BinaryExpr); ok {
					if temp3, ok := temp2.Y.(*ast.BasicLit); ok {
						fmt.Println("iota -> ", temp3.Value)
					}

				} else if temp2, ok := subv.(*ast.BasicLit); ok {
					fmt.Println("cst, ", temp2.Value)
				}
			}

		}
		fmt.Printf("%v, %T\n", k, v.Decl)
	}
	fmt.Println(m)

	for _, v := range file.Imports {
		name := ""
		if v.Name == nil {
			temp := strings.Split(v.Path.Value, "/")
			name = temp[len(temp)-1]
		} else {
			name = v.Name.Name
		}

		p.imports[strings.Trim(name, `"`)] = strings.Trim(v.Path.Value, `"`)
	}

	fileBuf, _ := ioutil.ReadFile(filename)
	for _, v := range file.Decls {
		// fmt.Println("2222, ", v.Pos(), v, v.End())

		begin := int(v.Pos())
		end := int(v.End())
		src := string(fileBuf[begin-1 : end])

		b := file.Comments[0].Pos()
		c := file.Comments[0].End()
		fmt.Println("dddd, ", string(fileBuf[b-1:c]))

		if !strings.Contains(src, "\n") && !strings.Contains(src, "import ") {
			// fmt.Println("1111,", src)
			temps := strings.Split(src, " ")
			// fmt.Println(temps[2])
			temps[2] = strings.TrimSpace(temps[2])
			// 	fmt.Println(temps[2])
			p.types[temps[1]] = temps[2]
		}

		// switch v.(type) {
		// case *ast.DeclStmt:
		// 	fmt.Println(11)
		// }
		// if d, ok := v.(*ast.DeclStmt); ok {

		// }
	}

	var collectStructs func(x ast.Node) bool
	collectStructs = func(x ast.Node) bool {
		if _, ok := x.(*ast.ValueSpec); ok {
			// fmt.Println("解析到值定义 -> ", ts.Names, ts.Type, ts.Values[0])
			return true
		}
		ts, ok := x.(*ast.TypeSpec)
		if !ok || ts.Type == nil {
			return true
		}

		s, ok := ts.Type.(*ast.StructType)
		if !ok {
			return true
		}

		// 获取结构体名称
		structName := ts.Name.Name

		// fmt.Println(p.mod, " -> 解析到结构体:", structName)

		if _, ok := p.structInfos[structName]; ok {
			return true
		}

		p.structInfos[structName] = s

		for _, field := range s.Fields.List {
			if field.Tag != nil && p.parseTag(field.Tag.Value) != "" {
				if p.isGolangPrimitiveType(field.Type) {
					// fmt.Println("	1111 -> ddddd, ", field.Names[0].Name, field.Type, field.Tag.Value, field.Doc.Text())
				} else {
					dependOuterFile, pkgName, _ := p.parseType(field.Type)
					// fmt.Println("	2222 -> ddddd, ", field.Names[0].Name, field.Type, field.Tag.Value)
					if dependOuterFile {
						// fmt.Println(field.Type, "依赖外部包")
						// flags := strings.Split(fmt.Sprintf("%v", field.Type), " ")
						// structName := strings.Split(flags[1], "}")
						// ParsePkg("parseStruct2/model2")
						if otherPkg := p.imports[pkgName]; otherPkg != "" {
							p.mgr.OnDetectNewPkg(otherPkg, p.isOutterProject(otherPkg))
						} else {
							fmt.Println(p.fullPath, field.Names[0].Name, field.Type)
							fmt.Printf("警告！！！，未找到%v包\n", pkgName)
							p.parseType(field.Type)
						}
						// if dd, ok := g_def[os.Getenv("GOPATH")+"/src/parseStruct2/model2"][structName[0]]; ok {
						// 	for _, v := range dd.Fields.List {
						// 		fmt.Println("2222222 -> dddddd", v.Names[0].Name, v.Type, v.Tag.Value, field.Doc.Text())
						// 	}
						// }

						// fmt.Println("---->", ok, g_def[os.Getenv("GOPATH")+"/src/parseStruct2/model2"])
					} else {

					}
				}
			}

		}

		return true
	}

	ast.Inspect(file, collectStructs)
}

func (p *pkg) isOutterProject(pkgName string) bool {
	ours := strings.Split(p.mod, "/")
	others := strings.Split(pkgName, "/")
	if ours[0] == others[0] {
		return false
	} else {
		return true
	}
}

// 返回该字段是否是外部包、包名、结构体名
func (p *pkg) parseType(typeExpr ast.Expr) (bool, string, string) {
	// fmt.Println("ttttttt, ", typeExpr)
	switch expr := typeExpr.(type) {
	case *ast.StructType:
		// ddd := typeExpr.(*ast.StructType)
		// fmt.Println("	ssssssss, ", ddd.Fields.NumFields())
		return false, p.mod, ""

	case *ast.Ident:
		//expr := typeExpr.(*ast.Ident)
		// fmt.Println("	sssss222222, ", expr.String())
		return false, p.mod, expr.String()
	case *ast.StarExpr:
		dd := typeExpr.(*ast.StarExpr)
		// fmt.Println(dd)
		return p.parseType(dd.X)
	case *ast.SliceExpr:
		dd := typeExpr.(*ast.SliceExpr)
		// fmt.Println(dd)
		return p.parseType(dd.X)
	case *ast.ArrayType:
		dd := typeExpr.(*ast.ArrayType)
		// fmt.Println(dd)
		return p.parseType(dd.Elt)
	case *ast.InterfaceType:
		return false, p.mod, "interface{}"
	case *ast.MapType:
		dd := typeExpr.(*ast.MapType)
		return p.parseType(dd.Value)

	default:
		// fmt.Println("	default   -> ", typeExpr)
		flag := fmt.Sprintf("%v", typeExpr)

		a := regexp.MustCompile(`^&{(\w+) (\w+)}`)
		f := a.FindAllStringSubmatch(flag, -1)
		// fmt.Println(f)
		return true, f[0][1], f[0][2]

	}

	panic("shirt")

	return true, p.mod, ""
}
