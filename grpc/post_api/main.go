package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	chzUI "github.com/rantav/go-grpc-channelz"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	chzService "google.golang.org/grpc/channelz/service"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	postStoragePB "grpc/post_storage/proto/autogen/post_storage"

	pb "grpc/post_api/proto/autogen/post_api"
	"grpc/post_api/src/core"
)

func main() {
	logrus.SetLevel(logrus.InfoLevel)

	// channelz
	chzAddr := ":50050"
	lis, err := net.Listen("tcp", chzAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer lis.Close()
	s := grpc.NewServer()
	chzService.RegisterChannelzServiceToServer(s)
	go s.Serve(lis)
	logrus.Info("starting channelz server")
	defer s.Stop()

	http.Handle("/", chzUI.CreateHandler("/", chzAddr))
	chzUIAddr := ":8080"
	lis, err = net.Listen("tcp", chzUIAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	go http.Serve(lis, nil) // head to localhost:8080/channelz/ (trailing slash is important)

	// post api server
	port := "8888"
	lis, err = net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	postStorageName := "post-storage-grpc"
	if noHost := os.Getenv("NO_HOST"); strings.ToLower(noHost) == "true" || noHost == "1" {
		postStorageName = ""
	}

	psClient := getPostStorageClientOrPanic(postStorageName, "8890")
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
		logrus.WithError(err).Error("failed to serve grpc")
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
	fmt.Printf("initializing health check with %d milli seconds interval...\n", interval.Milliseconds())
	conn, err := grpc.Dial(fmt.Sprintf(":%s", port),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.WithError(err).Error("failed to dial grpc")
		return
	}

	c := pb.NewPostAPIClient(conn)

	ticker := time.NewTicker(interval)
	for t := range ticker.C {
		go func(t time.Time) {
			_, err := c.ShowPost(context.TODO(), &pb.ShowPostRequest{Token: "post-api-health-check"})
			if err != nil {
				logrus.WithError(err).WithField("time", t).Error("health check failed")
			}
		}(t)
	}
}

func registerChannelzService(port string) {
	// port = port + "0"

}
