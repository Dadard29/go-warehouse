package repositories

import (
	"errors"
	"fmt"
	"github.com/Dadard29/go-warehouse/models"
	"io/ioutil"
	"os"
	"path"
)

const (
	baseDirStore = "store"
	Tmp = "tmp"
	mp3Extension = ".mp3"

)

func getFullFilePath(token string, tags models.Tags) string {
	fingerprint := token[:6]
	return path.Join(baseDirStore, tags.Artist, tags.Album, fingerprint + "_" + tags.Title + mp3Extension)
}

func GetFilePathForDownload(token string, tags models.Tags) (string, error) {
	var p string
	filePath := getFullFilePath(token, tags)

	if _, err := os.Stat(filePath); err != nil {
		return p, err
	}

	return filePath, nil
}

func AddFile(srcPath string, token string, tags models.Tags) (models.File, error) {

	var f models.File

	//if err := checkPlaceholder(token); err != nil {
	//	return f, err
	//}

	if err := checkArtist(tags.Artist); err != nil {
		return f, err
	}

	if err := checkAlbum(tags.Artist, tags.Album); err != nil {
		return f, err
	}

	if checkFileExist(token, tags) {
		return f, errors.New(fmt.Sprintf("file %s already exists", tags.Title))
	}


	outputPath := getFullFilePath(token, tags)
	if err := os.Rename(srcPath, outputPath); err != nil {
		return f, err
	}

	infos, _ := os.Stat(outputPath)

	return models.File{
		Filename: infos.Name(),
		AddedAt:  infos.ModTime(),
		Metadata: tags,
	}, nil
}

func RemoveFile(token string, tags models.Tags) (models.File, error) {
	var f models.File

	p := getFullFilePath(token, tags)

	infos, err := os.Stat(p)
	if err != nil {
		return f, err
	}

	err = os.Remove(p)
	if err != nil {
		return f, err
	}

	return models.File{
		Filename: infos.Name(),
		AddedAt:  infos.ModTime(),
		Metadata: tags,
	}, nil
}

func ListFiles(token string) ([]models.File, error) {
	var l = make([]models.File, 0)

	// read all artists
	artistDirs, err := ioutil.ReadDir(path.Join(baseDirStore))
	if err != nil {
		logger.Error("error listing artists")
		return nil, err
	}
	for _, ar := range artistDirs {
		if !ar.IsDir() {
			continue
		}

		// read all albums
		albumDirs, err := ioutil.ReadDir(path.Join(baseDirStore, ar.Name()))
		if err != nil {
			logger.Error("error listing album for artist " + ar.Name())
			return nil, err
		}
		for _, al := range albumDirs {

			// read all titles
			titleDirs, err := ioutil.ReadDir(path.Join(baseDirStore, ar.Name(), al.Name()))
			if err != nil {
				logger.Error("error listing titles for artist " + ar.Name() + " and album " + al.Name())
				return nil, err
			}
			for _, t := range titleDirs {
				tags, err := ReadTags(
					path.Join(baseDirStore, ar.Name(), al.Name(), t.Name()))
				if err != nil {
					logger.Error("error reading tags of file " + t.Name())
					logger.Error(err.Error())
					continue
				}

				l = append(l, models.File{
					Filename: t.Name(),
					AddedAt:  t.ModTime(),
					Metadata: tags,
				})
			}

		}
	}

	return l, nil
}
