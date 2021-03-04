package persistence

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strconv"
)

func ApplyPaginationSettings(c *gin.Context, tx *gorm.DB) *gorm.DB {
	limit, offset := getQuerySettings(c)
	return tx.Limit(limit).Offset(offset)
}

func ApplyFilterQueries(c *gin.Context, tx *gorm.DB, whitelistedFilters []string) *gorm.DB {
	queries := c.QueryMap("q")
	for key, value := range queries {
		for _, filter := range whitelistedFilters {
			if key == filter {
				tx = tx.Where(fmt.Sprintf("%s = ?", filter), value)
				break
			}
		}
	}
	return tx
}

func getQuerySettings(c *gin.Context) (limit, offset int) {
	limit = getQueryValueAsIntegerOrDefault(c, "limit", 20)
	offset = getQueryValueAsIntegerOrDefault(c, "offset", 0)
	return
}

func getQueryValueAsIntegerOrDefault(c *gin.Context, key string, def int) int {
	val := c.Query(key)
	if val == "" {
		return def
	}
	v, err := strconv.Atoi(val)
	if err != nil {
		return def
	}
	return v
}
