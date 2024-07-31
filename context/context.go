package context

import (
	"context"

	"github.com/Ali-Gorgani/lenslocked/models"
)

type key string

const (
	userKey key = "user"
)

func WithUser(ctx context.Context, user *models.User) context.Context {
	return context.WithValue(ctx, userKey, user)
}

func User(ctx context.Context) *models.User {
	val := ctx.Value(userKey)
	user, ok := val.(*models.User)
	if !ok {
		// This situation occurs when no user is found or an error occurs during user retrieval(like invalid value in context).
		return nil
	}
	return user
}
