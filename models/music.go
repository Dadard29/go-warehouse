package models

import "time"

// stored in db
type MusicEntity struct {
	Title       string `gorm:"type:varchar(70);index:title"`
	Artist      string `gorm:"type:varchar(70);index:artist"`
	Album       string `gorm:"type:varchar(70);index:album"`
	PublishedAt string `gorm:"type:varchar(20);index:published_at"`
	Genre       string `gorm:"type:varchar(40);index:genre"`
	ImageUrl    string `gorm:"type:varchar(200);index:image_url"`

	AddedAt time.Time `gorm:"type:datetime;index:added_at"`
	AddedBy string    `gorm:"type:varchar(70);index:added_by"`
}

func (MusicEntity) TableName() string {
	return "music"
}

func (m MusicEntity) ToDto() MusicDto {
	return MusicDto{
		Title:       m.Title,
		Artist:      m.Artist,
		Album:       m.Album,
		PublishedAt: m.PublishedAt,
		Genre:       m.Genre,
		ImageUrl:    m.ImageUrl,
		AddedAt:     m.AddedAt,
	}
}

// exposed
type MusicDto struct {
	Title       string `json:"title"`
	Artist      string `json:"artist"`
	Album       string `json:"album"`
	PublishedAt string `json:"published_at"`
	Genre       string `json:"genre"`
	ImageUrl    string `json:"image_url"`

	AddedAt time.Time `json:"added_at"`
}

// input
type MusicParam struct {
	ImageUrl string `json:"image_url"`
}

func (m MusicParam) CheckSanity() bool {
	return m.ImageUrl != ""
}
