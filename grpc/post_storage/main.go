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

	pb "post_storage/proto/autogen/post_storage"
	"post_storage/src/core"
)

func main() {
	logrus.SetLevel(logrus.InfoLevel)
	port := "8890"

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	reflection.Register(grpcServer)
	pb.RegisterPostStorageServer(grpcServer, core.NewService())

	intervalStr, ok := os.LookupEnv("POST_STORAGE_INTERVAL_MILLIS")

	interval := 2000 // millis
	if ok {
		interval, err = strconv.Atoi(intervalStr)
		if err != nil {
			panic(err)
		}
	}

	go healthCheck(port, time.Duration(interval)*time.Millisecond)

	logrus.Infof("starting post-storage-grpc server on port %s", port)
	err = grpcServer.Serve(lis)
	if err != nil {
		panic(err)
	}

	logrus.Info("shutting down post-storage-grpc server...")
}

func healthCheck(port string, interval time.Duration) {
	fmt.Println("initializing health check...")
	conn, err := grpc.Dial(fmt.Sprintf(":%s", port),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.WithError(err).Error("failed to dial grpc")
		return
	}

	c := pb.NewPostStorageClient(conn)

	ticker := time.NewTicker(interval)
	for t := range ticker.C {
		_, err := c.GetPost(context.TODO(), &pb.GetPostRequest{Token: "post-storage-health-check"})
		if err != nil {
			logrus.WithError(err).WithField("time", t).Error("health check failed")
		}
	}
}
