package middleware

import (
	"fmt"
	"github.com/U5K492/go-rest/framework/router"
	"time"
)

func TimeCost(ctx *router.Context) {
	now := time.Now()
	ctx.Next()
	fmt.Println("time cost:", time.Since(now).Milliseconds())
}
