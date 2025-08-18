package model

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ SequenceModel = (*customSequenceModel)(nil)

type (
	// SequenceModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSequenceModel.
	SequenceModel interface {
		sequenceModel
		withSession(session sqlx.Session) SequenceModel

		// 自定义方法
		ReplaceIntoByStub(ctx context.Context, stub string) (sql.Result, error)
	}

	customSequenceModel struct {
		*defaultSequenceModel
	}
)

// NewSequenceModel returns a model for the database table.
func NewSequenceModel(conn sqlx.SqlConn) SequenceModel {
	return &customSequenceModel{
		defaultSequenceModel: newSequenceModel(conn),
	}
}

func (m *customSequenceModel) withSession(session sqlx.Session) SequenceModel {
	return NewSequenceModel(sqlx.NewSqlConnFromSession(session))
}

// 自定义方法实现
func (m *customSequenceModel) ReplaceIntoByStub(ctx context.Context, stub string) (sql.Result, error) {
	query := fmt.Sprintf("REPLACE INTO %s (`stub`) VALUES (?)", m.table)
	return m.conn.ExecCtx(ctx, query, stub)
}
