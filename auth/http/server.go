package http

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
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
