package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"livelife/backend/internal/application"
	sqlitestore "livelife/backend/internal/infrastructure/sqlite"
	"livelife/backend/internal/interfaces/httpapi"
)

// cmd/server 只负责程序入口和依赖组装。
// 业务规则、HTTP 处理、数据库实现都放在 internal 下，避免多人协作时把所有逻辑继续堆进 main.go。
func main() {
	port := getenv("BACKEND_PORT", getenv("PORT", "8080"))
	server, err := NewFileServer(defaultDatabasePath())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("LIVE LIFE local server listening on http://localhost:%s", port)
	if err := http.ListenAndServe(":"+port, server); err != nil {
		log.Fatal(err)
	}
}

// NewServer 给测试使用内存数据库，避免污染本地 SQLite 文件。
func NewServer() http.Handler {
	server, err := newServer(":memory:", filepath.Join("..", "..", "static"))
	if err != nil {
		panic(err)
	}
	return server
}

// NewFileServer 给本地预览和部署环境使用文件数据库。
func NewFileServer(databasePath string) (http.Handler, error) {
	return newServer(databasePath, filepath.Join(".", "static"))
}

func newServer(databasePath string, staticDir string) (http.Handler, error) {
	store, err := sqlitestore.NewStore(databasePath)
	if err != nil {
		return nil, err
	}
	service := application.NewService(store)
	return httpapi.NewServer(service, staticDir), nil
}

func defaultDatabasePath() string {
	if path := os.Getenv("DATABASE_PATH"); path != "" {
		return path
	}
	return filepath.Join("data", "livelife.dev.db")
}

func getenv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
