package main

import (
	"fmt"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/austien/file-server/files"
	"github.com/austien/file-server/http"
	"github.com/austien/file-server/watcher"
)

func main() {
	logLevel, ok := os.LookupEnv("LOG_LEVEL")
	if ok {
		switch logLevel {
		case "info", "INFO":
			log.SetLevel(log.InfoLevel)
		case "debug", "DEBUG":
			log.SetLevel(log.DebugLevel)
		default:
			log.Warnf("Unknown log level %q. Defaulting to info", logLevel)
		}
	} else {
		log.SetLevel(log.InfoLevel)
	}

	rootFolder, ok := os.LookupEnv("FILE_SERVER_ROOT_FOLDER")
	if !ok {
		log.Fatal("missing envvar for media folder")
	}
	log.WithField("path", rootFolder).Info("Using root folder")

	host, ok := os.LookupEnv("FILE_SERVER_HOST")
	if !ok {
		log.Fatal("missing envvar for host")
	}
	log.WithField("host", host).Info("Using host")

	fileTypes, ok := os.LookupEnv("FILE_SERVER_WHITELISTED_EXTENSIONS")
	if ok {
		files.VaildFiletypes = strings.Split(fileTypes, ",")
	}
	log.WithField("extensions", files.VaildFiletypes).Info("Using whitelisted extensions")

	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8000"
	}
	log.WithField("port", port).Info("Using port")

	client := files.NewClient(fmt.Sprintf("%s:%s", host, port), rootFolder)

	if err := client.Init(rootFolder); err != nil {
		log.Fatal(err)
	}

	w, err := watcher.NewWatcher(rootFolder)
	if err != nil {
		log.Fatal(err)
	}
	defer w.Close()

	folderEvent := make(chan string)
	errChan := make(chan error)

	go w.Watch(folderEvent, errChan)

	go func() {
		for {
			select {
			case f := <-folderEvent:
				if err := client.HandleFolder(f); err != nil {
					log.Fatal(err)
				}
			case err := <-errChan:
				log.Fatal(err)
			}
		}
	}()

	log.Infof("Running file-server at port %s", port)

	panic(http.NewServer(fmt.Sprintf(":%s", port), client, rootFolder).ListenAndServe())
}
