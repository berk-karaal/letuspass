package bodybinder

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/berk-karaal/letuspass/backend/internal/schemas"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type validationErrorResponse struct {
	Errors []validationErrorResponseItem `json:"errors"`
}

type validationErrorResponseItem struct {
	Field  string `json:"field"`
	Reason string `json:"reason"`
}

// Bind runs gin.Context.ShouldBind() function underneath, returns true if binding is successful. If binding
// is unsuccessful, writes appropriate error responses to gin context and returns false. Controller functions
// typically just returns itself when this function returns false.
func Bind(obj any, c *gin.Context) (ok bool) {
	if err := c.ShouldBind(obj); err != nil {
		var valErrors validator.ValidationErrors
		if errors.As(err, &valErrors) {
			c.JSON(http.StatusUnprocessableEntity, validationErrorResponse{Errors: marshalValErrors(valErrors)})
			return false
		}
		c.JSON(http.StatusBadRequest, schemas.BadRequestResponse{Error: "Parsing request body failed."})
		return false
	}
	return true
}

// marshalValErrors returns list of fields, which failed on validation, and their fail reasons.
func marshalValErrors(valErrors validator.ValidationErrors) []validationErrorResponseItem {
	var errs []validationErrorResponseItem
	for _, f := range valErrors {
		err := f.ActualTag()
		if f.Param() != "" {
			err = fmt.Sprintf("%s=%s", err, f.Param())
		}
		errs = append(errs, validationErrorResponseItem{Field: f.Field(), Reason: err})
	}
	return errs
}
