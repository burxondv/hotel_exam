package v1

import (
	"net/http"
	"strconv"

	"github.com/burxondv/hotel_exam/api/models"
	"github.com/burxondv/hotel_exam/storage/repo"
	"github.com/gin-gonic/gin"
)

// @Security ApiKeyAuth
// @Router /room [post]
// @Summary Create a room
// @Description Create a room
// @Tags room
// @Accept json
// @Produce json
// @Param room body models.CreateRoomRequest true "Room"
// @Success 201 {object} models.Room
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) CreateRoom(c *gin.Context) {
	var (
		req models.CreateRoomRequest
	)

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	resp, err := h.storage.Room().Create(&repo.Room{
		ImageUrl: req.ImageUrl,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusCreated, models.Room{
		ID:        resp.ID,
		HotelID:   resp.HotelID,
		Status:    resp.Status,
		ImageUrl:  resp.ImageUrl,
		CreatedAt: resp.CreatedAt,
	})
}

// @Router /room/{id} [get]
// @Summary Get room by id
// @Description Get room by id
// @Tags room
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} models.Room
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) GetRoom(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	resp, err := h.storage.Room().Get(int64(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, models.Room{
		ID:        resp.ID,
		HotelID:   resp.HotelID,
		Status:    resp.Status,
		ImageUrl:  resp.ImageUrl,
		CreatedAt: resp.CreatedAt,
	})
}

// @Router /room [get]
// @Summary Get all room
// @Description Get all room
// @Tags room
// @Accept json
// @Produce json
// @Param filter query models.GetAllParams false "Filter"
// @Success 200 {object} models.GetAllRoomResponse
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) GetAllRoom(c *gin.Context) {
	req, err := validateGetAllRoomParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	result, err := h.storage.Room().GetAll(&repo.GetAllRoomsParams{
		Page:    req.Page,
		Limit:   req.Limit,
		HotelID: req.HotelID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, getRoomResponse(result))
}

// @Security ApiKeyAuth
// @Router /room/{id} [put]
// @Summary Update a room
// @Description Update a room
// @Tags room
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Param room body models.CreateRoomRequest true "Room"
// @Success 200 {object} models.Room
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) UpdateRoom(c *gin.Context) {
	var req repo.Room

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

	updated, err := h.storage.Room().Update(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusCreated, parseRoomModel(updated))
}

// @Security ApiKeyAuth
// @Router /room/{id} [delete]
// @Summary Delete a room
// @Description Delete a room
// @Tags room
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} models.ResponseOK
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) DeleteRoom(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	err = h.storage.Room().Delete(int64(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully deleted",
	})
}

func getRoomResponse(data *repo.GetAllRoomsResult) *models.GetAllRoomResponse {
	response := models.GetAllRoomResponse{
		Rooms: make([]*models.Room, 0),
		Count: data.Count,
	}

	for _, c := range data.Rooms {
		response.Rooms = append(response.Rooms, &models.Room{
			ID:        c.ID,
			HotelID:   c.HotelID,
			Status:    c.Status,
			ImageUrl:  c.ImageUrl,
			CreatedAt: c.CreatedAt,
		})
	}

	return &response
}

func parseRoomModel(room *repo.Room) models.Room {
	return models.Room{
		ID:           room.ID,
		HotelID:      room.HotelID,
		Status:       room.Status,
		ImageUrl:     room.ImageUrl,
		CreatedAt:    room.CreatedAt,
	}
}
