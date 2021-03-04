package persistence

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"sync"
)

type Transaction struct {
	once     sync.Once
	rollback bool
	Tx       *gorm.DB
}

func (t *Transaction) Close() {
	t.once.Do(func() {
		if t.rollback {
			t.Tx.Rollback()
		} else {
			t.Tx.Commit()
		}
	})
}

func (t *Transaction) Fail() {
	t.rollback = true
}

func NewTransaction(ctx *gin.Context, db *gorm.DB) (*gorm.DB, *Transaction) {
	tx := db.WithContext(ctx).Begin()
	return tx, &Transaction{Tx: tx}
}
