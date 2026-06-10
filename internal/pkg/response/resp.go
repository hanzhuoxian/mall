package response

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/hanzhuoxian/mall/pkg/log"
)

type ErrResponse struct {
	Code      int    `json:"code"`
	Message   string `json:"message"`
	Reference string `json:"reference,omitempty"`
}

func WriteResponse(c *gin.Context, err error, data any) {
	if err != nil {
		log.Errorf("%#+v", err)

		return
	}

	c.JSON(http.StatusOK, data)
}
