package repositories

import (
	"errors"
	"fmt"
	"github.com/Dadard29/go-warehouse/models"
	"github.com/bogem/id3v2"
	"github.com/h2non/filetype"
	"io/ioutil"
	"os"
	"path"
)

// create placeholder if needed
func checkPlaceholder(token string) error {
	path2check := path.Join(baseDirStore, token)
	f, err := os.Stat(path2check)
	if err != nil {
		if os.IsNotExist(err) {
			if err := os.Mkdir(path2check, 0755); err != nil {
				return err
			}
			return nil
		} else {
			return err
		}
	}

	if !f.IsDir() {
		return errors.New(fmt.Sprintf("placeholder %s is a file", token))
	}

	return nil
}

func checkArtist(token string, artist string) error {
	path2check := path.Join(baseDirStore, token, artist)
	f, err := os.Stat(path2check)
	if err != nil {
		if os.IsNotExist(err) {
			if err := os.Mkdir(path2check, 0755); err != nil {
				return err
			}
			return nil
		} else {
			return err
		}
	}

	if !f.IsDir() {
		return errors.New(fmt.Sprintf("artist placeholder %s is a file", artist))
	}

	return nil
}

func checkAlbum(token string, artist string, album string) error {
	path2check := path.Join(baseDirStore, token, artist, album)
	f, err := os.Stat(path2check)
	if err != nil {
		if os.IsNotExist(err) {
			if err := os.Mkdir(path2check, 0755); err != nil {
				return err
			}
			return nil
		} else {
			return err
		}
	}

	if !f.IsDir() {
		return errors.New(fmt.Sprintf("album placeholder %s is a file", album))
	}

	return nil
}


func CheckFileAudio(path string) bool {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return false
	}

	return filetype.IsAudio(buf)
}

func ReadTags(path string) (models.Tags, error) {
	var t models.Tags
	f, err := id3v2.Open(path, id3v2.Options{Parse: true})
	if err != nil {
		return t, err
	}

	if f.Title() == "" {
		return t, errors.New("title tag empty")
	}

	if f.Artist() == "" {
		return t, errors.New("artist tag empty")
	}

	if f.Album() == "" {
		return t, errors.New("album tag empty")
	}

	t.Title = f.Title()
	t.Artist = f.Artist()
	t.Album = f.Album()

	return t, nil
}

// return true if file exist
func checkFileExist(token string, tags models.Tags) bool {
	file2check := getFullFilePath(token, tags)
	_, err := os.Stat(file2check)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		} else if os.IsExist(err) {
			return true
		}
	}

	return true
}