package repositories

import (
	"errors"
	"github.com/Dadard29/go-warehouse/api"
	"github.com/Dadard29/go-warehouse/models"
	"time"
)

const (
	listLimit = 5
)

func musicExists(title string, artist string) bool {
	_, err := MusicGetFromTitle(title, artist)
	return err == nil
}

func MusicGetFromTitle(title string, artist string) (models.MusicEntity, error) {
	var f models.MusicEntity
	var m models.MusicEntity
	api.Api.Database.Orm.Where(&models.MusicEntity{
		Title:  title,
		Artist: artist,
	}).First(&m)

	if m.Title != title && m.Artist != artist {
		return f, errors.New("music not found")
	}

	return m, nil
}

func MusicDelete(title string, artist string) (models.MusicEntity, error) {
	var f models.MusicEntity
	if !musicExists(title, artist) {
		return f, errors.New("music not found")
	}

	m, err := MusicGetFromTitle(title, artist)
	if err != nil {
		return f, err
	}

	api.Api.Database.Orm.Delete(&m)

	if musicExists(title, artist) {
		return f, errors.New("error deleting music")
	}

	return m, nil
}

func MusicCreate(token string, mp models.MusicParam, t models.Tags) (models.MusicEntity, error) {
	var f models.MusicEntity

	if musicExists(t.Title, t.Artist) {
		return f, errors.New("music already exists")
	}

	var m = models.MusicEntity{
		Title:       t.Title,
		Artist:      t.Artist,
		Album:       t.Album,
		PublishedAt: t.PublishedAt,
		Genre:       t.Genre,
		ImageUrl:    mp.ImageUrl,
		AddedAt:     time.Now(),
		AddedBy:     token,
	}
	api.Api.Database.Orm.Create(&m)

	if !musicExists(t.Title, t.Artist) {
		return f, errors.New("error storing in DB")
	}

	return m, nil
}

func MusicList() ([]models.MusicEntity, error) {
	var l []models.MusicEntity
	api.Api.Database.Orm.Order("added_at desc").Limit(listLimit).Find(&l)

	return l, nil
}
