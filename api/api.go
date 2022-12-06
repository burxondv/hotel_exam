package api

import (
	v1 "github.com/burxondv/hotel_exam/api/v1"
	"github.com/burxondv/hotel_exam/config"
	"github.com/burxondv/hotel_exam/storage"
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware

	_ "github.com/burxondv/hotel_exam/api/docs"
)

type RouterOptions struct {
	Cfg      *config.Config
	Storage  storage.StorageI
	InMemory storage.InMemoryStorageI
}

// @title           Swagger for hotel api
// @version         1.0
// @description     This is a hotel service api.
// @host      localhost:8000
// @BasePath  /v1
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @Security ApiKeyAuth
func New(opt *RouterOptions) *gin.Engine {
	router := gin.Default()

	handlerV1 := v1.New(&v1.HandlerV1Options{
		Cfg:      opt.Cfg,
		Storage:  opt.Storage,
		InMemory: opt.InMemory,
	})

	router.Static("/media", "./media")

	apiV1 := router.Group("/v1")

	apiV1.POST("/auth-guest/register", handlerV1.RegisterGuest)
	apiV1.POST("/auth-guest/verify", handlerV1.VerifyGuest)
	apiV1.POST("/auth-guest/login", handlerV1.LoginGuest)
	apiV1.POST("/auth-guest/forgot-password", handlerV1.ForgotPasswordGuest)
	apiV1.POST("/auth-guest/verify-forgot-password", handlerV1.VerifyForgotPasswordGuest)
	apiV1.POST("/auth-guest/update-password", handlerV1.AuthMiddleware, handlerV1.UpdatePasswordGuest)

	apiV1.POST("/auth-owner/register", handlerV1.RegisterOwner)
	apiV1.POST("/auth-owner/verify", handlerV1.VerifyOwner)
	apiV1.POST("/auth-owner/login", handlerV1.LoginOwner)
	apiV1.POST("/auth-owner/forgot-password", handlerV1.ForgotPasswordOwner)
	apiV1.POST("/auth-owner/verify-forgot-password", handlerV1.VerifyForgotPasswordOwner)
	apiV1.POST("/auth-owner/update-password", handlerV1.AuthMiddleware, handlerV1.UpdatePasswordOwner)

	apiV1.POST("/guest", handlerV1.CreateGuest)
	apiV1.GET("/guest/:id", handlerV1.GetGuest)
	apiV1.GET("/guest", handlerV1.GetAllGuests)
	apiV1.PUT("/guest/:id", handlerV1.AuthMiddleware, handlerV1.UpdateGuest)
	apiV1.DELETE("/guest/:id", handlerV1.AuthMiddleware, handlerV1.DeleteGuest)


	apiV1.POST("/owner", handlerV1.CreateOwner)
	apiV1.GET("/owner/:id", handlerV1.GetOwner)
	apiV1.GET("/owner", handlerV1.GetAllOwner)
	apiV1.PUT("/owner/:id", handlerV1.AuthMiddleware, handlerV1.UpdateOwner)
	apiV1.DELETE("/owner/:id", handlerV1.AuthMiddleware, handlerV1.DeleteOwner)

	apiV1.POST("hotel", handlerV1.CreateHotel)
	apiV1.GET("hotel/:id", handlerV1.GetHotel)
	apiV1.GET("hotel", handlerV1.GetAllHotels)
	apiV1.PUT("hotel/:id", handlerV1.AuthMiddleware, handlerV1.UpdateHotel)
	apiV1.DELETE("hotel/:id", handlerV1.AuthMiddleware, handlerV1.DeleteHotel)
	apiV1.POST("hotel/file-upload", handlerV1.AuthMiddleware, handlerV1.UploadFileHotel)

	apiV1.POST("room", handlerV1.CreateRoom)
	apiV1.GET("room/:id", handlerV1.GetRoom)
	apiV1.GET("room", handlerV1.GetAllRoom)
	apiV1.PUT("room/:id", handlerV1.AuthMiddleware, handlerV1.UpdateRoom)
	apiV1.DELETE("room/:id", handlerV1.AuthMiddleware, handlerV1.DeleteRoom)
	apiV1.POST("room/file-upload", handlerV1.AuthMiddleware, handlerV1.UploadFileRoom)


	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}
