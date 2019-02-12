package BossBot

type WeeklyReportItem struct {
	Id          int    `bb_data:"id"`
	Year        int    `bb_data:"year"`
	WeekOfYear  int    `bb_data:"week_of_year"`
	UserSlackId string `bb_data:"user_slack_id"`
	Done        string `bb_data:"done"`
	OnGoing     string `bb_data:"ongoing"`
}
