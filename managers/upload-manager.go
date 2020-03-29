package managers

import (
	"errors"
	"fmt"
	"github.com/Dadard29/go-warehouse/models"
	"github.com/Dadard29/go-warehouse/repositories"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path"
)

const (
	mimeMp3 = "audio/mpeg"

	maxMegaBytes = 10
	maxSize = maxMegaBytes << (10 * 2)

	maxFilesNumber = 10
)

func cleanTempFile(path string) {
	err := os.Remove(path)
	if err != nil {
		logger.Error(err.Error())
	}
}

func FileListManager(token string) (models.FileListJson, error) {
	var flJson models.FileListJson

	files, err := repositories.ListFiles(token)
	if err != nil {
		logger.Error("error listing files")
		return flJson, err
	}

	flJson.List = files
	flJson.Quota = maxFilesNumber
	flJson.Current = len(files)

	return flJson, err
}

func FileDeleteManager(token string, tags models.Tags) (models.File, error) {
	var f models.File

	fileDeleted, err := repositories.RemoveFile(token, tags)
	if err != nil {
		logger.Error(err.Error())
		return f, errors.New("error while deleting file")
	}

	return fileDeleted, nil
}

func FileStoreManager(token string, file multipart.File, headers *multipart.FileHeader) (models.File, error) {
	var f models.File

	defer file.Close()

	// check quota
	storedFiles, err := repositories.ListFiles(token)
	if err != nil {
		logger.Error("error checking existing stored files")
		return f, err
	}

	storedFilesCount := len(storedFiles)
	if storedFilesCount > maxFilesNumber {
		return f, errors.New(fmt.Sprintf("file number quota already reached (%d/%d)",
			storedFilesCount, maxFilesNumber))
	}

	// check size
	if headers.Size > maxSize {
		return f, errors.New(fmt.Sprintf(
			"file too big: maximum allowed is %d Mb", maxMegaBytes))
	}

	// check mime
	if headers.Header.Get("Content-Type") != mimeMp3 {
		return f, errors.New("bad mime")
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		logger.Error("error reading file")
		return f, err
	}

	tempFileName := "tempfile.mp3"
	tempFilePath := path.Join(repositories.Tmp, tempFileName)
	err = ioutil.WriteFile(tempFilePath, data, 0644)
	if err != nil {
		logger.Error("error writing file")
		return f, err
	}

	// check is audio
	if !repositories.CheckFileAudio(tempFilePath) {
		cleanTempFile(tempFilePath)

		msg := "not an audio file"
		logger.Error(msg)
		return f, errors.New(msg)
	}

	// check if mp3 by reading ID3V2 tag
	tags, err := repositories.ReadTags(tempFilePath)
	if err != nil {
		cleanTempFile(tempFilePath)

		logger.Error(err.Error())
		return f, errors.New("error reading id3v2 tags")
	}

	var fileAdded models.File
	if fileAdded, err = repositories.AddFile(tempFilePath, token, tags); err != nil {
		cleanTempFile(tempFilePath)

		logger.Error(err.Error())
		return f, errors.New("error storing file in library")
	}

	return fileAdded, nil
}
