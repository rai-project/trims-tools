package assets

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/Unknwon/com"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"

	"github.com/rai-project/downloadmanager"
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

type ModelManifest struct {
	Name              string            `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty" yaml:"name,omitempty"`
	Version           string            `protobuf:"bytes,2,opt,name=version,proto3" json:"version,omitempty" yaml:"version,omitempty"`
	Description       string            `protobuf:"bytes,5,opt,name=description,proto3" json:"description,omitempty" yaml:"description,omitempty"`
	Reference         []string          `protobuf:"bytes,6,rep,name=reference" json:"reference,omitempty" yaml:"references,omitempty"`
	License           string            `protobuf:"bytes,7,opt,name=license,proto3" json:"license,omitempty" yaml:"license,omitempty"`
	BeforePreprocess  string            `protobuf:"bytes,10,opt,name=before_preprocess,json=beforePreprocess,proto3" json:"before_preprocess,omitempty" yaml:"before_preprocess,omitempty"`
	Preprocess        string            `protobuf:"bytes,11,opt,name=preprocess,proto3" json:"preprocess,omitempty" yaml:"preprocess,omitempty"`
	AfterPreprocess   string            `protobuf:"bytes,12,opt,name=after_preprocess,json=afterPreprocess,proto3" json:"after_preprocess,omitempty" yaml:"after_preprocess,omitempty"`
	BeforePostprocess string            `protobuf:"bytes,13,opt,name=before_postprocess,json=beforePostprocess,proto3" json:"before_postprocess,omitempty" yaml:"before_postprocess,omitempty"`
	Postprocess       string            `protobuf:"bytes,14,opt,name=postprocess,proto3" json:"postprocess,omitempty" yaml:"postprocess,omitempty"`
	AfterPostprocess  string            `protobuf:"bytes,15,opt,name=after_postprocess,json=afterPostprocess,proto3" json:"after_postprocess,omitempty" yaml:"after_postprocess,omitempty"`
	Model             ModelAssets       `protobuf:"bytes,16,opt,name=model" json:"model,omitempty" yaml:"model,omitempty"`
	Attributes        map[string]string `protobuf:"bytes,17,rep,name=attributes" json:"attributes,omitempty" yaml:"attributes,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Hidden            bool              `protobuf:"varint,18,opt,name=hidden,proto3" json:"hidden,omitempty" yaml:"hidden,omitempty"`
}

type ModelManifests []ModelManifest

var (
	Models ModelManifests
)

func (model ModelManifest) baseURL() string {
	if model.Model.BaseUrl == "" {
		return ""
	}
	return strings.TrimSuffix(model.Model.BaseUrl, "/") + "/"
}

func (model ModelManifest) WorkDir() string {
	baseDir := filepath.Join(Config.BasePath, model.Name)
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

func (model ModelManifest) Download(ctx context.Context) error {
	if model.Model.IsArchive {
		baseURL := model.Model.BaseUrl
		_, err := downloadmanager.DownloadInto(baseURL, model.WorkDir(), downloadmanager.Context(ctx))
		if err != nil {
			return errors.Wrapf(err, "failed to download model archive from %v", model.Model.BaseUrl)
		}
		return nil
	}
	checksum := model.Model.GraphChecksum
	if checksum == "" {
		return errors.New("Need graph file checksum in the model manifest")
	}
	if _, err := downloadmanager.DownloadFile(model.GetGraphUrl(), model.GetGraphPath(), downloadmanager.MD5Sum(checksum)); err != nil {
		return err
	}

	checksum = model.Model.WeightsChecksum
	if checksum == "" {
		return errors.New("Need weights file checksum in the model manifest")
	}
	if _, err := downloadmanager.DownloadFile(model.GetWeightsUrl(), model.GetWeightsPath(), downloadmanager.MD5Sum(checksum)); err != nil {
		return err
	}

	return nil
}

func (ms ModelManifests) Download(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)
	for _, model := range ms {
		g.Go(func() error {
			return model.Download(ctx)
		})
	}
	// Wait for all downloads to complete.
	err := g.Wait()
	if err != nil {
		return err
	}
	log.Info("Successfully downloaded all models.")
	return nil
}

func initModels() error {
	assets, err := AssetDir("builtin_models")
	if err != nil {
		return err
	}
	for _, asset := range assets {
		ext := filepath.Ext(asset)
		if ext != ".yml" && ext != ".yaml" {
			return err
		}

		bts, err := Asset(asset)
		if err != nil {
			log.WithField("asset", asset).Error("failed to get asset bytes")
			return err
		}

		var model ModelManifest
		if err := yaml.Unmarshal(bts, &model); err != nil {
			log.WithField("asset", asset).WithError(err).Error("failed to unmarshal model")
			return err
		}
		if model.Name == "" {
			log.WithField("asset", asset).WithField("name", model.Name).Error("empty model name")
			return errors.New("empty model name found")
		}
		Models = append(Models, model)
	}
	return nil
}

func init() {
	initModels()
}