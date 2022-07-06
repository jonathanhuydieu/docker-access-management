package main

import (
	// "encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
)

func main() {
	var socketPath string
	flag.StringVar(&socketPath, "socket", "/run/guest/volumes-service.sock", "Unix domain socket to listen on")
	flag.Parse()

	os.RemoveAll(socketPath)
	
	logrus.New().Infof("Starting listening on %s\n", socketPath)
	router := echo.New()
	router.HideBanner = true

	startURL := ""

	ln, err := listen(socketPath)
	if err != nil {
		log.Fatal(err)
	}
	router.Listener = ln

	router.GET("/repo", repo)

	log.Fatal(router.Start(startURL))
}

func listen(path string) (net.Listener, error) {
	return net.Listen("unix", path)
}

type Repository struct {
	User            string `json:"user,omitempty"`
	Name            string `json:"name"`
	Namespace       string `json:"namespace"`
	Description     string `json:"description"`
	Private         bool   `json:"is_private"`
	PullCount       int    `json:"pull_count"`
	FullDescription string `json:"full_description,omitempty"`
}
type Repositories struct {
	User         string                 `json:"user,omitempty"`
	Repositories map[string]interface{} `json:"results"`
	MaxResults   int
}

func repo (ctx echo.Context)  error {
	c := NewClient("https://hub-stage.docker.com/v2", "ryanhristovski", "Hackathon2022")
	repository := Repository{}
	err := c.sendRequest(ctx.Request().Context(), "GET", fmt.Sprintf("/repositories/ryanhristovski/personal-repo-demo/"), nil, &repository)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, repository)
}

type HTTPMessageBody struct {
	Message string
}