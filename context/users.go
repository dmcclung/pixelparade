package context

import (
	"context"

	"github.com/dmcclung/pixelparade/models"
)


type ctxKey string

const (
	userKey ctxKey = "user"
)

func WithUser(ctx context.Context, user *models.User) context.Context {
	return context.WithValue(ctx, userKey, user)
}

func User(ctx context.Context) *models.User {
	val := ctx.Value(userKey)
	user, ok := val.(*models.User)
	if !ok {
		return nil
	}

	return user
}