package middlewares

import (
	"net/http"
	"time"

	"github.com/berk-karaal/letuspass/backend/internal/logging"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

// LogHandler middleware logs requests
func LogHandler(logger *logging.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		rawQuery := c.Request.URL.RawQuery
		if rawQuery != "" {
			path = path + "?" + rawQuery
		}

		c.Next()

		processTime := time.Now().Sub(start)

		msg := "Request"
		if len(c.Errors) > 0 {
			// TODO: improve ?
			msg = c.Errors.String()
		}

		var event *zerolog.Event

		switch {
		case c.Writer.Status() >= http.StatusBadRequest && c.Writer.Status() < http.StatusInternalServerError:
			event = logger.NewEvent(zerolog.WarnLevel)
		case c.Writer.Status() >= http.StatusInternalServerError:
			event = logger.NewEvent(zerolog.ErrorLevel)
		default:
			event = logger.NewEvent(zerolog.InfoLevel)
		}
		event.
			Int("status", c.Writer.Status()).
			Str("method", c.Request.Method).
			Str("path", path).
			Str("ip", c.ClientIP()).
			Dur("process_time", processTime).
			Str("user_agent", c.Request.UserAgent()).
			Int("body_size", c.Writer.Size()).
			Str("request_id", requestid.Get(c)).
			Msg(msg)
	}
}
