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

// LocalizedText 用来保存一个字段的三语言展示文案。
//
// 现在产品默认展示中文，同时支持日语和英语。后端直接返回三语言文案，
// 前端只负责根据当前语言取对应字段，这样以后接数据库时也能沿用同一套结构。
type LocalizedText map[string]string

// Event 是首页和演出情报页使用的活动数据结构。
//
// 说明：
// - Title/Summary 保留英文基础值，方便没有语言包时兜底展示。
// - TitleI18n/SummaryI18n 是对外展示字段，必须使用 LIVE LIFE 当前支持的三语言。
// - Owned 用来区分“我们自己的演出”和“推荐演出”。前端会把 Owned=true 的活动固定在最上面。
// - TicketNote 暂时只放票务说明，不直接做购票；因为你提到演出票站有代理逻辑，后续再接外部票站。
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

// CatalogItem 表示 CD 页面和 Shop 页面共用的商品/内容卡片。
//
// 当前阶段先不做登录、购物车和支付，只把页面结构和展示数据拆清楚：
// CD 页面负责音乐发行物，Shop 页面负责周边或其他商品。等购买流程确定后，
// 可以在这个结构上补库存、SKU、支付状态和订单字段。
type CatalogItem struct {
	ID          string        `json:"id"`
	Brand       string        `json:"brand"`
	Kind        string        `json:"kind"`
	Title       string        `json:"title"`
	TitleI18n   LocalizedText `json:"titleI18n"`
	Summary     string        `json:"summary"`
	SummaryI18n LocalizedText `json:"summaryI18n"`
	Status      string        `json:"status"`
	Price       string        `json:"price"`
	ImageURL    string        `json:"imageUrl"`
}

// ContentItem 是首页下方的编辑型内容入口。
//
// 它不直接等同于商品或活动，而是用于放文章、播放列表、照片记录、
// SNS 聚合等内容。这样首页不会把“演出、CD、Shop、售后联系”混在一起。
type ContentItem struct {
	ID          string        `json:"id"`
	Brand       string        `json:"brand"`
	Type        string        `json:"type"`
	Title       string        `json:"title"`
	TitleI18n   LocalizedText `json:"titleI18n"`
	Summary     string        `json:"summary"`
	SummaryI18n LocalizedText `json:"summaryI18n"`
}

// ConnectRequest 是页面底部 Connect 表单提交的数据。
//
// 这个表单不是“Join us”。它更像一个统一联系入口：可以问票务、CD/Shop
// 发货、付款后未收到货、投稿、合作等问题。现在只做本地 API 接收验证，
// 后续可以接邮件、数据库或客服系统。
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
	mux.HandleFunc("GET /api/contents", s.handleContents)
	mux.HandleFunc("GET /api/cd-items", s.handleCDItems)
	mux.HandleFunc("GET /api/shop-items", s.handleShopItems)
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
		"status":  "ok",
		"service": "LIVE LIFE API",
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

func (s *Server) handleContents(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"brand":    brandName,
		"contents": seedContents(),
	})
}

func (s *Server) handleCDItems(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"brand": brandName,
		"items": seedCDItems(),
	})
}

func (s *Server) handleShopItems(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"brand": brandName,
		"items": seedShopItems(),
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
				"en": "LIVE LIFE presents Red Hair Boy Murder Case in Tokyo",
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
				"en": "On July 10 and July 14, LIVE LIFE presents two Tokyo shows by Red Hair Boy Murder Case. The July 14 Shimokitazawa show features Tokyo alternative rock band Rusantiman and Onomichi DIY project Osoroshia Kakumei.",
			},
			Lineup:     []string{"紅髪少年殺人事件", "ルサンチマン", "おそロシア革命"},
			TicketNote: "Ticket links are pending. Keep the external ticket agency flow separate from the LIVE LIFE site.",
			TicketNoteI18n: LocalizedText{
				"zh": "票务链接待确认。演出票站可能由外部代理处理，LIVE LIFE 站内先展示情报，不直接结算。",
				"ja": "チケットリンクは確認中です。外部プレイガイドの導線は LIVE LIFE サイト内決済とは分けて設計します。",
				"en": "Ticket links are pending. The ticketing agency flow should stay separate from LIVE LIFE checkout.",
			},
			MapURL:     "https://maps.google.com/?q=GRIT+at+Shibuya",
			ImageURL:   "/assets/events/redhair-2026-july.jpg",
			SourceNote: "Details are based on the user-provided poster and copy. Public official event pages were not found during the first search pass.",
			SourceNoteI18n: LocalizedText{
				"zh": "活动细节来自你提供的海报和文本。第一轮公开检索暂未找到稳定的官方活动页。",
				"ja": "イベント詳細は提供されたフライヤーと本文に基づきます。初回の公開検索では安定して参照できる公式イベントページは見つかっていません。",
				"en": "Details are based on the provided poster and copy. A stable official event page was not found in the first public search pass.",
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
				"en": "Wednesday Wonderland archive visual",
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
				"en": "This poster is kept as a LIVE LIFE archive and visual reference, separate from upcoming live recommendations.",
			},
			Lineup: []string{"Wednesday Wonderland", "TiDE", "できないみらい", "いきものコーナー"},
			TicketNoteI18n: LocalizedText{
				"zh": "历史活动，仅用于样式参考。",
				"ja": "過去イベントのため、スタイル参考のみ。",
				"en": "Past event, used only as a style reference.",
			},
			MapURL:   "https://maps.google.com/?q=BASEMENT+BAR+Shimokitazawa",
			ImageURL: "/assets/events/wednesday-wonderland-2025-08-21.jpg",
			SourceNoteI18n: LocalizedText{
				"zh": "信息来自你提供的海报。",
				"ja": "提供されたフライヤーに基づく情報です。",
				"en": "Information is based on the provided poster.",
			},
		},
	}
}

