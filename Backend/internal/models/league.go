package models

type League struct {
	Base
	Name          string  `gorm:"not null" json:"name"`
	ShortName     string  `json:"short_name"`
	Country       string  `json:"country"`
	LogoURL       string  `json:"logo_url"`
	Type          string  `gorm:"type:varchar(10);default:league;not null" json:"type"` // "league" หรือ "cup"
	SortOrder     int     `gorm:"default:0;not null" json:"sort_order"`
	ExternalAPIID *string `gorm:"uniqueIndex" json:"-"` // id สำหรับซิงก์ external API (ไม่แสดงใน JSON)
}
