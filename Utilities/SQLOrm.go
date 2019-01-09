package Utilities

import (
	"database/sql"
	"github.com/pkg/errors"
	"reflect"
)

func RowsToStruct(rows *sql.Rows, to interface{}) error {
	v := reflect.ValueOf(to)
	if v.Elem().Type().Kind() != reflect.Struct {
		return errors.New("Expect a struct")
	}

	var scanDest []interface{}
	columnNames, _ := rows.Columns()

	addrByColumnName := map[string]interface{}{}

	for i := 0; i < v.Elem().NumField(); i++ {
		oneValue := v.Elem().Field(i)
		columnName := v.Elem().Type().Field(i).Tag.Get("bb_data")
		if columnName == "" {
			columnName = oneValue.Type().Name()
		}
		put := oneValue.Addr().Interface()
		addrByColumnName[columnName] = put
	}

	for _, columnName := range columnNames {
		scanDest = append(scanDest, addrByColumnName[columnName])
	}

	return rows.Scan(scanDest...)

}
