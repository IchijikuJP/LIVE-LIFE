package domain

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

const BrandName = "LIVE LIFE"

// LocalizedText 保存对外展示文案的多语言版本。
// 当前固定支持中文、日语、英语。API 一次返回完整文案，前端只负责按当前语言取值。
type LocalizedText map[string]string

// Event 是演出情报的业务实体。
// Owned=true 代表 LIVE LIFE 自主演出，前端和 API 都会把它固定放在演出情报最上方。
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

// CatalogItem 是 CD 严选里的单品实体。
// 当前不设置顶层 Shop；购买路径固定为 CD 严选 -> 单品 -> 外部 Shop。
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

// ContentItem 是 Archive、设计备注、资料页摘要等内容入口。
// 它和演出、商品分开，避免历史资料影响演出和购买流程。
type ContentItem struct {
	ID          string        `json:"id"`
	Brand       string        `json:"brand"`
	Type        string        `json:"type"`
	Title       string        `json:"title"`
	TitleI18n   LocalizedText `json:"titleI18n"`
	Summary     string        `json:"summary"`
	SummaryI18n LocalizedText `json:"summaryI18n"`
}

// ConnectRequest 是 Connect 表单提交的入参。
// 它统一承接票务、购买未收到货、发货、合作、投稿等联系问题。
type ConnectRequest struct {
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Topic    string `json:"topic"`
	Message  string `json:"message"`
}

type ConnectMessage struct {
	ID        string
	Brand     string
	Nickname  string
	Email     string
	Topic     string
	Message   string
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NormalizeConnectRequest(req ConnectRequest) ConnectRequest {
	req.Nickname = strings.TrimSpace(req.Nickname)
	req.Email = strings.TrimSpace(req.Email)
	req.Topic = strings.TrimSpace(req.Topic)
	req.Message = strings.TrimSpace(req.Message)
	return req
}

func ValidateConnectRequest(req ConnectRequest) error {
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

func NewConnectMessageID(now time.Time) string {
	return fmt.Sprintf("conn_%d", now.UTC().UnixNano())
}
