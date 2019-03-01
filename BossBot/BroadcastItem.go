package BossBot

import (
	"Utilities"
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"time"
)

type BroadcastItem struct {
	Id            int            `bb_data:"id" json:"id"`
	StartDate     mysql.NullTime `bb_data:"start" json:"start"`
	EndDate       mysql.NullTime `bb_data:"end" json:"end"`
	Message       string         `bb_data:"message" json:"message"`
	MessageId     int            `bb_data:"message_id" json:"message_id"`
	Day           sql.NullInt64  `bb_data:"day_in_month" json:"day_in_month"`
	WeekDay       sql.NullInt64  `bb_data:"day_in_week" json:"day_in_week"`
	BroadcastTime string         `bb_data:"broadcast_time" json:"broadcast_time"`
	ChannelName   string         `bb_data:"channel_name" json:"channel_name"`
	Active        int            `bb_data:"active" json:"is_active"`
}

func (bc *BroadcastItem) String() string {
	var ret string
	var activeSymbol string

	if bc.Active == 1 {
		activeSymbol = "o"
	} else {
		activeSymbol = "x"
	}
	ret += fmt.Sprintf("(%s) - Message : %s(%d) ", activeSymbol, bc.Message, bc.MessageId)
	if bc.StartDate.Valid {
		ret += fmt.Sprintf("from %s ", bc.StartDate.Time)
	}

	if bc.EndDate.Valid {
		ret += fmt.Sprintf("until %s ", bc.EndDate.Time)
	}

	ret += fmt.Sprintf("at %s ", bc.BroadcastTime)

	if bc.WeekDay.Valid {
		ret += fmt.Sprintf("at every %s ", time.Weekday(bc.Day.Int64))
	}

	if bc.Day.Valid {
		ret += fmt.Sprintf("at every %d day of the month ", bc.Day.Int64)
	}

	ret += fmt.Sprintf("to channel : %s.", bc.ChannelName)

	return ret
}

func (bc *BroadcastItem) StringForSlackItem() string {
	var ret string

	ret += fmt.Sprintf("Scheduler ID : %d\n", bc.Id)
	ret += fmt.Sprintf("Message : %s\n", bc.Message)
	ret += fmt.Sprintf("Message ID : %d\n", bc.MessageId)

	if bc.StartDate.Valid {
		ret += fmt.Sprintf("From %s\n", bc.StartDate.Time.Format(time.RFC822))
	}

	if bc.EndDate.Valid {
		ret += fmt.Sprintf("Until %s\n", bc.EndDate.Time.Format(time.RFC822))
	}

	ret += fmt.Sprintf("Broadcast at : %s\n", bc.BroadcastTime)

	if bc.WeekDay.Valid {
		ret += fmt.Sprintf("Recursive in weekday : %s\n", time.Weekday(bc.WeekDay.Int64))
	}

	if bc.Day.Valid {
		ret += fmt.Sprintf("Recursive Day in Month : %d\n", bc.Day.Int64)
	}

	ret += fmt.Sprintf("Target channel : %s\n", bc.ChannelName)

	return ret
}

func GetBroadcastList() ([]BroadcastItem, error) {
	db := GetConfiguration().ServiceContext.DBObject.GetDB()
	rows, err := db.Query(
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
