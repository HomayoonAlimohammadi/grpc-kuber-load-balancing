package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"

	"http/post_storage/api"
)

func main() {
	logrus.SetLevel(logrus.InfoLevel)
	port := "8090"
	var err error

	srv := api.NewServer()

	http.HandleFunc("/posts/get", srv.GetPost)

	intervalStr, ok := os.LookupEnv("POST_STORAGE_INTERVAL_MILLIS")

	interval := 20000 // milis
	if ok {
		interval, err = strconv.Atoi(intervalStr)
		if err != nil {
			panic(err)
		}
	}

	go healthCheck(port, time.Duration(interval)*time.Millisecond)

	logrus.Infof("starting post-storage-http server on port %s", port)
	err = http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		panic(err)
	}

	logrus.Info("shutting down post-storage-http server...")
}

func healthCheck(port string, interval time.Duration) {
	fmt.Printf("initializing health check with %d milliseconds interval...\n", interval.Milliseconds())

	ticker := time.NewTicker(interval)
	for range ticker.C {
		go func() {
			res, err := http.Get(fmt.Sprintf("http://localhost:%s/posts/get?token=post-storage-health-check", port))
			if err != nil {
				logrus.WithError(err).Error("failed to health check")
				return
			}

			if res.StatusCode != 200 {
				logrus.WithError(err).Error("failed to health check")
			}
		}()
	}
}
