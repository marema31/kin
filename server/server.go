package server

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/marema31/kin/cache"
	"github.com/sirupsen/logrus"
)

var (
	base     string
	database *cache.Cache
	log      *logrus.Entry
	root     string
)

func preParseRequest(request *http.Request) (string, *logrus.Entry) {
	path := request.URL.Path
	path = "/site" + strings.TrimPrefix(path, base)

	logReq := log.WithField("path", path)

	if strings.HasSuffix(path, "/") {
		path += "index.html"
	}

	return path, logReq
}

//Run start the kin http server.
func Run(applog *logrus.Entry, db *cache.Cache, baseURL string, rootPath string, port int) error {
	applog.Infof("Starting web server on port %d for %s", port, base)

	base = strings.TrimSuffix(baseURL, "/")
	database = db
	log = applog
	root = rootPath

	if root != "" {
		http.HandleFunc(baseURL, handleResponseFile)
	} else {
		http.HandleFunc(baseURL, handleResponsePkger)
	}

	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
