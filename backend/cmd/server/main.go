package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Event struct {
	ID          string            `json:"id"`
	Title       string            `json:"title"`
	TitleI18n   map[string]string `json:"titleI18n"`
	Date        string            `json:"date"`
	Time        string            `json:"time"`
	Venue       string            `json:"venue"`
	Area        string            `json:"area"`
	Price       string            `json:"price"`
	Tags        []string          `json:"tags"`
	Summary     string            `json:"summary"`
	SummaryI18n map[string]string `json:"summaryI18n"`
	MapURL      string            `json:"mapUrl"`
	ImageHint   string            `json:"imageHint"`
}

type ContentItem struct {
	ID          string            `json:"id"`
	Type        string            `json:"type"`
	Title       string            `json:"title"`
	TitleI18n   map[string]string `json:"titleI18n"`
	Summary     string            `json:"summary"`
	SummaryI18n map[string]string `json:"summaryI18n"`
}

type JoinRequest struct {
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	Message  string `json:"message"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type Server struct {
	mux    *http.ServeMux
	static http.Handler
}

func main() {
	port := getenv("BACKEND_PORT", getenv("PORT", "8080"))
	server := NewServer()

	log.Printf("LiveLife local server listening on http://localhost:%s", port)
	if err := http.ListenAndServe(":"+port, server); err != nil {
		log.Fatal(err)
	}
}

func NewServer() *Server {
	mux := http.NewServeMux()
	staticDir := filepath.Join(".", "static")

	s := &Server{
		mux:    mux,
		static: http.FileServer(http.Dir(staticDir)),
	}

	mux.HandleFunc("GET /api/health", s.handleHealth)
	mux.HandleFunc("GET /api/events", s.handleEvents)
	mux.HandleFunc("GET /api/contents", s.handleContents)
	mux.HandleFunc("POST /api/join", s.handleJoin)
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
	writeJSON(w, http.StatusOK, map[string]any{
		"status":  "ok",
		"service": "livelife-api",
		"time":    time.Now().UTC().Format(time.RFC3339),
	})
}

func (s *Server) handleEvents(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"events": seedEvents(),
	})
}

func (s *Server) handleContents(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"contents": seedContents(),
	})
}

func (s *Server) handleJoin(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req JoinRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	req.Nickname = strings.TrimSpace(req.Nickname)
	req.Email = strings.TrimSpace(req.Email)
	req.Role = strings.TrimSpace(req.Role)
	req.Message = strings.TrimSpace(req.Message)

	if err := validateJoinRequest(req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusAccepted, map[string]any{
		"accepted": true,
		"message":  "Thanks. The local API received your join request.",
	})
}

func (s *Server) handleStatic(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, "/api/") {
		writeError(w, http.StatusNotFound, "api route not found")
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/")
	if path == "" {
		http.ServeFile(w, r, filepath.Join(".", "static", "index.html"))
		return
	}

	s.static.ServeHTTP(w, r)
}

func validateJoinRequest(req JoinRequest) error {
	if req.Nickname == "" {
		return errors.New("nickname is required")
	}
	if req.Email == "" {
		return errors.New("email is required")
	}
	if !strings.Contains(req.Email, "@") {
		return errors.New("email is invalid")
	}
	return nil
}

func seedEvents() []Event {
	return []Event{
		{
			ID:    "tokyo-loop-001",
			Title: "Tokyo Loop Night",
			TitleI18n: map[string]string{
				"zh": "东京 Loop 夜",
				"ja": "東京 Loop Night",
			},
			Date:    "2026-07-05",
			Time:    "18:30",
			Venue:   "FEVER",
			Area:    "Shindaita",
			Price:   "JPY 2,500 + 1D",
			Tags:    []string{"livehouse", "indie", "recommended"},
			Summary: "A compact livehouse night for discovering bands, DJs, and small-label releases around Tokyo.",
			SummaryI18n: map[string]string{
				"zh": "面向东京周边乐队、DJ 和小厂牌发行物的轻量 livehouse 夜场。",
				"ja": "東京周辺のバンド、DJ、小規模レーベルのリリースに出会うためのコンパクトなライブハウス企画。",
			},
			MapURL:    "https://maps.google.com/?q=FEVER+Shindaita",
			ImageHint: "poster-placeholder",
		},
		{
			ID:    "shimokita-cd-002",
			Title: "CD Shop Listening Hour",
			TitleI18n: map[string]string{
				"zh": "CD 店试听时间",
				"ja": "CDショップ試聴会",
			},
			Date:    "2026-07-12",
			Time:    "15:00",
			Venue:   "Basement Bar",
			Area:    "Shimokitazawa",
			Price:   "Free entry",
			Tags:    []string{"cd", "shop", "listening"},
			Summary: "A small listening session for new CD arrivals, zines, and artist recommendations.",
			SummaryI18n: map[string]string{
				"zh": "围绕新到 CD、zine 和音乐人推荐展开的小型试听会。",
				"ja": "新入荷のCD、ZINE、アーティスト推薦を中心にした小さな試聴セッション。",
			},
			MapURL:    "https://maps.google.com/?q=Shimokitazawa+Basement+Bar",
			ImageHint: "cd-placeholder",
		},
	}
}

func seedContents() []ContentItem {
	return []ContentItem{
		{
			ID:    "why-we-do",
			Type:  "article",
			Title: "Why We Do",
			TitleI18n: map[string]string{
				"zh": "我们为什么做",
				"ja": "なぜやるのか",
			},
			Summary: "Notes on why small live spaces, handmade releases, and local scenes still matter.",
			SummaryI18n: map[string]string{
				"zh": "关于小型现场空间、手工发行物和本地场景为什么仍然重要的笔记。",
				"ja": "小さなライブスペース、手作りのリリース、ローカルシーンが今も大切な理由についてのメモ。",
			},
		},
		{
			ID:    "next-playlist",
			Type:  "playlist",
			Title: "Next Playlist",
			TitleI18n: map[string]string{
				"zh": "下一张播放列表",
				"ja": "次のプレイリスト",
			},
			Summary: "A rotating playlist for upcoming event discovery.",
			SummaryI18n: map[string]string{
				"zh": "用于发现即将到来的活动和音乐人的轮换播放列表。",
				"ja": "これからのイベントやアーティストを見つけるためのローテーション型プレイリスト。",
			},
		},
		{
			ID:    "connect",
			Type:  "sns",
			Title: "Connect",
			TitleI18n: map[string]string{
				"zh": "连接入口",
				"ja": "つながる",
			},
			Summary: "Links for articles, photos, SNS, videos, and submissions.",
			SummaryI18n: map[string]string{
				"zh": "文章、照片、SNS、视频和投稿入口会集中在这里。",
				"ja": "記事、写真、SNS、動画、投稿入口をここにまとめます。",
			},
		},
	}
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

func getenv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
