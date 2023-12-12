package handler

import "github.com/U5K492/go-rest/framework/router"

func NameHandler() func(ctx *router.Context) {
	return func(ctx *router.Context) {
		name := ctx.GetParam("name")

		type Res struct {
			Name string `json:"name"`
		}
		res := &Res{
			Name: name,
		}
		ctx.Json(res)
	}
}
