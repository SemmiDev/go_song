package server

import (
	"fmt"
	"github.com/SemmiDev/go-song/common/password"
	"github.com/SemmiDev/go-song/common/random"
	"strings"
	"time"

	db "github.com/SemmiDev/go-song/db/datastore"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (s *Server) RegisterHandler(c *fiber.Ctx) error {
	var req RegisterRequest
	err := c.BodyParser(&req)
	if err != nil {
		return c.Render("register", req)
	}

	req.validate()
	if req.Error != "" {
		return c.Render("register", req)
	}

	hashedPassword, _ := password.Hash(req.Password)
	arg := db.CreateAccountParams{
		ID:       uuid.NewString(),
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
	}

	account, err := s.datastore.CreateAccount(c.Context(), arg)
	if err != nil {
		if strings.Index(err.Error(), "duplicate") >= 0 {
			req.Error = ErrEmailAlreadyExists
			return c.Render("register", req)
		}
		req.Error = ErrInternalServer
		return c.Render("register", req)
	}

	verificationCode := random.VerificationCode(6)
	verificationID := random.UniqueString(2)

	msg := fmt.Sprintf("Kode verifikasi anda adalah: %s \n Kode ini valid selama 30 menit", verificationCode)
	go s.memstore.Set(verificationID, verificationCode, time.Minute*30)
	s.mail.Send(arg.Email, "Verifikasi Akun", msg)

	return c.Render("register-code-verification", fiber.Map{
		"CodeID":    verificationID,
		"AccountID": account.ID,
	})
}

func (s *Server) RegisterCodeVerificationHandler(c *fiber.Ctx) error {
	var req RegisterCodeVerificationRequest
	err := c.BodyParser(&req)
	if err != nil {
		return c.Render("register-code-verification", req)
	}

	req.validate()
	if req.Error != "" {
		return c.Render("register-code-verification", req)
	}

	getCodeID, _ := s.memstore.Get(req.CodeID)
	if getCodeID == nil {
		req.Error = ErrCodeIsNotValid
		return c.Render("register-code-verification", req)
	}

	if string(getCodeID) != req.Code {
		req.Error = ErrWrongCode
		return c.Render("register-code-verification", req)
	}

	go s.memstore.Delete(req.CodeID)
	arg := db.UpdateAccountEmailVerificationByIDParams{
		ID:              req.AccountID,
		IsEmailVerified: true,
	}

	go s.datastore.UpdateAccountEmailVerificationByID(c.Context(), arg)
	return c.Redirect("/login")
}

func (s *Server) LoginHandler(c *fiber.Ctx) error {
	var req LoginAccountReq
	err := c.BodyParser(&req)
	if err != nil {
		return c.Render("login", req)
	}

	req.validate()
	if req.Error != "" {
		return c.Render("login", req)
	}

	user, err := s.datastore.GetAccountByEmail(c.Context(), req.Email)
	if err != nil {
		if strings.Index(err.Error(), "no") >= 0 {
			req.Error = ErrEmailUnregistered
			return c.Render("login", req)
		}
		req.Error = ErrInternalServer
		return c.Render("login", req)
	}

	err = password.Check(req.Password, user.Password)
	if err != nil {
		req.Error = ErrWrongPassword
		return c.Render("login", req)
	}

	if !user.IsEmailVerified {
		req.Error = ErrEmailUnverification
		return c.Render("login", req)
	}

	s.SetCookie(c, user, req.Remember)
	return c.Redirect("/")
}

func (s *Server) ResendEmailVerificationHandler(c *fiber.Ctx) error {
	var req ResendEmailVerificationReq
	err := c.BodyParser(&req)
	if err != nil {
		return c.Render("resend-email-verification", req)
	}

	req.validate()
	if req.Error != "" {
		return c.Render("resend-email-verification", req)
	}

	account, err := s.datastore.GetAccountByEmail(c.Context(), req.Email)
	if err != nil {
		if strings.Index(err.Error(), "no") >= 0 {
			req.Error = ErrEmailUnregistered
			return c.Render("resend-email-verification", req)
		}
		req.Error = ErrInternalServer
		return c.Render("resend-email-verification", req)
	}

	verificationCode := random.VerificationCode(6)
	verificationID := random.UniqueString(2)

	go s.memstore.Set(verificationID, verificationCode, time.Minute*30)
	return c.Render("register-code-verification", fiber.Map{
		"CodeID":    verificationID,
		"AccountID": account.ID,
	})
}

func (s *Server) LogoutHandler(c *fiber.Ctx) error {
	s.DeleteCookie(c)
	return c.Redirect("/login")
}

func (s *Server) ForgotPasswordHandler(c *fiber.Ctx) error {
	email := c.FormValue("email")
	if email == "" {
		return c.Render("forgot", fiber.Map{"Error": ErrEmailIsEmpty})
	}
	if !mailRegex.MatchString(email) {
		return c.Render("forgot", fiber.Map{"Error": ErrInvalidEmail})
	}

	user, err := s.datastore.GetAccountByEmail(c.Context(), email)
	if err != nil {
		if strings.Index(err.Error(), "no") >= 0 {
			return c.Render("forgot", fiber.Map{"Error": ErrEmailUnregistered})
		}
		return c.Render("forgot", fiber.Map{"Error": ErrInternalServer})
	}

	if !user.IsEmailVerified {
		return c.Render("forgot", fiber.Map{"Error": ErrAccountUnVerification})
	}

	identifier := random.UniqueString(2) + "~" + email
	url := fmt.Sprintf("%s/reset-password?id=%s", s.config.Hostname, identifier)

	go s.memstore.Set(identifier, []byte(url), time.Minute*30)

	return c.Render("forgot", fiber.Map{
		"Success": "Sukses, silahkan cek email anda",
	})
}

func (s *Server) ResetPasswordPageHandler(c *fiber.Ctx) error {
	url := c.Query("id")
	data, _ := s.memstore.Get(url)
	if data != nil {
		email := strings.Split(string(data), "~")[1]
		return c.Render("reset-password", fiber.Map{
			"Email": email,
		})
	}
	return c.Redirect("/forgot")
}

func (s *Server) ResetPasswordHandler(c *fiber.Ctx) error {
	var req ResetPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Render("reset-password", req)
	}

	req.validate()
	if req.Error != "" {
		return c.Render("reset-password", req)
	}

	hashedPassword, _ := password.Hash(req.Password)
	arg := db.UpdateAccountPasswordByEmailParams{
		Email:    req.Email,
		Password: hashedPassword,
	}

	go s.datastore.UpdateAccountPasswordByEmail(c.Context(), arg)

	return c.Render("reset-password", fiber.Map{
		"Success": "Password telah berhasil diubah",
	})
}
