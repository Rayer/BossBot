package BossBot

import (
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
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
