package config

import "github.com/spf13/viper"

func MinioEndpoint() string {
	return viper.GetString("MINIO_ENDPOINT")
}

func MinioAccessKey() string {
	return viper.GetString("MINIO_ACCESS_KEY")
}

func MinioSecretKey() string {
	return viper.GetString("MINIO_SECRET_KEY")
}

func MinioBucketName() string {
	return viper.GetString("MINIO_BUCKET_NAME")
}

func MinioUseSSL() bool {
	return viper.GetBool("MINIO_USE_SSL")
}
