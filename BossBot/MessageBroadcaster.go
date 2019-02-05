package BossBot

import (
	"Utilities"
	"database/sql"
	"fmt"
	"github.com/nlopes/slack"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"strconv"
	"time"
)

type MessageBroadcaster struct {
	config Configuration
}

func (mb *MessageBroadcaster) Processing() error {

	db := mb.config.ServiceContext.DBObject

	log.Debugln("Start process handling routine....")
	conn := db.GetDB()

	//This query includes : 1. active = 1 2. At least start or end date is assigned 3. date < end 4. date > start
	queryString := `select bs.id, bs.start_date as start, bs.end_date as end, bm.message as message, 
bs.message_id as message_id, CONVERT_TZ(bs.last_run, '+00:00', '+08:00') as last_run, 
bs.recursive_day_in_month as day_in_month, bs.recursive_day_in_week as day_in_week, 
bs.broadcast_time, bbc.webhook, bm.format
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
	broadcast_time, err := time.Parse("2006 Jan _2 15:04:05 -0700 MST", broadcast)
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
		lastTime := broadcastItem["last_run"].(time.Time)
		//lastTime, err := time.Parse("2006-01-02 15:04:05", last)
		if err != nil {
			return 0, errors.Wrap(err, "Error parsing last_run!")
		}

		if currentTime.Sub(lastTime).Hours() < 24 {
			return 0, nil
		}
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
	webhookUrl := getContextAsString(broadcastItem, "webhook")
	message := getContextAsString(broadcastItem, "message")
	log.Println("Posting " + string(broadcastItem["message"].([]byte)) + " to " + string(broadcastItem["webhook"].([]byte)))

	outgoing := slack.WebhookMessage{
		Text: message,
	}
	err = slack.PostWebhook(webhookUrl, &outgoing)

	if err != nil {
		return 0, errors.Wrap(err, "Error posting to channel!")
	}

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
	log.Infoln("Successfully posted " + getContextAsString(broadcastItem, "message") + " to " + getContextAsString(broadcastItem, "webhook"))
	return 1, nil
}

func (mb *MessageBroadcaster) InvokeBroadcast(id int) (int, error) {
	db := mb.config.ServiceContext.DBObject.GetDB()

	res, err := db.Query(`select bs.id, bm.message as message, bbc.webhook, bm.id as message_id
	from bb_broadcast_schedule as bs
	inner join bb_broadcast_msg as bm on bs.message_id = bm.id
	inner join bb_broadcast_channel as bbc on bs.channel_id = bbc.id
	where bs.id = ?`, id)

	if err != nil {
		return 0, errors.Wrap(err, "Error invoke broadcast!")
	}

	type ib_item struct {
		Id        int    `bb_data:"id"`
		Message   string `bb_data:"message"`
		Webhook   string `bb_data:"webhook"`
		MessageId int    `bb_data:"message_id"`
	}

	ib := ib_item{}
	res.Next()
	err = Utilities.RowsToStruct("bb_data", res, &ib)
	if err != nil {
		return 0, errors.Wrap(err, "Error parsing return data")
	}

	log.Debugf("Get ib item : %+v\n", ib)

	webhookUrl := ib.Webhook
	message := ib.Message

	outgoing := slack.WebhookMessage{
		Text: message,
	}
	err = slack.PostWebhook(webhookUrl, &outgoing)

	if err != nil {
		return 0, errors.Wrap(err, "Error posting to channel!")
	}

	lastRunSql := fmt.Sprintf("update bb_broadcast_schedule set last_run = CONVERT_TZ(Now(), '+00:00', '+08:00') where id = %d;", id)
	_, err = db.Exec(lastRunSql)
	if err != nil {
		return 1, errors.Wrap(err, "Message is sent but fail to update last_run in bb_broadcast_schedule!")
	}

	lastRunSql = fmt.Sprintf("update bb_broadcast_msg set last_broadcast = CONVERT_TZ(Now(), '+00:00', '+08:00'), broadcast_count = if(broadcast_count is null, 1, broadcast_count + 1) where id = %d;", ib.MessageId)
	_, err = db.Exec(lastRunSql)
	if err != nil {
		return 1, errors.Wrap(err, "Message is sent but fail to update last_run in bb_broadcast_msg!")
	}

	return 1, nil
}

func (mb *MessageBroadcaster) SetActive(schId int, active bool) error {
	db := mb.config.ServiceContext.DBObject.GetDB()

	var isActive int
	if active {
		isActive = 1
	} else {
		isActive = 0
	}

	_, err := db.Exec("update bb_broadcast_schedule set active = ? where id = ?", isActive, schId)
	if err != nil {
		return errors.Wrap(err, "Error updating schedule table!")
	}
	return nil
}

func StartBroadcaster(conf Configuration) {
	log.Println("Initialization for Broadcaster...")
	for {
		ticker := time.NewTicker(time.Minute)
		select {
		case <-ticker.C:
			broadcaster := MessageBroadcaster{conf}
			err := broadcaster.Processing()
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
