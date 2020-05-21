package server

import (
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"text/template"

	"github.com/markbates/pkger"
	"github.com/markbates/pkger/pkging"
	"github.com/sirupsen/logrus"
)

func openPkger(logReq *logrus.Entry, path string) (pkging.File, os.FileInfo, error) {
	fi, err := pkger.Stat(path)
	if err != nil {
		logReq.Errorf("Error getting stats of file: %v", err)
		return nil, fi, err
	}

	file, err := pkger.Open(path)
	if err != nil {
		logReq.Errorf("Error opening file: %v", err)
		return nil, fi, err
	}

	return file, fi, err
}

func templatedResponsePkger(logReq *logrus.Entry, response http.ResponseWriter, tmplPath string) {
	tmplName := path.Base(tmplPath)

	file, _, err := openPkger(logReq, tmplPath)
	if err != nil {
		http.Error(response, "Internal error", http.StatusInternalServerError)
		return
	}

	defer file.Close()

	tmplContent, err := ioutil.ReadAll(file)
	if err != nil {
		logReq.Errorf("Error reading template: %v", err)
		http.Error(response, "Internal error", http.StatusInternalServerError)

		return
	}

	tmplt, err := template.New(tmplName).Parse(string(tmplContent))
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

func rawResponsePkger(logReq *logrus.Entry, response http.ResponseWriter, request *http.Request, path string) {
	logReq.Debug("raw")

	file, fi, err := openPkger(logReq, path)
	if err != nil {
		http.Error(response, "Internal error", http.StatusInternalServerError)
		return
	}

	defer file.Close()

	http.ServeContent(response, request, filepath.Base(path), fi.ModTime(), file)
}

func handleResponsePkger(response http.ResponseWriter, request *http.Request) {
	path, logReq := preParseRequest(request)

	fi, err := pkger.Stat(path)
	if err == nil && !fi.IsDir() {
		rawResponsePkger(logReq, response, request, path)
		return
	}

	fi, err = pkger.Stat(path + ".tpl")
	if err == nil && !fi.IsDir() {
		templatedResponsePkger(logReq, response, path+".tpl")
		return
	}

	logReq.Error("not found")
	http.NotFound(response, request)
}
