package server

import (
	"time"

	db "github.com/SemmiDev/go-song/db/datastore"
	"github.com/gofiber/fiber/v2"
)

func (s *Server) SetCookie(c *fiber.Ctx, user db.Account, isRememberMe bool) {
	duration := time.Hour * 24 // 1 day
	token, _ := s.tokenMaker.CreateToken(user.Email, duration)
	cookie := &fiber.Cookie{
		Name:     s.config.CookieName,
		Value:    token,
		Path:     "/",
		Expires:  time.Now().Add(duration),
		Secure:   false,
		HTTPOnly: true,
	}

	if isRememberMe {
		duration = time.Hour * 24 * 7 // 7 days
		token, _ = s.tokenMaker.CreateToken(user.Email, duration)
		cookie = &fiber.Cookie{
			Name:     s.config.CookieName,
			Value:    token,
			Path:     "/",
			Expires:  time.Now().Add(duration),
			Secure:   false,
			HTTPOnly: true,
		}
		c.Cookie(cookie)
		return
	}
	c.Cookie(cookie)
}

func (s *Server) DeleteCookie(c *fiber.Ctx) {
	cookie := &fiber.Cookie{
		Name:     s.config.CookieName,
		Value:    "",
		Path:     "/",
		Expires:  time.Now().Add(-time.Hour),
		Secure:   false,
		HTTPOnly: true,
	}
	c.Cookie(cookie)
}
