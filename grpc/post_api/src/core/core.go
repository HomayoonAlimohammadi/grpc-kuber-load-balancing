package core

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"

	pb "grpc/post_api/proto/autogen/post_api"
	"grpc/post_storage/proto/autogen/post_storage"
)

type service struct {
	pb.UnimplementedPostAPIServer

	postStorageClient post_storage.PostStorageClient
}

func NewService(
	postStorageClient post_storage.PostStorageClient,
) pb.PostAPIServer {
	return &service{
		postStorageClient: postStorageClient,
	}
}

func (s *service) ShowPost(ctx context.Context, req *pb.ShowPostRequest) (*pb.ShowPostResponse, error) {
	logrus.Info("got show post request with token:", req.Token)

	psResp, err := s.postStorageClient.GetPost(ctx, &post_storage.GetPostRequest{Token: req.Token})
	if err != nil {
		return nil, fmt.Errorf("failed to get post from post storage: %w", err)
	}

	logrus.Info("got post from post storage with token:", req.Token)

	if psResp.Post == nil {
		return nil, fmt.Errorf("post can not be nil")
	}

	return &pb.ShowPostResponse{
		Title:       psResp.Post.Title,
		Description: psResp.Post.Description,
	}, nil
}
