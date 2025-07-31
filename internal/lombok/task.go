package lombok

import (
	"fmt"
	"github.com/heyuuu/go-lombok/internal/utils/mapkit"
	strkit2 "github.com/heyuuu/go-lombok/internal/utils/strkit"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
)

const generateExt = ".properties.go"

type statistic struct {
	unchanged int
	updated   int
	deleted   int
}

type taskType int

const (
	TaskGenerate taskType = iota
	TaskClear
)

func RunTask(task taskType, dir string) {
	basePkg := getNameFromModFile(dir)
	fmt.Println(basePkg)

	var stat statistic
	eachGoDir(dir, func(dirPath string, filePaths []string) {
		var dirPkg string
		if basePkg != "" && strings.HasPrefix(dirPath, dir) {
			dirPkg = basePkg + dirPath[len(dir):]
		}

		err := handleDir(dirPkg, dirPath, filePaths, task, &stat)
		if err != nil {
			log.Fatalln(err)
		}
	})
	log.Printf("处理完成. 共有更新文件 %d, 未变更文件 %d, 移除文件 %d\n", stat.updated, stat.unchanged, stat.deleted)
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

// 处理单个文件夹，不递归
func handleDir(dirPkg string, dirPath string, filePaths []string, task taskType, stat *statistic) error {
	// records old files
	srcFiles := make([]string, 0, len(filePaths))
	oldFiles := make(map[string]bool)
	for _, filePath := range filePaths {
		if isGenerateFile(filePath) {
			oldFiles[filePath] = true
		} else {
			srcFiles = append(srcFiles, filePath)
		}
	}

	// generate files
	if task == TaskGenerate { // parse pkg info
		pkg, err := ScanPkgInfo(dirPkg, srcFiles)
		if err != nil {
			return err
		}

		// show pkg info
		showPkgInfo(pkg)

		// generate file
		genCode, ok := GenFileCode(pkg)
		if ok {
			genFile := filepath.Join(dirPath, pkg.Name+generateExt)
			isChanged, err := writeFileIfChanged(genFile, genCode)
			if err != nil {
				return err
			}

			delete(oldFiles, genFile)
			if isChanged {
				stat.updated++
				log.Println("Update file: " + genFile)
			} else {
				stat.unchanged++
			}
		}
	}

	// remove old file
	for oldFile, _ := range oldFiles {
		err := os.Remove(oldFile)
		if err != nil {
			return err
		}
		stat.deleted++
		log.Println("Remove File: " + oldFile)
	}
	return nil
}

func showPkgInfo(pkg *PkgInfo) {
	var first = true
	for _, typ := range pkg.SortedTypes() {
		guessTags := make(map[string]string)
		properties := slices.DeleteFunc(typ.Properties(), func(prop *Property) bool {
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
				strkit2.PadRight(typ.Name, 20, ' '),
				strkit2.PadRight(prop.Name, 20, ' '),
				strkit2.PadRight(tag, 20, ' '),
				guessTag)
		}
	}
}

func tryGuessTag(prop *Property) (string, bool) {
	var getterMode int // 0: 未匹配，1: 默认模式, 2: 自定义模式, 3: 'Get' 前缀模式
	var setterMode int // 0: 未匹配，1: 默认模式, 2: 自定义模式
	var getterTag, setterTag string

	ucName := strkit2.UpperCamelCase(prop.Name)
	if prop.ExistGetters[ucName] {
		getterMode, getterTag = 1, `get:""`
	} else if prop.ExistGetters["Get"+ucName] {
		getterMode, getterTag = 3, `get:"@"`
	} else if getter, ok := mapkit.FirstKey(prop.ExistGetters); ok {
		getterMode, getterTag = 2, fmt.Sprintf(`get:"%s"`, getter)
	}

	if prop.ExistSetters["Set"+ucName] {
		setterMode, setterTag = 1, `set:""`
	} else if setter, ok := mapkit.FirstKey(prop.ExistSetters); ok {
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

func isGenerateFile(filepath string) bool {
	return strings.HasSuffix(filepath, generateExt)
}

func eachGoDir(dir string, handler func(dirPath string, filePaths []string)) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return
	}
	var filePaths []string
	for _, file := range files {
		name := file.Name()
		if strings.HasPrefix(name, "_") || strings.HasPrefix(name, ".") {
			continue
		}

		path := filepath.Join(dir, name)
		if file.IsDir() {
			eachGoDir(path, handler)
		} else if strings.HasSuffix(name, ".go") && !strings.HasSuffix(name, "_test.go") {
			filePaths = append(filePaths, path)
		}
	}
	if len(filePaths) > 0 {
		handler(dir, filePaths)
	}
}

func writeFileIfChanged(fileName string, content string) (changed bool, err error) {
	existContent, err := os.ReadFile(fileName)
	if err == nil && string(existContent) == content {
		return false, nil
	}

	err = os.WriteFile(fileName, []byte(content), 0644)
	return true, err
}
