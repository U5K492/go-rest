package handler

import "github.com/U5K492/go-rest/framework/router"

func HealthHandler() func(ctx *router.Context) {
	return func(ctx *router.Context) {
		type Res struct {
			Status string `json:"status"`
		}
		res := &Res{
			Status: "ok",
		}
		ctx.Json(res)
	}
}
