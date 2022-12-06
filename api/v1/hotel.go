package v1

import (
	"net/http"
	"strconv"

	"github.com/burxondv/hotel_exam/api/models"
	"github.com/burxondv/hotel_exam/storage/repo"
	"github.com/gin-gonic/gin"
)

// @Security ApiKeyAuth
// @Router /hotel [post]
// @Summary Create a hotel
// @Description Create a hotel
// @Tags hotel
// @Accept json
// @Produce json
// @Param hotel body models.CreateHotelRequest true "Hotel"
// @Success 201 {object} models.Hotel
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) CreateHotel(c *gin.Context) {
	var (
		req models.CreateHotelRequest
	)

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	resp, err := h.storage.Hotel().Create(&repo.Hotel{
		Name:         req.Name,
		Address:      req.Address,
		StarsPopular: req.StarsPopular,
		ImageUrl:     req.ImageUrl,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusCreated, models.Hotel{
		ID:           resp.ID,
		Name:         resp.Name,
		StarsPopular: resp.StarsPopular,
		ImageUrl:     resp.ImageUrl,
		CreatedAt:    resp.CreatedAt,
	})
}

// @Router /hotel/{id} [get]
// @Summary Get hotel by id
// @Description Get hotel by id
// @Tags hotel
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} models.Hotel
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) GetHotel(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	resp, err := h.storage.Hotel().Get(int64(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, models.Hotel{
		ID:           resp.ID,
		Name:         resp.Name,
		Address:      resp.Address,
		StarsPopular: resp.StarsPopular,
		ImageUrl:     resp.ImageUrl,
		CreatedAt:    resp.CreatedAt,
	})
}

// @Router /hotel [get]
// @Summary Get all hotel
// @Description Get all hotel
// @Tags hotel
// @Accept json
// @Produce json
// @Param filter query models.GetAllParams false "Filter"
// @Success 200 {object} models.GetAllHotelResponse
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) GetAllHotels(c *gin.Context) {
	req, err := validateGetAllParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	result, err := h.storage.Hotel().GetAll(&repo.GetAllHotelsParams{
		Page:   req.Page,
		Limit:  req.Limit,
		Search: req.Search,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, getHotelResponse(result))
}

// @Security ApiKeyAuth
// @Router /hotel/{id} [put]
// @Summary Update a hotel
// @Description Update a hotel
// @Tags hotel
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Param hotel body models.CreateHotelRequest true "Hotel"
// @Success 200 {object} models.Hotel
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) UpdateHotel(c *gin.Context) {
	var req repo.Hotel

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

	updated, err := h.storage.Hotel().Update(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusCreated, parseHotelModel(updated))
}

// @Security ApiKeyAuth
// @Router /hotel/{id} [delete]
// @Summary Delete a hotel
// @Description Delete a hotel
// @Tags hotel
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} models.ResponseOK
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) DeleteHotel(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	err = h.storage.Hotel().Delete(int64(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully deleted",
	})
}

func getHotelResponse(data *repo.GetAllHotelsResult) *models.GetAllHotelResponse {
	response := models.GetAllHotelResponse{
		Hotels: make([]*models.Hotel, 0),
		Count:  data.Count,
	}

	for _, c := range data.Hotels {
		response.Hotels = append(response.Hotels, &models.Hotel{
			ID:           c.ID,
			Name:         c.Name,
			Address:      c.Address,
			StarsPopular: c.StarsPopular,
			ImageUrl:     c.ImageUrl,
			CreatedAt:    c.CreatedAt,
		})
	}

	return &response
}

func parseHotelModel(hotel *repo.Hotel) models.Hotel {
	return models.Hotel{
		ID:           hotel.ID,
		Name:         hotel.Name,
		Address:      hotel.Address,
		StarsPopular: hotel.StarsPopular,
		ImageUrl:     hotel.ImageUrl,
		CreatedAt:    hotel.CreatedAt,
	}
}
