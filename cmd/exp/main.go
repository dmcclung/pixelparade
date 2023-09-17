package main

import (
	stdctx "context"
	"fmt"

	"github.com/dmcclung/pixelparade/models"
	"github.com/dmcclung/pixelparade/context"
)

type ctxKey string

const (
	someKey ctxKey = "someKey"
)


func main() {
	ctx := stdctx.Background()

	user := models.User{
		Email: "admin@pixelparade",
	}

	ctx = context.WithUser(ctx, &user)

	val := context.User(ctx)
	fmt.Printf("User email is %v\n", val.Email)
}
