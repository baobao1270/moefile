package dto

import (
	"bytes"
	"encoding/xml"
	"text/template"
	"time"

	"moefile/dist"
	"moefile/res"
)

type ListBucketResult struct {
	DirInfo
	ServerTZ string `xml:"ServerTimezoneOffset"`
	XSLT     string `xml:",innerxml"`
}

func (i *DirInfo) ToS3XMLWithXSLT(indent bool, xslt string) (data []byte, err error) {
	result := ListBucketResult{
		DirInfo:  *i,
		ServerTZ: time.Now().Format("-07:00"),
		XSLT:     xslt,
	}

	if indent {
		data, err = xml.MarshalIndent(result, "", "\t")
	} else {
		data, err = xml.Marshal(result)
	}
	return
}

func (i *DirInfo) ToS3XMLWithoutXSLT(indent bool) (data []byte, err error) {
	return i.ToS3XMLWithXSLT(indent, "")
}

func (i *DirInfo) ToS3XML(indent bool) (data []byte, err error) {
	t, err := template.New("xslt").Parse(res.XSLT)
	if err != nil {
		return nil, err
	}

	c, err := dist.Embed.ReadFile("index.html")
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	err = t.Execute(buf, string(c))
	if err != nil {
		return nil, err
	}

	data, err = i.ToS3XMLWithXSLT(indent, buf.String())
	if err != nil {
		return nil, err
	}
	data = append(res.XMLHeader, data...)
	return data, err
}
