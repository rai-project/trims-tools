package assets

import (
	"os"
)

func assetInfo(path string) (os.FileInfo, error) {
	return os.Stat(path)
}
