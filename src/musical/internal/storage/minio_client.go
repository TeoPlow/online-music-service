package storage

import (
	"context"
	"io"
	"log/slog"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/TeoPlow/online-music-service/src/musical/internal/config"
	"github.com/TeoPlow/online-music-service/src/musical/internal/logger"
)

type MinIOClient struct {
	client *minio.Client
	bucket string
}

func NewMinIOClient(ctx context.Context, bucket string) (*MinIOClient, error) {
	client, err := minio.New(config.Config.Minio.Endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(config.Config.Minio.AccessKey,
			config.Config.Minio.SecretKey, "",
		),
		Secure: config.Config.Minio.UseSSL,
	})
	if err != nil {
		return nil, err
	}

	exists, err := client.BucketExists(ctx, bucket)
	if err != nil {
		return nil, err
	}
	if !exists {
		if err := client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{}); err != nil {
			return nil, err
		}
	}

	return &MinIOClient{
		client: client,
		bucket: bucket,
	}, nil
}

func (c *MinIOClient) Upload(ctx context.Context, key string, r io.Reader, size int64) error {
	_, err := c.client.PutObject(ctx, c.bucket, key, r, size, minio.PutObjectOptions{})
	return err
}

func (c *MinIOClient) Download(ctx context.Context, key string) (io.ReadCloser, error) {
	return c.client.GetObject(ctx, c.bucket, key, minio.GetObjectOptions{})
}

func (c *MinIOClient) Delete(ctx context.Context, key string) error {
	return c.client.RemoveObject(ctx, c.bucket, key, minio.RemoveObjectOptions{})
}

func (c *MinIOClient) Truncate(ctx context.Context) {
	objectCh := make(chan minio.ObjectInfo)

	go func() {
		defer close(objectCh)
		for object := range c.client.ListObjects(ctx, c.bucket, minio.ListObjectsOptions{Recursive: true}) {
			if object.Err != nil {
				logger.Logger.Error("error listing object", slog.String("error", object.Err.Error()))
				continue
			}
			objectCh <- object
		}
	}()

	for r := range c.client.RemoveObjects(ctx, c.bucket, objectCh, minio.RemoveObjectsOptions{}) {
		if r.Err != nil {
			logger.Logger.Error("error deleting object",
				slog.String("object", r.ObjectName),
				slog.String("error", r.Err.Error()))
		}
	}
}
