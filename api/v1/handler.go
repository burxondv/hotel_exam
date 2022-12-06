package v1

import (
	"strconv"

	"github.com/burxondv/hotel_exam/api/models"
	"github.com/burxondv/hotel_exam/config"
	"github.com/burxondv/hotel_exam/storage"
	"github.com/burxondv/hotel_exam/storage/repo"
	"github.com/gin-gonic/gin"
)

type handlerV1 struct {
	cfg      *config.Config
	storage  storage.StorageI
	inMemory storage.InMemoryStorageI
}

type HandlerV1Options struct {
	Cfg      *config.Config
	Storage  storage.StorageI
	InMemory storage.InMemoryStorageI
}

func New(options *HandlerV1Options) *handlerV1 {
	return &handlerV1{
		cfg:      options.Cfg,
		storage:  options.Storage,
		inMemory: options.InMemory,
	}
}

func errorResponse(err error) *models.ErrorResponse {
	return &models.ErrorResponse{
		Error: err.Error(),
	}
}

func validateGetAllParams(c *gin.Context) (*models.GetAllParams, error) {
	var (
		limit int = 10
		page  int = 1
		err   error
	)

	if c.Query("limit") != "" {
		limit, err = strconv.Atoi(c.Query("limit"))
		if err != nil {
			return nil, err
		}
	}

	if c.Query("page") != "" {
		page, err = strconv.Atoi(c.Query("page"))
		if err != nil {
			return nil, err
		}
	}

	return &models.GetAllParams{
		Limit:  int32(limit),
		Page:   int32(page),
		Search: c.Query("search"),
	}, nil
}

func validateGetAllRoomParams(c *gin.Context) (*repo.GetAllRoomsParams, error) {
	var (
		limit    int = 10
		page     int = 1
		hotel_id int
		err      error
	)

	if c.Query("limit") != "" {
		limit, err = strconv.Atoi(c.Query("limit"))
		if err != nil {
			return nil, err
		}
	}

	if c.Query("page") != "" {
		page, err = strconv.Atoi(c.Query("page"))
		if err != nil {
			return nil, err
		}
	}

	if c.Query("hotel_id") != "" {
		hotel_id, err = strconv.Atoi(c.Query("hotel_id"))
		if err != nil {
			return nil, err
		}
	}

	return &repo.GetAllRoomsParams{
		Limit:   int32(limit),
		Page:    int32(page),
		HotelID: int32(hotel_id),
	}, nil
}
