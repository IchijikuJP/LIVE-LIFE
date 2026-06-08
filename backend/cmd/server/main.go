package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

const brandName = "LIVE LIFE"

// LocalizedText 保存对外展示文案的三语言版本。
// 当前默认语言是中文，同时支持日本语和 English。
// 后端保留完整三语言数据，前端根据当前界面语言选择对应字段。
type LocalizedText map[string]string

// Event 是 /api/events 对外返回的数据结构。
// 注意它是 API DTO，不是数据库表模型。这样前端视觉从 V2 切到 V3 时，
// 后端数据库模型可以保持稳定，API 也不用因为布局变化而改变。
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

// CatalogItem 是 /api/cd-items 对外返回的数据结构。
// 这里覆盖 CD 严选里的 CD 和黑胶。购买按钮只跳外部 Shop，不做站内订单。
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

// ContentItem 是 /api/contents 对外返回的数据结构。
// Archive、设计备注、历史海报说明等都先归到这里，避免混进演出或购买流程。
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
// 它承接票务、外部购买未收到货、发货、合作、投稿等问题。
type ConnectRequest struct {
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Topic    string `json:"topic"`
	Message  string `json:"message"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

// 下面这些 Model 是数据库表模型。
// 目前用 GORM + SQLite 做完整数据链路：数据库 -> ORM -> API DTO -> 前端。
// 它们不直接暴露给前端，避免未来数据库调整影响 API。
type EventModel struct {
	ID           string `gorm:"primaryKey"`
	Brand        string `gorm:"not null;default:LIVE LIFE"`
	Owned        bool   `gorm:"not null;default:false;index"`
	Category     string `gorm:"not null"`
	Title        string `gorm:"not null"`
	Date         string
	Time         string
	Venue        string
	Area         string
	Price        string
	Summary      string
	TicketNote   string
	MapURL       string
	ImageURL     string
	SourceNote   string
	DisplayOrder int
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Translations []EventTranslationModel `gorm:"foreignKey:EventID"`
	Lineup       []EventLineupModel      `gorm:"foreignKey:EventID"`
	Tags         []EventTagModel         `gorm:"foreignKey:EventID"`
}

type EventTranslationModel struct {
	EventID    string `gorm:"primaryKey"`
	Lang       string `gorm:"primaryKey"`
	Title      string
	Summary    string
	TicketNote string
	SourceNote string
}

type EventLineupModel struct {
	ID           uint `gorm:"primaryKey"`
	EventID      string
	Name         string `gorm:"not null"`
	Role         string
	DisplayOrder int
}

type EventTagModel struct {
	ID           uint `gorm:"primaryKey"`
	EventID      string
	Tag          string `gorm:"not null"`
	DisplayOrder int
}

type CatalogItemModel struct {
	ID           string `gorm:"primaryKey"`
	Brand        string `gorm:"not null;default:LIVE LIFE"`
	Format       string `gorm:"not null;index"`
	Artist       string `gorm:"not null"`
	Title        string `gorm:"not null"`
	Summary      string
	Status       string
	Price        string
	ImageURL     string
	PurchaseURL  string
	DisplayOrder int
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Translations []CatalogItemTranslationModel `gorm:"foreignKey:ItemID"`
	Tracks       []CatalogItemTrackModel       `gorm:"foreignKey:ItemID"`
}

type CatalogItemTranslationModel struct {
	ItemID       string `gorm:"primaryKey"`
	Lang         string `gorm:"primaryKey"`
	Title        string
	Summary      string
	PurchaseText string
}

type CatalogItemTrackModel struct {
	ID           uint `gorm:"primaryKey"`
	ItemID       string
	SideLabel    string
	Position     int
	Title        string `gorm:"not null"`
	DurationText string
}

type ContentItemModel struct {
	ID           string `gorm:"primaryKey"`
	Brand        string `gorm:"not null;default:LIVE LIFE"`
	Type         string `gorm:"not null;index"`
	Title        string `gorm:"not null"`
	Summary      string
	DisplayOrder int
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Translations []ContentTranslationModel `gorm:"foreignKey:ContentID"`
}

type ContentTranslationModel struct {
	ContentID string `gorm:"primaryKey"`
	Lang      string `gorm:"primaryKey"`
	Title     string
	Summary   string
}

type ConnectMessageModel struct {
	ID           string `gorm:"primaryKey"`
	Brand        string `gorm:"not null;default:LIVE LIFE"`
	Nickname     string `gorm:"not null"`
	Email        string `gorm:"not null"`
	Topic        string `gorm:"not null;index"`
	Message      string
	Status       string `gorm:"not null;default:new;index"`
	UserLocale   string
	SourcePage   string
	InternalNote string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	HandledAt    *time.Time
}

type Server struct {
	mux    *http.ServeMux
	static http.Handler
	db     *gorm.DB
}

func main() {
	port := getenv("BACKEND_PORT", getenv("PORT", "8080"))
	server, err := NewFileServer(defaultDatabasePath())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%s local server listening on http://localhost:%s", brandName, port)
	if err := http.ListenAndServe(":"+port, server); err != nil {
		log.Fatal(err)
	}
}

// NewServer 保持测试兼容：默认使用内存 SQLite，测试不会污染本地文件。
func NewServer() *Server {
	server, err := newServer(":memory:")
	if err != nil {
		panic(err)
	}
	return server
}

// NewFileServer 给本地预览和未来部署使用，会把数据落到 SQLite 文件。
func NewFileServer(databasePath string) (*Server, error) {
	return newServer(databasePath)
}

func newServer(databasePath string) (*Server, error) {
	db, err := openDatabase(databasePath)
	if err != nil {
		return nil, err
	}
	if err := migrateDatabase(db); err != nil {
		return nil, err
	}
	if err := seedDatabase(db); err != nil {
		return nil, err
	}

	mux := http.NewServeMux()
	staticDir := filepath.Join(".", "static")

	s := &Server{
		mux:    mux,
		static: http.FileServer(http.Dir(staticDir)),
		db:     db,
	}

	mux.HandleFunc("GET /api/health", s.handleHealth)
	mux.HandleFunc("GET /api/events", s.handleEvents)
	mux.HandleFunc("GET /api/cd-items", s.handleCDItems)
	mux.HandleFunc("GET /api/contents", s.handleContents)
	mux.HandleFunc("POST /api/connect", s.handleConnect)
	mux.HandleFunc("POST /api/join", s.handleConnect)
	mux.HandleFunc("/", s.handleStatic)

	return s, nil
}

func openDatabase(databasePath string) (*gorm.DB, error) {
	if databasePath != ":memory:" {
		if err := os.MkdirAll(filepath.Dir(databasePath), 0o755); err != nil {
			return nil, err
		}
	}

	db, err := gorm.Open(sqlite.Open(databasePath), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if err := db.Exec("PRAGMA foreign_keys = ON").Error; err != nil {
		return nil, err
	}
	return db, nil
}

func migrateDatabase(db *gorm.DB) error {
	return db.AutoMigrate(
		&EventModel{},
		&EventTranslationModel{},
		&EventLineupModel{},
		&EventTagModel{},
		&CatalogItemModel{},
		&CatalogItemTranslationModel{},
		&CatalogItemTrackModel{},
		&ContentItemModel{},
		&ContentTranslationModel{},
		&ConnectMessageModel{},
	)
}

func seedDatabase(db *gorm.DB) error {
	var count int64
	if err := db.Model(&EventModel{}).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	return db.Transaction(func(tx *gorm.DB) error {
		for index, event := range seedEvents() {
			if err := createSeedEvent(tx, event, index); err != nil {
				return err
			}
		}
		for index, item := range seedCDItems() {
			if err := createSeedCatalogItem(tx, item, index); err != nil {
				return err
			}
		}
		for index, item := range seedContents() {
			if err := createSeedContent(tx, item, index); err != nil {
				return err
			}
		}
		return nil
	})
}

func createSeedEvent(tx *gorm.DB, event Event, order int) error {
	model := EventModel{
		ID:           event.ID,
		Brand:        event.Brand,
		Owned:        event.Owned,
		Category:     event.Category,
		Title:        event.Title,
		Date:         event.Date,
		Time:         event.Time,
		Venue:        event.Venue,
		Area:         event.Area,
		Price:        event.Price,
		Summary:      event.Summary,
		TicketNote:   event.TicketNote,
		MapURL:       event.MapURL,
		ImageURL:     event.ImageURL,
		SourceNote:   event.SourceNote,
		DisplayOrder: order,
	}
	if err := tx.Create(&model).Error; err != nil {
		return err
	}
	for lang, title := range event.TitleI18n {
		translation := EventTranslationModel{
			EventID:    event.ID,
			Lang:       lang,
			Title:      title,
			Summary:    event.SummaryI18n[lang],
			TicketNote: event.TicketNoteI18n[lang],
			SourceNote: event.SourceNoteI18n[lang],
		}
		if err := tx.Create(&translation).Error; err != nil {
			return err
		}
	}
	for index, name := range event.Lineup {
		lineup := EventLineupModel{EventID: event.ID, Name: name, DisplayOrder: index}
		if err := tx.Create(&lineup).Error; err != nil {
			return err
		}
	}
	for index, tag := range event.Tags {
		tagModel := EventTagModel{EventID: event.ID, Tag: tag, DisplayOrder: index}
		if err := tx.Create(&tagModel).Error; err != nil {
			return err
		}
	}
	return nil
}

func createSeedCatalogItem(tx *gorm.DB, item CatalogItem, order int) error {
	model := CatalogItemModel{
		ID:           item.ID,
		Brand:        item.Brand,
		Format:       item.Format,
		Artist:       item.Artist,
		Title:        item.Title,
		Summary:      item.Summary,
		Status:       item.Status,
		Price:        item.Price,
		ImageURL:     item.ImageURL,
		PurchaseURL:  item.PurchaseURL,
		DisplayOrder: order,
	}
	if err := tx.Create(&model).Error; err != nil {
		return err
	}
	for lang, title := range item.TitleI18n {
		translation := CatalogItemTranslationModel{
			ItemID:       item.ID,
			Lang:         lang,
			Title:        title,
			Summary:      item.SummaryI18n[lang],
			PurchaseText: item.PurchaseText[lang],
		}
		if err := tx.Create(&translation).Error; err != nil {
			return err
		}
	}
	for index, track := range item.Tracks {
		trackModel := CatalogItemTrackModel{ItemID: item.ID, Position: index + 1, Title: track}
		if err := tx.Create(&trackModel).Error; err != nil {
			return err
		}
	}
	return nil
}

func createSeedContent(tx *gorm.DB, item ContentItem, order int) error {
	model := ContentItemModel{
		ID:           item.ID,
		Brand:        item.Brand,
		Type:         item.Type,
		Title:        item.Title,
		Summary:      item.Summary,
		DisplayOrder: order,
	}
	if err := tx.Create(&model).Error; err != nil {
		return err
	}
	for lang, title := range item.TitleI18n {
		translation := ContentTranslationModel{
			ContentID: item.ID,
			Lang:      lang,
			Title:     title,
			Summary:   item.SummaryI18n[lang],
		}
		if err := tx.Create(&translation).Error; err != nil {
			return err
		}
	}
	return nil
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
	events, err := s.listEvents()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load events")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"brand":             brandName,
		"events":            events,
		"ownedEvents":       filterEvents(events, true),
		"recommendedEvents": filterEvents(events, false),
	})
}

func (s *Server) handleCDItems(w http.ResponseWriter, r *http.Request) {
	items, err := s.listCatalogItems()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load catalog items")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"brand": brandName,
		"items": items,
		"cd":    filterCatalog(items, "cd"),
		"vinyl": filterCatalog(items, "vinyl"),
	})
}

func (s *Server) handleContents(w http.ResponseWriter, r *http.Request) {
	contents, err := s.listContents()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load contents")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"brand":    brandName,
		"contents": contents,
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

	message := ConnectMessageModel{
		ID:        fmt.Sprintf("conn_%d", time.Now().UTC().UnixNano()),
		Brand:     brandName,
		Nickname:  req.Nickname,
		Email:     req.Email,
		Topic:     req.Topic,
		Message:   req.Message,
		Status:    "new",
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	if err := s.db.Create(&message).Error; err != nil {
		writeError(w, http.StatusInternalServerError, "failed to save connect message")
		return
	}

	writeJSON(w, http.StatusAccepted, map[string]any{
		"accepted":  true,
		"brand":     brandName,
		"message":   "LIVE LIFE received your message.",
		"messageId": message.ID,
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

func (s *Server) listEvents() ([]Event, error) {
	var models []EventModel
	err := s.db.
		Preload("Translations").
		Preload("Lineup").
		Preload("Tags").
		Order("owned DESC").
		Order("display_order ASC").
		Order("id ASC").
		Find(&models).Error
	if err != nil {
		return nil, err
	}

	events := make([]Event, 0, len(models))
	for _, model := range models {
		sort.Slice(model.Lineup, func(i, j int) bool {
			return model.Lineup[i].DisplayOrder < model.Lineup[j].DisplayOrder
		})
		sort.Slice(model.Tags, func(i, j int) bool {
			return model.Tags[i].DisplayOrder < model.Tags[j].DisplayOrder
		})

		event := Event{
			ID:             model.ID,
			Brand:          model.Brand,
			Owned:          model.Owned,
			Category:       model.Category,
			Title:          model.Title,
			TitleI18n:      LocalizedText{},
			Date:           model.Date,
			Time:           model.Time,
			Venue:          model.Venue,
			Area:           model.Area,
			Price:          model.Price,
			Tags:           []string{},
			Summary:        model.Summary,
			SummaryI18n:    LocalizedText{},
			Lineup:         []string{},
			TicketNote:     model.TicketNote,
			TicketNoteI18n: LocalizedText{},
			MapURL:         model.MapURL,
			ImageURL:       model.ImageURL,
			SourceNote:     model.SourceNote,
			SourceNoteI18n: LocalizedText{},
		}
		for _, translation := range model.Translations {
			event.TitleI18n[translation.Lang] = translation.Title
			event.SummaryI18n[translation.Lang] = translation.Summary
			event.TicketNoteI18n[translation.Lang] = translation.TicketNote
			event.SourceNoteI18n[translation.Lang] = translation.SourceNote
		}
		for _, lineup := range model.Lineup {
			event.Lineup = append(event.Lineup, lineup.Name)
		}
		for _, tag := range model.Tags {
			event.Tags = append(event.Tags, tag.Tag)
		}
		events = append(events, event)
	}
	return events, nil
}

func (s *Server) listCatalogItems() ([]CatalogItem, error) {
	var models []CatalogItemModel
	err := s.db.
		Preload("Translations").
		Preload("Tracks").
		Order("display_order ASC").
		Order("id ASC").
		Find(&models).Error
	if err != nil {
		return nil, err
	}

	items := make([]CatalogItem, 0, len(models))
	for _, model := range models {
		sort.Slice(model.Tracks, func(i, j int) bool {
			return model.Tracks[i].Position < model.Tracks[j].Position
		})

		item := CatalogItem{
			ID:           model.ID,
			Brand:        model.Brand,
			Format:       model.Format,
			Artist:       model.Artist,
			Title:        model.Title,
			TitleI18n:    LocalizedText{},
			Summary:      model.Summary,
			SummaryI18n:  LocalizedText{},
			Tracks:       []string{},
			Status:       model.Status,
			Price:        model.Price,
			ImageURL:     model.ImageURL,
			PurchaseURL:  model.PurchaseURL,
			PurchaseText: LocalizedText{},
		}
		for _, translation := range model.Translations {
			item.TitleI18n[translation.Lang] = translation.Title
			item.SummaryI18n[translation.Lang] = translation.Summary
			item.PurchaseText[translation.Lang] = translation.PurchaseText
		}
		for _, track := range model.Tracks {
			item.Tracks = append(item.Tracks, track.Title)
		}
		items = append(items, item)
	}
	return items, nil
}

func (s *Server) listContents() ([]ContentItem, error) {
	var models []ContentItemModel
	err := s.db.
		Preload("Translations").
		Order("display_order ASC").
		Order("id ASC").
		Find(&models).Error
	if err != nil {
		return nil, err
	}

	contents := make([]ContentItem, 0, len(models))
	for _, model := range models {
		item := ContentItem{
			ID:          model.ID,
			Brand:       model.Brand,
			Type:        model.Type,
			Title:       model.Title,
			TitleI18n:   LocalizedText{},
			Summary:     model.Summary,
			SummaryI18n: LocalizedText{},
		}
		for _, translation := range model.Translations {
			item.TitleI18n[translation.Lang] = translation.Title
			item.SummaryI18n[translation.Lang] = translation.Summary
		}
		contents = append(contents, item)
	}
	return contents, nil
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
	if req.Topic == "" {
		return errors.New("topic is required")
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
				"ja": "ホームのテクスチャは実在バンド名を文化的座標として使い、コードやバイナリの代わりに抽象的な音軌、波形、サンプルトラックを使います。",
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
