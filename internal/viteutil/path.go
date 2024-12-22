package viteutil

import (
	"fmt"
	"os"
)

const (
	ViteSrcDir    = "./src"
	ViteOutputDir = "./dist"
)

func GetHTMLFileName(name string, isDev bool) string {
	return fmt.Sprintf("%s.%s.html", name, map[bool]string{true: "dev", false: "prod"}[isDev])
}

func WriteAppSrc(name string, isDev bool, appData any) {
	fileName := fmt.Sprintf("%s/%s/%s.html", ViteSrcDir, name, name)
	fileContent := []byte(RAppSrc(name, isDev, appData))
	println(fmt.Sprintf("Writing: %s [%d bytes]", fileName, len(fileContent)))
	err := os.WriteFile(fileName, fileContent, 0644)
	if err != nil {
		panic(err)
	}
}
