package router

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
)

type Context struct {
	W          http.ResponseWriter
	R          *http.Request
	params     map[string]string
	pool       map[string]any
	mut        sync.RWMutex
	hasTimeout bool
	handlers   []func(ctx *Context)
	index      int
	basicWriter
}

func NewRouteContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		W:      w,
		R:      r,
		params: map[string]string{},
		mut:    sync.RWMutex{},
		index:  -1,
	}
}

func (ctx *Context) Get(key string) any {
	ctx.mut.RLock()
	defer ctx.mut.RUnlock()
	if ctx.pool == nil {
		return ""
	}
	res, ok := ctx.pool[key]
	if !ok {
		return ""
	}
	return res
}

func (ctx *Context) Set(key string, value string) {
	ctx.mut.Lock()
	defer ctx.mut.Unlock()
	if ctx.pool == nil {
		ctx.pool = make(map[string]any)
	}

	ctx.pool[key] = value
}

func (ctx *Context) SetHasTimeout(timeout bool) {
	ctx.hasTimeout = timeout
}

func (ctx *Context) HasTimeout() bool {
	return ctx.hasTimeout
}

func (ctx *Context) SetHandlers(handlers []func(ctx *Context)) {
	ctx.handlers = handlers
}

func (ctx *Context) Next() {
	ctx.index++
	for ctx.index < len(ctx.handlers) {
		ctx.handlers[ctx.index](ctx)
		ctx.index++
	}
}

func (ctx *Context) BindJson(data any) error {
	byteData, err := io.ReadAll(ctx.R.Body)
	if err != nil {
		return err
	}

	ctx.R.Body = io.NopCloser(bytes.NewBuffer(byteData))

	return json.Unmarshal(byteData, data)
}

type basicWriter struct {
	http.ResponseWriter
	wroteHeader bool
	code        int
	// bytes       int
	// tee         io.Writer
}

func (bw *basicWriter) WriteHeader(code int) {
	if !bw.wroteHeader {
		bw.code = code
		bw.wroteHeader = true
		bw.WriteHeader(code)
	}
}

func (ctx *Context) WriteString(data string) {
	if ctx.hasTimeout {
		ctx.basicWriter.WriteHeader(http.StatusInternalServerError)
		return
	}
	ctx.basicWriter.WriteHeader(http.StatusOK)
	fmt.Fprint(ctx.W, data)
}

func (ctx *Context) Json(data any) {
	if ctx.hasTimeout {
		ctx.basicWriter.WriteHeader(http.StatusInternalServerError)
		return
	}
	resData, err := json.Marshal(data)
	if err != nil {
		ctx.basicWriter.WriteHeader(http.StatusInternalServerError)
		return
	}
	ctx.W.Header().Set("Content-Type", "application/json")
	ctx.basicWriter.WriteHeader(http.StatusOK)

	ctx.W.Write(resData)
}

func (ctx *Context) Status() int {
	return ctx.basicWriter.code
}

func (ctx *Context) QueryAll() map[string][]string {
	return ctx.R.URL.Query()
}

func (ctx *Context) QueryKey(key string, defaultValue string) string {
	values := ctx.QueryAll()

	if target, ok := values[key]; ok {
		if len(target) == 0 {
			return defaultValue
		}

		return target[len(target)-1]
	}

	return defaultValue
}

func (ctx *Context) SetParams(params map[string]string) {
	ctx.params = params
}

func (ctx *Context) GetParam(key string) string {
	params := ctx.params

	if v, ok := params[":"+key]; ok {
		return v
	}

	return ""
}
