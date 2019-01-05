package BossBot

import (
	"Utilities"
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"reflect"
)

type DataManager struct {
	config *Configuration
	db     *Utilities.DBObject
}

var instantiated *DataManager = nil

func GetDataManager(conf *Configuration) (*DataManager, error) {
	if instantiated == nil {
		instantiated = new(DataManager)
		instantiated.config = conf
		dbObject, err := Utilities.CreateDBObject(conf.SqlHost, conf.SqlAcc, conf.SqlPass)
		if err != nil {
			return nil, errors.Wrap(err, "Error creating DBObject in DataManager!")
		}

		instantiated.db = dbObject
	}
	return instantiated, nil
}

type BroadcastItem struct {
	Id            int            `bb_data:"id"`
	StartDate     mysql.NullTime `bb_data:"start"`
	EndDate       mysql.NullTime `bb_data:"end"`
	Message       string         `bb_data:"message"`
	MessageId     int            `bb_data:"message_id"`
	Day           sql.NullInt64  `bb_data:"day_in_month"`
	WeekDay       sql.NullInt64  `bb_data:"day_in_week"`
	BroadcastTime mysql.NullTime `bb_data:"broadcast_time"`
	ChannelName   string         `bb_data:"channel_name"`
}

func ToStruct(rows *sql.Rows, to interface{}) error {
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

func (dm *DataManager) GetBroadcastList() ([]BroadcastItem, error) {
	rows, err := dm.db.GetConnection().Query(
		`select bs.id, bs.start_date as start, bs.end_date as end, bm.message as message, bs.message_id as message_id, bs.recursive_day_in_month as day_in_month, bs.recursive_day_in_week as day_in_week, bs.broadcast_time, bbc.name as channel_name
				from bb_broadcast_schedule as bs
				inner join bb_broadcast_msg as bm on bs.message_id = bm.id
				inner join bb_broadcast_channel as bbc on bs.channel_id = bbc.id;`)

	if err != nil {
		return nil, errors.Wrap(err, "In GetBroadcastList()")
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			log.Warningf("Error closing rows in GetBroadcastList")
		}
	}()

	var ret []BroadcastItem

	if err != nil {
		return nil, errors.Wrap(err, "In GetBroadcastList()")
	}

	for rows.Next() {
		//First, create struct, and fill pointer list according to column name list
		bi := BroadcastItem{}
		err = ToStruct(rows, &bi)
		if err != nil {
			log.Errorln(err)
		}
		ret = append(ret, bi)
	}

	return ret, nil
}
