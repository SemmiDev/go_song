package server

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	db "github.com/SemmiDev/go-song/db/datastore"
	"github.com/SemmiDev/go-song/util"
	"github.com/gofiber/fiber/v2"
)

var mailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)

type RegisterAccountRequest struct {
	Name     string `form:"name"`
	Email    string `form:"email"`
	Password string `form:"password"`
	Error    string
}

func (r RegisterAccountRequest) validate() string {
	if r.Name == "" {
		return "Nama tidak boleh kosong"
	}
	if r.Email == "" {
		return "Email tidak boleh kosong"
	}
	if r.Password == "" {
		return "Password tidak boleh kosong"
	}
	if len(r.Password) < 6 {
		return "Password minimal 6 karakter"
	}
	if !mailRegex.MatchString(r.Email) {
		return "Email tidak valid"
	}
	return ""
}

func (s *Server) WebRegisterProcess(c *fiber.Ctx) error {
	var req RegisterAccountRequest
	if err := c.BodyParser(&req); err != nil {
		req.Error = err.Error()
		return c.Render("register", req)
	}

	if req.Error = req.validate(); req.Error != "" {
		return c.Render("register", req)
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		req.Error = err.Error()
		return c.Render("register", req)
	}

	arg := db.CreateAccountParams{
		ID:       util.RandomUUID(1),
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
	}

	account, err := s.datastore.CreateAccount(c.Context(), arg)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			req.Error = "Email telah terdaftar"
			return c.Render("register", req)
		}
		req.Error = "Terjadi masalah di server, tetaplah santuy"
		return c.Render("register", req)
	}

	verificationCode := util.RandomVerificationCode(6)
	verificationID := util.RandomUUID(2)

	log.Println(string(verificationCode))

	s.memorystore.Set(verificationID, verificationCode, time.Minute*30)

	return c.Render("register-code-verification", fiber.Map{
		"CodeID":    verificationID,
		"AccountID": account.ID,
	})
}

type RegisterCodeVerificationRequest struct {
	AccountID string `form:"account_id"`
	CodeID    string `form:"code_id"`
	Code      string `form:"code"`
	Error     string
}

func (r RegisterCodeVerificationRequest) validate() string {
	if r.Code == "" {
		return "Code tidak boleh kosong"
	}
	if _, err := strconv.Atoi(r.Code); err != nil {
		return "Code harus berupa angka"
	}
	if len(r.Code) != 6 {
		return "Code harus 6 karakter"
	}
	return ""
}

func (s *Server) WebRegisterCodeVerificationProcess(c *fiber.Ctx) error {
	var req RegisterCodeVerificationRequest
	if err := c.BodyParser(&req); err != nil {
		req.Error = err.Error()
		return c.Render("register-code-verification", req)
	}

	log.Println(req)
	if req.Error = req.validate(); req.Error != "" {
		return c.Render("register-code-verification", req)
	}

	getCodeID, _ := s.memorystore.Get(req.CodeID)
	if getCodeID == nil {
		req.Error = "Code tidak valid atau mungkin sudah kadaluarsa"
		return c.Render("register-code-verification", req)
	}

	if string(getCodeID) != req.Code {
		return c.Render("register-code-verification", fiber.Map{
			"Error": "Code salah",
		})
	}

	s.memorystore.Delete(req.CodeID)

	arg := db.UpdateAccountEmailVerificationByIDParams{
		ID:              req.AccountID,
		IsEmailVerified: true,
	}

	s.datastore.UpdateAccountEmailVerificationByID(c.Context(), arg)

	return c.Redirect("/login")
}

type LoginAccountReq struct {
	Email    string `form:"email"`
	Password string `form:"password"`
	Remember bool   `form:"remember"`
	Error    string
}

func (r LoginAccountReq) validate() string {
	if r.Email == "" {
		return "Email tidak boleh kosong"
	}
	if r.Password == "" {
		return "Password tidak boleh kosong"
	}
	if len(r.Password) < 6 {
		return "Password minimal 6 karakter"
	}
	if !mailRegex.MatchString(r.Email) {
		return "Email tidak valid"
	}
	return ""
}

