package api

import (
	"github.com/gin-gonic/gin"
	"server/internal/api/handlers"
)

type Server interface {
	Run() error
}

type httpServer struct {
	handler handlers.Handler
	router  *gin.Engine
	port    string
}

func NewHTTPServer(router *gin.Engine, handler handlers.Handler, port string) Server {
	return &httpServer{handler: handler, router: router, port: port}
}

func (s *httpServer) Run() error {
	s.router.POST("/login", s.handler.Login)
	s.router.POST("/register", s.handler.Register)

	posts := s.router.Group("/posts").Use(s.handler.LoginMiddleware)
	{
		posts.POST("/", s.handler.WritePost)
		posts.GET("/", s.handler.GetPosts)
		posts.PATCH("/:id", s.handler.EditPost)
		posts.GET("/:id", s.handler.GetPost)
		posts.DELETE("/:id", s.handler.DeletePost)
	}

	comments := s.router.Group("/comments").Use(s.handler.LoginMiddleware)
	{
		comments.POST("/", s.handler.WriteComment)
		comments.GET("/", s.handler.GetComments)
		comments.DELETE("/:id", s.handler.DeleteComment)
	}

	users := s.router.Group("/users").Use(s.handler.LoginMiddleware)
	{
		users.GET("/:mail", s.handler.GetUser)
	}

	return s.router.Run(s.port)
}
