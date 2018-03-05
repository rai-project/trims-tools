package assets

import (
	"context"
	"io/ioutil"
	"path/filepath"

	"github.com/rai-project/micro18-tools/pkg/config"
)

var (
	Features = map[string]string{}
)

func DownloadFeatues(ctx context.Context) error {
	prefix := "pkg/assets/builtin_features"
	assets, err := AssetDir(prefix)
	if err != nil {
		return err
	}
	baseDir := config.Config.BasePath
	for _, asset := range assets {
		bts, err := Asset(prefix + "/" + asset)
		if err != nil {
			log.WithField("asset", asset).Error("failed to get asset bytes")
			return err
		}
		ioutil.WriteFile(filepath.Join(baseDir, asset), bts, 0644)
		Features[asset] = string(bts)
	}
	return nil
}
