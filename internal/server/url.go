package server

import (
	"path"
	"path/filepath"

	"moefile/internal/log"
)

type urlInfo struct {
	requestURL string
	relPath    string
	ok         bool
}

func resolve(url string) (result urlInfo) {
	log.T("url").Dbgf("HTTP_URL_REQUEST:  %s", url)

	url = path.Clean("./" + url)
	log.T("url").Dbgf(" ->  URL_CLEAN:    %s", url)

	url, err := filepath.Localize(url)
	log.T("url").Dbgf(" ->  URL_PLATFORM: %s", url)
	if err != nil {
		log.T("url").Errf("Invalid URL in request <%s>: %s", url, err)
		return
	}

	cleanURL := filepath.Join("/", url)
	log.T("url").Dbgf(" ->  URL_REAL:     %s", cleanURL)

	result = urlInfo{
		requestURL: cleanURL,
		relPath:    url,
		ok:         true,
	}
	return
}
