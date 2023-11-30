package main

import (
	"fmt"
	"net/http"
	"os"
	"post_api/api"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetLevel(logrus.InfoLevel)
	port := "8080"
	psPort := "8090"

	// set this to post-storage-http for kubernetes
	psServiceName := "post-storage-http"
	var err error

	psClient := api.GetPostStorageClient(psServiceName, psPort)
	srv := api.NewServer(psClient)

	http.HandleFunc("/posts/show", srv.ShowPost)

	intervalStr, ok := os.LookupEnv("POST_API_INTERVAL_MILLIS")

	interval := 2000 // milis
	if ok {
		interval, err = strconv.Atoi(intervalStr)
		if err != nil {
			panic(err)
		}
	}

	go healthCheck(port, time.Duration(interval)*time.Millisecond)

	logrus.Infof("starting post-api-http server on port %s", port)
	err = http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		panic(err)
	}

	logrus.Info("shutting down post-api-http server...")
}

func healthCheck(port string, interval time.Duration) {
	fmt.Println("initializing health check...")

	ticker := time.NewTicker(interval)
	for range ticker.C {
		res, err := http.Get(fmt.Sprintf("http://localhost:%s/posts/show?token=post-api-health-check", port))
		if err != nil {
			logrus.WithError(err).Error("failed to health check")
			continue
		}
		if res.StatusCode != 200 {
			logrus.WithError(err).Error("failed to health check")
		}
	}
}
