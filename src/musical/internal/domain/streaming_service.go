package domain

import (
	"bytes"
	"context"
	"io"

	"github.com/google/uuid"
)

type FileStorage interface {
	Upload(ctx context.Context, key string, r io.Reader, size int64) error
	Download(ctx context.Context, key string) (io.ReadCloser, error)
	Delete(ctx context.Context, key string) error
}

type StreamingService struct {
	storage FileStorage
}

func NewStreamingService(s FileStorage) *StreamingService {
	return &StreamingService{s}
}

func (s *StreamingService) SaveTrack(ctx context.Context, id uuid.UUID, b *bytes.Buffer) error {
	size := b.Len()
	return s.storage.Upload(ctx, id.String(), b, int64(size))
}

func (s *StreamingService) DownloadTrack(ctx context.Context, userID, trackID uuid.UUID) (io.ReadCloser, error) {
	SendListeningEvent(ctx, userID, trackID)
	return s.storage.Download(ctx, trackID.String())
}

func (s *StreamingService) DeleteTrack(ctx context.Context, id uuid.UUID) error {
	return s.storage.Delete(ctx, id.String())
}
