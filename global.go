package main

import (
	"fmt"
	"os"
	"strings"
)

type global struct {
	// 所有包
	pkgs  map[string]*pkg
	goMod map[string]string
}

func (g *global) OnDetectNewPkg(mod string, isOutterProject bool) {
	if !isOutterProject {
		fullPath := os.Getenv("GOPATH") + "/src/" + mod
		p := NewPkg(g, mod, fullPath)
		g.pkgs[fullPath] = p
		p.Parse()
	} else {
		for k, v := range g.goMod {
			if strings.HasPrefix(mod, k) {
				modName := strings.Split(v, "@")
				temp1 := strings.TrimPrefix(mod, modName[0])
				fullPath := os.Getenv("GOPATH") + "/pkg/mod/" + v + temp1
				p := NewPkg(g, mod, fullPath)
				g.pkgs[fullPath] = p
				fmt.Println(fullPath)
				p.Parse()

				return
			}
		}

		fmt.Printf("警告！！！未能找到%v包, %v\n", mod, isOutterProject)

	}
}

func (g *global) GetPkg(mod string, isOutterProject bool) *pkg {
	if !isOutterProject {
		p, ok := g.pkgs[os.Getenv("GOPATH")+"/src/"+mod]
		if !ok {
			fmt.Printf("警告！！！未能找到%v包,%v\n", mod, isOutterProject)
			return nil
		}

		return p
	} else {
		for k, v := range g.goMod {
			if strings.HasPrefix(mod, k) {
				modName := strings.Split(v, "@")
				temp1 := strings.TrimPrefix(mod, modName[0])
				fullPath := os.Getenv("GOPATH") + "/pkg/mod/" + v + temp1
				p, ok := g.pkgs[fullPath]
				if !ok {
					fmt.Printf("警告！！！未能找到%v包,%v\n", mod, isOutterProject)
					return nil
				}

				return p
			}
		}

		fmt.Printf("警告！！！未能找到%v包,%v\n", mod, isOutterProject)
		return nil
	}
}
