package httping

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func NewHttpServer(port int) *HttpServer {
	engine := gin.Default()
	engine.HandleMethodNotAllowed = true
	server := &http.Server{
		Addr:    ":" + strconv.Itoa(port),
		Handler: engine,
	}
	return &HttpServer{server: server, engine: engine}
}

type HttpServer struct {
	server     *http.Server
	engine     *gin.Engine
	middleware MiddlewareFunc
}

func (server *HttpServer) NewRoute(baseRoute *Route, path string) *Route {
	if baseRoute != nil {
		g := baseRoute.route.Group(path)
		return &Route{route: g}
	}
	g := server.engine.Group(path)
	return &Route{route: g, middleware: server.middleware}
}

func (server *HttpServer) RunServer() (ServerCloseFunc, chan error) {
	chErr := make(chan error)
	go func(server *http.Server) {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			chErr <- err
		}
	}(server.server)
	return func() error {
		return server.server.Close()
	}, chErr
}

func (server *HttpServer) SetMiddleware(middleware MiddlewareFunc) *HttpServer {
	server.middleware = middleware
	return server
}

type ServerCloseFunc func() error
