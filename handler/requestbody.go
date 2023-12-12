package handler

import "github.com/U5K492/go-rest/framework/router"

type RequestBody struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func RequestBodyHandler() func(ctx *router.Context) {
	return func(ctx *router.Context) {
		requestBody := &RequestBody{}
		if err := ctx.BindJson(requestBody); err != nil {
			ctx.WriteHeader(500)
			return
		}

		ctx.Json(requestBody)
	}
}
