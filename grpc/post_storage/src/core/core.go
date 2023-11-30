package core

import (
	"context"

	"github.com/sirupsen/logrus"

	pb "post_storage/proto/autogen/post_storage"
)

type service struct {
	pb.UnimplementedPostStorageServer
}

func NewService() pb.PostStorageServer {
	return &service{}
}

func (s *service) GetPost(ctx context.Context, req *pb.GetPostRequest) (*pb.GetPostResponse, error) {
	logrus.Info("got GetPost request with token:", req.Token)

	return &pb.GetPostResponse{
		Post: &pb.Post{
			Title:       "some title",
			Description: "some description",
			Phone:       "09120000001",
		},
	}, nil
}
