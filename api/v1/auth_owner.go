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

// @Router /auth-owner/register [post]
// @Summary Register a owner
// @Description Register a owner
// @Tags auth-owner
// @Accept json
// @Produce json
// @Param data body models.RegisterOwnerRequest true "Data"
// @Success 200 {object} models.ResponseOK
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) RegisterOwner(c *gin.Context) {
	var (
		req models.RegisterOwnerRequest
	)

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	_, err = h.storage.Owner().GetByEmail(req.Email)
	if !errors.Is(err, sql.ErrNoRows) {
		c.JSON(http.StatusBadRequest, errorResponse(ErrEmailExists))
		return
	}

	hashedPassword, err := pkg.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	owner := repo.Owner{
		HotelID:   req.HotelID,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Password:  hashedPassword,
	}

	userData, err := json.Marshal(owner)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = h.inMemory.Set("owner_"+owner.Email, string(userData), 10*time.Minute)
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

// @Router /auth-owner/verify [post]
// @Summary Verify owner
// @Description Verify owner
// @Tags auth-owner
// @Accept json
// @Produce json
// @Param data body models.VerifyRequest true "Data"
// @Success 200 {object} models.AuthOwnerResponse
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) VerifyOwner(c *gin.Context) {
	var (
		req models.VerifyRequest
	)

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	userData, err := h.inMemory.Get("owner_" + req.Email)
	if err != nil {
		c.JSON(http.StatusForbidden, errorResponse(err))
		return
	}

	var owner repo.Owner
	err = json.Unmarshal([]byte(userData), &owner)
	if err != nil {
		c.JSON(http.StatusForbidden, errorResponse(err))
		return
	}

	code, err := h.inMemory.Get(RegisterCodeKey + owner.Email)
	if err != nil {
		c.JSON(http.StatusForbidden, errorResponse(ErrCodeExpired))
		return
	}

	if req.Code != code {
		c.JSON(http.StatusForbidden, errorResponse(ErrIncorrectCode))
		return
	}

	result, err := h.storage.Owner().Create(&owner)
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

	c.JSON(http.StatusCreated, models.AuthOwnerResponse{
		ID:          result.ID,
		HotelID:     result.HotelID,
		FirstName:   result.FirstName,
		LastName:    result.LastName,
		Email:       result.Email,
		CreatedAt:   result.CreatedAt,
		AccessToken: token,
	})
}

// @Router /auth-owner/login [post]
// @Summary Login owner
// @Description Login owner
// @Tags auth-owner
// @Accept json
// @Produce json
// @Param data body models.LoginRequest true "Data"
// @Success 200 {object} models.AuthOwnerResponse
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) LoginOwner(c *gin.Context) {
	var (
		req models.LoginRequest
	)

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	result, err := h.storage.Owner().GetByEmail(req.Email)
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

	c.JSON(http.StatusCreated, models.AuthOwnerResponse{
		ID:          result.ID,
		HotelID:     result.HotelID,
		FirstName:   result.FirstName,
		LastName:    result.LastName,
		Email:       result.Email,
		CreatedAt:   result.CreatedAt,
		AccessToken: token,
	})
}

// @Router /auth-owner/forgot-password [post]
// @Summary Forgot password
// @Description Forgot password
// @Tags auth-owner
// @Accept json
// @Produce json
// @Param data body models.ForgotPasswordRequest true "Data"
// @Success 200 {object} models.ResponseOK
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) ForgotPasswordOwner(c *gin.Context) {
	var (
		req models.ForgotPasswordRequest
	)

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	_, err = h.storage.Owner().GetByEmail(req.Email)
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

// @Router /auth-owner/verify-forgot-password [post]
// @Summary Verify forgot password
// @Description Verify forgot password
// @Tags auth-owner
// @Accept json
// @Produce json
// @Param data body models.VerifyRequest true "Data"
// @Success 200 {object} models.AuthOwnerResponse
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) VerifyForgotPasswordOwner(c *gin.Context) {
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

	result, err := h.storage.Owner().GetByEmail(req.Email)
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

	c.JSON(http.StatusCreated, models.AuthOwnerResponse{
		ID:          result.ID,
		HotelID:     result.HotelID,
		FirstName:   result.FirstName,
		LastName:    result.LastName,
		Email:       result.Email,
		CreatedAt:   result.CreatedAt,
		AccessToken: token,
	})
}

// @Security ApiKeyAuth
// @Router /auth-owner/update-password [post]
// @Summary Update password
// @Description Update password
// @Tags auth-owner
// @Accept json
// @Produce json
// @Param data body models.UpdatePasswordRequest true "Data"
// @Success 200 {object} models.ResponseOK
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) UpdatePasswordOwner(c *gin.Context) {
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

	err = h.storage.Owner().UpdateOwnerPassword(&repo.UpdateOwnerPassword{
		OwnerID:  payload.UserID,
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