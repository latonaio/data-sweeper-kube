package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Msg is request message from client
type Msg struct {
	DirPath               string   `json:"dir_path" form:"dir_path" query:"dir_path"`
	ExcludeFiles          []string `json:"exclude_files" form:"exclude_files" query:"exclude_files"`
	ExcludeFileExtensions []string `json:"exclude_file_extensions" from:"exclude_file_extensions" query:"exclude_file_extensions"`
	IsRecursive           bool     `json:"is_recursive" form:"is_recursive" query:"is_recursive"`
}

// EchoServer is API server
type EchoServer struct {
	server *echo.Echo
	config *http.Server
}

// Start method is starting server
func (s *EchoServer) Start() {
	s.server.Logger.Fatal(s.server.StartServer(s.config))
}

// Stop method is stopping server
func (s *EchoServer) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	if err := s.server.Shutdown(ctx); err != nil {
		s.server.Logger.Fatal(err)
	}
}

// NewServer is generating new EchoServer
func NewServer(host string, port int) *EchoServer {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
	}))

	e.POST("/sweeper", sweep)

	config := &http.Server{
		Addr: host + ":" + strconv.Itoa(port),
	}

	return &EchoServer{
		server: e,
		config: config,
	}
}

func generateSet(keys []string) map[string]struct{} {
	sets := make(map[string]struct{})
	for _, key := range keys {
		sets[key] = struct{}{}
	}
	return sets
}

func deleteFiles(dirPath string, excludeFiles map[string]struct{}, excludeFileExtensions map[string]struct{}, isRecursive bool) error {
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return err
	}
	for _, file := range files {
		fmt.Println("find file: " + file.Name())
		if file.IsDir() {
			if isRecursive {
				deleteFiles(filepath.Join(dirPath, file.Name()), excludeFiles, excludeFileExtensions, isRecursive)
			}
		} else if _, ignore := excludeFiles[file.Name()]; ignore {
			fmt.Println("ignore file: " + file.Name())
		} else if _, ignore := excludeFileExtensions[filepath.Ext(file.Name())]; ignore {
			fmt.Println("ignore file: " + file.Name())
		} else {
			deleteFile(filepath.Join(dirPath, file.Name()))
		}
	}
	return nil
}

func deleteFile(filePath string) error {
	if err := os.Remove(filePath); err != nil {
		return err
	}
	fmt.Println("deleted file: " + filePath)
	return nil
}

func sweep(c echo.Context) error {
	m := new(Msg)
	if err := c.Bind(m); err != nil {
		fmt.Println("cannot bind data")
		return c.String(http.StatusBadRequest, "invalid data format")
	}
	deleteFiles(m.DirPath, generateSet(m.ExcludeFiles), generateSet(m.ExcludeFileExtensions), m.IsRecursive)
	return c.String(http.StatusOK, "delete files")
}
