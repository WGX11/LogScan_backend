package pkg

import (
	sql_ "database/sql"
	"log"
	"logscan/sql"
)

type Alarm struct {
	Id        int      `json:"id"`
	Input     string   `json:"input"`
	MatchType string   `json:"matchType"`
	Logs      []string `json:"logs"`
}

var alarmMonitor = make(map[int]Alarm, 0)

// 查询数据库，启动报警监控
func StartAlarmMonitor() {
	db := sql.Db
	query1 := `
	SELECT alarm_id,  alarm_input, alarm_match_mode FROM alarm_group
	`
	query2 := `
	SELECT ag.alarm_id, al.log_content FROM
	alarm_group ag
	LEFT JOIN alarm_logs al 
	ON ag.alarm_id = al.alarm_id
	`
	rows1, err := db.Query(query1)
	if err != nil {
		log.Println("Failed to start alarm monitor:", err)
		return
	}
	defer rows1.Close()
	for rows1.Next() {
		var alarm Alarm
		err := rows1.Scan(&alarm.Id, &alarm.Input, &alarm.MatchType)
		if err != nil {
			log.Println("Failed to scan alarm info:", err)
			return
		}
		alarmMonitor[alarm.Id] = alarm
	}
	row2, err := db.Query(query2)
	if err != nil {
		log.Println("Failed to start alarm monitor:", err)
		return
	}
	for row2.Next() {
		var (
			id         int
			logContent sql_.NullString
		)
		err := row2.Scan(&id, &logContent)
		if err != nil {
			log.Println("Failed to scan alarm log:", err)
			return
		}
		if !logContent.Valid {
			continue
		}
		alarm, ok := alarmMonitor[id]
		if !ok {
			log.Println("Failed to find alarm:", id)
			continue
		}
		alarm.Logs = append(alarm.Logs, logContent.String)
	}
}
