package v1

import (
	"errors"
	"net/http"

	"github.com/burxondv/hotel_exam/pkg"
	"github.com/gin-gonic/gin"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationPayloadKey = "authorization_payload"
)

func (h *handlerV1) AuthMiddleware(c *gin.Context) {
	accessToken := c.GetHeader(authorizationHeaderKey)

	if len(accessToken) == 0 {
		err := errors.New("authorization header is not provided")
		c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	payload, err := pkg.VerifyToken(h.cfg, accessToken)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	c.Set(authorizationPayloadKey, payload)
	c.Next()
}

func (m *handlerV1) GetAuthPayload(ctx *gin.Context) (*pkg.Payload, error) {
	i, exists := ctx.Get(authorizationPayloadKey)
	if !exists {
		return nil, errors.New("")
	}

	payload, ok := i.(*pkg.Payload)
	if !ok {
		return nil, errors.New("unknown user")
	}
	return payload, nil
}
