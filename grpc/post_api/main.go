package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	postStoragePB "post_storage/proto/autogen/post_storage"

	pb "post_api/proto/autogen/post_api"
	"post_api/src/core"
)

func main() {
	logrus.SetLevel(logrus.InfoLevel)
	port := "8888"

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)
	psClient := getPostStorageClientOrPanic("post-storage-grpc", "8890")
	pb.RegisterPostAPIServer(grpcServer, core.NewService(psClient))

	intervalStr, ok := os.LookupEnv("POST_API_INTERVAL_MILLIS")

	interval := 2000 // millis
	if ok {
		interval, err = strconv.Atoi(intervalStr)
		if err != nil {
			panic(err)
		}
	}

	go healthCheck(port, time.Duration(interval)*time.Millisecond)

	logrus.Infof("starting post-api-grpc server on port %s", port)
	err = grpcServer.Serve(lis)
	if err != nil {
		panic(err)
	}

	logrus.Info("shutting down post-api-grpc server...")
}

func getPostStorageClientOrPanic(serviceName, port string) postStoragePB.PostStorageClient {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", serviceName, port),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	return postStoragePB.NewPostStorageClient(conn)
}

func healthCheck(port string, interval time.Duration) {
	fmt.Println("initializing health check...")
	conn, err := grpc.Dial(fmt.Sprintf(":%s", port),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.WithError(err).Error("failed to dial grpc")
		return
	}

	c := pb.NewPostAPIClient(conn)

	ticker := time.NewTicker(interval)
	for t := range ticker.C {
		_, err := c.ShowPost(context.TODO(), &pb.ShowPostRequest{Token: "post-api-health-check"})
		if err != nil {
			logrus.WithError(err).WithField("time", t).Error("health check failed")
		}
	}
}
