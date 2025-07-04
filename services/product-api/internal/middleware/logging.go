package middleware

import (
	"context"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

// RequestIDMiddleware adds a unique request ID to each request
func RequestIDMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			requestID := strconv.FormatInt(time.Now().UnixNano(), 36)
			c.Response().Header().Set("X-Request-ID", requestID)
			c.Set("requestId", requestID)
			
			// Add request ID to context for service layer
			ctx := context.WithValue(c.Request().Context(), "requestId", requestID)
			c.SetRequest(c.Request().WithContext(ctx))
			
			return next(c)
		}
	}
}

// StructuredLoggingMiddleware provides structured request logging
func StructuredLoggingMiddleware(logger *logrus.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			requestID := c.Get("requestId").(string)
			
			// Log request start
			logger.WithFields(logrus.Fields{
				"method":     c.Request().Method,
				"uri":        c.Request().RequestURI,
				"remote_ip":  c.RealIP(),
				"user_agent": c.Request().UserAgent(),
				"request_id": requestID,
				"event":      "request_start",
			}).Info("🌐 Request started")
			
			// Process request
			err := next(c)
			
			// Calculate duration
			duration := time.Since(start)
			
			// Determine log level based on status
			status := c.Response().Status
			logLevel := logrus.InfoLevel
			emoji := "✅"
			
			if status >= 400 && status < 500 {
				logLevel = logrus.WarnLevel
				emoji = "⚠️"
			} else if status >= 500 {
				logLevel = logrus.ErrorLevel
				emoji = "💥"
			}
			
			// Log request completion
			logEntry := logger.WithFields(logrus.Fields{
				"method":        c.Request().Method,
				"uri":           c.Request().RequestURI,
				"status":        status,
				"duration_ms":   duration.Milliseconds(),
				"duration":      duration.String(),
				"bytes_out":     c.Response().Size,
				"remote_ip":     c.RealIP(),
				"request_id":    requestID,
				"event":         "request_complete",
			})
			
			if err != nil {
				logEntry = logEntry.WithError(err)
			}
			
			logEntry.Log(logLevel, emoji+" Request completed")
			
			return err
		}
	}
}

// MetricsMiddleware collects basic request metrics
func MetricsMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			
			err := next(c)
			
			// You could collect metrics here and send to a metrics service
			// For now, we'll just add timing headers
			duration := time.Since(start)
			c.Response().Header().Set("X-Response-Time", duration.String())
			
			return err
		}
	}
}

// ErrorHandlingMiddleware provides centralized error handling
func ErrorHandlingMiddleware(logger *logrus.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := next(c)
			if err != nil {
				requestID := c.Get("requestId")
				
				logger.WithFields(logrus.Fields{
					"error":      err.Error(),
					"method":     c.Request().Method,
					"uri":        c.Request().RequestURI,
					"request_id": requestID,
				}).Error("💥 Request error")
				
				// Return the error to let Echo handle the response
				return err
			}
			return nil
		}
	}
}