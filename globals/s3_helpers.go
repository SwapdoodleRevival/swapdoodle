package globals

import (
	"context"
	"fmt"
	"strings"

	"github.com/PretendoNetwork/nex-go/v2/types"
	"github.com/minio/minio-go/v7"
)

func S3StatObject(bucket string, key string) (minio.ObjectInfo, error) {
	return MinIOClient.StatObject(context.TODO(), bucket, key, minio.StatObjectOptions{})
}

func S3SetFileContent(bucket, key string, content string) (minio.ObjectInfo, error) {
	MinIOClient.PutObject(context.TODO(), bucket, key, strings.NewReader(content), int64(len(content)), minio.PutObjectOptions{})
	return MinIOClient.StatObject(context.TODO(), bucket, key, minio.StatObjectOptions{})
}

func S3GetNotificationKey(pid types.PID) string {
	return fmt.Sprintf("notifications/%d", pid)
}
