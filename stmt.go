package orm

import (
	"fmt"
	"strings"
)

type Wherer struct {
	col  string
	args []any
}

func Where(col string, args ...any) *Wherer {
	if len(args) > 1 {
		panic("orm: Where can only have one argument")
	}
	return &Wherer{
		col:  col,
		args: args,
	}
}

func WhereStmt(stmt string, wheres []*Wherer, args ...any) (string, []any) {
	if len(args) == 0 {
		args = make([]any, 0, len(wheres))
	}

	whereStmt := make([]string, 0, len(wheres))

	for _, where := range wheres {
		if len(where.args) == 0 {
			whereStmt = append(whereStmt, where.col)
		} else {
			whereStmt = append(whereStmt, fmt.Sprintf("%s = $%d", where.col, len(args)+1))
			args = append(args, where.args[0])
		}
	}

	return fmt.Sprintf("%s WHERE %s", stmt, strings.Join(whereStmt, " AND ")), args
}

func OrderStmt(stmt string, orderBy ...string) string {
	if len(orderBy) == 0 {
		return stmt
	}

	return fmt.Sprintf("%s ORDER BY %s", stmt, strings.Join(orderBy, ", "))
}

type setter struct {
	col  string
	args []any
}

func Set(col string, args ...any) *setter {
	if len(args) > 1 {
		panic("orm: setter can only have one argument")
	}

	return &setter{
		col:  col,
		args: args,
	}
}

func UpdateStmt(stmt string, setters ...*setter) (string, []any) {
	args := make([]any, 0, len(setters))
	setStmt := make([]string, 0, len(setters))

	for _, setter := range setters {
		if len(setter.args) == 0 {
			setStmt = append(setStmt, setter.col)
		} else {
			setStmt = append(setStmt, fmt.Sprintf("%s = $%d", setter.col, len(args)+1))
			args = append(args, setter.args[0])
		}
	}

	return fmt.Sprintf(
		`%s SET %s`,
		stmt,
		strings.Join(setStmt, ", "),
	), args
}
