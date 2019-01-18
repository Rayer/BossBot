package BossBot

import (
	"Utilities"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
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

func (dm *DataManager) GetBroadcastList() ([]BroadcastItem, error) {
	rows, err := dm.db.GetDB().Query(
		`select bs.id, bs.start_date as start, bs.end_date as end, bm.message as message, bs.message_id as message_id, bs.recursive_day_in_month as day_in_month, bs.recursive_day_in_week as day_in_week, bs.broadcast_time, bbc.name as channel_name, bs.active as active
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
		err = Utilities.RowsToStruct("bb_data", rows, &bi)
		if err != nil {
			log.Errorln(err)
		}
		ret = append(ret, bi)
	}

	return ret, nil
}
