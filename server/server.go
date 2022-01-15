package server

import (
	"fmt"
	"github.com/SemmiDev/go-song/common/config"
	"github.com/SemmiDev/go-song/common/mail"
	"github.com/SemmiDev/go-song/common/token"
	"github.com/SemmiDev/go-song/db/memorystore"
	"github.com/gofiber/template/html"

	db "github.com/SemmiDev/go-song/db/datastore"
	"github.com/gofiber/fiber/v2"
)

type Server struct {
	config     config.Config
	tokenMaker token.Maker
	mail       mail.Mailer
	datastore  db.Store
	memstore   *memorystore.Storage
	router     *fiber.App
}

func New(config config.Config, datastore db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:     config,
		datastore:  datastore,
		tokenMaker: tokenMaker,
		mail:       mail.NewSMTP(),
		memstore:   memorystore.New(),
	}

	templateEngine := html.New("./web/templates", ".html")
	fiberConfig := fiber.Config{Views: templateEngine}
	server.router = fiber.New(fiberConfig)

	server.router.Static("/", "./web/assets")
	server.CommonMiddleware()

	// render pages
	server.router.Get("/", server.AuthMiddleware(), server.RenderHomePage)
	server.router.Get("/register", server.RenderRegisterPage)
	server.router.Get("/login", server.RenderLoginPage)
	server.router.Get("/resend-email-verification", server.RenderResendEmailVerificationPage)
	server.router.Get("/forgot", server.RenderForgotPasswordPage)
	server.router.Get("/reset-password", server.ResetPasswordPageHandler)

	// process handler
	server.router.Post("/register", server.RegisterHandler)
	server.router.Post("/register-code-verification", server.RegisterCodeVerificationHandler)
	server.router.Post("/login", server.LoginHandler)
	server.router.Post("/resend-email-verification", server.ResendEmailVerificationHandler)
	server.router.Post("/forgot", server.ForgotPasswordHandler)
	server.router.Post("/reset-password", server.ResetPasswordHandler)
	server.router.Get("/logout", server.AuthMiddleware(), server.LogoutHandler)

	return server, nil
}

func (s *Server) Start(addr string) error {
	return s.router.Listen(addr)
}
