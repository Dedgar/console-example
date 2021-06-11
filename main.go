package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/dedgar/console-example/downloader"
	"github.com/dedgar/console-example/routers"
	"github.com/dedgar/console-example/tree"
)

var (
	certFile     = os.Getenv("CERT_FILE")
	keyFile      = os.Getenv("KEY_FILE")
	downloadURL  = os.Getenv("DOWNLOAD_URL")
	filePath     = os.Getenv("TLS_FILE_PATH")
	startPort    = os.Getenv("TLS_START_PORT")
	localTesting = os.Getenv("LOCAL_TESTING")
)

func main() {
	e := routers.Routers

	if localTesting != "" {
		fmt.Println("localtesting is set to:", localTesting)
		go tree.ExampleTree()

		e.Logger.Info(e.Start(":" + localTesting))
	} else {
		err := downloader.FileFromURL(downloadURL, filePath, certFile, keyFile)
		if err != nil {
			fmt.Println(err)
		}

		go func() {
			time.Sleep(24 * 60 * time.Hour)
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			if err := e.Shutdown(ctx); err != nil {
				e.Logger.Info(err)
			}
		}()

		if _, err := os.Stat(filePath + certFile); os.IsNotExist(err) {
			fmt.Println("Cert file does not exist:", err)
		}

		e.Logger.Info(e.StartTLS(startPort, filePath+certFile, filePath+keyFile))
	}
}
