package helper

import (
	"context"
	"database/sql"
	"fmt"
)

func WithTransaction(ctx context.Context, db *sql.DB, fn func(*sql.Tx) (interface{}, error)) (interface{}, error) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				err = fmt.Errorf("error rolling back transaction: %v (original error : %w)", rbErr, err)
			}
		} else {
			cerr := tx.Commit()
			if cerr != nil {
				err = fmt.Errorf("error committing transaction: %v (original error : %w)", cerr, err)
			}
		}
	}()

	res, err := fn(tx)

	return res, err
}
