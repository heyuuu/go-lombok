package lombok

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
)

const genFileName = "properties.gen.go"

// 通过代码目录遍历go包并调用回调
func iterPkgFiles(root string, handler func(dir string, srcFiles []string)) {
	files, err := os.ReadDir(root)
	if err != nil {
		return
	}

	var srcFiles []string
	for _, file := range files {
		name := file.Name()
		if name == "" || name[0] == '_' || name[0] == '.' {
			continue
		}

		path := filepath.Join(root, name)
		if file.IsDir() {
			iterPkgFiles(path, handler)
		} else if strings.HasSuffix(name, ".go") && !strings.HasSuffix(name, "_test.go") && !strings.HasSuffix(name, ".gen.go") {
			srcFiles = append(srcFiles, path)
		}
	}
	if len(srcFiles) > 0 {
		handler(root, srcFiles)
	}
}

// GenerateByCode 基于代码字符串的扫描和生成，主要用于单元测试
func GenerateByCode(pkgName string, code string) (string, error) {
	pkg, err := ScanCode(pkgName, code)
	if err != nil {
		return "", err
	}

	result := GenFileCode(pkg)
	return result, nil
}

// Clear 清理生成文件
func Clear(root string) {
	deleted := 0
	iterPkgFiles(root, func(dir string, _ []string) {
		genFile := filepath.Join(dir, genFileName)
		exists, err := deleteFileIfExists(genFile)
		if err != nil {
			log.Fatalln(err)
		}

		if exists {
			deleted++
			log.Println("Remove File: " + genFile)
		}
	})
	log.Printf("处理完成. 移除文件 %d\n", deleted)
}

type statistic struct {
	unchanged int
	updated   int
	deleted   int
}

// Generate 基于代码目录的扫描、生成、清理
func Generate(root string) {
	basePkg := getNameFromModFile(root)
	fmt.Println(basePkg)

	var stat statistic
	iterPkgFiles(root, func(dir string, srcFiles []string) {
		var dirPkg string
		if basePkg != "" && strings.HasPrefix(dir, root) {
			dirPkg = basePkg + dir[len(root):]
		}

		err := handlePkg(dirPkg, dir, srcFiles, &stat)
		if err != nil {
			log.Fatalln(err)
		}
	})
	log.Printf("处理完成. 共有更新文件 %d, 未变更文件 %d, 移除文件 %d\n", stat.updated, stat.unchanged, stat.deleted)
}

// 处理单个包(即单个文件夹)，不处理子包
func handlePkg(pkgName string, dir string, srcFiles []string, stat *statistic) error {
	// 扫描源代码文件，生成目标代码
	genCode, err := genPkgCode(pkgName, srcFiles)
	if err != nil {
		return err
	}

	// 更新或删除生成文件
	genFile := filepath.Join(dir, genFileName)
	if genCode != "" { // 有生成代码时，创建或更新文件
		changed, err := writeFileIfChanged(genFile, genCode)
		if err != nil {
			return fmt.Errorf("写入文件异常: file=%s, err=%w", genFile, err)
		}

		if changed {
			stat.updated++
			log.Println("Update file: " + genFile)
		} else {
			stat.unchanged++
		}
	} else { // genCode == ""，没有生成代码时，尝试删除文件
		exists, err := deleteFileIfExists(genFile)
		if err != nil {
			return fmt.Errorf("删除文件异常: file=%s, err=%w", genFile, err)
		}

		if exists {
			stat.deleted++
			log.Println("Remove File: " + genFile)
		}
	}
	return nil
}

func genPkgCode(pkgName string, srcFiles []string) (string, error) {
	// 扫描包信息
	pkg, err := ScanPkgInfo(pkgName, srcFiles)
	if err != nil {
		return "", err
	}

	// show pkg info
	showPkgInfo(pkg)

	// 生成文件代码
	code := GenFileCode(pkg)
	return code, nil
}

func getNameFromModFile(dir string) string {
	modFile := filepath.Join(dir, "go.mod")
	data, err := os.ReadFile(modFile)
	if err != nil {
		return ""
	}

	if match := regexp.MustCompile(`module ([\w./]+)`).FindSubmatch(data); len(match) > 0 {
		return string(match[1])
	}
	return ""
}

func showPkgInfo(pkg *PkgInfo) {
	var first = true
	for _, typ := range pkg.SortedTypes() {
		guessTags := make(map[string]string)
		properties := slices.DeleteFunc(slices.Collect(typ.Properties()), func(prop *Property) bool {
			guessTag, _ := tryGuessTag(prop)
			guessTags[prop.Name] = guessTag
			return guessTag == "" || guessTag == prop.Tag
		})
		if len(properties) == 0 {
			continue
		}
		if first {
			first = false
			fmt.Printf("package %s\n", pkg.Pkg)
		}

		fmt.Printf("type %s: recv=%s\n", typ.Name, typ.RecvName)
		for _, prop := range properties {
			guessTag := guessTags[prop.Name]
			tag := prop.Tag
			if tag == "" {
				tag = "-"
			}
			fmt.Printf("    %s.%s %s => %s\n",
				padRight(typ.Name, 20, ' '),
				padRight(prop.Name, 20, ' '),
				padRight(tag, 20, ' '),
				guessTag)
		}
	}
}

func tryGuessTag(prop *Property) (string, bool) {
	var getterMode int // 0: 未匹配，1: 默认模式, 2: 自定义模式, 3: 'Get' 前缀模式
	var setterMode int // 0: 未匹配，1: 默认模式, 2: 自定义模式
	var getterTag, setterTag string

	ucName := pascalCase(prop.Name)
	if prop.ExistsGetter(ucName) {
		getterMode, getterTag = 1, `get:""`
	} else if prop.ExistsGetter("Get" + ucName) {
		getterMode, getterTag = 3, `get:"@"`
	} else if getter, ok := firstOf(prop.ExistingGetters()); ok {
		getterMode, getterTag = 2, fmt.Sprintf(`get:"%s"`, getter)
	}

	if prop.ExistsSetter("Set" + ucName) {
		setterMode, setterTag = 1, `set:""`
	} else if setter, ok := firstOf(prop.ExistingSetters()); ok {
		setterMode, setterTag = 2, fmt.Sprintf(`set:"%s"`, setter)
	}

	var tag string
	if getterMode == 0 && setterMode == 0 {
		return "", false
	} else if getterMode == 1 && setterMode == 1 {
		tag = `prop:""`
	} else if getterMode == 3 && setterMode == 1 {
		tag = `prop:"@"`
	} else {
		if getterTag == "" {
			tag = setterTag
		} else if setterTag == "" {
			tag = getterTag
		} else {
			tag = getterTag + " " + setterTag
		}
	}

	return "`" + tag + "`", true
}
