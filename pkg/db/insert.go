package db

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"
)

type Transfer struct {
	From  string
	To    string
	Value string
}

func (c *DBClient) InsertUSDCTransfer(ctx context.Context, transfers *[]Transfer) (sql.Result, error) {
	sqlStr := `INSERT INTO usdc_transfer('from', 'to', 'value') VALUES `
	vals := []interface{}{}

	for _, row := range *transfers {
		sqlStr += "(?, ?, ?),"
		vals = append(vals, row.From, row.To, row.Value)
	}
	sqlStr = sqlStr[0 : len(sqlStr)-1]
	stmt, err := c.db.Prepare(sqlStr)
	if err != nil {
		return nil, errors.Wrap(err, "failed to c.db.Prepare(sqlStr)")
	}

	res, err := stmt.ExecContext(ctx, vals...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to insertUSDCTransfer")
	}

	return res, nil
}
