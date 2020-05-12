package server

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"
	"text/template"

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

func templatedResponse(logReq *logrus.Entry, response http.ResponseWriter, tmplPath string) {
	//TODO: query the data1
	data := []containerInfo{
		{Name: "Mon Site 1", URL: "http://localhost/1"},
		{Name: "Mon Site 2", URL: "http://localhost/2"},
		{Name: "Mon Site 3", URL: "http://localhost/3"},
	}

	tmplName := path.Base(tmplPath)
	tmplt, err := template.New(tmplName).ParseFiles(tmplPath)
	if err != nil {
		logReq.Errorf("Error parsing template: %v", err)
		http.Error(response, "Internal error", http.StatusInternalServerError)

		return
	}

	err = tmplt.Execute(response, data)
	if err != nil {
		logReq.Errorf("Error rendering template: %v", err)
		http.Error(response, "Internal error", http.StatusInternalServerError)

		return
	}

	logReq.Debug("templated")
}

func rawResponse(logReq *logrus.Entry, response http.ResponseWriter, request *http.Request, path string) {
	logReq.Debug("raw")
	http.ServeFile(response, request, path)
}

func handleResponse(response http.ResponseWriter, request *http.Request) {
	path := request.URL.Path
	path = strings.TrimPrefix(path, base)

	logReq := log.WithField("path", path)

	if strings.HasSuffix(path, "/") {
		path += "index.html"
	}

	path = root + path

	fi, err := os.Stat(path)
	if err == nil && !fi.IsDir() {
		rawResponse(logReq, response, request, path)
		return
	}

	fi, err = os.Stat(path + ".tpl")
	if err == nil && !fi.IsDir() {
		templatedResponse(logReq, response, path+".tpl")
		return
	}

	logReq.Error("not found")
	http.NotFound(response, request)
}

func Run(applog *logrus.Entry, baseURL string, rootPath string, port int) error {
	applog.Infof("Starting web server on port %d for %s", port, base)

	base = strings.TrimSuffix(baseURL, "/")
	log = applog
	root = rootPath

	http.HandleFunc(baseURL, handleResponse)

	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
