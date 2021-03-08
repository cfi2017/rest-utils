package persistence

import (
	"errors"
	"github.com/cfi2017/rest-utils/pkg/util"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

/*
PersistenceMiddleware provides a transaction per request
*/
func PersistenceMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		_, t := NewTransaction(db.WithContext(c))
		defer t.Close()
		c.Set("tx", t)
		c.Next()
	}
}

func GetTx(ctx *gin.Context) *gorm.DB {
	return ctx.MustGet("tx").(*Transaction).Tx
}

func FailTx(ctx *gin.Context) {
	t := ctx.MustGet("tx").(*Transaction)
	t.Fail()
}

type StaticHandler func(c *gin.Context, tx *gorm.DB, id string) (entity interface{}, err error)

type EntityMiddlewareOpts struct {
	Preloads        []string
	ContinueOnError bool
	StaticHandler   StaticHandler
	StaticPaths     []string
}

func EntityMiddleware(c *gin.Context, id string, entity interface{}, opts *EntityMiddlewareOpts) {
	tx := GetTx(c)
	var continueOnError bool
	if opts != nil {
		if opts.Preloads != nil {
			for _, p := range opts.Preloads {
				tx.Preload(p)
			}
		}
		continueOnError = opts.ContinueOnError
		if opts.StaticPaths != nil {
			for _, path := range opts.StaticPaths {
				if path == id {
					entity, err := opts.StaticHandler(c, tx, id)
					if err != nil {
						_ = c.Error(err)
						if !opts.ContinueOnError {
							c.Abort()
							return
						}
						return
					}
					c.Set("entity", entity)
					return
				}
			}
		}
	}
	if err := tx.First(entity, id).Error; err != nil && !continueOnError {
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
