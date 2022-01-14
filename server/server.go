package server

import (
	"fmt"

	"github.com/SemmiDev/go-song/db/memorystore"
	"github.com/SemmiDev/go-song/mail"
	"github.com/gofiber/template/html"

	db "github.com/SemmiDev/go-song/db/datastore"
	"github.com/SemmiDev/go-song/token"
	"github.com/SemmiDev/go-song/util"
	"github.com/gofiber/fiber/v2"
)

type Server struct {
	env         util.Env
	tokenMaker  token.Maker
	mailer      mail.Mailer
	datastore   db.Store
	memorystore *memorystore.Storage
	router      *fiber.App
}

func New(env util.Env, datastore db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(env.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	s := &Server{
		env:         env,
		datastore:   datastore,
		tokenMaker:  tokenMaker,
		mailer:      mail.NewSTMP(),
		memorystore: memorystore.New(),
	}

	templateEngine := html.New("./web/templates", ".html")
	fiberConfig := fiber.Config{Views: templateEngine}
	s.router = fiber.New(fiberConfig)

	s.CommonMiddleware()

	s.router.Static("/", "./web/assets")

	s.router.Get("/", s.AuthMiddleware(), s.HomePage)
	s.router.Get("/register", s.RegisterPage)
	s.router.Get("/login", s.LoginPage)
	s.router.Get("/resend-email-verification", s.ResendEmailVerificationPage)
	s.router.Get("/forgot", s.ForgotPasswordPage)
	s.router.Get("/logout", s.AuthMiddleware(), s.WebLogoutProcess)
	s.router.Get("/reset-password", s.ResetPasswordPage)

	s.router.Post("/register", s.WebRegisterProcess)
	s.router.Post("/register-code-verification", s.WebRegisterCodeVerificationProcess)
	s.router.Post("/login", s.WebLoginProcess)
	s.router.Post("/resend-email-verification", s.WebResendEmailVerification)
	s.router.Post("/forgot", s.WebForgotPasswordProcess)
	s.router.Post("/reset-password", s.WebResetPasswordProcess)

	return s, nil
}

func (s *Server) Start(addr string) error {
	return s.router.Listen(addr)
}
