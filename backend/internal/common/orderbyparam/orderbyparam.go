package orderbyparam

import (
	"errors"

	"github.com/gin-gonic/gin"
)

const orderingQueryParam = "ordering"

// GenerateOrdering used to generate a sql "order by" statement from request query parameters.
func GenerateOrdering(c *gin.Context, mapping map[string]string, defaultQueryParamValue string) (string, error) {
	orderingParam := c.Query(orderingQueryParam)
	if orderingParam == "" {
		return mapQueryParamValueToOrdering(defaultQueryParamValue, mapping)
	}
	return mapQueryParamValueToOrdering(orderingParam, mapping)
}

func mapQueryParamValueToOrdering(queryParamValue string, mapping map[string]string) (string, error) {
	if queryParamValue == "" {
		return "", errors.New("ordering query param value is empty")
	}

	isDesc := string(queryParamValue[0]) == "-"

	var orderingField string
	var ok bool
	if isDesc {
		orderingField, ok = mapping[queryParamValue[1:]]
	} else {
		orderingField, ok = mapping[queryParamValue]
	}
	if !ok {
		return "", errors.New("ordering query param value is not valid")
	}

	if isDesc {
		return orderingField + " DESC", nil
	} else {
		return orderingField, nil
	}
}
