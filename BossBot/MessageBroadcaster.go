package BossBot

import (
	"Utilities"
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func Processing() error {

	db, err := Utilities.CreateDBObject("node.rayer.idv.tw", "acc", "12qw34er")
	if err != nil {
		return errors.Wrap(err, "Error creating db object")
	}

	log.Println("Start process handling routine....")
	conn := db.GetConnection()
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Fatalf("Error closing connection!")
		}
	}()

	//This query includes : 1. active = 1 2. At least start or end date is assigned 3. date < end 4. date > start
	queryString := `select bs.id, bs.start_date as start, bs.end_date as end, bm.message as message, bs.message_id as message_id, CONVERT_TZ(bs.last_run, '+00:00', '+08:00') as last_run, bs.recursive_day_in_month as day_in_month, bs.recursive_day_in_week as day_in_week, bs.broadcast_time, bbc.webhook
from bb_broadcast_schedule as bs
inner join bb_broadcast_msg as bm on bs.message_id = bm.id
inner join bb_broadcast_channel as bbc on bs.channel_id = bbc.id
where bs.active = 1 and (bs.last_run is NULL or bs.last_run < CONVERT_TZ(NOW(), '+00:00', '+08:00'))
and (bs.start_date is null or bs.start_date < CONVERT_TZ(NOW(), '+00:00', '+08:00'))
and (bs.end_date is null or bs.end_date > CONVERT_TZ(NOW(), '+00:00', '+08:00'));`

	result, err := Utilities.QueryToMap(conn, queryString)
	if err != nil {
		return errors.Wrap(err, "Error in QueryToMap")
	}
	log.Debugf("Evaluating %d items...", len(result))
	for _, entry := range result {
		//output := "Evaluating : "
		//for key, value := range entry {
		//	//fmt.Printf("%s : %s\n", key, value)
		//	output += fmt.Sprintf("%s : %s ", key, value)
		//}
		//log.Println(output)

		_, err = tryBroadcast(entry, conn)
		if err != nil {
			return errors.Wrap(err, "Broadcast fail!")
		}
	}

	return nil
}

func getContextAsString(result Utilities.RowResult, column string) string {
	return string(result[column].([]byte))
}

//return code : 0 = not broadcast, 1 = broadcast
func tryBroadcast(broadcastItem Utilities.RowResult, conn *sql.DB) (int, error) {

	loc, err := time.LoadLocation("Asia/Taipei")
	if err != nil {
		return 0, errors.Wrap(err, "Fail to load time.LoadLocation")
	}
	currentTime := time.Now().In(loc)
	//log.Println("Using local time : " + currentTime.String())
	y, m, d := currentTime.Date()
	broadcast := fmt.Sprintf("%d %s %d %s +0800 CST", y, m.String()[0:3], d, string(broadcastItem["broadcast_time"].([]byte)))
	broadcast_time, err := time.Parse("2006 Jan 02 15:04:05 -0700 MST", broadcast)
	if err != nil {
		return 0, errors.Wrap(err, "Error parsing time")
	}
	//log.Println("Broadcast time : " + broadcast_time.String())
	if broadcast_time.After(currentTime) {
		return 0, nil
	}

	//Basic concept is :
	//If it is executed in 24 hours, skip. It means today's broadcast have been done
	if broadcastItem["last_run"] != nil {
		last := fmt.Sprintf(string(broadcastItem["last_run"].([]byte)))
		lastTime, err := time.Parse("2006-01-02 15:04:05", last)
		if err != nil {
			return 0, errors.Wrap(err, "Error parsing last_run!")
		}

		if currentTime.Sub(lastTime).Hours() < 24 {
			return 0, nil
		}
		println(currentTime.Sub(lastTime).Hours())
	}
	//check day of month and day of week
	if broadcastItem["day_in_week"] != nil {

		res, err := strconv.Atoi(getContextAsString(broadcastItem, "day_in_week"))
		if err != nil {
			return 0, errors.Wrap(err, "Error parsing day_in_week!")
		}

		if res != int(currentTime.Weekday()) {
			return 0, nil
		}
	}

	if broadcastItem["day_in_month"] != nil {

		res, err := strconv.Atoi(getContextAsString(broadcastItem, "day_in_month"))
		if err != nil {
			return 0, errors.Wrap(err, "Error parsing day_in_month!")
		}

		if res != currentTime.Day() {
			return 0, nil
		}
	}

	//do broadcast
	log.Println("Posting " + string(broadcastItem["message"].([]byte)) + " to " + string(broadcastItem["webhook"].([]byte)))
	outgoing := "{\"text\":\"" + string(broadcastItem["message"].([]byte)) + "\"}"
	response, err := http.Post(getContextAsString(broadcastItem, "webhook"), "application/json", strings.NewReader(outgoing))
	if err != nil {
		return 0, errors.Wrap(err, "Error posting to channel!")
	}
	defer func() {
		err := response.Body.Close()
		if err != nil {
			log.Errorf("Error while closing response!")
		}
	}()
	//fmt.Println(ioutil.ReadAll(response.Body))

	//update last run in both schedule and msg
	id := getContextAsString(broadcastItem, "id")
	last_run_sql := fmt.Sprintf("update bb_broadcast_schedule set last_run = CONVERT_TZ(Now(), '+00:00', '+08:00') where id = %s;", id)
	_, err = conn.Exec(last_run_sql)
	if err != nil {
		return 1, errors.Wrap(err, "Message is sent but fail to update last_run in bb_broadcast_schedule!")
	}

	id = getContextAsString(broadcastItem, "message_id")
	last_run_sql = fmt.Sprintf("update bb_broadcast_msg set last_broadcast = CONVERT_TZ(Now(), '+00:00', '+08:00'), broadcast_count = if(broadcast_count is null, 1, broadcast_count + 1) where id = %s;", id)
	_, err = conn.Exec(last_run_sql)
	if err != nil {
		return 1, errors.Wrap(err, "Message is sent but fail to update last_run in bb_broadcast_msg!")
	}
	log.Println("Successfully posted " + getContextAsString(broadcastItem, "message") + " to " + getContextAsString(broadcastItem, "webhook"))
	return 1, nil
}

func StartBroadcaster() {
	log.Println("Initialization for Broadcaster...")
	for {
		ticker := time.NewTicker(time.Minute)
		select {
		case <-ticker.C:
			err := Processing()
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
