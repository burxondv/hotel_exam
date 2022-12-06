package v1

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/burxondv/hotel_exam/api/models"
	"github.com/burxondv/hotel_exam/pkg"
	"github.com/burxondv/hotel_exam/storage/repo"
	"github.com/gin-gonic/gin"
)

var (
	ErrWrongEmailOrPass = errors.New("wrong email or password")
	ErrEmailExists      = errors.New("email already exists")
	ErrUserNotVerified  = errors.New("user not verified")
	ErrIncorrectCode    = errors.New("incorrect verification code")
	ErrCodeExpired      = errors.New("verification code has been expired")
)

const (
	RegisterCodeKey   = "register_code_"
	ForgotPasswordKey = "forgot_password_code_"
)

// @Router /auth-guest/register [post]
// @Summary Register a guest
// @Description Register a guest
// @Tags auth-guest
// @Accept json
// @Produce json
// @Param data body models.RegisterGuestRequest true "Data"
// @Success 200 {object} models.ResponseOK
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) RegisterGuest(c *gin.Context) {
	var (
		req models.RegisterGuestRequest
	)

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	_, err = h.storage.Guest().GetByEmail(req.Email)
	if !errors.Is(err, sql.ErrNoRows) {
		c.JSON(http.StatusBadRequest, errorResponse(ErrEmailExists))
		return
	}

	hashedPassword, err := pkg.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	guest := repo.Guest{
		HotelID:   req.HotelID,
		RoomID:    req.RoomID,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Password:  hashedPassword,
	}

	userData, err := json.Marshal(guest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = h.inMemory.Set("guest_"+guest.Email, string(userData), 10*time.Minute)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	go func() {
		err := h.sendVerificationCode(RegisterCodeKey, req.Email)
		if err != nil {
			fmt.Printf("failed to send verification code: %v", err)
		}
	}()

	c.JSON(http.StatusCreated, models.ResponseOK{
		Message: "Verification code has been sent!",
	})
}

// @Router /auth-guest/verify [post]
// @Summary Verify guest
// @Description Verify guest
// @Tags auth-guest
// @Accept json
// @Produce json
// @Param data body models.VerifyRequest true "Data"
// @Success 200 {object} models.AuthGuestResponse
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) VerifyGuest(c *gin.Context) {
	var (
		req models.VerifyRequest
	)

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	userData, err := h.inMemory.Get("guest_" + req.Email)
	if err != nil {
		c.JSON(http.StatusForbidden, errorResponse(err))
		return
	}

	var guest repo.Guest
	err = json.Unmarshal([]byte(userData), &guest)
	if err != nil {
		c.JSON(http.StatusForbidden, errorResponse(err))
		return
	}

	code, err := h.inMemory.Get(RegisterCodeKey + guest.Email)
	if err != nil {
		c.JSON(http.StatusForbidden, errorResponse(ErrCodeExpired))
		return
	}

	if req.Code != code {
		c.JSON(http.StatusForbidden, errorResponse(ErrIncorrectCode))
		return
	}

	result, err := h.storage.Guest().Create(&guest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	token, _, err := pkg.CreateToken(h.cfg, &pkg.TokenParams{
		UserID:   result.ID,
		Email:    result.Email,
		Duration: time.Hour * 24,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusCreated, models.AuthGuestResponse{
		ID:          result.ID,
		HotelID:     result.HotelID,
		RoomID:      result.RoomID,
		FirstName:   result.FirstName,
		LastName:    result.LastName,
		Email:       result.Email,
		CreatedAt:   result.CreatedAt,
		AccessToken: token,
	})
}

// @Router /auth-guest/login [post]
// @Summary Login guest
// @Description Login guest
// @Tags auth-guest
// @Accept json
// @Produce json
// @Param data body models.LoginRequest true "Data"
// @Success 200 {object} models.AuthGuestResponse
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) LoginGuest(c *gin.Context) {
	var (
		req models.LoginRequest
	)

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	result, err := h.storage.Guest().GetByEmail(req.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusForbidden, errorResponse(ErrWrongEmailOrPass))
			return
		}

		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = pkg.CheckPassword(req.Password, result.Password)
	if err != nil {
		c.JSON(http.StatusForbidden, errorResponse(ErrWrongEmailOrPass))
		return
	}

	token, _, err := pkg.CreateToken(h.cfg, &pkg.TokenParams{
		UserID:   result.ID,
		Email:    result.Email,
		Duration: time.Hour * 24,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusCreated, models.AuthGuestResponse{
		ID:          result.ID,
		HotelID:     result.HotelID,
		RoomID:      result.RoomID,
		FirstName:   result.FirstName,
		LastName:    result.LastName,
		Email:       result.Email,
		CreatedAt:   result.CreatedAt,
		AccessToken: token,
	})
}

// @Router /auth-guest/forgot-password [post]
// @Summary Forgot password
// @Description Forgot password
// @Tags auth-guest
// @Accept json
// @Produce json
// @Param data body models.ForgotPasswordRequest true "Data"
// @Success 200 {object} models.ResponseOK
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) ForgotPasswordGuest(c *gin.Context) {
	var (
		req models.ForgotPasswordRequest
	)

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	_, err = h.storage.Guest().GetByEmail(req.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	go func() {
		err := h.sendVerificationCode(ForgotPasswordKey, req.Email)
		if err != nil {
			fmt.Printf("failed to send verification code: %v", err)
		}
	}()

	c.JSON(http.StatusCreated, models.ResponseOK{
		Message: "Verification code has been sent!",
	})
}

// @Router /auth-guest/verify-forgot-password [post]
// @Summary Verify forgot password
// @Description Verify forgot password
// @Tags auth-guest
// @Accept json
// @Produce json
// @Param data body models.VerifyRequest true "Data"
// @Success 200 {object} models.AuthGuestResponse
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) VerifyForgotPasswordGuest(c *gin.Context) {
	var (
		req models.VerifyRequest
	)

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	code, err := h.inMemory.Get(ForgotPasswordKey + req.Email)
	if err != nil {
		c.JSON(http.StatusForbidden, errorResponse(ErrCodeExpired))
		return
	}

	if req.Code != code {
		c.JSON(http.StatusForbidden, errorResponse(ErrIncorrectCode))
		return
	}

	result, err := h.storage.Guest().GetByEmail(req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	token, _, err := pkg.CreateToken(h.cfg, &pkg.TokenParams{
		UserID:   result.ID,
		Email:    result.Email,
		Duration: time.Minute * 30,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusCreated, models.AuthGuestResponse{
		ID:          result.ID,
		HotelID:     result.HotelID,
		RoomID:      result.RoomID,
		FirstName:   result.FirstName,
		LastName:    result.LastName,
		Email:       result.Email,
		CreatedAt:   result.CreatedAt,
		AccessToken: token,
	})
}

// @Security ApiKeyAuth
// @Router /auth-guest/update-password [post]
// @Summary Update password
// @Description Update password
// @Tags auth-guest
// @Accept json
// @Produce json
// @Param data body models.UpdatePasswordRequest true "Data"
// @Success 200 {object} models.ResponseOK
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) UpdatePasswordGuest(c *gin.Context) {
	var (
		req models.UpdatePasswordRequest
	)

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payload, err := h.GetAuthPayload(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	hashedPassword, err := pkg.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = h.storage.Guest().UpdateGuestPassword(&repo.UpdateGuestPassword{
		GuestID:  payload.UserID,
		Password: hashedPassword,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusCreated, models.ResponseOK{
		Message: "Password has been updated!",
	})
}

func (h *handlerV1) sendVerificationCode(key, email string) error {
	code, err := pkg.GenerateRandomCode(6)
	if err != nil {
		return err
	}

	err = h.inMemory.Set(key+email, code, time.Minute)
	if err != nil {
		return err
	}

	err = pkg.SendEmail(h.cfg, &pkg.SendEmailRequest{
		To:      []string{email},
		Subject: "Verification email",
		Body: map[string]string{
			"code": code,
		},
		Type: pkg.VerificationEmail,
	})
	if err != nil {
		return err
	}

	return nil
}