func seedCDItems() []CatalogItem {
	return []CatalogItem{
		{
			ID:    "redhair-cd-placeholder",
			Brand: brandName,
			Kind:  "cd",
			Title: "LIVE LIFE CD selection",
			TitleI18n: LocalizedText{
				"zh": "LIVE LIFE CD 选择",
				"ja": "LIVE LIFE CD セレクション",
				"en": "LIVE LIFE CD selection",
			},
			Summary: "A placeholder for CDs and releases curated around upcoming shows.",
			SummaryI18n: LocalizedText{
				"zh": "这里先独立成 CD 页面，用来放我们推荐或发行相关的 CD、唱片和试听信息。",
				"ja": "ここはCD専用ページとして、推薦盤・関連リリース・試聴情報を分けて掲載します。",
				"en": "A separate CD page for curated releases, records, and listening notes.",
			},
			Status: "planning",
			Price:  "TBD",
		},
	}
}

func seedShopItems() []CatalogItem {
	return []CatalogItem{
		{
			ID:    "live-life-shop-placeholder",
			Brand: brandName,
			Kind:  "shop",
			Title: "LIVE LIFE Shop",
			TitleI18n: LocalizedText{
				"zh": "LIVE LIFE Shop",
				"ja": "LIVE LIFE Shop",
				"en": "LIVE LIFE Shop",
			},
			Summary: "A placeholder for merch and future checkout discussion.",
			SummaryI18n: LocalizedText{
				"zh": "Shop 页面先单独列出。是否需要登录注册、购物车、订单和支付，等购买流程确定后再做。",
				"ja": "Shopページは独立させます。ログイン、カート、注文、決済は購入フロー確定後に設計します。",
				"en": "The Shop page is separated. Login, cart, orders, and checkout can be designed once the buying flow is decided.",
			},
			Status: "discussion",
			Price:  "TBD",
		},
	}
}

func seedContents() []ContentItem {
	return []ContentItem{
		{
			ID:    "homepage-positioning",
			Brand: brandName,
			Type:  "article",
			Title: "What LIVE LIFE is building",
			TitleI18n: LocalizedText{
				"zh": "LIVE LIFE 正在做什么",
				"ja": "LIVE LIFE が作っているもの",
				"en": "What LIVE LIFE is building",
			},
			Summary: "LIVE LIFE is a Tokyo-facing music entry point for shows, CD selections, shop items, and support messages.",
			SummaryI18n: LocalizedText{
				"zh": "LIVE LIFE 会先成为东京演出情报、我们自己的演出、CD 选择、Shop 和售后/合作联系的统一入口。",
				"ja": "LIVE LIFE は東京のライブ情報、自主公演、CDセレクト、Shop、問い合わせをまとめる入口として始めます。",
				"en": "LIVE LIFE starts as one entry point for Tokyo shows, owned events, CD picks, shop items, and support messages.",
			},
		},
		{
			ID:    "public-research-note",
			Brand: brandName,
			Type:  "research",
			Title: "Research note",
			TitleI18n: LocalizedText{
				"zh": "公开资料备注",
				"ja": "公開情報メモ",
				"en": "Research note",
			},
			Summary: "Band and venue descriptions should distinguish public sources from user-provided poster details.",
			SummaryI18n: LocalizedText{
				"zh": "共演乐队和场地资料可引用公开来源；具体演出安排在没有官方页面前，先标注为来自海报和你提供的文本。",
				"ja": "共演者と会場情報は公開ソースを参照できます。公演詳細は公式ページ確認前はフライヤーと提供文面由来として扱います。",
				"en": "Band and venue notes can cite public sources. Event details should be marked as poster/user-provided until official pages are confirmed.",
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
