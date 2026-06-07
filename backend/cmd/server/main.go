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

const brandName = "LIVE LIFE"

// LocalizedText 保存对外展示文案的三语言版本。
//
// 当前产品默认中文，同时支持日语和英语。后端直接返回三语言文案，
// 前端只需要按照当前界面语言读取对应字段，后续接数据库时也能沿用这个结构。
type LocalizedText map[string]string

// Event 是演出情报数据结构。
//
// Owned=true 表示 LIVE LIFE 自主演出，前端会把它固定展示在演出情报最上方。
// TicketNote 用来说明票务状态；目前演出票务先跳外部票站，不做站内支付。
type Event struct {
	ID             string        `json:"id"`
	Brand          string        `json:"brand"`
	Owned          bool          `json:"owned"`
	Category       string        `json:"category"`
	Title          string        `json:"title"`
	TitleI18n      LocalizedText `json:"titleI18n"`
	Date           string        `json:"date"`
	Time           string        `json:"time"`
	Venue          string        `json:"venue"`
	Area           string        `json:"area"`
	Price          string        `json:"price"`
	Tags           []string      `json:"tags"`
	Summary        string        `json:"summary"`
	SummaryI18n    LocalizedText `json:"summaryI18n"`
	Lineup         []string      `json:"lineup"`
	TicketNote     string        `json:"ticketNote"`
	TicketNoteI18n LocalizedText `json:"ticketNoteI18n"`
	MapURL         string        `json:"mapUrl"`
	ImageURL       string        `json:"imageUrl"`
	SourceNote     string        `json:"sourceNote"`
	SourceNoteI18n LocalizedText `json:"sourceNoteI18n"`
}

// CatalogItem 是 CD 严选里的单品数据结构。
//
// 注意：现在没有顶层 Shop 页面。商业路径是：
// CD 严选列表 -> 单品详情 -> “点击此处购买”按钮 -> 外部 Shop，比如 BASE。
// Format 用来区分 CD / vinyl，PurchaseURL 是外部购买链接。
type CatalogItem struct {
	ID           string        `json:"id"`
	Brand        string        `json:"brand"`
	Format       string        `json:"format"`
	Artist       string        `json:"artist"`
	Title        string        `json:"title"`
	TitleI18n    LocalizedText `json:"titleI18n"`
	Summary      string        `json:"summary"`
	SummaryI18n  LocalizedText `json:"summaryI18n"`
	Tracks       []string      `json:"tracks"`
	Status       string        `json:"status"`
	Price        string        `json:"price"`
	ImageURL     string        `json:"imageUrl"`
	PurchaseURL  string        `json:"purchaseUrl"`
	PurchaseText LocalizedText `json:"purchaseText"`
}

// ContentItem 是 Archive / 内容摘要入口。
//
// Archive 用来承接历史海报、照片、文章、公开资料备注等内容，
// 避免把它们混进演出情报或 CD 严选购买流程。
type ContentItem struct {
	ID          string        `json:"id"`
	Brand       string        `json:"brand"`
	Type        string        `json:"type"`
	Title       string        `json:"title"`
	TitleI18n   LocalizedText `json:"titleI18n"`
	Summary     string        `json:"summary"`
	SummaryI18n LocalizedText `json:"summaryI18n"`
}

// ConnectRequest 是 Connect 表单提交的数据。
//
// 这个表单是统一联系入口，可以处理票务、外部购买后未收到货、发货、
// 合作、投稿等问题。当前阶段只做本地 API 验证，后续可接邮件或客服系统。
type ConnectRequest struct {
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Topic    string `json:"topic"`
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

	log.Printf("%s local server listening on http://localhost:%s", brandName, port)
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
	writeJSON(w, http.StatusOK, map[string]any{
		"brand":   brandName,
		"service": "LIVE LIFE API",
		"status":  "ok",
		"time":    time.Now().UTC().Format(time.RFC3339),
	})
}

func (s *Server) handleEvents(w http.ResponseWriter, r *http.Request) {
	events := seedEvents()
	writeJSON(w, http.StatusOK, map[string]any{
		"brand":             brandName,
		"events":            events,
		"ownedEvents":       filterEvents(events, true),
		"recommendedEvents": filterEvents(events, false),
	})
}

