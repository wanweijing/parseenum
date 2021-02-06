package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type fieldInfo struct {
	Name         string      `json:"名字"`
	Type         string      `json:"类型"`
	Tag          string      `json:"tag"`
	Comment      string      `json:"注释"`
	SubFieldInfo []fieldInfo `json:"子字段"`
}

func print(offset int, info []fieldInfo) {
	prefix := ""
	for i := 0; i < offset; i++ {
		prefix += "	"
	}

	for _, v := range info {
		fmt.Printf("%s 名字：%v, 类型：%v, tag：%v, 注释:%v\n", prefix, v.Name, v.Type, v.Tag, v.Comment)
		if len(v.SubFieldInfo) > 0 {
			print(offset+1, v.SubFieldInfo)
		}
	}
}

func fmtFieldInfo(offset int, info []fieldInfo) string {
	if len(info) <= 0 {
		return ""
	}

	ret := ""
	prefix := "#### "
	for i := 0; i < offset; i++ {
		prefix += "&emsp;"
	}

	// ret += prefix

	for _, v := range info {
		ret += prefix + fmt.Sprintf("字段名：%v, 类型：%v, 说明:%v\n", v.Tag, v.Type, v.Comment)
		if len(v.SubFieldInfo) > 0 {
			ret += prefix + "{\n"
			ret += fmtFieldInfo(offset+1, v.SubFieldInfo)
			ret += prefix + "}\n"
		}
	}

	// buf, _ := json.MarshalIndent(info, "	", "")

	return ret
}

func parseMod(pName string) map[string]string {
	temps := strings.Split(pName, "/")
	path := os.Getenv("GOPATH") + "/src/" + temps[0] + "/go.mod"
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	strs := strings.Split(string(buf), "require (\n")
	strs = strings.Split(strs[1], "\n)")

	strs = strings.Split(strs[0], "\n")
	rets := make(map[string]string)
	for _, v := range strs {
		temps := strings.Split(v, "//")
		temp := strings.TrimSpace(temps[0])

		temps = strings.Split(temp, " ")
		rets[temps[0]] = strings.Replace(temp, " ", "@", -1)
	}

	return rets
}

func getAllModPath(dirName string) []string {

	files, err := ioutil.ReadDir(dirName)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	var fileList []string
	for _, file := range files {
		if file.IsDir() {
			if file.Name() == "model" {
				fileList = append(fileList, dirName+string(os.PathSeparator)+file.Name())
			}

			fileList = append(fileList, getAllModPath(dirName+string(os.PathSeparator)+file.Name())...)
		}
	}

	var rets []string
	for _, v := range fileList {
		prefix := strings.Replace(os.Getenv("GOPATH")+"/src/", "\\", "/", -1)
		temp := strings.TrimPrefix(v, prefix)
		rets = append(rets, temp)
	}

	return rets

}

func getAllProjectName() []string {
	dirName := os.Getenv("GOPATH") + "/src"
	files, err := ioutil.ReadDir(dirName)
	if err != nil {
		fmt.Println(err)
	}
	var fileList []string
	for _, file := range files {
		if file.IsDir() {
			fileList = append(fileList, file.Name())
		}
	}

	return fileList
}

type pkgInfo struct {
	pkgName string
	buf     string
}

var g_structStr = map[string]pkgInfo{}

func parseProject(proName string) {
	pkgs := getAllModPath(fmt.Sprintf("C:/work/go/src/%v/pkg", proName))
	if len(pkgs) <= 0 {
		return
	}
	// 	pkgs := []string{"parseStruct2/model1"}

	for _, waitParse := range pkgs {
		g := &global{
			pkgs:  make(map[string]*pkg),
			goMod: parseMod(waitParse),
		}

		fullPath := os.Getenv("GOPATH") + "/src/" + waitParse
		p := NewPkg(g, waitParse, fullPath)
		g.pkgs[fullPath] = p
		p.Parse()

		structNames := p.GetAllStructName()

		fmt.Printf("包%v\n", waitParse)
		for _, v := range structNames {

			// if v != "RoleLimitItem" && v != "Menu" {
			fmt.Println("	============", v)
			if v == "StayTime" {
				fmt.Println(1)
			}
			info := p.GetStructInfoByName(v)
			// print(1, info)
			if buf := fmtFieldInfo(1, info); buf != "" {
				g_structStr[v] = pkgInfo{
					pkgName: waitParse,
					buf:     buf,
				}
			}

			// fmt.Println(buf)
			// }
		}
	}
}

func main() {

	// dirs := getAllProjectName()
	// for _, v := range dirs {
	// 	parseProject(v)
	// }

	pkgs := []string{"parseenum/model1"}

	for _, waitParse := range pkgs {
		g := &global{
			pkgs:  make(map[string]*pkg),
			goMod: parseMod(waitParse),
		}

		fullPath := os.Getenv("GOPATH") + "/src/" + waitParse
		p := NewPkg(g, waitParse, fullPath)
		g.pkgs[fullPath] = p
		p.Parse()

		structNames := p.GetAllStructName()

		fmt.Printf("包%v\n", waitParse)
		for _, v := range structNames {
			fmt.Println("	============", v)
			// info := p.GetStructInfoByName(v)
			// print(1, info)

			info := p.GetStructInfoByName(v)
			print(1, info)
			if buf := fmtFieldInfo(1, info); buf != "" {
				g_structStr[v] = pkgInfo{
					pkgName: waitParse,
					buf:     buf,
				}
			}
		}
	}

	fileBuf := ""

	fmt.Println("解析结果:")
	index := 0
	for k, v := range g_structStr {
		tabBuf := "## " + k + fmt.Sprintf(" (%v)", strings.Replace(v.pkgName, "\\", "/", -1))
		tabBuf = fmt.Sprintln(tabBuf)
		tabBuf += "#### {\n"
		// fmt.Printf("-> 结构体%v\n", k)
		// fmt.Println(v)
		tabBuf += v.buf
		tabBuf += "#### }"
		tabBuf += fmt.Sprintln("")
		for i := 0; i < 3; i++ {
			tabBuf += fmt.Sprintln("<br>")
		}
		tabBuf += fmt.Sprintln("")
		fileBuf += tabBuf
		index++
		// if index >= 2 {
		// 	break
		// }
	}

	fmt.Println("总表数: ", index)

	ioutil.WriteFile("model.md", []byte(fileBuf), 0644)
}
