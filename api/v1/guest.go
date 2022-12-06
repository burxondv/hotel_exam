package v1

import (
	"net/http"
	"strconv"

	"github.com/burxondv/hotel_exam/api/models"
	"github.com/burxondv/hotel_exam/storage/repo"
	"github.com/gin-gonic/gin"
)

// @Security ApiKeyAuth
// @Router /guest [post]
// @Summary Create a guest
// @Description Create a guest
// @Tags guest
// @Accept json
// @Produce json
// @Param guest body models.CreateGuestRequest true "Guest"
// @Success 201 {object} models.Guest
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) CreateGuest(c *gin.Context) {
	var (
		req models.CreateGuestRequest
	)

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	resp, err := h.storage.Guest().Create(&repo.Guest{
		HotelID:     req.HotelID,
		RoomID:      req.RoomID,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		PhoneNumber: req.PhoneNumber,
		Email:       req.Email,
		Password:    req.Password,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusCreated, parseGuestModel(resp))
}

// @Router /guest/{id} [get]
// @Summary Get guest by id
// @Description Get guest by id
// @Tags guest
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} models.Guest
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) GetGuest(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	resp, err := h.storage.Guest().Get(int64(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, parseGuestModel(resp))
}

// @Router /guest [get]
// @Summary Get all guest
// @Description Get all guest
// @Tags guest
// @Accept json
// @Produce json
// @Param filter query models.GetAllParams false "Filter"
// @Success 200 {object} models.GetAllGuestsResponse
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) GetAllGuests(c *gin.Context) {
	req, err := validateGetAllParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	result, err := h.storage.Guest().GetAll(&repo.GetAllGuestsParams{
		Page:   req.Page,
		Limit:  req.Limit,
		Search: req.Search,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, getGuestResponse(result))
}

// @Security ApiKeyAuth
// @Router /guest/{id} [put]
// @Summary Update a guest
// @Description Update a guest
// @Tags guest
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Param guest body models.CreateGuestRequest true "Guest"
// @Success 200 {object} models.Guest
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) UpdateGuest(c *gin.Context) {
	var req repo.Guest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	req.ID = int64(id)

	updated, err := h.storage.Guest().Update(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusCreated, parseGuestModel(updated))
}

// @Security ApiKeyAuth
// @Router /guest/{id} [delete]
// @Summary Delete a guest
// @Description Delete a guest
// @Tags guest
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} models.ResponseOK
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) DeleteGuest(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	err = h.storage.Guest().Delete(int64(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully deleted",
	})
}

func getGuestResponse(data *repo.GetAllGuestsResult) *models.GetAllGuestsResponse {
	response := models.GetAllGuestsResponse{
		Guests: make([]*models.Guest, 0),
		Count: data.Count,
	}

	for _, user := range data.Guests {
		u := parseGuestModel(user)
		response.Guests = append(response.Guests, &u)
	}

	return &response
}

func parseGuestModel(guest *repo.Guest) models.Guest {
	return models.Guest{
		ID:          guest.ID,
		FirstName:   guest.FirstName,
		LastName:    guest.LastName,
		PhoneNumber: guest.PhoneNumber,
		Email:       guest.Email,
		CreatedAt:   guest.CreatedAt,
	}
}
