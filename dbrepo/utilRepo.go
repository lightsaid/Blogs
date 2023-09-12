package dbrepo

import (
	"context"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

// utilRepo 工具类，提供一些通用的方法,可嵌套在表 Repo中使用
type utilRepo struct{}

// execTx 定义个执行事务公共的方法
func (utilRepo) execTx(ctx context.Context, qb Queryable, fn func(*Repository) error) error {
	// qb => sql.DB/sql.Tx
	db, ok := qb.(*sqlx.DB)
	if !ok {
		return errors.New("DB is not *sqlx.DB")
	}

	// 开启事务
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	repository := NewRepository(tx)

	if err = fn(repository); err != nil {
		// 执行不成功，则 Rollback
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	// 执行成功 Commit
	return tx.Commit()
}
