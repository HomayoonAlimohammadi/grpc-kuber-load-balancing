package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
)

type PostStorageServer interface {
	GetPost(w http.ResponseWriter, r *http.Request)
}

type server struct{}

func NewServer() PostStorageServer {
	return &server{}
}

type GetPostResponse struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Phone       string `json:"phone"`
}

func (s *server) GetPost(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")

	logrus.Info("got get post request with token: ", token)

	resp := GetPostResponse{
		Title:       "some title",
		Description: "some description",
		Phone:       "09120000001",
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
