package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/sirupsen/logrus"

	postStorage "post_storage/api"
)

type PostAPIServer interface {
	ShowPost(w http.ResponseWriter, r *http.Request)
}

type server struct {
	postStorageClient PostStorageClient
}

func NewServer(psClient PostStorageClient) PostAPIServer {
	return &server{
		postStorageClient: psClient,
	}
}

type ShowPostResponse struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (s *server) ShowPost(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")

	logrus.Info("got show post request with token: ", token)

	psResp, err := s.postStorageClient.GetPost(r.Context(), token)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Errorf("failed to get post: %w", err).Error()))
		fmt.Println("got error:", err.Error())
		return
	}

	logrus.Info("got post from post storage with token: ", token)

	resp := ShowPostResponse{
		Title:       psResp.Title,
		Description: psResp.Description,
	}

	mResp, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Errorf("failed to marshal response to json: %w", err).Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(mResp)
}

type PostStorageClient interface {
	GetPost(ctx context.Context, token string) (*postStorage.GetPostResponse, error)
}

type postStorageClient struct {
	serviceName string
	port        string
}

func GetPostStorageClient(serviceName, port string) PostStorageClient {
	return &postStorageClient{
		serviceName: serviceName,
		port:        port,
	}
}

func (psc *postStorageClient) GetPost(ctx context.Context, token string) (*postStorage.GetPostResponse, error) {
	res, err := http.Get(fmt.Sprintf("http://%s:%s/posts/get?token=%s", psc.serviceName, psc.port, token))
	if err != nil {
		return nil, fmt.Errorf("failed to get post from post storage: %w", err)
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	defer res.Body.Close()

	var pResp postStorage.GetPostResponse
	err = json.Unmarshal(b, &pResp)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal data to post storage response: %w", err)
	}

	return &pResp, nil
}
