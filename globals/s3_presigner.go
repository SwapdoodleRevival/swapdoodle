package globals

import (
	"context"
	"net/url"
	"time"

	common_globals "github.com/PretendoNetwork/nex-protocols-common-go/v2/globals"
	"github.com/minio/minio-go/v7"
)

type S3Presigner struct{}

func (S3Presigner) GetObject(bucket, key string, lifetime time.Duration) (*common_globals.S3GetObjectData, error) {
	reqParams := make(url.Values)

	url, err := MinIOClient.PresignedGetObject(context.Background(), bucket, key, lifetime, reqParams)
	if err != nil {
		return nil, err
	}

	stat, err := S3StatObject(bucket, key)
	if err != nil {
		return nil, err
	}

	return &common_globals.S3GetObjectData{
		URL:  url,
		Size: uint32(stat.Size),
	}, nil
}

func (S3Presigner) PostObject(bucket, key string, lifetime time.Duration) (*common_globals.S3PostObjectData, error) {
	policy := minio.NewPostPolicy()

	err := policy.SetBucket(bucket)
	if err != nil {
		return nil, err
	}

	err = policy.SetKey(key)
	if err != nil {
		return nil, err
	}

	err = policy.SetExpires(time.Now().UTC().Add(lifetime).UTC())
	if err != nil {
		return nil, err
	}

	url, formData, err := MinIOClient.PresignedPostPolicy(context.Background(), policy)
	if err != nil {
		return nil, err
	}

	return &common_globals.S3PostObjectData{
		URL:      url,
		FormData: formData,
	}, nil
}

func NewS3Presigner(minioClient *minio.Client) S3Presigner {
	return S3Presigner{}
}
