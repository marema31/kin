package server

import (
	"net/http"
	"os"
	"path"
	"text/template"

	"github.com/sirupsen/logrus"
)

func templatedResponseFile(logReq *logrus.Entry, response http.ResponseWriter, tmplPath string) {
	tmplName := path.Base(tmplPath)
	tmplt, err := template.New(tmplName).ParseFiles(tmplPath)

	if err != nil {
		logReq.Errorf("Error parsing template: %v", err)
		http.Error(response, "Internal error", http.StatusInternalServerError)

		return
	}

	ci, err := prepareData(logReq)
	if err != nil {
		http.Error(response, "Internal error", http.StatusInternalServerError)
		return
	}

	err = tmplt.Execute(response, ci)
	if err != nil {
		logReq.Errorf("Error rendering template: %v", err)
		http.Error(response, "Internal error", http.StatusInternalServerError)

		return
	}

	logReq.Debug("templated")
}

func rawResponseFile(logReq *logrus.Entry, response http.ResponseWriter, request *http.Request, path string) {
	logReq.Debug("raw")
	http.ServeFile(response, request, path)
}

func handleResponseFile(response http.ResponseWriter, request *http.Request) {
	path, logReq := preParseRequest(request)

	path = root + path

	fi, err := os.Stat(path)
	if err == nil && !fi.IsDir() {
		rawResponseFile(logReq, response, request, path)
		return
	}

	fi, err = os.Stat(path + ".tpl")
	if err == nil && !fi.IsDir() {
		templatedResponseFile(logReq, response, path+".tpl")
		return
	}

	logReq.Error("not found")
	http.NotFound(response, request)
}
