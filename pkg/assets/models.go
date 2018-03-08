package assets

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/Unknwon/com"
	"github.com/pkg/errors"

	"github.com/rai-project/downloadmanager"
	"github.com/rai-project/micro18-tools/pkg/config"
	yaml "gopkg.in/yaml.v2"
)

type ModelAssets struct {
	BaseUrl         string `protobuf:"bytes,1,opt,name=base_url,json=baseUrl,proto3" json:"base_url,omitempty" yaml:"base_url,omitempty"`
	WeightsPath     string `protobuf:"bytes,2,opt,name=weights_path,json=weightsPath,proto3" json:"weights_path,omitempty" yaml:"weights_path,omitempty"`
	GraphPath       string `protobuf:"bytes,3,opt,name=graph_path,json=graphPath,proto3" json:"graph_path,omitempty" yaml:"graph_path,omitempty"`
	IsArchive       bool   `protobuf:"varint,4,opt,name=is_archive,json=isArchive,proto3" json:"is_archive,omitempty" yaml:"is_archive,omitempty"`
	WeightsChecksum string `protobuf:"bytes,5,opt,name=weights_checksum,json=weightsChecksum,proto3" json:"weights_checksum,omitempty" yaml:"weights_checksum,omitempty"`
	GraphChecksum   string `protobuf:"bytes,6,opt,name=graph_checksum,json=graphChecksum,proto3" json:"graph_checksum,omitempty" yaml:"graph_checksum,omitempty"`
}

type ModelManifest_Type struct {
	Type        string                 `protobuf:"bytes,1,opt,name=type,proto3" json:"type,omitempty" yaml:"type,omitempty"`
	Description string                 `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty" yaml:"description,omitempty"`
	Parameters  map[string]interface{} `protobuf:"bytes,3,rep,name=parameters" json:"parameters,omitempty" yaml:"parameters,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value"`
}

