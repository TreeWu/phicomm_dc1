package transaction

import (
	"context"
	"gorm.io/gorm"
)

var contextTransactionKey struct{}

// ExecTx 开启一个事务，数据库操作要写在fn里面，并且通过GetDB获取session，否则事务会失效
func ExecTx(ctx context.Context, db *gorm.DB, fn func(ctx2 context.Context) error) error {
	if ctx.Value(contextTransactionKey) != nil {
		return fn(ctx)
	}
	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		ctx := context.WithValue(ctx, contextTransactionKey, tx)
		return fn(ctx)
	})
}

// GetDB 获取事务操作session对象，如果没有开启事务，就会返回普通session
func GetDB(ctx context.Context, db *gorm.DB) *gorm.DB {
	if ctx.Value(contextTransactionKey) != nil {
		return ctx.Value(contextTransactionKey).(*gorm.DB)
	}
	return db.WithContext(ctx)
}
