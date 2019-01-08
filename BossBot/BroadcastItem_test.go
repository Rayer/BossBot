package BossBot

import (
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"testing"
	"time"
)

func TestBroadcastItem_String(t *testing.T) {
	testItem := BroadcastItem{
		1,
		mysql.NullTime{Time: time.Now()},
		mysql.NullTime{Time: time.Now(), Valid: true},
		"This is a testing message!",
		5,
		sql.NullInt64{Int64: 4, Valid: true},
		sql.NullInt64{Int64: 2, Valid: true},
		"13:00:00",
		"Channel9",
		1,
	}

	t.Log(testItem.String())
	t.Log(testItem.StringForSlackItem())
}
