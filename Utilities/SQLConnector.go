package Utilities

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type DBObject struct {
	Host        string
	User        string
	Passwd      string
	conn_string string
	conn        *sql.DB
}

func CreateDBObject(host string, user string, passwd string) (*DBObject, error) {
	ret := DBObject{
		Host:   host,
		User:   user,
		Passwd: passwd,
	}

	ret.conn_string = fmt.Sprintf("%s:%s@tcp(%s)/apps?charset=utf8&loc=Asia%%2FTaipei&parseTime=true", user, passwd, host)

	db, err := sql.Open("mysql", ret.conn_string)
	log.Debugf("Attempting logging in with Server String : %s", ret.conn_string)
	if err != nil {
		wrapped := errors.Wrap(err, "Error while initialization with sql string : "+ret.conn_string)
		return nil, wrapped
	}
	log.Debugf("Attempting pinging : %s", ret.conn_string)
	err = db.Ping()
	if err != nil {
		return nil, errors.Wrap(err, "Fail to ping server : "+ret.conn_string)
	}

	ret.conn = db

	return &ret, nil
}

func (db *DBObject) GetDB() *sql.DB {
	return db.conn
}

type RowResult map[string]interface{}

//map[string]interface{} waste more memory then map[string][]interface{}, but let use it first
func QueryToMap(conn *sql.DB, queryString string) ([]RowResult, error) {
	result, err := conn.Query(queryString)
	if err != nil {
		return nil, errors.Wrap(err, "Error when querying : "+queryString)
	}
	//defer result.Close()
	defer func() {
		err := result.Close()
		if err != nil {
			log.Errorf("Error closing result set!")
		}
	}()

	columns, err := result.Columns()
	ret := make([]RowResult, 0)

	if err != nil {
		return nil, errors.Wrap(err, "Error when trying getting columns for "+queryString)
	}

	for i := 0; result.Next(); {
		columnInstance := make([]interface{}, len(columns))
		columnPointers := make([]interface{}, len(columns))

		for i := range columns {
			columnPointers[i] = &columnInstance[i]
		}

		row := make(RowResult)
		err = result.Scan(columnPointers...)
		if err != nil {
			return nil, errors.Wrap(err, "Error reading from column points, query string is "+queryString)
		}
		for j := 0; j < len(columns); {
			row[columns[j]] = columnInstance[j]
			j++
		}
		ret = append(ret, row)
		i++
	}

	return ret, nil
}
