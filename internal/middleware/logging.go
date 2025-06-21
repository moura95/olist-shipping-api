package middleware

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func RequestLogMiddleware(logger *zap.SugaredLogger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		startTime := time.Now()

		// Gera trace ID único para a requisição
		traceID := generateTraceID()

		// Adiciona trace ID ao contexto
		ctx.Set("trace_id", traceID)

		var requestBody string
		if ctx.Request.Body != nil {
			bodyBytes, err := io.ReadAll(ctx.Request.Body)
			if err == nil {
				requestBody = string(bodyBytes)
				ctx.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			}
		}

		// Log da requisição
		logger.Infow("Request",
			"trace_id", traceID,
			"method", ctx.Request.Method,
			"path", ctx.Request.URL.Path,
			"query", ctx.Request.URL.RawQuery,
			"user_agent", ctx.Request.UserAgent(),
			"client_ip", ctx.ClientIP(),
			"headers", ctx.Request.Header,
			"body", requestBody,
			"timestamp", startTime.Format(time.RFC3339),
		)

		enrichedLogger := logger.With("trace_id", traceID)
		ctx.Set("logger", enrichedLogger)

		ctx.Next()
	}
}

func generateTraceID() string {
	// Gera UUID v4
	uuid := uuid.New()

	// Cria hash MD5 para encurtar
	hash := md5.Sum([]byte(uuid.String()))

	return fmt.Sprintf("%x", hash)[:8]
}

func GetLoggerFromContext(ctx *gin.Context) *zap.SugaredLogger {
	if logger, exists := ctx.Get("logger"); exists {
		if enrichedLogger, ok := logger.(*zap.SugaredLogger); ok {
			return enrichedLogger
		}
	}
	// Fallback para logger padrão
	return zap.NewNop().Sugar()
}
