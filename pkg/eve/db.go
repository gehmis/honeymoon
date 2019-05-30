package eve

import (
	"database/sql"
	"fmt"
)

type column struct {
	value interface{}
}

type row map[string]interface{}

func (c *column) Scan(v interface{}) error {
	switch t := v.(type) {
	case bool, float64, int64, nil:
		c.value = t
	case []uint8:
		c.value = string(t)
	default:
		return fmt.Errorf("unknown type: %T", t)
	}
	return nil
}

func query(format string, args ...interface{}) ([]row, error) {
	db, err := sql.Open("sqlite3", "./sqlite-latest.sqlite")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query(fmt.Sprintf(format, args...))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	rs := []row{}

	for rows.Next() {
		ir := make([]interface{}, len(columns))
		for i := range ir {
			ir[i] = &column{}
		}

		if err := rows.Scan(ir...); err != nil {
			return nil, err
		}

		r := row{}

		for i := range ir {
			r[columns[i]] = ir[i].(*column).value
		}

		rs = append(rs, r)
	}

	return rs, nil
}
