package sqlite

import "time"

// EventModel 是演出情报主表。
// LIVE LIFE 自主演出和推荐演出共用这张表，用 owned 字段决定前端是否固定展示在顶部。
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

// EventTranslationModel 保存演出多语言文案。
// 一条演出对应 zh/ja/en 三行翻译，前端通过当前语言从 titleI18n、summaryI18n 等字段取展示文本。
type EventTranslationModel struct {
	EventID    string `gorm:"primaryKey"`
	Lang       string `gorm:"primaryKey"`
	Title      string
	Summary    string
	TicketNote string
	SourceNote string
}

// EventLineupModel 保存演出阵容。
// 阵容独立成表，是为了后续能追加出演顺序、角色说明、乐队主页链接等字段。
type EventLineupModel struct {
	ID           uint `gorm:"primaryKey"`
	EventID      string
	Name         string `gorm:"not null"`
	Role         string
	DisplayOrder int
}

// EventTagModel 保存演出标签。
// 标签用于前端做快速扫描，也方便以后做 Archive 或 Shows 的筛选。
type EventTagModel struct {
	ID           uint `gorm:"primaryKey"`
	EventID      string
	Tag          string `gorm:"not null"`
	DisplayOrder int
}

// CatalogItemModel 是 CD 严选的单品主表。
// format 固定承接当前需求里的 cd / vinyl；购买入口不在顶层 Shop，而在单品 purchase_url。
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

// CatalogItemTranslationModel 保存 CD/黑胶单品的多语言标题、介绍和购买按钮文本。
type CatalogItemTranslationModel struct {
	ItemID       string `gorm:"primaryKey"`
	Lang         string `gorm:"primaryKey"`
	Title        string
	Summary      string
	PurchaseText string
}

// CatalogItemTrackModel 保存曲目。
// CD 和黑胶都可以使用同一张表；黑胶后续可用 side_label 区分 A/B 面。
type CatalogItemTrackModel struct {
	ID           uint `gorm:"primaryKey"`
	ItemID       string
	SideLabel    string
	Position     int
	Title        string `gorm:"not null"`
	DurationText string
}

// ContentItemModel 保存 Archive、设计说明、资料页入口等轻量内容。
// 它和 Event/Catalog 分开，避免历史资料影响演出和购买链路。
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

// ContentTranslationModel 保存 Archive/资料内容的多语言文案。
type ContentTranslationModel struct {
	ContentID string `gorm:"primaryKey"`
	Lang      string `gorm:"primaryKey"`
	Title     string
	Summary   string
}

// ConnectMessageModel 保存 Connect 表单消息。
// 这里先统一收票务、购买未收到货、合作、投稿等问题；后台管理上线后再按 topic/status 分流。
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
