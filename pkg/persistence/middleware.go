package persistence

import (
	"errors"
	"github.com/cfi2017/rest-utils/pkg/util"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strconv"
)

/*
PersistenceMiddleware provides a transaction per request
*/
func PersistenceMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		_, t := NewTransaction(db)
		defer t.Close()
		ctx.Set("tx", t)
		ctx.Next()
	}
}

func GetTx(ctx *gin.Context) *gorm.DB {
	return ctx.MustGet("tx").(*Transaction).Tx
}

func FailTx(ctx *gin.Context) {
	t := ctx.MustGet("tx").(*Transaction)
	t.Fail()
}

type EntityMiddlewareOpts struct {
	Preloads        []string
	ContinueOnError bool
	StaticHandler   func(c *gin.Context, tx *gorm.DB, id string) (entity interface{}, matches bool, err error)
}

func EntityMiddleware(param string, entity interface{}, opts *EntityMiddlewareOpts) gin.HandlerFunc {
	return func(c *gin.Context) {
		tx := GetTx(c)
		id := c.Param(param)
		var continueOnError bool
		if opts != nil {
			if opts.Preloads != nil {
				for _, p := range opts.Preloads {
					tx = tx.Preload(p)
				}
			}
			continueOnError = opts.ContinueOnError
			if opts.StaticHandler != nil {
				e, ok, err := opts.StaticHandler(c, tx, id)
				if err != nil && !continueOnError {
					_ = c.AbortWithError(500, err)
					return
				}
				if ok {
					c.Set("entity", e)
					return
				}
			}
		}
		idx, err := strconv.Atoi(id)
		if err != nil {
			c.AbortWithStatusJSON(400, util.NewErrorResponse("invalid id", util.ErrorCodeBadRequest))
			return
		}
		if err := tx.First(entity, idx).Error; err != nil && !continueOnError {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.AbortWithStatusJSON(404, util.NewErrorResponse("could not find object", util.ErrorCodeNotFound))
			} else {
				_ = c.AbortWithError(500, err)
				return
			}
		} else {
			c.Set("entity", entity)
			c.Next()
		}
	}
}

func MustGetEntity(c *gin.Context) interface{} {
	entity, ok := GetEntity(c)
	if !ok {
		panic("entity requested but entity middleware not used")
	}
	return entity
}

func GetEntity(c *gin.Context) (interface{}, bool) {
	return c.Get("entity")
}
