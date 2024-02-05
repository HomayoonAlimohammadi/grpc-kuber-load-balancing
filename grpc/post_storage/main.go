package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	chzUI "github.com/rantav/go-grpc-channelz"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	chzService "google.golang.org/grpc/channelz/service"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	pb "grpc/post_storage/proto/autogen/post_storage"
	"grpc/post_storage/src/core"
)

func main() {
	logrus.SetLevel(logrus.InfoLevel)

	// channelz
	chzAddr := ":50051"
	lis, err := net.Listen("tcp", chzAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer lis.Close()
	s := grpc.NewServer()
	chzService.RegisterChannelzServiceToServer(s)
	go s.Serve(lis)
	logrus.Info("started channelz server")
	defer s.Stop()

	http.Handle("/", chzUI.CreateHandler("/", chzAddr))
	chzUIAddr := ":8081"
	lis, err = net.Listen("tcp", chzUIAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	go http.Serve(lis, nil) // head to localhost:8081/channelz/ (trailing slash is important)

	// post storage server
	port := "8890"
	lis, err = net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	reflection.Register(grpcServer)
	chzService.RegisterChannelzServiceToServer(grpcServer)
	pb.RegisterPostStorageServer(grpcServer, core.NewService())

	intervalStr, ok := os.LookupEnv("POST_STORAGE_INTERVAL_MILLIS")

	interval := 20000 // millis
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
		logrus.WithError(err).Error("failed to serve grpc")
	}

	logrus.Info("shutting down post-storage-grpc server...")
}

func healthCheck(port string, interval time.Duration) {
	fmt.Printf("initializing health check with %d milli seconds interval...\n", interval.Milliseconds())
	conn, err := grpc.Dial(fmt.Sprintf(":%s", port),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.WithError(err).Error("failed to dial grpc")
		return
	}

	c := pb.NewPostStorageClient(conn)

	ticker := time.NewTicker(interval)
	for t := range ticker.C {
		go func(t time.Time) {
			_, err := c.GetPost(context.TODO(), &pb.GetPostRequest{Token: "post-storage-health-check"})
			if err != nil {
				logrus.WithError(err).WithField("time", t).Error("health check failed")
			}
		}(t)
	}
}
