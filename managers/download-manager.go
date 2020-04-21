package managers

import (
	"github.com/Dadard29/go-warehouse/models"
	"github.com/Dadard29/go-warehouse/repositories"
)

func DownloadGetManager(tags models.Tags) (string, error) {
	return repositories.GetFilePathForDownload(tags)
}
