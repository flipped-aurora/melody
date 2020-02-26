package plugin

import (
	"io/ioutil"
	"strings"
)

// Scan ... 读取插件
func Scan(folder, pattern string) ([]string, error) {
	files, erro := ioutil.ReadDir(folder)
	if erro != nil {
		return []string{}, erro
	}

	plugins := []string{}
	for _, file := range files {
		if !file.IsDir() && strings.Contains(file.Name(), pattern) {
			plugins = append(plugins, folder+file.Name())
		}
	}

	return plugins, nil
}
