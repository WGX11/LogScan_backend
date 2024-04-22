package sql

import (
	"fmt"
	"strings"
)

type AlarmInfo struct {
	AlarmName        string `json:"alarm_name"`
	AlarmId          int    `json:"alarm_id"`
	AlarmDescription string `json:"alarm_description"`
}

type DashBoard struct {
	AlarmName string `json:"alarm_name"`
	AlarmId   int    `json:"alarm_id"`
	AlarmNum  int    `json:"alarm_num"`
}

func AddAlarm(userId int, alarmName, alarmInput, alarmMatchMode,
	alarmDescription, alarmEmail, alarmPhoneNumber string,
	alarmEmailNotification, alarmPhoneNotification bool, logs []string) error {
	db := Db
	query := `INSERT INTO alarm_group
		(user_id, alarm_name, alarm_input, alarm_match_mode, alarm_description, 
		alarm_email, alarm_phone, sms_notification, 
		email_notification)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	result, err := db.Exec(query, userId, alarmName, alarmInput, alarmMatchMode,
		alarmDescription, alarmEmail, alarmPhoneNumber, alarmEmailNotification,
		alarmPhoneNotification)
	if err != nil {
		fmt.Println("Failed to insert alarm group:", err)
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		fmt.Println("Failed to get last insert id:", err)
	}
	fmt.Println(id)
	valuesStrings := make([]string, 0, len(logs))
	valueArgs := make([]interface{}, 0, len(logs)*2)
	for _, log := range logs {
		valuesStrings = append(valuesStrings, "(?, ?)")
		valueArgs = append(valueArgs, id, log)
	}
	stmt := fmt.Sprintf("INSERT INTO alarm_logs (alarm_id, log_content) VALUES %s", strings.Join(valuesStrings, ","))

	_, err = db.Exec(stmt, valueArgs...)
	if err != nil {
		fmt.Println("Failed to insert alarm logs:", err)
		return err
	}
	return nil
}

func GetAlarmList(userId int) []AlarmInfo {
	db := Db
	query := `SELECT alarm_id, alarm_name, alarm_description FROM alarm_group WHERE user_id = ?`
	rows, err := db.Query(query, userId)
	if err != nil {
		fmt.Println("Failed to get alarm list:", err)
		return nil
	}
	defer rows.Close()
	alarmList := make([]AlarmInfo, 0)
	for rows.Next() {
		var alarm AlarmInfo
		err := rows.Scan(&alarm.AlarmId, &alarm.AlarmName, &alarm.AlarmDescription)
		if err != nil {
			fmt.Println("Failed to scan alarm info:", err)
			return nil
		}
		alarmList = append(alarmList, alarm)
	}
	return alarmList
}

func DeleteAlarm(alarmId string) error {
	db := Db
	query := `DELETE FROM alarm_group WHERE alarm_id = ?`
	_, err := db.Exec(query, alarmId)
	if err != nil {
		fmt.Println("Failed to delete alarm:", err)
		return err
	}
	return nil
}

func GetAlarmData15Days(userId int) []DashBoard {
	db := Db
	query := `
	SELECT ag.alarm_id, ag.alarm_name, COUNT(ai.alarm_item_id) AS count
	FROM alarm_group ag
	LEFT JOIN alarm_item ai ON ag.alarm_id = ai.alarm_id AND ai.alarm_time > DATE_SUB(NOW(), INTERVAL 15 DAY)
	WHERE ag.user_id = ? 
	GROUP BY ag.alarm_id
	`
	rows, err := db.Query(query, userId)
	if err != nil {
		fmt.Println("Failed to get alarm list:", err)
		return nil
	}
	defer rows.Close()
	dashBoards := make([]DashBoard, 0)
	for rows.Next() {
		var dashBoard DashBoard
		err := rows.Scan(&dashBoard.AlarmId, &dashBoard.AlarmName, &dashBoard.AlarmNum)
		if err != nil {
			fmt.Println("Failed to scan alarm info:", err)
			return nil
		}
		dashBoards = append(dashBoards, dashBoard)
	}
	return dashBoards
}
