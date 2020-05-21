package server

import (
	"context"
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
	server   *http.Server
	root     string
)

type tplData map[string][]cache.ContainerInfo

func prepareData(log *logrus.Entry) (tplData, error) {
	data := tplData{}

	ci, err := database.RetrieveData(log)
	if err != nil {
		return data, err
	}

	for _, container := range ci {
		if _, ok := data[container.Group]; !ok {
			data[container.Group] = make([]cache.ContainerInfo, 0)
		}

		data[container.Group] = append(data[container.Group], container)
	}

	return data, nil
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
func Run(ctx context.Context, applog *logrus.Entry, db *cache.Cache, baseURL string, rootPath string, port int) error {
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

	server = &http.Server{Addr: fmt.Sprintf(":%d", port)}

	return server.ListenAndServe()
}

//Shutdown stop the kin http server.
func Shutdown(ctx context.Context) error {
	return server.Shutdown(ctx)
}
