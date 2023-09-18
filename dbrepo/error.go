package dbrepo

import (
	"errors"
)

// 执行数SQL，没有影响的行
var ErrNotRowsAffected = errors.New("not rows affected")

func handleRowsAffected(rowsAff int64, err error) error {
	if err != nil {
		return err
	}

	if rowsAff == 0 {
		return ErrNotRowsAffected
	}

	return nil
}
