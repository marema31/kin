package server

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
)

type containerInfo struct {
	Name string
	URL  string
}

var (
	base string
	log  *logrus.Entry
	root string
)

func retrieveData() []containerInfo {
	//TODO: query the data1
	return []containerInfo{
		{Name: "Mon Site 1", URL: "http://localhost/1"},
		{Name: "Mon Site 2", URL: "http://localhost/2"},
		{Name: "Mon Site 3", URL: "http://localhost/3"},
	}
}

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
func Run(applog *logrus.Entry, baseURL string, rootPath string, port int) error {
	applog.Infof("Starting web server on port %d for %s", port, base)

	base = strings.TrimSuffix(baseURL, "/")
	log = applog
	root = rootPath

	if root != "" {
		http.HandleFunc(baseURL, handleResponseFile)
	} else {
		http.HandleFunc(baseURL, handleResponsePkger)
	}

	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
