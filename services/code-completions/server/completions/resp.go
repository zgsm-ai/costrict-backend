package completions

import (
	"code-completion/pkg/completions"
	"code-completion/pkg/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func respCompletion(c *gin.Context, clientId string, rsp *completions.CompletionResponse) {
	if rsp.Status != model.StatusSuccess {
		zap.L().Warn("completion failed", zap.String("completionID", rsp.ID),
			zap.String("clientID", clientId),
			zap.String("status", string(rsp.Status)),
			zap.Any("response", rsp))
	} else {
		zap.L().Info("completion succeeded", zap.String("completionID", rsp.ID),
			zap.String("clientID", clientId),
			zap.Any("response", rsp))
	}
	statusCode := http.StatusOK
	switch rsp.Status {
	case model.StatusSuccess:
		statusCode = http.StatusOK
	case model.StatusEmpty:
		statusCode = http.StatusNoContent
	case model.StatusCanceled:
		statusCode = 499 // Client Closed Request
	case model.StatusTimeout:
		statusCode = http.StatusGatewayTimeout
	case model.StatusBusy:
		statusCode = http.StatusTooManyRequests
	case model.StatusReqError:
	case model.StatusRejected:
		statusCode = http.StatusBadRequest
	default:
		statusCode = http.StatusInternalServerError
	}
	c.JSON(statusCode, rsp)
}
