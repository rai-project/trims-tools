package micro

var Features map[string]string

func initFeatures() error {
	assets, err := AssetDir("builtin_features")
	if err != nil {
		return err
	}
	for _, asset := range assets {
		bts, err := Asset(asset)
		if err != nil {
			log.WithField("asset", asset).Error("failed to get asset bytes")
			return err
		}
		Features[asset] = string(bts)
	}
	return nil
}

func init() {
	initFeatures()
}
