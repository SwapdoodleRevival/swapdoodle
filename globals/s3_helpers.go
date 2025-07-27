package globals

import (
	"context"

	"github.com/minio/minio-go/v7"
)

func S3StatObject(bucket string, key string) (minio.ObjectInfo, error) {
	return MinIOClient.StatObject(context.TODO(), bucket, key, minio.StatObjectOptions{})
}
