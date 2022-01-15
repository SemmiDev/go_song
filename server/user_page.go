package server

import "github.com/gofiber/fiber/v2"

func (s *Server) RenderHomePage(c *fiber.Ctx) error {
	return c.Render("home", nil)
}

func (s *Server) RenderRegisterPage(c *fiber.Ctx) error {
	return c.Render("register", nil)
}

func (s *Server) RenderLoginPage(c *fiber.Ctx) error {
	return c.Render("login", nil)
}

func (s *Server) RenderResendEmailVerificationPage(c *fiber.Ctx) error {
	return c.Render("resend-email-verification", nil)
}

func (s *Server) RenderForgotPasswordPage(c *fiber.Ctx) error {
	return c.Render("forgot", nil)
}
