package http

import (
	"log"

	"github.com/gistsapp/api/auth/core"
	"github.com/gofiber/fiber/v2"
)

type Server struct {
	listen_addr string
	app         *fiber.App
}

func NewServer(listen_addr string) *Server {
	return &Server{
		listen_addr: listen_addr,
		app:         fiber.New(),
	}
}

func (s *Server) Setup(handlers ...Handler) {
	s.app.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowOrigins:     "*",
	}))

	s.app.Use(logger.New())

	for _, handler := range handlers {
		handler.Register(s.app)
	}
}

func (s *Server) Ignite() {
	log.Fatal(s.app.Listen(s.listen_addr))
}

func NewServer(jwtService core.JWTService) *fiber.App {
    app := fiber.New()

    handler := NewHandler(jwtService)
    handler.Register(app)

    return app
}