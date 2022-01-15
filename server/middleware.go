package server

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func (s *Server) CommonMiddleware() {
	s.router.Use(limiter.New(limiter.Config{
		Max:        60,
		Expiration: 1 * time.Minute,
	}))
	s.router.Use(logger.New())
}

func (s *Server) AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		cookie := c.Cookies(s.config.CookieName)
		if cookie == "" {
			return c.Redirect("/login")
		}

		_, err := s.tokenMaker.VerifyToken(cookie)
		if err != nil {
			return c.Redirect("/login")
		}

		return c.Next()
	}
}
