package trace

import (
	"bytes"
	"encoding/json"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/rai-project/aws"
	"github.com/rai-project/micro18-tools/pkg/config"
	"github.com/rai-project/store"
	"github.com/rai-project/store/s3"
)

func (tr Trace) Upload() error {
	js, err := json.Marshal(tr)
	if err != nil {
		return errors.Wrap(err, "failed to marshal trace as json")
	}

	session, err := aws.NewSession(
		aws.Region(aws.Config.Region),
		aws.AccessKey(aws.Config.AccessKey),
		aws.SecretKey(aws.Config.SecretKey),
	)
	if err != nil {
		return err
	}

	st, err := s3.New(
		s3.Session(session),
		store.Bucket(config.Config.UploadBucketName),
		store.BaseURL(config.Config.BaseBucketURL),
	)
	if err != nil {
		return err
	}

	uploadKey := config.Config.UploadBucketName + "/" + tr.ID + ".json"

	key, err := st.UploadFrom(
		bytes.NewBuffer(js),
		uploadKey,
		s3.Metadata(map[string]interface{}{
			"id":         tr.ID,
			"type":       "profile_upload",
			"created_at": time.Now(),
		}),
		s3.ContentType("application/json"),
		store.UploadProgressOutput(os.Stdout),
	)
	if err != nil {
		return err
	}

	log.WithField("key", key).Info("profile uploaded")

	return nil
}