type ModelManifest struct {
	Name              string                `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty" yaml:"name,omitempty"`
	Version           string                `protobuf:"bytes,2,opt,name=version,proto3" json:"version,omitempty" yaml:"version,omitempty"`
	Description       string                `protobuf:"bytes,5,opt,name=description,proto3" json:"description,omitempty" yaml:"description,omitempty"`
	Reference         []string              `protobuf:"bytes,6,rep,name=reference" json:"reference,omitempty" yaml:"references,omitempty"`
	License           string                `protobuf:"bytes,7,opt,name=license,proto3" json:"license,omitempty" yaml:"license,omitempty"`
	Inputs            []*ModelManifest_Type `protobuf:"bytes,8,rep,name=inputs" json:"inputs,omitempty" yaml:"inputs,omitempty"`
	Output            *ModelManifest_Type   `protobuf:"bytes,9,opt,name=output" json:"output,omitempty" yaml:"output,omitempty"`
	BeforePreprocess  string                `protobuf:"bytes,10,opt,name=before_preprocess,json=beforePreprocess,proto3" json:"before_preprocess,omitempty" yaml:"before_preprocess,omitempty"`
	Preprocess        string                `protobuf:"bytes,11,opt,name=preprocess,proto3" json:"preprocess,omitempty" yaml:"preprocess,omitempty"`
	AfterPreprocess   string                `protobuf:"bytes,12,opt,name=after_preprocess,json=afterPreprocess,proto3" json:"after_preprocess,omitempty" yaml:"after_preprocess,omitempty"`
	BeforePostprocess string                `protobuf:"bytes,13,opt,name=before_postprocess,json=beforePostprocess,proto3" json:"before_postprocess,omitempty" yaml:"before_postprocess,omitempty"`
	Postprocess       string                `protobuf:"bytes,14,opt,name=postprocess,proto3" json:"postprocess,omitempty" yaml:"postprocess,omitempty"`
	AfterPostprocess  string                `protobuf:"bytes,15,opt,name=after_postprocess,json=afterPostprocess,proto3" json:"after_postprocess,omitempty" yaml:"after_postprocess,omitempty"`
	Model             ModelAssets           `protobuf:"bytes,16,opt,name=model" json:"model,omitempty" yaml:"model,omitempty"`
	Attributes        map[string]string     `protobuf:"bytes,17,rep,name=attributes" json:"attributes,omitempty" yaml:"attributes,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Hidden            bool                  `protobuf:"varint,18,opt,name=hidden,proto3" json:"hidden,omitempty" yaml:"hidden,omitempty"`
}

type ModelManifests []ModelManifest

var (
	Models = ModelManifests{}
)

func (model ModelManifest) baseURL() string {
	if model.Model.BaseUrl == "" {
		return ""
	}
	return strings.TrimSuffix(model.Model.BaseUrl, "/") + "/"
}

func (model ModelManifest) CanonicalName() (string, error) {
	modelName := strings.ToLower(model.Name)
	if modelName == "" {
		return "", errors.New("model name must not be empty")
	}
	modelVersion := model.Version
	if modelVersion == "" {
		modelVersion = "latest"
	}
	cannonicalName := modelName + ":" + modelVersion
	cannonicalName = strings.Replace(cannonicalName, ":", "_", -1)
	cannonicalName = strings.Replace(cannonicalName, " ", "_", -1)
	cannonicalName = strings.Replace(cannonicalName, "-", "_", -1)
	return cannonicalName, nil
}

func (model ModelManifest) MustCanonicalName() string {
	cannonicalName, err := model.CanonicalName()
	if err != nil {
		panic(err)
	}
	return cannonicalName
}

func (model ModelManifest) WorkDir() string {
	cannonicalName, err := model.CanonicalName()
	if err != nil {
		return ""
	}

	baseDir := filepath.Join(config.Config.BasePath, cannonicalName)
	if !com.IsDir(baseDir) {
		err := os.MkdirAll(baseDir, os.ModePerm)
		if err != nil {
			log.WithError(err).WithField("path", baseDir).WithField("model", model.Name).Panic("failed to create directory for model's graph path")
		}
	}
	return baseDir
}

func (model ModelManifest) GetGraphPath() string {
	baseDir := model.WorkDir()
	return filepath.Join(baseDir, "model.symbol")
}

func (model ModelManifest) GetWeightsPath() string {
	baseDir := model.WorkDir()
	return filepath.Join(baseDir, "model.params")
}

func (model ModelManifest) GetFeaturesPath() string {
	baseDir := model.WorkDir()
	return filepath.Join(baseDir, "features.txt")
}

func (model ModelManifest) GetWeightsUrl() string {
	if model.Model.IsArchive {
		return model.Model.BaseUrl
	}
	return model.baseURL() + model.Model.WeightsPath
}

func (model ModelManifest) GetGraphUrl() string {
	if model.Model.IsArchive {
		return model.Model.BaseUrl
	}
	return model.baseURL() + model.Model.GraphPath
}

func (model ModelManifest) GetFeaturesUrl() string {
	if model.Output == nil {
		return ""
	}
	params := model.Output.Parameters
	pfeats, ok := params["features_url"]
	if !ok {
		return ""
	}
	return pfeats.(string)
}

func (model ModelManifest) GetImageDimensions() ([]uint32, error) {
	modelInputs := model.Inputs
	typeParameters := modelInputs[0].Parameters
	if typeParameters == nil {
		return nil, errors.New("invalid type parameters")
	}
	pdims0, ok := typeParameters["dimensions"]
	if !ok {
		return nil, errors.New("expecting image type dimensions")
	}
	pdims, ok := pdims0.(string)
	if !ok {
		return nil, errors.New("expecting image type string")
	}

	var dims []uint32
	if err := yaml.Unmarshal([]byte(pdims), &dims); err != nil {
		return nil, errors.Errorf("unable to get image dimensions %v as an integer slice", pdims)
	}
	if len(dims) != 3 {
		return nil, errors.Errorf("expecting a dimensions size of 3, but got %v. do not put the batch size in the input dimensions.", len(dims))
	}
	return dims, nil
}

func (model ModelManifest) GetMeanImage() ([]float32, error) {
	modelInputs := model.Inputs
	typeParameters := modelInputs[0].Parameters
	if typeParameters == nil {
		return nil, errors.New("invalid type parameters")
	}
	pmean, ok := typeParameters["mean"]
	if !ok {
		log.Debug("using 0,0,0 as the mean image")
		return []float32{0, 0, 0}, nil
	}
	pmeanVal, ok := pmean.(string)
	if !ok {
		return nil, errors.New("expecting parameters type string")
	}

	var vals []float32
	if err := yaml.Unmarshal([]byte(pmeanVal), &vals); err == nil {
		return vals, nil
	}
	var val float32
	if err := yaml.Unmarshal([]byte(pmeanVal), &val); err != nil {
		return nil, errors.Errorf("unable to get image mean %v as a float or slice", pmeanVal)
	}

	return []float32{val, val, val}, nil
}

func (model ModelManifest) Download(ctx context.Context) (err error) {
	defer func() {
		if err != nil {
			log.WithError(err).WithField("name", model.Name).Error("failed to download model")
			return
		}
		log.WithField("name", model.Name).Info("model downloaded")
	}()
	if model.Model.IsArchive {
		baseURL := model.Model.BaseUrl
		_, err = downloadmanager.DownloadInto(baseURL, model.WorkDir(), downloadmanager.Context(ctx))
		if err != nil {
			err = errors.Wrapf(err, "failed to download model archive from %v", model.Model.BaseUrl)
			return
		}
		return
	}
	checksum := model.Model.GraphChecksum
	if checksum == "" {
		err = errors.New("Need graph file checksum in the model manifest")
		return err
	}
	if _, err = downloadmanager.DownloadFile(model.GetGraphUrl(), model.GetGraphPath(), downloadmanager.MD5Sum(checksum)); err != nil {
		return
	}

	checksum = model.Model.WeightsChecksum
	if checksum == "" {
		return errors.New("Need weights file checksum in the model manifest")
	}
	if _, err = downloadmanager.DownloadFile(model.GetWeightsUrl(), model.GetWeightsPath(), downloadmanager.MD5Sum(checksum)); err != nil {
		return
	}
	if featuresURL := model.GetFeaturesUrl(); featuresURL != "" {
		if _, err = downloadmanager.DownloadFile(featuresURL, model.GetFeaturesPath()); err != nil {
			return
		}
	}

	return
}

func (ms ModelManifests) Download(ctx context.Context) error {
	for _, model := range ms {
		log.WithField("name", model.Name).Info("downloading model")
		if err := model.Download(ctx); err != nil {
			log.WithError(err).WithField("name", model.Name).Info("failed to download model")
		}
	}
	log.Info("Successfully downloaded all models.")
	return nil
}

func DownloadModels(ctx context.Context) error {
	return Models.Download(ctx)
}

func init() {
	prefix := "pkg/assets/builtin_models"
	assets, err := AssetDir(prefix)
	if err != nil {
		return
	}
	for _, asset := range assets {
		ext := filepath.Ext(asset)
		if ext != ".yml" && ext != ".yaml" {
			continue
		}

		bts, err := Asset(prefix + "/" + asset)
		if err != nil {
			log.WithField("asset", asset).Error("failed to get asset bytes")
			continue
		}

		var model ModelManifest
		if err := yaml.Unmarshal(bts, &model); err != nil {
			log.WithField("asset", asset).WithError(err).Error("failed to unmarshal model")
			continue
		}
		if model.Name == "" {
			log.WithField("asset", asset).WithField("name", model.Name).Error("empty model name")
			continue
		}
		Models = append(Models, model)
	}
}
