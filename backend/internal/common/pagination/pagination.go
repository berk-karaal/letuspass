package pagination

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const maxPageSize = 100

func Paginate(c *gin.Context) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		page, _ := strconv.Atoi(c.Query("page"))
		if page <= 0 {
			page = 1
		}

		pageSize, _ := strconv.Atoi(c.Query("page_size"))
		switch {
		case pageSize > maxPageSize:
			pageSize = maxPageSize
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

type StandardPaginationResponse[T any] struct {
	Results []T `json:"results" binding:"required"`
	Count   int `json:"count" binding:"required"`
}
