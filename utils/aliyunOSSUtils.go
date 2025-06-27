package utils

import (
	"volunteer-system-backend/config"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/google/uuid"
	"mime/multipart"
	"path/filepath"
)

func AliOSSUtils(file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	defer file.Close()
	AliyunOSSConfig := config.ProjectConfig.AliyunOSS
	//accessKeySecret := AliyunOSSConfig.AccessKeySecret
	//objectKey := AliyunOSSConfig.ObjectKey
	//bucketName := AliyunOSSConfig.BucketName
	//endpoint := AliyunOSSConfig.Endpoint
	//area := AliyunOSSConfig.Area

	client, err := oss.New(AliyunOSSConfig.Endpoint, AliyunOSSConfig.AccessKeyId, AliyunOSSConfig.AccessKeySecret)
	if err != nil {
		return "", err
	}

	bucket, err := client.Bucket(AliyunOSSConfig.BucketName)
	if err != nil {
		return "", err
	}

	// 上传文件内容
	ext := filepath.Ext(fileHeader.Filename)
	newObjectKey := fmt.Sprintf("%s/image_%s%s", AliyunOSSConfig.ObjectKey, uuid.New().String(), ext)
	err = bucket.PutObject(newObjectKey, file)
	if err != nil {
		return "", err
	}

	// 返回图片 URL（前提是文件可公开访问）
	imageURL := fmt.Sprintf("https://%s.oss-cn-%s.aliyuncs.com/%s", AliyunOSSConfig.BucketName, AliyunOSSConfig.Area, newObjectKey)
	//imageURL := fmt.Sprintf("https://avatar-minshenyao.oss-cn-beijing.aliyuncs.com/%s", newObjectKey)
	return imageURL, nil
}