func (s *Server) WebLoginProcess(c *fiber.Ctx) error {
	var req LoginAccountReq
	if err := c.BodyParser(&req); err != nil {
		req.Error = err.Error()
		return c.Render("login", req)
	}

	if req.Error = req.validate(); req.Error != "" {
		return c.Render("login", req)
	}

	user, err := s.datastore.GetAccountByEmail(c.Context(), req.Email)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			req.Error = "E-mail belum terdaftar, silahkan mendaftar terlebih dahulu"
			return c.Render("login", req)
		}
		req.Error = "Terjadi masalah di server, tetaplah santuy"
		return c.Render("register", req)
	}

	err = util.CheckPassword(req.Password, user.Password)
	if err != nil {
		req.Error = "Password salah"
		return c.Render("login", req)
	}

	if !user.IsEmailVerified {
		req.Error = `email belum terverifikasi`
		return c.Render("login", req)
	}

	if req.Remember {
		return c.Redirect("/")
	}

	s.SetCookie(c, user, false)
	return c.Redirect("/")
}

type ResendEmailVerificationReq struct {
	Email string `form:"email"`
	Error string
}

func (r ResendEmailVerificationReq) validate() string {
	if r.Email == "" {
		return "Email tidak boleh kosong"
	}
	if !mailRegex.MatchString(r.Email) {
		return "Email tidak valid"
	}
	return ""
}

func (s *Server) WebResendEmailVerification(c *fiber.Ctx) error {
	var req ResendEmailVerificationReq
	if err := c.BodyParser(&req); err != nil {
		req.Error = err.Error()
		return c.Render("resend-email-verification", req)
	}

	if req.Error = req.validate(); req.Error != "" {
		return c.Render("resend-email-verification", req)
	}

	account, err := s.datastore.GetAccountByEmail(c.Context(), req.Email)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			req.Error = "E-mail belum terdaftar, silahkan mendaftar terlebih dahulu"
			return c.Render("resend-email-verification", req)
		}
		req.Error = "Terjadi masalah di server, tetaplah santuy"
		return c.Render("resend-email-verification", req)
	}

	verificationCode := util.RandomVerificationCode(6)
	verificationID := util.RandomUUID(2)
	log.Println(string(verificationCode))

	s.memorystore.Set(verificationID, verificationCode, time.Minute*30)

	return c.Render("register-code-verification", fiber.Map{
		"CodeID":    verificationID,
		"AccountID": account.ID,
	})
}

func (s *Server) WebLogoutProcess(c *fiber.Ctx) error {
	s.DeleteCookie(c)
	return c.Redirect("/login")
}

func (s *Server) WebForgotPasswordProcess(c *fiber.Ctx) error {
	email := c.FormValue("email")
	if email == "" {
		return c.Render("forgot", fiber.Map{"Error": "Email tidak boleh kosong"})
	}
	if !mailRegex.MatchString(email) {
		return c.Render("forgot", fiber.Map{"Error": "Email tidak valid"})
	}

	identifier := util.RandomUUID(2) + "~" + email
	url := fmt.Sprintf("https://42bc-114-125-55-70.ngrok.io/reset-password?id=%s", identifier)

	s.memorystore.Set(identifier, []byte(url), time.Minute*30)

	log.Println(url)

	return c.Render("forgot", fiber.Map{
		"Success": "Sukses, silahkan cek email anda",
	})
}

func (s *Server) ResetPasswordPage(c *fiber.Ctx) error {
	url := c.Query("id")
	data, _ := s.memorystore.Get(url)
	if data != nil {
		email := strings.Split(string(data), "~")[1]
		return c.Render("reset-password", fiber.Map{
			"Email": email,
		})
	}
	return c.Redirect("/forgot")
}

type ResetPasswordRequest struct {
	Email        string `form:"email"`
	Password     string `form:"password"`
	PasswordSame string `form:"password_same"`
	Error        string
}

func (r ResetPasswordRequest) validate() string {
	if r.Email == "" {
		return "Email tidak boleh kosong"
	}
	if r.Password == "" || r.PasswordSame == "" {
		return "Password tidak boleh kosong"
	}
	if len(r.Password) < 6 || r.PasswordSame == "" {
		return "Password minimal 6 karakter"
	}
	if r.Password != r.PasswordSame {
		return "Password tidak sama"
	}
	if !mailRegex.MatchString(r.Email) {
		return "Email tidak valid"
	}
	return ""
}

func (s *Server) WebResetPasswordProcess(c *fiber.Ctx) error {
	var req ResetPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		req.Error = err.Error()
		return c.Render("reset-password", req)
	}

	if req.Error = req.validate(); req.Error != "" {
		return c.Render("reset-password", req)
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		req.Error = err.Error()
		return c.Render("reset-password", req)
	}

	arg := db.UpdateAccountPasswordByEmailParams{
		Email:    req.Email,
		Password: hashedPassword,
	}

	s.datastore.UpdateAccountPasswordByEmail(c.Context(), arg)
	return c.Render("reset-password", fiber.Map{
		"Success": "Email telah berhasil diubah",
	})
}
