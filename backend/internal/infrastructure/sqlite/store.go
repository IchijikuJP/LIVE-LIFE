package sqlite

import (
	"context"
	"os"
	"path/filepath"
	"sort"

	sqlitedriver "github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"livelife/backend/internal/domain"
)

type Store struct {
	db *gorm.DB
}

// NewStore 打开 SQLite、执行自动迁移，并在空库时写入本地开发种子数据。
// 生产环境仍然可以继续用同一套表结构；只是数据来源会逐步换成后台管理或导入脚本。
func NewStore(databasePath string) (*Store, error) {
	if databasePath != ":memory:" {
		if err := os.MkdirAll(filepath.Dir(databasePath), 0o755); err != nil {
			return nil, err
		}
	}

	db, err := gorm.Open(sqlitedriver.Open(databasePath), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if err := db.Exec("PRAGMA foreign_keys = ON").Error; err != nil {
		return nil, err
	}
	if err := migrateDatabase(db); err != nil {
		return nil, err
	}
	if err := seedDatabase(db); err != nil {
		return nil, err
	}
	return &Store{db: db}, nil
}

func (s *Store) ListEvents(ctx context.Context) ([]domain.Event, error) {
	var models []EventModel
	err := s.db.WithContext(ctx).
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

	events := make([]domain.Event, 0, len(models))
	for _, model := range models {
		sort.Slice(model.Lineup, func(i, j int) bool {
			return model.Lineup[i].DisplayOrder < model.Lineup[j].DisplayOrder
		})
		sort.Slice(model.Tags, func(i, j int) bool {
			return model.Tags[i].DisplayOrder < model.Tags[j].DisplayOrder
		})

		event := domain.Event{
			ID:             model.ID,
			Brand:          model.Brand,
			Owned:          model.Owned,
			Category:       model.Category,
			Title:          model.Title,
			TitleI18n:      domain.LocalizedText{},
			Date:           model.Date,
			Time:           model.Time,
			Venue:          model.Venue,
			Area:           model.Area,
			Price:          model.Price,
			Tags:           []string{},
			Summary:        model.Summary,
			SummaryI18n:    domain.LocalizedText{},
			Lineup:         []string{},
			TicketNote:     model.TicketNote,
			TicketNoteI18n: domain.LocalizedText{},
			MapURL:         model.MapURL,
			ImageURL:       model.ImageURL,
			SourceNote:     model.SourceNote,
			SourceNoteI18n: domain.LocalizedText{},
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

func (s *Store) ListCatalogItems(ctx context.Context) ([]domain.CatalogItem, error) {
	var models []CatalogItemModel
	err := s.db.WithContext(ctx).
		Preload("Translations").
		Preload("Tracks").
		Order("display_order ASC").
		Order("id ASC").
		Find(&models).Error
	if err != nil {
		return nil, err
	}

	items := make([]domain.CatalogItem, 0, len(models))
	for _, model := range models {
		sort.Slice(model.Tracks, func(i, j int) bool {
			return model.Tracks[i].Position < model.Tracks[j].Position
		})

		item := domain.CatalogItem{
			ID:           model.ID,
			Brand:        model.Brand,
			Format:       model.Format,
			Artist:       model.Artist,
			Title:        model.Title,
			TitleI18n:    domain.LocalizedText{},
			Summary:      model.Summary,
			SummaryI18n:  domain.LocalizedText{},
			Tracks:       []string{},
			Status:       model.Status,
			Price:        model.Price,
			ImageURL:     model.ImageURL,
			PurchaseURL:  model.PurchaseURL,
			PurchaseText: domain.LocalizedText{},
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

func (s *Store) ListContents(ctx context.Context) ([]domain.ContentItem, error) {
	var models []ContentItemModel
	err := s.db.WithContext(ctx).
		Preload("Translations").
		Order("display_order ASC").
		Order("id ASC").
		Find(&models).Error
	if err != nil {
		return nil, err
	}

	contents := make([]domain.ContentItem, 0, len(models))
	for _, model := range models {
		item := domain.ContentItem{
			ID:          model.ID,
			Brand:       model.Brand,
			Type:        model.Type,
			Title:       model.Title,
			TitleI18n:   domain.LocalizedText{},
			Summary:     model.Summary,
			SummaryI18n: domain.LocalizedText{},
		}
		for _, translation := range model.Translations {
			item.TitleI18n[translation.Lang] = translation.Title
			item.SummaryI18n[translation.Lang] = translation.Summary
		}
		contents = append(contents, item)
	}
	return contents, nil
}

func (s *Store) SaveConnectMessage(ctx context.Context, message domain.ConnectMessage) error {
	model := ConnectMessageModel{
		ID:        message.ID,
		Brand:     message.Brand,
		Nickname:  message.Nickname,
		Email:     message.Email,
		Topic:     message.Topic,
		Message:   message.Message,
		Status:    message.Status,
		CreatedAt: message.CreatedAt,
		UpdatedAt: message.UpdatedAt,
	}
	return s.db.WithContext(ctx).Create(&model).Error
}

func (s *Store) CountConnectMessages(ctx context.Context) (int64, error) {
	var count int64
	err := s.db.WithContext(ctx).Model(&ConnectMessageModel{}).Count(&count).Error
	return count, err
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
