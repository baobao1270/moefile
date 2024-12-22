package server

import (
	"bytes"
	"encoding/json"
	"io"
	"io/fs"
	"net/url"
	"path/filepath"
	"slices"
	"strings"
	"text/template"

	"github.com/baobao1270/slang"

	"moefile/dist"
	"moefile/internal/log"
	"moefile/pkg/dto"
)

func (s *serverConfig) readFSDir(name string) ([]fs.DirEntry, error) {
	vfs, ok := s.rootFS.(fs.ReadDirFS)
	if !ok {
		log.T("server/xml").Errf("Unsupported filesystemc cast to fs.ReadDirFS: %T", s.rootFS)
	}

	return vfs.ReadDir(name)
}

func (s *serverConfig) statFSDir(name string) ([]fs.FileInfo, error) {
	files := make([]fs.FileInfo, 0)
	entries, err := s.readFSDir(name)
	if err != nil {
		return files, err
	}

	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			return files, err
		}
		files = append(files, info)
	}

	return files, nil
}

func (s *serverConfig) createS3XMLFromFSDir(url, path string) ([]byte, error) {
	res := dto.NewFSDirInfo(s.app.ServerName, strings.TrimPrefix(url, "/"))
	dir, err := s.statFSDir(path)
	if err != nil {
		return nil, err
	}

	for _, stat := range dir {
		res.AddFSFile(stat)
	}

	buf, err := res.ToS3XML(s.app.XMLIndent)
	if err != nil {
		log.T("server/xml").Errf("Unable to marshal ListBucketResult: %s", err)
		return nil, err
	}

	return buf, nil
}

func (s *serverConfig) searchPlayerData(requestURL string) (dto.PlayerData, error) {
	log.T("server/player/search").Dbgf("-------- enter searchPlayerData --------")
	if !strings.HasPrefix(requestURL, "/") {
		requestURL = "/" + requestURL
	}

	requestURL, err := url.PathUnescape(requestURL)
	if err != nil {

		return dto.PlayerData{}, err
	}

	path := strings.TrimPrefix(requestURL, "/")
	urlDir := strings.TrimSuffix(filepath.Dir(requestURL), "/")
	pathDir := strings.TrimSuffix(filepath.Dir(path), "/")
	baseName := filepath.Base(requestURL)
	title := strings.TrimSuffix(baseName, filepath.Ext(baseName))
	data := dto.PlayerData{
		Subs: make([]dto.PlayerSub, 0),
	}
	dir, err := s.readFSDir(pathDir)
	if err != nil {
		return data, nil
	}

	for _, entry := range dir {
		entryURL := filepath.Join(urlDir, entry.Name())
		if entry.IsDir() {
			continue
		}

		// danmaku: match <title>.<ext>[.danmaku|.danmuku|.comment|.comments].xml
		if matchPrefixSuffixList(entry.Name(), baseName,
			".danmaku.xml", ".danmuku.xml", ".comment.xml", ".comments.xml", ".xml") {
			data.DanmakuURL = entryURL
			continue
		}

		// danmaku: match <title>[.danmaku|.danmuku|.comment|.comments].xml
		if matchPrefixSuffixList(entry.Name(), title,
			".danmaku.xml", ".danmuku.xml", ".comment.xml", ".comments.xml", ".xml") {
			data.DanmakuURL = entryURL
			continue
		}

		// danmaku: match danmaku|danmuku|comment|comments.xml
		if slices.Contains([]string{"danmaku.xml", "danmuku.xml", "comment.xml", "comments.xml"}, entry.Name()) {
			data.DanmakuURL = entryURL
			continue
		}

		// subtitles: match <title>.<ext>*<lang>.srt
		if matchPrefixSuffixList(entry.Name(), baseName, ".srt") {
			data.Subs = append(data.Subs, createSub(baseName, entry.Name(), entryURL))
			continue
		}

		// subtitles: match <title>*<lang>.srt
		if matchPrefixSuffixList(entry.Name(), title, ".srt") {
			data.Subs = append(data.Subs, createSub(title, entry.Name(), entryURL))
			continue
		}
	}

	log.T("server/player/search").Dbgf("PlayerData: %+v", data)
	log.T("server/player/search").Dbgf("-------- leave searchPlayerData --------")
	return data, nil
}

func matchSuffixList(name string, suffixList ...string) bool {
	for _, suffix := range suffixList {
		if strings.HasSuffix(name, suffix) {
			return true
		}
	}
	return false
}

func matchPrefixSuffixList(name string, prefix string, suffixList ...string) bool {
	if !strings.HasPrefix(name, prefix) {
		return false
	}
	return matchSuffixList(name, suffixList...)
}

var subLangFilterOut = []rune{
	'?', '@', '|', ' ',
	'[', ']', '(', ')', '{', '}', '<', '>',
	'？', '｜',
	'【', '】', '（', '）', '「', '」', '『', '』',
}

func defaultSub(url string) dto.PlayerSub {
	return dto.PlayerSub{
		Lang:     "und",
		URL:      url,
		LangName: "CC",
		LangInfo: slang.Lang{
			Name:       "und",
			BCP47:      "zz",
			WinID:      "ZZZ",
			ISO639Set1: "zz",
			ISO639Set2: "und",
			ISO639Set3: "und",
		},
	}
}

func getSubLangCode(prefix, name string) string {
	result := strings.TrimPrefix(name, prefix)
	result = strings.TrimSuffix(result, ".srt")
	result = strings.Trim(result, ".- ")
	code := ""
	for _, c := range result {
		if c == '.' || c == '_' {
			c = '-'
		}
		if slices.Contains(subLangFilterOut, c) {
			continue
		}
		code += string(c)
	}
	code = strings.Trim(code, " ")
	return code
}

func createSub(prefix, name, url string) dto.PlayerSub {
	langCode := getSubLangCode(prefix, name)

	parser, err := slang.NewParser()
	if err != nil || langCode == "" {
		return defaultSub(url)
	}

	lang := parser.Parse(langCode)
	if lang == nil {
		return defaultSub(url)
	}

	langName := lang.WinID
	if !lang.IsValidWinID() {
		langName = strings.ToUpper(lang.ISO639Set3)
	}

	return dto.PlayerSub{
		Lang:     lang.BCP47,
		URL:      url,
		LangName: strings.ToUpper(langName),
		LangInfo: *lang,
	}
}

func renderPlayerData(data dto.PlayerData) (io.Reader, error) {
	tmtmplBuf, err := dist.Embed.ReadFile("player.html")
	if err != nil {
		log.T("server/player").Errf("Unable to read player.html: %s", err)
		return nil, err
	}

	dataBuf, err := json.Marshal(data)
	if err != nil {
		log.T("server/player").Errf("Unable to marshal danmaku and subtitles (PlayerData): %v", err)
	}

	tmpl, err := template.New("player.html").Parse(string(tmtmplBuf))
	if err != nil {
		log.T("server/player").Errf("Unable to create template with player.html: %v", err)
		return nil, err
	}

	outBuf := new(bytes.Buffer)
	err = tmpl.Execute(outBuf, string(dataBuf))
	if err != nil {
		log.T("server/player").Errf("Unable to render player.html with data <%s>: %v", string(dataBuf), err)
		return nil, err
	}

	return outBuf, nil
}
