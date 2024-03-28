package files

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"

	libstr "github.com/ZYallers/golib/funcs/strings"
)

func CurrentMethodName() string {
	pc, _, _, _ := runtime.Caller(1)
	name := runtime.FuncForPC(pc).Name()
	return libstr.StrFirstToLower(name[strings.LastIndex(name, ".")+1:])
}

func CurrentFileName() string {
	_, path, _, _ := runtime.Caller(1)
	file := path[strings.LastIndex(path, "/")+1:]
	return file[0:strings.Index(file, ".")]
}

func FileIsExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func GetFileNameWithoutExt(filePath string) string {
	// 获取文件名（包含扩展名）
	fileName := filepath.Base(filePath)
	// 分割文件名和扩展名
	fileNameWithoutExt := strings.TrimSuffix(fileName, filepath.Ext(fileName))
	return fileNameWithoutExt
}
