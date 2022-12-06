package v1

import (
	"net/http"
	"strconv"

	"github.com/burxondv/hotel_exam/api/models"
	"github.com/burxondv/hotel_exam/storage/repo"
	"github.com/gin-gonic/gin"
)

// @Security ApiKeyAuth
// @Router /owner [post]
// @Summary Create a owner
// @Description Create a owner
// @Tags owner
// @Accept json
// @Produce json
// @Param owner body models.CreateOwnerRequest true "Owner"
// @Success 201 {object} models.Owner
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) CreateOwner(c *gin.Context) {
	var (
		req models.CreateGuestRequest
	)

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	resp, err := h.storage.Owner().Create(&repo.Owner{
		HotelID:     req.HotelID,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Email:       req.Email,
		Password:    req.Password,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusCreated, parseOwnerModel(resp))
}

// @Router /owner/{id} [get]
// @Summary Get owner by id
// @Description Get owner by id
// @Tags owner
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} models.Owner
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) GetOwner(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	resp, err := h.storage.Owner().Get(int64(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, parseOwnerModel(resp))
}

// @Router /owner [get]
// @Summary Get all owner
// @Description Get all owner
// @Tags owner
// @Accept json
// @Produce json
// @Param filter query models.GetAllParams false "Filter"
// @Success 200 {object} models.GetAllOwnersResponse
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) GetAllOwner(c *gin.Context) {
	req, err := validateGetAllParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	result, err := h.storage.Owner().GetAll(&repo.GetAllOwnersParams{
		Page:   req.Page,
		Limit:  req.Limit,
		Search: req.Search,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, getOwnerResponse(result))
}

// @Security ApiKeyAuth
// @Router /owner/{id} [put]
// @Summary Update a owner
// @Description Update a owner
// @Tags owner
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Param owner body models.CreateOwnerRequest true "Owner"
// @Success 200 {object} models.Owner
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) UpdateOwner(c *gin.Context) {
	var req repo.Owner

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

	updated, err := h.storage.Owner().Update(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusCreated, parseOwnerModel(updated))
}

// @Security ApiKeyAuth
// @Router /owner/{id} [delete]
// @Summary Delete a owner
// @Description Delete a owner
// @Tags owner
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} models.ResponseOK
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) DeleteOwner(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	err = h.storage.Owner().Delete(int64(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully deleted",
	})
}

func getOwnerResponse(data *repo.GetAllOwnersResult) *models.GetAllOwnersResponse {
	response := models.GetAllOwnersResponse{
		Owners: make([]*models.Owner, 0),
		Count:  data.Count,
	}

	for _, owner := range data.Owners {
		u := parseOwnerModel(owner)
		response.Owners = append(response.Owners, &u)
	}

	return &response
}

func parseOwnerModel(owner *repo.Owner) models.Owner {
	return models.Owner{
		ID:          owner.ID,
		FirstName:   owner.FirstName,
		LastName:    owner.LastName,
		Email:       owner.Email,
		CreatedAt:   owner.CreatedAt,
	}
}
