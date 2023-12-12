package middleware

import (
	"fmt"
	"github.com/U5K492/go-rest/framework/router"
)

func Logger(ctx *router.Context) {
	ctx.Next()
	fmt.Printf("Method: %s\npath: %s\nstatus: %d\n", ctx.R.Method, ctx.R.URL, ctx.Status())
}
