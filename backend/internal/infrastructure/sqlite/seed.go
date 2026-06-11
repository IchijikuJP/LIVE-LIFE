package sqlite

import (
	"gorm.io/gorm"

	"livelife/backend/internal/domain"
)

func seedDatabase(db *gorm.DB) error {
	return db.Transaction(func(tx *gorm.DB) error {
		for index, event := range seedEvents() {
			if err := createSeedEventIfMissing(tx, event, index); err != nil {
				return err
			}
		}
		for index, item := range seedCDItems() {
			if err := createSeedCatalogItemIfMissing(tx, item, index); err != nil {
				return err
			}
		}
		for index, item := range seedContents() {
			if err := createSeedContentIfMissing(tx, item, index); err != nil {
				return err
			}
		}
		return nil
	})
}

func createSeedEventIfMissing(tx *gorm.DB, event domain.Event, order int) error {
	var count int64
	if err := tx.Model(&EventModel{}).Where("id = ?", event.ID).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return updateSeedEvent(tx, event, order)
	}
	return createSeedEvent(tx, event, order)
}

func createSeedEvent(tx *gorm.DB, event domain.Event, order int) error {
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
	return replaceSeedEventChildren(tx, event)
}

func updateSeedEvent(tx *gorm.DB, event domain.Event, order int) error {
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
	if err := tx.Save(&model).Error; err != nil {
		return err
	}
	if err := tx.Where("event_id = ?", event.ID).Delete(&EventTranslationModel{}).Error; err != nil {
		return err
	}
	if err := tx.Where("event_id = ?", event.ID).Delete(&EventLineupModel{}).Error; err != nil {
		return err
	}
	if err := tx.Where("event_id = ?", event.ID).Delete(&EventTagModel{}).Error; err != nil {
		return err
	}
	return replaceSeedEventChildren(tx, event)
}

func replaceSeedEventChildren(tx *gorm.DB, event domain.Event) error {
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

func createSeedCatalogItemIfMissing(tx *gorm.DB, item domain.CatalogItem, order int) error {
	var count int64
	if err := tx.Model(&CatalogItemModel{}).Where("id = ?", item.ID).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return updateSeedCatalogItem(tx, item, order)
	}
	return createSeedCatalogItem(tx, item, order)
}

func createSeedCatalogItem(tx *gorm.DB, item domain.CatalogItem, order int) error {
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
	return replaceSeedCatalogItemChildren(tx, item)
}

func updateSeedCatalogItem(tx *gorm.DB, item domain.CatalogItem, order int) error {
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
	if err := tx.Save(&model).Error; err != nil {
		return err
	}
	if err := tx.Where("item_id = ?", item.ID).Delete(&CatalogItemTranslationModel{}).Error; err != nil {
		return err
	}
	if err := tx.Where("item_id = ?", item.ID).Delete(&CatalogItemTrackModel{}).Error; err != nil {
		return err
	}
	return replaceSeedCatalogItemChildren(tx, item)
}

func replaceSeedCatalogItemChildren(tx *gorm.DB, item domain.CatalogItem) error {
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

func createSeedContentIfMissing(tx *gorm.DB, item domain.ContentItem, order int) error {
	var count int64
	if err := tx.Model(&ContentItemModel{}).Where("id = ?", item.ID).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return updateSeedContent(tx, item, order)
	}
	return createSeedContent(tx, item, order)
}

func createSeedContent(tx *gorm.DB, item domain.ContentItem, order int) error {
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
	return replaceSeedContentChildren(tx, item)
}

func updateSeedContent(tx *gorm.DB, item domain.ContentItem, order int) error {
	model := ContentItemModel{
		ID:           item.ID,
		Brand:        item.Brand,
		Type:         item.Type,
		Title:        item.Title,
		Summary:      item.Summary,
		DisplayOrder: order,
	}
	if err := tx.Save(&model).Error; err != nil {
		return err
	}
	if err := tx.Where("content_id = ?", item.ID).Delete(&ContentTranslationModel{}).Error; err != nil {
		return err
	}
	return replaceSeedContentChildren(tx, item)
}

