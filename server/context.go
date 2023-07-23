package server

import (
	"context"
	"errors"
)

type contextKey string

const (
	contextKeyUser contextKey = "user"
)

// User represents the user stored in the request context.
// authenticator middleware should add this after authentication.
type User struct {
	Name     string
	Username string
}

// userToContext add provided User to context.
func userToContext(ctx context.Context, key contextKey, user User) context.Context {
	ctx = context.WithValue(ctx, key, user)
	return ctx
}

// userFromContext get user from context.
func userFromContext(ctx context.Context, key contextKey) (User, error) {
	val := ctx.Value(key)
	if val == nil {
		return User{}, errors.New("incorrect context")
	}
	user, ok := val.(User)
	if !ok {
		return User{}, errors.New("incorrect type")
	}
	return user, nil
}
