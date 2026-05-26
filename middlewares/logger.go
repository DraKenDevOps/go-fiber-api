package middleware

import (
	"encoding/json"
	"slices"
	"time"

	"go-fiber-api/zplogger"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func LogRequestResponse(ignorePaths []string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if slices.Contains(ignorePaths, c.Path()) {
			return c.Next()
		}

		requestID := uuid.New().String()
		c.Set("X-Request-ID", requestID)

		// ipaddr := c.Get("X-Real-Ip")
		clientIP := c.IP()
		forwardedFor := c.Get("X-Forwarded-For")
		userAgent := c.Get("User-Agent")

		var reqBody any
		bodyBytes := c.Body()
		if len(bodyBytes) > 0 {
			if err := json.Unmarshal(bodyBytes, &reqBody); err != nil {
				reqBody = string(bodyBytes)
			}
		}

		zplogger.Logger.Info(
			"Request",
			zap.String("requestId", requestID),
			zap.String("method", c.Method()),
			zap.String("url", c.OriginalURL()),
			zap.String("ip", clientIP),
			zap.String("forwardedFor", forwardedFor),
			zap.String("userAgent", userAgent),
			zap.Any("reqbody", truncateData(reqBody, 500)),
		)

		startTime := time.Now()

		err := c.Next()

		duration := time.Since(startTime)

		var resBody any
		resBodyBytes := c.Response().Body()
		if len(resBodyBytes) > 0 {
			if err := json.Unmarshal(resBodyBytes, &resBody); err != nil {
				resBody = string(resBodyBytes)
			}
		}
		zplogger.Logger.Info(
			"Response",
			zap.String("requestId", requestID),
			zap.Int("statusCode", c.Response().StatusCode()),
			zap.Duration("duration", duration),
			zap.String("path", c.Path()),
			zap.String("method", c.Method()),
			zap.Any("resbody", truncateData(resBody, 100)),
		)

		return err
	}
}

func LogRequestResponseSimple(ignorePaths []string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if slices.Contains(ignorePaths, c.Path()) {
			return c.Next()
		}

		requestID := uuid.New().String()
		c.Set("X-Request-ID", requestID)

		zplogger.Logger.Info(
			"Request",
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.String("requestId", requestID),
			zap.String("ip", c.IP()),
		)

		startTime := time.Now()
		err := c.Next()
		duration := time.Since(startTime)

		zplogger.Logger.Info(
			"Response",
			zap.String("requestId", requestID),
			zap.Int("statusCode", c.Response().StatusCode()),
			zap.Duration("duration", duration),
		)

		return err
	}
}

func ErrorLoggingMiddleware(c *fiber.Ctx) error {
	err := c.Next()

	if err != nil {
		requestID := c.Get("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		zplogger.Logger.Error(
			"Request error",
			zap.String("requestId", requestID),
			zap.String("path", c.Path()),
			zap.String("method", c.Method()),
			zap.Int("statusCode", c.Response().StatusCode()),
			zap.Error(err),
		)

		return c.Status(c.Response().StatusCode()).JSON(fiber.Map{
			"status":    "error",
			"message":   err.Error(),
			"requestId": requestID,
		})
	}

	return nil
}

func truncateData(data interface{}, maxLength int) interface{} {
	if maxLength <= 0 {
		maxLength = 500
	}

	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return data
	}

	if len(jsonBytes) > maxLength {
		return string(jsonBytes[:maxLength]) + "... (truncated)"
	}

	return data
}

func GetRequestID(c *fiber.Ctx) string {
	return c.Get("X-Request-ID")
}

func LogWithRequestID(c *fiber.Ctx, level string, message string, fields ...zap.Field) {
	fields = append(fields, zap.String("requestId", GetRequestID(c)))

	switch level {
	case "info":
		zplogger.Logger.Info(message, fields...)
	case "warn":
		zplogger.Logger.Warn(message, fields...)
	case "error":
		zplogger.Logger.Error(message, fields...)
	case "debug":
		zplogger.Logger.Debug(message, fields...)
	default:
		zplogger.Logger.Info(message, fields...)
	}
}
