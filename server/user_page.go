package server

import "github.com/gofiber/fiber/v2"

func (s *Server) HomePage(c *fiber.Ctx) error {
	return c.Render("home", nil)
}

func (s *Server) RegisterPage(c *fiber.Ctx) error {
	return c.Render("register", nil)
}

func (s *Server) LoginPage(c *fiber.Ctx) error {
	return c.Render("login", nil)
}

func (s *Server) ResendEmailVerificationPage(c *fiber.Ctx) error {
	return c.Render("resend-email-verification", nil)
}

func (s *Server) ForgotPasswordPage(c *fiber.Ctx) error {
	return c.Render("forgot", nil)
}
