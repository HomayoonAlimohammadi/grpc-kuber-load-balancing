package core

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"

	pb "grpc/post_storage/proto/autogen/post_storage"
)

type service struct {
	pb.UnimplementedPostStorageServer
}

func NewService() pb.PostStorageServer {
	return &service{}
}

func (s *service) GetPost(ctx context.Context, req *pb.GetPostRequest) (*pb.GetPostResponse, error) {
	logrus.Info("got GetPost request with token:", req.Token)

	sleepSeconds := 1000 // milli seconds
	if tStr := os.Getenv("POST_STORAGE_SLEEP_TIME_MILLIS"); tStr != "" {
		if t, err := strconv.Atoi(tStr); err == nil {
			sleepSeconds = t
		}
	}

	logrus.Infof("sleeping for %d milli seconds", sleepSeconds)

	// sleep to mimic an expensive calculation
	time.Sleep(time.Duration(sleepSeconds) * time.Millisecond)

	logrus.Info("sending response...")

	return &pb.GetPostResponse{
		Post: &pb.Post{
			Title:       "some title",
			Description: "some description",
			Phone:       "09120000001",
		},
	}, nil
}
