package viteutil

import (
	"bytes"
	"text/template"

	"moefile/res/tmpl"
)

var RootPathDevMapping = map[bool]string{
	true:  "/",
	false: "https://cdn.local/?_/",
}

type LayoutData struct {
	HeaderContent string
	RootPath      string
	AppName       string
}

func TmplInsert(tmplFile string, innerData any) string {
	tmplText, err := tmpl.TmplFS.ReadFile(tmplFile)
	if err != nil {
		panic(err)
	}

	tmpl, err := template.New(tmplFile).Parse(string(tmplText))
	if err != nil {
		panic(err)
	}

	buf := new(bytes.Buffer)
	err = tmpl.Execute(buf, innerData)
	if err != nil {
		panic(err)
	}
	return buf.String()
}

func RLayout(data LayoutData) string {
	return TmplInsert("_layout.html", data)
}

func RAppSrc(name string, isDev bool, appData any) string {
	return RLayout(LayoutData{
		HeaderContent: TmplInsert(GetHTMLFileName(name, isDev), appData),
		RootPath:      RootPathDevMapping[isDev],
		AppName:       name,
	})
}
