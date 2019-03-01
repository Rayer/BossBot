package BossBot

import (
	"Utilities"
	"github.com/pkg/errors"
)

type WeeklyReportItem struct {
	Id          int    `bb_data:"id"`
	Year        int    `bb_data:"year"`
	WeekOfYear  int    `bb_data:"week_of_year"`
	UserSlackId string `bb_data:"user_id"`
	Done        string `bb_data:"done"`
	OnGoing     string `bb_data:"ongoing"`
}

func GetWeeklyReports(year int, weekOfYear int) ([]WeeklyReportItem, error) {
	//year, week := time.Now().ISOWeek()
	db := GetConfiguration().ServiceContext.DBObject.GetDB()
	res, err := db.Query("select * from bb_weekly_report where week_of_year = ? and year = ?", weekOfYear, year)
	if err != nil {
		return nil, errors.Wrap(err, "Error fetching database!")
	}
	var reports []WeeklyReportItem
	for res.Next() {
		wri := WeeklyReportItem{}
		err := Utilities.RowsToStruct("bb_data", res, &wri)
		if err != nil {
			return nil, errors.Wrap(err, "Error marshaling WeeklyReportItem!")
		}
		reports = append(reports, wri)
	}

	return reports, nil
}
