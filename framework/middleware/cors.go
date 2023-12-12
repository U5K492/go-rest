package middleware

import "github.com/U5K492/go-rest/framework/router"

func Cors(ctx *router.Context) {
	ctx.Next()

	ctx.W.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
	ctx.W.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	ctx.W.Header().Set("Access-Control-Allow-Credentials", "true")
	ctx.W.Header().Set("Access-Control-Allow-Methods", "GET,PUT,POST,DELETE,UPDATE,OPTIONS")
	ctx.W.Header().Set("Content-Type", "application/json")
}
