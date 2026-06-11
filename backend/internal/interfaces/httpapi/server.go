package httpapi

import (
	"encoding/json"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"livelife/backend/internal/application"
	"livelife/backend/internal/domain"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type Server struct {
	mux       *http.ServeMux
	static    http.Handler
	staticDir string
	service   *application.Service
}

func NewServer(service *application.Service, staticDir string) *Server {
	mux := http.NewServeMux()
	s := &Server{
		mux:       mux,
		static:    http.FileServer(http.Dir(staticDir)),
		staticDir: staticDir,
		service:   service,
	}

	// API 路由保持轻量：这里只做 HTTP 方法、JSON 编解码和状态码转换。
	// 具体的业务校验、数据读取、数据库保存都交给 application/service 和 repository 实现。
	mux.HandleFunc("GET /api/health", s.handleHealth)
	mux.HandleFunc("GET /api/events", s.handleEvents)
	mux.HandleFunc("GET /api/cd-items", s.handleCDItems)
	mux.HandleFunc("GET /api/contents", s.handleContents)
	mux.HandleFunc("POST /api/connect", s.handleConnect)
	mux.HandleFunc("POST /api/join", s.handleConnect)
	mux.HandleFunc("/", s.handleStatic)

	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("X-Content-Type-Options", "nosniff")
	if strings.HasPrefix(r.URL.Path, "/api/") {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	}
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	s.mux.ServeHTTP(w, r)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	// health 返回品牌名，方便前端、本地脚本、上线巡检确认当前服务确实是 LIVE LIFE API。
	writeJSON(w, http.StatusOK, map[string]any{
		"brand":   domain.BrandName,
		"service": "LIVE LIFE API",
		"status":  "ok",
		"time":    time.Now().UTC().Format(time.RFC3339),
	})
}

func (s *Server) handleEvents(w http.ResponseWriter, r *http.Request) {
	events, err := s.service.ListEvents(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load events")
		return
	}
	// events 保留完整列表，同时额外返回 ownedEvents / recommendedEvents。
	// 这样前端可以直接把 LIVE LIFE 自主演出固定在上方，而不用在页面组件里重复写分组逻辑。
	writeJSON(w, http.StatusOK, map[string]any{
		"brand":             domain.BrandName,
		"events":            events,
		"ownedEvents":       filterEvents(events, true),
		"recommendedEvents": filterEvents(events, false),
	})
}

func (s *Server) handleCDItems(w http.ResponseWriter, r *http.Request) {
	items, err := s.service.ListCatalogItems(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load catalog items")
		return
	}
	// CD 严选内部拆成 CD / vinyl 两类。这里没有恢复顶层 Shop API，
	// 购买行为只通过单品里的 purchaseUrl 跳到 BASE 等外部商店。
	writeJSON(w, http.StatusOK, map[string]any{
		"brand": domain.BrandName,
		"items": items,
		"cd":    filterCatalog(items, "cd"),
		"vinyl": filterCatalog(items, "vinyl"),
	})
}

func (s *Server) handleContents(w http.ResponseWriter, r *http.Request) {
	contents, err := s.service.ListContents(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load contents")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"brand":    domain.BrandName,
		"contents": contents,
	})
}

func (s *Server) handleConnect(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req domain.ConnectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	result, err := s.service.SubmitConnect(r.Context(), req)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Connect 是统一联系入口：票务、CD 未收到货、合作、投稿都先进入同一张消息表。
	// 后续如果做后台管理，可以按 topic/status 分流处理。
	writeJSON(w, http.StatusAccepted, map[string]any{
		"accepted":  true,
		"brand":     result.Brand,
		"message":   result.Message,
		"messageId": result.MessageID,
	})
}

func (s *Server) handleStatic(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, "/api/") {
		writeError(w, http.StatusNotFound, "api route not found")
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/")
	if path == "" {
		http.ServeFile(w, r, filepath.Join(s.staticDir, "index.html"))
		return
	}

	s.static.ServeHTTP(w, r)
}

func filterEvents(events []domain.Event, owned bool) []domain.Event {
	filtered := make([]domain.Event, 0, len(events))
	for _, event := range events {
		if event.Owned == owned {
			filtered = append(filtered, event)
		}
	}
	return filtered
}

func filterCatalog(items []domain.CatalogItem, format string) []domain.CatalogItem {
	filtered := make([]domain.CatalogItem, 0, len(items))
	for _, item := range items {
		if item.Format == format {
			filtered = append(filtered, item)
		}
	}
	return filtered
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(value); err != nil {
		log.Printf("write json: %v", err)
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, ErrorResponse{Error: message})
}
