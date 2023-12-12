package router

import (
	"context"
	"fmt"
	"github.com/U5K492/go-rest/config"
	"github.com/rs/cors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

type Engine struct {
	Router  *Router
	Context *Context
}

func NewEngine() *Engine {
	return &Engine{
		Router: &Router{
			routingTable: map[string]*TreeNode{
				"GET":    Constructor(),
				"POST":   Constructor(),
				"PUT":    Constructor(),
				"DELETE": Constructor(),
			},
		},
	}
}

type Router struct {
	routingTable   map[string]*TreeNode
	middlewareList []func(ctx *Context)
}

func (r *Router) Use(middleware func(ctx *Context)) {
	r.middlewareList = append(r.middlewareList, middleware)
}

func (r *Router) register(method string, pathName string, handler func(ctx *Context)) error {
	routingTable := r.routingTable[method]
	pathName = strings.TrimSuffix(pathName, "/")
	existedHandler := routingTable.walk(pathName)

	if existedHandler != nil {
		panic("already exists")
	}
	routingTable.Insert(pathName, handler)
	return nil
}

func (r *Router) Get(pathName string, handler func(ctx *Context)) error {
	return r.register("GET", pathName, handler)
}

func (r *Router) Post(pathName string, handler func(ctx *Context)) error {
	return r.register("POST", pathName, handler)
}

func (r *Router) Put(pathName string, handler func(ctx *Context)) error {
	return r.register("PUT", pathName, handler)
}

func (r *Router) Delete(pathName string, handler func(ctx *Context)) error {
	return r.register("DELETE", pathName, handler)
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	ctx := NewRouteContext(w, r)

	routingTable := e.Router.routingTable[r.Method]
	pathName := r.URL.Path
	pathName = strings.TrimSuffix(pathName, "/")
	targetNode := routingTable.walk(pathName)

	if targetNode == nil || targetNode.handler == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	paramDict := targetNode.ParseParams(r.URL.Path)
	ctx.SetParams(paramDict)

	handlers := append(e.Router.middlewareList, targetNode.handler)
	ctx.SetHandlers(handlers)
	ctx.Next()
}

func (e *Engine) Run(port uint64) {
	ch := make(chan os.Signal, 1)
	signal.Ignore()
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	c := cors.New(cors.Options{
		AllowedOrigins:   config.CorsConfig().AllowedOrigins,
		AllowedMethods:   config.CorsConfig().AllowedMethods,
		AllowedHeaders:   config.CorsConfig().AllowedHeaders,
		AllowCredentials: config.CorsConfig().AllowCredentials,
	})

	server := &http.Server{Addr: fmt.Sprintf(":%d", port), Handler: c.Handler(e)}

	go func() {
		server.ListenAndServe()
		log.Println("Server Started")
	}()

	<-ch
	fmt.Println("shutdown...")

	if err := server.Shutdown(context.Background()); err != nil {
		fmt.Println("error occurred at shutdown")
	}
	fmt.Println("shutdown completed")
}
