package assets

import (
	"context"
	"os"
)

func assetInfo(path string) (os.FileInfo, error) {
	return os.Stat(path)
}

func Download(ctx context.Context) error {
	if err := DownloadFeatues(ctx); err != nil {
		return err
	}
	if err := DownloadModels(ctx); err != nil {
		return err
	}
	return nil
}
