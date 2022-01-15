package server

import (
	"regexp"
	"strconv"
)

var mailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)

type RegisterRequest struct {
	Name                 string `form:"name"`
	Email                string `form:"email"`
	Password             string `form:"password"`
	PasswordConfirmation string `form:"password_confirmation"`
	Error                string
}

const (
	ErrNameIsEmpty           = "Nama tidak boleh kosong"
	ErrEmailIsEmpty          = "Email tidak boleh kosong"
	ErrPasswordIsEmpty       = "Password tidak boleh kosong"
	ErrPasswordMin6Chars     = "Password minimal 6 karakter"
	ErrPasswordNotSame       = "Password tidak sama"
	ErrInvalidEmail          = "Email tidak valid"
	ErrEmailAlreadyExists    = "Email telah terdaftar"
	ErrInternalServer        = "Sistem sedang bermasalah, tetaplah santuy"
	ErrCodeIsEmpty           = "Kode tidak boleh kosong"
	ErrCodeMustNumber        = "Kode harus berupa angka"
	ErrCodeMust6Chars        = "Kode harus 6 karakter"
	ErrCodeIsNotValid        = "Kode tidak valid atau mungkin sudah kadaluarsa"
	ErrWrongCode             = "Kode salah"
	ErrEmailUnregistered     = "Email belum terdaftar, silahkan mendaftar terlebih dahulu"
	ErrWrongPassword         = "Password salah"
	ErrEmailUnverification   = "Email belum terverifikasi"
	ErrAccountUnVerification = "Akun belum terverifikasi"
)

func (r *RegisterRequest) validate() {
	if r.Name == "" {
		r.Error = ErrNameIsEmpty
		return
	}
	if r.Email == "" {
		r.Error = ErrEmailIsEmpty
		return
	}
	if r.Password == "" || r.PasswordConfirmation == "" {
		r.Error = ErrPasswordIsEmpty
		return
	}
	if len(r.Password) < 6 || len(r.PasswordConfirmation) < 6 {
		r.Error = ErrPasswordMin6Chars
		return
	}
	if r.Password != r.PasswordConfirmation {
		r.Error = ErrPasswordNotSame
		return
	}
	if !mailRegex.MatchString(r.Email) {
		r.Error = ErrInvalidEmail
		return
	}
}

type RegisterCodeVerificationRequest struct {
	AccountID string `form:"account_id"`
	CodeID    string `form:"code_id"`
	Code      string `form:"code"`
	Error     string
}

func (r *RegisterCodeVerificationRequest) validate() {
	if r.Code == "" {
		r.Error = ErrCodeIsEmpty
		return
	}
	if _, err := strconv.Atoi(r.Code); err != nil {
		r.Error = ErrCodeMustNumber
		return
	}
	if len(r.Code) != 6 {
		r.Error = ErrCodeMust6Chars
		return
	}
}

type LoginAccountReq struct {
	Email    string `form:"email"`
	Password string `form:"password"`
	Remember bool   `form:"remember"`
	Error    string
}

func (r *LoginAccountReq) validate() {
	if r.Email == "" {
		r.Error = ErrEmailIsEmpty
	}
	if r.Password == "" {
		r.Error = ErrPasswordIsEmpty
	}
	if len(r.Password) < 6 {
		r.Error = ErrPasswordMin6Chars
	}
	if !mailRegex.MatchString(r.Email) {
		r.Error = ErrInvalidEmail
	}
}

type ResendEmailVerificationReq struct {
	Email string `form:"email"`
	Error string
}

func (r *ResendEmailVerificationReq) validate() {
	if r.Email == "" {
		r.Error = ErrEmailIsEmpty
	}
	if !mailRegex.MatchString(r.Email) {
		r.Error = ErrInvalidEmail
	}
}

type ResetPasswordRequest struct {
	Email        string `form:"email"`
	Password     string `form:"password"`
	PasswordSame string `form:"password_same"`
	Error        string
}

func (r *ResetPasswordRequest) validate() {
	if r.Email == "" {
		r.Error = ErrEmailIsEmpty
	}
	if r.Password == "" || r.PasswordSame == "" {
		r.Error = ErrPasswordIsEmpty
	}
	if len(r.Password) < 6 || r.PasswordSame == "" {
		r.Error = ErrPasswordMin6Chars
	}
	if r.Password != r.PasswordSame {
		r.Error = ErrPasswordNotSame
	}
	if !mailRegex.MatchString(r.Email) {
		r.Error = ErrInvalidEmail
	}
}