func replaceSeedContentChildren(tx *gorm.DB, item domain.ContentItem) error {
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

func seedEvents() []domain.Event {
	return []domain.Event{
		{
			ID:       "redhair-japan-2026-july",
			Brand:    domain.BrandName,
			Owned:    true,
			Category: "own-live",
			Title:    "LIVE LIFE presents Red Hair Boy Murder Case in Tokyo",
			TitleI18n: domain.LocalizedText{
				"zh": "LIVE LIFE presents 紅髮少年殺人事件 东京双日演出",
				"ja": "LIVE LIFE presents 紅髪少年殺人事件 東京2公演",
				"en": "LIVE LIFE PRESENTS RED HAIR BOY MURDER CASE IN TOKYO",
			},
			Date:  "2026-07-10 / 2026-07-14",
			Time:  "7/10 OPEN 18:45 START 19:30; 7/14 OPEN 19:00 START 19:30",
			Venue: "GRIT at Shibuya / Shimokitazawa THREE",
			Area:  "Shibuya / Shimokitazawa",
			Price: "7/10 adv ¥5,000 + 1D, door ¥5,500 + 1D; 7/14 adv ¥4,000 + 1D, door ¥4,500 + 1D",
			Tags: []string{
				"LIVE LIFE",
				"own live",
				"alternative rock",
				"Tokyo",
				"向井秀徳アコースティック&エレクトリック",
				"ルサンチマン",
				"おそロシア革命",
				"dj:じん",
			},
			Summary: "Two Tokyo shows by 紅髪少年殺人事件. The July 10 Shibuya guest is 向井秀徳アコースティック&エレクトリック. The July 14 Shimokitazawa show features ルサンチマン, おそロシア革命, and dj:じん.",
			SummaryI18n: domain.LocalizedText{
				"zh": "7月10日和7月14日，LIVE LIFE 将在东京分别呈现两场 紅髮少年殺人事件 演出。7月10日涩谷场嘉宾为「向井秀徳アコースティック&エレクトリック」；7月14日下北泽场共演为「ルサンチマン」「おそロシア革命」，DJ 为「dj:じん」。",
				"ja": "7月10日と7月14日、LIVE LIFE は東京で 紅髪少年殺人事件 の2公演を行います。7月10日の渋谷公演には「向井秀徳アコースティック&エレクトリック」、7月14日の下北沢公演には「ルサンチマン」「おそロシア革命」、DJ として「dj:じん」が出演します。",
				"en": "ON JULY 10 AND JULY 14, LIVE LIFE PRESENTS TWO TOKYO SHOWS BY RED HAIR BOY MURDER CASE. THE JULY 10 SHIBUYA GUEST IS MUKAI SHUTOKU ACOUSTIC & ELECTRIC. THE JULY 14 SHIMOKITAZAWA SHOW FEATURES RUSANTIMAN, OSOROSHIA KAKUMEI, AND DJ:JIN.",
			},
			Lineup: []string{
				"紅髪少年殺人事件",
				"向井秀徳アコースティック&エレクトリック",
				"ルサンチマン",
				"おそロシア革命",
				"dj:じん",
			},
			TicketNote: "Ticket links are pending. Keep the external ticket agency flow separate from the LIVE LIFE site.",
			TicketNoteI18n: domain.LocalizedText{
				"zh": "票务链接待确认。演出票站可能由外部代理处理，LIVE LIFE 站内先展示情报，不直接结算。",
				"ja": "チケットリンクは確認中です。外部プレイガイドの導線と LIVE LIFE サイト内決済は分けて設計します。",
				"en": "TICKET LINKS ARE PENDING. THE TICKETING AGENCY FLOW STAYS OUTSIDE LIVE LIFE CHECKOUT.",
			},
			MapURL:   "https://maps.google.com/?q=GRIT+at+Shibuya",
			ImageURL: "/assets/events/redhair-2026-july.jpg",
			SourceNoteI18n: domain.LocalizedText{
				"zh": "活动细节来自你提供的海报和文本。第一轮公开检索暂未找到稳定的官方活动页。",
				"ja": "イベント詳細は提供されたフライヤーと本文に基づきます。初回の公開検索では安定して参照できる公式イベントページは見つかっていません。",
				"en": "DETAILS ARE BASED ON THE PROVIDED POSTER AND COPY. A STABLE OFFICIAL EVENT PAGE WAS NOT FOUND IN THE FIRST PUBLIC SEARCH PASS.",
			},
		},
		{
			ID:       "wednesday-wonderland-archive-2025",
			Brand:    domain.BrandName,
			Owned:    false,
			Category: "archive-reference",
			Title:    "Wednesday Wonderland archive visual",
			TitleI18n: domain.LocalizedText{
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
			SummaryI18n: domain.LocalizedText{
				"zh": "这张图先作为 LIVE LIFE 活动档案和视觉气质参考，不作为即将发生的演出推荐。",
				"ja": "この画像は LIVE LIFE のイベントアーカイブとビジュアル参考として置き、直近公演の推薦情報とは分けます。",
				"en": "THIS POSTER IS KEPT AS A LIVE LIFE ARCHIVE AND VISUAL REFERENCE, SEPARATE FROM UPCOMING LIVE RECOMMENDATIONS.",
			},
			Lineup: []string{"Wednesday Wonderland", "TiDE", "できないみらい", "いきものコーナー"},
			TicketNoteI18n: domain.LocalizedText{
				"zh": "历史活动，仅用于样式参考。",
				"ja": "過去イベントのため、スタイル参考のみ。",
				"en": "PAST EVENT, USED ONLY AS A STYLE REFERENCE.",
			},
			MapURL:   "https://maps.google.com/?q=BASEMENT+BAR+Shimokitazawa",
			ImageURL: "/assets/events/wednesday-wonderland-2025-08-21.jpg",
			SourceNoteI18n: domain.LocalizedText{
				"zh": "信息来自你提供的海报。",
				"ja": "提供されたフライヤーに基づく情報です。",
				"en": "INFORMATION IS BASED ON THE PROVIDED POSTER.",
			},
		},
	}
}

func seedCDItems() []domain.CatalogItem {
	return []domain.CatalogItem{
		{
			ID:     "redhair-demo-cd",
			Brand:  domain.BrandName,
			Format: "cd",
			Artist: "紅髪少年殺人事件",
			Title:  "Selected CD placeholder",
			TitleI18n: domain.LocalizedText{
				"zh": "紅髮少年殺人事件 CD 严选占位",
				"ja": "紅髪少年殺人事件 CD厳選仮枠",
				"en": "RED HAIR BOY MURDER CASE HAND-PICKED CD",
			},
			Summary: "A hand-picked CD slot for show-related releases. Purchase goes to an external shop.",
			SummaryI18n: domain.LocalizedText{
				"zh": "这里作为 CD 严选的 CD 分类占位。单品详情页会放「点击此处购买」按钮，跳转到 BASE 等外部 Shop。",
				"ja": "CD厳選内のCDカテゴリ仮枠です。詳細ページの購入ボタンから BASE など外部ショップへ遷移します。",
				"en": "A HAND-PICKED CD SLOT. THE DETAIL PAGE PURCHASE BUTTON LINKS TO AN EXTERNAL SHOP SUCH AS BASE.",
			},
			Tracks:      []string{"A-side reference", "Live note", "Shop link pending"},
			Status:      "external shop",
			Price:       "TBD",
			ImageURL:    "/assets/events/redhair-2026-july.jpg",
			PurchaseURL: "https://thebase.com/",
			PurchaseText: domain.LocalizedText{
				"zh": "点击此处购买",
				"ja": "こちらから購入",
				"en": "BUY HERE",
			},
		},
		{
			ID:     "live-life-vinyl-placeholder",
			Brand:  domain.BrandName,
			Format: "vinyl",
			Artist: "LIVE LIFE SELECT",
			Title:  "Selected vinyl placeholder",
			TitleI18n: domain.LocalizedText{
				"zh": "LIVE LIFE 黑胶严选占位",
				"ja": "LIVE LIFE ヴァイナルセレクト仮枠",
				"en": "LIVE LIFE VINYL SELECT",
			},
			Summary: "A vinyl slot for future selected records.",
			SummaryI18n: domain.LocalizedText{
				"zh": "这里作为黑胶分类占位，后续放精选黑胶、推荐语和外部购买链接。",
				"ja": "ヴァイナルカテゴリの仮枠です。今後、推薦盤、紹介文、外部購入リンクを掲載します。",
				"en": "A VINYL CATEGORY SLOT FOR SELECTED RECORDS, NOTES, AND EXTERNAL PURCHASE LINKS.",
			},
			Tracks:      []string{"Side A", "Side B", "Listening note pending"},
			Status:      "curating",
			Price:       "TBD",
			PurchaseURL: "https://thebase.com/",
			PurchaseText: domain.LocalizedText{
				"zh": "点击此处购买",
				"ja": "こちらから購入",
				"en": "BUY HERE",
			},
		},
	}
}

func seedContents() []domain.ContentItem {
	return []domain.ContentItem{
		{
			ID:    "homepage-positioning",
			Brand: domain.BrandName,
			Type:  "archive",
			Title: "What LIVE LIFE is building",
			TitleI18n: domain.LocalizedText{
				"zh": "LIVE LIFE 正在做什么",
				"ja": "LIVE LIFE が作っているもの",
				"en": "WHAT LIVE LIFE IS BUILDING",
			},
			Summary: "LIVE LIFE is a Tokyo-facing music entry point for shows, CD select, archive, and support messages.",
			SummaryI18n: domain.LocalizedText{
				"zh": "LIVE LIFE 会先成为东京演出情报、LIVE LIFE 自主演出、CD 严选、档案馆和 Connect 的统一入口。",
				"ja": "LIVE LIFE は東京のライブ情報、自主公演、CD厳選、Archive、Connect をまとめる入口として始めます。",
				"en": "LIVE LIFE STARTS AS ONE ENTRY POINT FOR TOKYO LIVES, OWNED EVENTS, HAND-PICKED CD, ARCHIVE, AND CONNECT.",
			},
		},
		{
			ID:    "visual-system-note",
			Brand: domain.BrandName,
			Type:  "design",
			Title: "Visual system note",
			TitleI18n: domain.LocalizedText{
				"zh": "视觉系统备注",
				"ja": "ビジュアルシステムメモ",
				"en": "VISUAL SYSTEM NOTE",
			},
			Summary: "The homepage uses real band names as cultural texture and abstract music production tracks as data flow.",
			SummaryI18n: domain.LocalizedText{
				"zh": "首页纹理使用真实乐队名作为文化坐标，同时用抽象音轨、波形和采样轨道替代代码/二进制。",
				"ja": "ホームのテクスチャは実在バンド名を文化的座標として使い、コードやバイナリの代わりに抽象的な音軌、波形、サンプルトラックを使います。",
				"en": "THE HOMEPAGE USES REAL BAND NAMES AS CULTURAL TEXTURE, WITH ABSTRACT TRACK LANES AND WAVEFORM-STYLE DATA INSTEAD OF CODE OR BINARY.",
			},
		},
	}
}
