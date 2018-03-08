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
	config.Config.Wait()
	baseDir := config.Config.BasePath
	for feature, data := range Features {
		ioutil.WriteFile(filepath.Join(baseDir, feature), []byte(data), 0644)
	}
	return nil
}

func init() {
	prefix := "pkg/assets/builtin_features"
	assets, err := AssetDir(prefix)
	if err != nil {
		return
	}
	for _, asset := range assets {
		bts, err := Asset(prefix + "/" + asset)
		if err != nil {
			log.WithField("asset", asset).Error("failed to get asset bytes")
			continue
		}
		Features[asset] = string(bts)
	}
}
