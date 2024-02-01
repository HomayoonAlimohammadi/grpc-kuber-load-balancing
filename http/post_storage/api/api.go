package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

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

	sleepTime := 1000 // millis
	if tStr := os.Getenv("POST_STORAGE_SLEEP_TIME_MILLIS"); tStr != "" {
		if t, err := strconv.Atoi(tStr); err == nil {
			sleepTime = t
		}
	}

	logrus.Infof("sleeping for %d milliseconds", sleepTime)
	// sleep to mimic an expensive calculation
	time.Sleep(time.Duration(sleepTime) * time.Millisecond)

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

	logrus.Info("sending response...")

	w.Header().Set("Content-Type", "application/json")
	w.Write(mResp)
}