func (s *Server) handleCDItems(w http.ResponseWriter, r *http.Request) {
	items := seedCDItems()
	writeJSON(w, http.StatusOK, map[string]any{
		"brand": brandName,
		"items": items,
		"cd":    filterCatalog(items, "cd"),
		"vinyl": filterCatalog(items, "vinyl"),
	})
}

func (s *Server) handleContents(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"brand":    brandName,
		"contents": seedContents(),
	})
}

func (s *Server) handleConnect(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req ConnectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	req.Nickname = strings.TrimSpace(req.Nickname)
	req.Email = strings.TrimSpace(req.Email)
	req.Topic = strings.TrimSpace(req.Topic)
	req.Message = strings.TrimSpace(req.Message)

	if err := validateConnectRequest(req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusAccepted, map[string]any{
		"accepted": true,
		"brand":    brandName,
		"message":  "LIVE LIFE received your message.",
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

func validateConnectRequest(req ConnectRequest) error {
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

func filterEvents(events []Event, owned bool) []Event {
	filtered := make([]Event, 0, len(events))
	for _, event := range events {
		if event.Owned == owned {
			filtered = append(filtered, event)
		}
	}
	return filtered
}

func filterCatalog(items []CatalogItem, format string) []CatalogItem {
	filtered := make([]CatalogItem, 0, len(items))
	for _, item := range items {
		if item.Format == format {
			filtered = append(filtered, item)
		}
	}
	return filtered
}

func seedEvents() []Event {
	return []Event{
		{
			ID:       "redhair-japan-2026-july",
			Brand:    brandName,
			Owned:    true,
			Category: "own-live",
			Title:    "LIVE LIFE presents Red Hair Boy Murder Case in Tokyo",
			TitleI18n: LocalizedText{
				"zh": "LIVE LIFE presents 紅髪少年殺人事件 东京双日演出",
				"ja": "LIVE LIFE presents 紅髪少年殺人事件 東京2公演",
				"en": "LIVE LIFE PRESENTS RED HAIR BOY MURDER CASE IN TOKYO",
			},
			Date:    "2026-07-10 / 2026-07-14",
			Time:    "7/10 OPEN 18:45 START 19:30; 7/14 OPEN 19:00 START 19:30",
			Venue:   "GRIT at Shibuya / Shimokitazawa THREE",
			Area:    "Shibuya / Shimokitazawa",
			Price:   "7/10 adv ¥5,000 + 1D, door ¥5,500 + 1D; 7/14 adv ¥4,000 + 1D, door ¥4,500 + 1D",
			Tags:    []string{"LIVE LIFE", "own live", "alternative rock", "Tokyo"},
			Summary: "Two Tokyo shows by 紅髪少年殺人事件. The July 14 Shimokitazawa show features ルサンチマン and おそロシア革命.",
			SummaryI18n: LocalizedText{
				"zh": "7月10日和7月14日，LIVE LIFE 将在东京分别呈现两场 紅髪少年殺人事件 演出。7月14日下北泽场共演为东京新生代另类摇滚乐队「ルサンチマン」和来自广岛尾道、极具个人特色的 DIY 独立音乐企划「おそロシア革命」。",
				"ja": "7月10日と7月14日、LIVE LIFE は東京で 紅髪少年殺人事件 の2公演を行います。7月14日の下北沢公演には、東京発のオルタナティブロックバンド「ルサンチマン」と、広島・尾道発のDIYインディー企画「おそロシア革命」が出演します。",
				"en": "ON JULY 10 AND JULY 14, LIVE LIFE PRESENTS TWO TOKYO SHOWS BY RED HAIR BOY MURDER CASE. THE JULY 14 SHIMOKITAZAWA SHOW FEATURES RUSANTIMAN AND OSOROSHIA KAKUMEI.",
			},
			Lineup:     []string{"紅髪少年殺人事件", "ルサンチマン", "おそロシア革命"},
			TicketNote: "Ticket links are pending. Keep the external ticket agency flow separate from the LIVE LIFE site.",
			TicketNoteI18n: LocalizedText{
				"zh": "票务链接待确认。演出票站可能由外部代理处理，LIVE LIFE 站内先展示情报，不直接结算。",
				"ja": "チケットリンクは確認中です。外部プレイガイドの導線は LIVE LIFE サイト内決済とは分けて設計します。",
				"en": "TICKET LINKS ARE PENDING. THE TICKETING AGENCY FLOW STAYS OUTSIDE LIVE LIFE CHECKOUT.",
			},
			MapURL:   "https://maps.google.com/?q=GRIT+at+Shibuya",
			ImageURL: "/assets/events/redhair-2026-july.jpg",
			SourceNoteI18n: LocalizedText{
				"zh": "活动细节来自你提供的海报和文本。第一轮公开检索暂未找到稳定的官方活动页。",
				"ja": "イベント詳細は提供されたフライヤーと本文に基づきます。初回の公開検索では安定して参照できる公式イベントページは見つかっていません。",
				"en": "DETAILS ARE BASED ON THE PROVIDED POSTER AND COPY. A STABLE OFFICIAL EVENT PAGE WAS NOT FOUND IN THE FIRST PUBLIC SEARCH PASS.",
			},
		},
		{
			ID:       "wednesday-wonderland-archive-2025",
			Brand:    brandName,
			Owned:    false,
			Category: "archive-reference",
			Title:    "Wednesday Wonderland archive visual",
			TitleI18n: LocalizedText{
				"zh": "Wednesday Wonderland 活动视觉档案",
				"ja": "Wednesday Wonderland ビジュアルアーカイブ",
				"en": "WEDNESDAY WONDERLAND ARCHIVE VISUAL",
			},
			Date:    "2025-08-21",
			Time:    "OPEN / START 18:00",
			Venue:   "BASEMENT BAR",
			Area:    "Shimokitazawa",
			Price:   "adv ¥2,500 / student ¥1,900 / door ¥3,000 + drink ¥600",
			Tags:    []string{"archive", "poster", "Shimokitazawa"},
			Summary: "A past poster kept as a visual reference for the LIVE LIFE event archive style.",
			SummaryI18n: LocalizedText{
				"zh": "这张图先作为 LIVE LIFE 活动档案和视觉气质参考，不作为即将发生的演出推荐。",
				"ja": "この画像は LIVE LIFE のイベントアーカイブとビジュアル参考として置き、直近公演の推薦情報とは分けます。",
				"en": "THIS POSTER IS KEPT AS A LIVE LIFE ARCHIVE AND VISUAL REFERENCE, SEPARATE FROM UPCOMING LIVE RECOMMENDATIONS.",
			},
			Lineup: []string{"Wednesday Wonderland", "TiDE", "できないみらい", "いきものコーナー"},
			TicketNoteI18n: LocalizedText{
				"zh": "历史活动，仅用于样式参考。",
				"ja": "過去イベントのため、スタイル参考のみ。",
				"en": "PAST EVENT, USED ONLY AS A STYLE REFERENCE.",
			},
			MapURL:   "https://maps.google.com/?q=BASEMENT+BAR+Shimokitazawa",
			ImageURL: "/assets/events/wednesday-wonderland-2025-08-21.jpg",
			SourceNoteI18n: LocalizedText{
				"zh": "信息来自你提供的海报。",
				"ja": "提供されたフライヤーに基づく情報です。",
				"en": "INFORMATION IS BASED ON THE PROVIDED POSTER.",
			},
		},
	}
}

func seedCDItems() []CatalogItem {
	return []CatalogItem{
		{
			ID:     "redhair-demo-cd",
			Brand:  brandName,
			Format: "cd",
			Artist: "紅髪少年殺人事件",
			Title:  "Selected CD placeholder",
			TitleI18n: LocalizedText{
				"zh": "紅髪少年殺人事件 CD 严选占位",
				"ja": "紅髪少年殺人事件 CD セレクト仮枠",
				"en": "RED HAIR BOY MURDER CASE CD SELECT",
			},
			Summary: "A curated CD slot for show-related releases. Purchase goes to an external shop.",
			SummaryI18n: LocalizedText{
				"zh": "这里作为 CD 严选的 CD 分类占位。单品详情页会放「点击此处购买」按钮，跳转到 BASE 等外部 Shop。",
				"ja": "CDセレクト内のCDカテゴリ仮枠です。詳細ページの購入ボタンから BASE など外部ショップへ遷移します。",
				"en": "A CD SELECT SLOT. THE DETAIL PAGE PURCHASE BUTTON LINKS TO AN EXTERNAL SHOP SUCH AS BASE.",
			},
			Tracks:      []string{"A-side reference", "Live note", "Shop link pending"},
			Status:      "external shop",
			Price:       "TBD",
			ImageURL:    "/assets/events/redhair-2026-july.jpg",
			PurchaseURL: "https://thebase.com/",
			PurchaseText: LocalizedText{
				"zh": "点击此处购买",
				"ja": "こちらから購入",
				"en": "BUY HERE",
			},
		},
		{
			ID:     "live-life-vinyl-placeholder",
			Brand:  brandName,
			Format: "vinyl",
			Artist: "LIVE LIFE SELECT",
			Title:  "Selected vinyl placeholder",
			TitleI18n: LocalizedText{
				"zh": "LIVE LIFE 黑胶严选占位",
				"ja": "LIVE LIFE ヴァイナルセレクト仮枠",
				"en": "LIVE LIFE VINYL SELECT",
			},
			Summary: "A vinyl slot for future selected records.",
			SummaryI18n: LocalizedText{
				"zh": "这里作为黑胶分类占位，后续放精选黑胶、推荐语和外部购买链接。",
				"ja": "ヴァイナルカテゴリの仮枠です。今後、推薦盤、紹介文、外部購入リンクを掲載します。",
				"en": "A VINYL CATEGORY SLOT FOR SELECTED RECORDS, NOTES, AND EXTERNAL PURCHASE LINKS.",
			},
			Tracks:      []string{"Side A", "Side B", "Listening note pending"},
			Status:      "curating",
			Price:       "TBD",
			PurchaseURL: "https://thebase.com/",
			PurchaseText: LocalizedText{
				"zh": "点击此处购买",
				"ja": "こちらから購入",
				"en": "BUY HERE",
			},
		},
	}
}

func seedContents() []ContentItem {
	return []ContentItem{
		{
			ID:    "homepage-positioning",
			Brand: brandName,
			Type:  "archive",
			Title: "What LIVE LIFE is building",
			TitleI18n: LocalizedText{
				"zh": "LIVE LIFE 正在做什么",
				"ja": "LIVE LIFE が作っているもの",
				"en": "WHAT LIVE LIFE IS BUILDING",
			},
			Summary: "LIVE LIFE is a Tokyo-facing music entry point for shows, CD select, archive, and support messages.",
			SummaryI18n: LocalizedText{
				"zh": "LIVE LIFE 会先成为东京演出情报、LIVE LIFE 自主演出、CD 严选、Archive 和 Connect 的统一入口。",
				"ja": "LIVE LIFE は東京のライブ情報、自主公演、CDセレクト、Archive、Connect をまとめる入口として始めます。",
				"en": "LIVE LIFE STARTS AS ONE ENTRY POINT FOR TOKYO SHOWS, OWNED EVENTS, CD SELECT, ARCHIVE, AND CONNECT.",
			},
		},
		{
			ID:    "visual-system-note",
			Brand: brandName,
			Type:  "design",
			Title: "Visual system note",
			TitleI18n: LocalizedText{
				"zh": "视觉系统备注",
				"ja": "ビジュアルシステムメモ",
				"en": "VISUAL SYSTEM NOTE",
			},
			Summary: "The homepage uses real band names as cultural texture and abstract music production tracks as data flow.",
			SummaryI18n: LocalizedText{
				"zh": "首页纹理使用真实乐队名作为文化坐标，同时用抽象音轨、波形和采样轨道替代代码/二进制。",
				"ja": "ホームのテクスチャは実在バンド名を文化的座標として使い、コードやバイナリの代わりに抽象的な音轨、波形、サンプルトラックを使います。",
				"en": "THE HOMEPAGE USES REAL BAND NAMES AS CULTURAL TEXTURE, WITH ABSTRACT TRACK LANES AND WAVEFORM-STYLE DATA INSTEAD OF CODE OR BINARY.",
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
