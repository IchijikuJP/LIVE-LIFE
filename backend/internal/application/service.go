package application

import (
	"context"
	"time"

	"livelife/backend/internal/domain"
)

// Repository 是应用层依赖的数据端口。
// 这里不写 GORM、SQL 或 SQLite 细节，方便后续从 SQLite 换到 PostgreSQL/MySQL 时保持业务逻辑不变。
type Repository interface {
	ListEvents(ctx context.Context) ([]domain.Event, error)
	ListCatalogItems(ctx context.Context) ([]domain.CatalogItem, error)
	ListContents(ctx context.Context) ([]domain.ContentItem, error)
	SaveConnectMessage(ctx context.Context, message domain.ConnectMessage) error
}

type ConnectResult struct {
	MessageID string
	Brand     string
	Message   string
}

type Service struct {
	repo Repository
	now  func() time.Time
}

func NewService(repo Repository) *Service {
	return NewServiceWithClock(repo, time.Now)
}

func NewServiceWithClock(repo Repository, now func() time.Time) *Service {
	return &Service{repo: repo, now: now}
}

func (s *Service) ListEvents(ctx context.Context) ([]domain.Event, error) {
	return s.repo.ListEvents(ctx)
}

func (s *Service) ListCatalogItems(ctx context.Context) ([]domain.CatalogItem, error) {
	return s.repo.ListCatalogItems(ctx)
}

func (s *Service) ListContents(ctx context.Context) ([]domain.ContentItem, error) {
	return s.repo.ListContents(ctx)
}

// SubmitConnect 是 Connect 表单的用例入口。
// 表单校验和消息生成放在应用层，HTTP 层只负责解析请求，数据库层只负责保存。
func (s *Service) SubmitConnect(ctx context.Context, req domain.ConnectRequest) (ConnectResult, error) {
	req = domain.NormalizeConnectRequest(req)
	if err := domain.ValidateConnectRequest(req); err != nil {
		return ConnectResult{}, err
	}

	now := s.now().UTC()
	message := domain.ConnectMessage{
		ID:        domain.NewConnectMessageID(now),
		Brand:     domain.BrandName,
		Nickname:  req.Nickname,
		Email:     req.Email,
		Topic:     req.Topic,
		Message:   req.Message,
		Status:    "new",
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := s.repo.SaveConnectMessage(ctx, message); err != nil {
		return ConnectResult{}, err
	}

	return ConnectResult{
		MessageID: message.ID,
		Brand:     domain.BrandName,
		Message:   "LIVE LIFE received your message.",
	}, nil
}
