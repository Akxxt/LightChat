package utils

import (
	"context"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var minioClient *minio.Client

// InitMinIO 初始化MinIO连接
func InitMinIO() {
	// 从环境变量获取配置
	endpoint := getEnv("MINIO_ENDPOINT", "localhost:9000")
	accessKey := getEnv("MINIO_ACCESS_KEY", "minioadmin")
	secretKey := getEnv("MINIO_SECRET_KEY", "minioadmin")
	useSSL := getEnv("MINIO_USE_SSL", "false") == "true"

	// 创建MinIO客户端
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		panic("MinIO连接失败: " + err.Error())
	}

	minioClient = client

	// 创建默认存储桶
	bucketName := "lightchat"
	exists, err := client.BucketExists(context.Background(), bucketName)
	if err != nil {
		panic("检查存储桶失败: " + err.Error())
	}
	if !exists {
		err = client.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{})
		if err != nil {
			panic("创建存储桶失败: " + err.Error())
		}
	}
}

// GetMinIOClient 获取MinIO客户端
func GetMinIOClient() *minio.Client {
	if minioClient == nil {
		InitMinIO()
	}
	return minioClient
}

// UploadFile 上传文件到MinIO
func UploadFile(bucketName, objectName, filePath string) (string, error) {
	client := GetMinIOClient()

	// 上传文件
	_, err := client.FPutObject(
		context.Background(),
		bucketName,
		objectName,
		filePath,
		minio.PutObjectOptions{},
	)
	if err != nil {
		return "", err
	}

	// 返回文件访问URL
	return "http://" + getEnv("MINIO_ENDPOINT", "localhost:9000") + "/" + bucketName + "/" + objectName, nil
}