package testutils

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"
)

func ContextWithUserID(ctx context.Context, userID uuid.UUID) context.Context {
	md := metadata.Pairs("user_id", userID.String())
	return metadata.NewOutgoingContext(ctx, md)
}
