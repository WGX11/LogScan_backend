package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"logscan/sql"
)

type Alarm struct {
	Name              string   `json:"name"`
	Input             string   `json:"input"`
	MatchType         string   `json:"matchType"`
	Description       string   `json:"description"`
	Email             string   `json:"email"`
	EmailNotification bool     `json:"emailNotification"`
	PhoneNumber       string   `json:"phoneNumber"`
	PhoneNotification bool     `json:"phoneNotification"`
	Logs              []string `json:"logs"`
}

func HandleAddAlarm(ctx *gin.Context) {
	var alarm Alarm
	if error := ctx.BindJSON(&alarm); error != nil {
		ctx.JSON(400, gin.H{"error": error.Error()})
		fmt.Println(error)
	}
	fmt.Println(alarm)
	err := sql.AddAlarm(1, alarm.Name, alarm.Input, alarm.MatchType, alarm.Description, alarm.Email, alarm.PhoneNumber,
		alarm.EmailNotification, alarm.PhoneNotification, alarm.Logs)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		fmt.Println("Failed to add alarm: ", err)
	}
	return
}

func HandleDeleteAlarm(ctx *gin.Context) {
	alarmId := ctx.Param("id")
	err := sql.DeleteAlarm(alarmId)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		fmt.Println("Failed to delete alarm: ", err)
	}
	return
}

func HandleGetAlarmList(ctx *gin.Context) {
	userId := 1
	alarmList := sql.GetAlarmList(userId)
	if alarmList == nil {
		ctx.JSON(500, gin.H{"error": "Failed to get alarm list"})
		return
	}
	ctx.JSON(200, alarmList)
	return
}

func HandleGetAlarmDashBoard(ctx *gin.Context) {
	userId := 1
	alarmList := sql.GetAlarmData15Days(userId)
	if alarmList == nil {
		ctx.JSON(500, gin.H{"error": "Failed to get alarm dashboard"})
		return
	}
	ctx.JSON(200, alarmList)
	return
}
