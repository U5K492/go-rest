package middleware

import (
	"context"
	"fmt"
	"github.com/U5K492/go-rest/framework/router"
	"time"
)

func TimeOut(ctx *router.Context) {
	successCh := make(chan struct{})
	panicCh := make(chan struct{})
	durationContext, cancel := context.WithTimeout(ctx.R.Context(), time.Second*5)

	defer cancel()

	go func() {

		defer func() {
			if err := recover(); err != nil {
				panicCh <- struct{}{}
			}
		}()
		ctx.Next()
		successCh <- struct{}{}
	}()

	select {
	case <-durationContext.Done():
		ctx.WriteString("timeout")
		ctx.SetHasTimeout(true)
	case <-panicCh:
		ctx.WriteString("panic")
	case <-successCh:
		fmt.Println("success")
	}
}
