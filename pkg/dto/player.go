package dto

import "github.com/baobao1270/slang"

type PlayerData struct {
	DanmakuURL string      `json:"danmaku"`
	Subs       []PlayerSub `json:"subtitles"`
}

type PlayerSub struct {
	Lang     string     `json:"lang"`
	LangName string     `json:"lang_name"`
	LangInfo slang.Lang `json:"lang_info"`
	URL      string     `json:"url"`
}
