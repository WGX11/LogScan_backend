package pkg

import (
	"context"
	sql_ "database/sql"
	"fmt"
	"log"
	"logscan/sql"
	"regexp"
	"strings"
	"sync"
	"time"
)

type Alarm struct {
	Id        int      `json:"id"`
	Input     string   `json:"input"`
	MatchType string   `json:"matchType"`
	Logs      []string `json:"logs"`
}

type Monitor struct {
	mu       sync.Mutex //防止数据竞争
	monitors map[int]context.CancelFunc
}

func NewMonitor() *Monitor {
	return &Monitor{
		monitors: make(map[int]context.CancelFunc),
	}
}

func (m *Monitor) AddMonitor(alarm Alarm) {
	id := alarm.Id
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.monitors[id]; ok {
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	m.monitors[id] = cancel
	//订阅日志消息
	logChan := make(chan string)
	LogProducerInstance.AddConsumer(id, logChan)
	//检查匹配模式
	mode := alarm.MatchType
	toCheckLogs := alarm.Logs
	if mode == "accurate" {
		go m.AccurateMonitorRoutine(ctx, id, logChan, toCheckLogs)
	} else if mode == "approximate" {
		go m.ApproximateMonitorRoutine(ctx, id, logChan, toCheckLogs)
	}
}

func (m *Monitor) RemoveMonitor(id int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if cancel, ok := m.monitors[id]; ok {
		cancel()
		//关闭信道
		LogProducerInstance.DeleteConsumer(id)
		delete(m.monitors, id)
	}
}

// 创建 KMP 部分匹配表
func createKMPTable(patterns []string) []int {
	table := make([]int, len(patterns))
	j := 0
	for i := 1; i < len(patterns); {
		if logCompare(patterns[i], patterns[j]) {
			j++
			table[i] = j
			i++
		} else if j > 0 {
			j = table[j-1]
		} else {
			table[i] = 0
			i++
		}
	}
	return table
}

// 使用 KMP 算法匹配日志序列
func KMPMatch(log string, patterns []string, table []int, state *int) bool {
	// 检查当前日志是否与当前状态指向的模式匹配
	if logCompare(log, patterns[*state]) {
		*state++

		if *state == len(patterns) {
			*state = 0
			return true
		}
	} else {

		if *state > 0 {
			*state = table[*state-1]                     // 使用部分匹配表调整状态
			return KMPMatch(log, patterns, table, state) // 递归检查新状态
		}
	}
	return false
}

// 精确匹配，即严格按照出现的顺序进行匹配 例如：1234只能匹配1234，不能匹配15234，利用kmp加速
func (m *Monitor) AccurateMonitorRoutine(ctx context.Context, id int, logChan chan string, toCheckLogs []string) {
	table := createKMPTable(toCheckLogs) // 创建 KMP 表
	state := 0                           // 初始化 KMP 状态

	for {
		select {
		case <-ctx.Done():
			return
		case logMessage := <-logChan:
			if KMPMatch(logMessage, toCheckLogs, table, &state) {
				err := sql.AddAlarmItem(id)
				if err != nil {
					log.Println("Failed to add alarm item:", err)
					return
				}
			}
		}
	}
}

// 模糊匹配，即只要包含了指定的字符串就可以匹配 例如：1234可以匹配1234，也可以匹配15234
func (m *Monitor) ApproximateMonitorRoutine(ctx context.Context, id int, logChan chan string, toCheckLogs []string) {
	begin := 0
	end := len(toCheckLogs)
	for {
		select {
		case <-ctx.Done():
			fmt.Println("monitor routine stop")
			return
		case logMessage := <-logChan:
			if logCompare(logMessage, toCheckLogs[begin]) {
				begin++
			}
			if begin == end {
				begin = 0
				err := sql.AddAlarmItem(id)
				if err != nil {
					log.Println("Failed to add alarm item:", err)
					return
				}
			}
		}
	}
}

// 查询数据库，启动报警监控
func StartAlarmMonitor() {
	db := sql.Db
	var alarmMonitor = make(map[int]Alarm, 0)
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
	defer row2.Close()
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

// 去除字母符号
func logPreprocess(log string) string {
	hexRegex := regexp.MustCompile(`(?i)0x[a-f0-9]+`)
	log = hexRegex.ReplaceAllString(log, " ")

	nonAlphaRegex := regexp.MustCompile(`[^a-zA-Z\s]+`)
	log = nonAlphaRegex.ReplaceAllString(log, " ")

	spaceRegex := regexp.MustCompile(`\s+`)
	log = spaceRegex.ReplaceAllString(log, " ")

	log = strings.ToLower(log)
	log = strings.TrimSpace(log)

	return log
}

// 计算jaccard相似度，越大越相似
func jaccardDistance(log1, log2 string) float64 {
	words1 := strings.Fields(log1)
	words2 := strings.Fields(log2)
	//使用struct{}的原因时空结构体不占内存
	set1 := make(map[string]struct{})
	set2 := make(map[string]struct{})
	for _, word := range words1 {
		if _, ok := set1[word]; !ok {
			set1[word] = struct{}{}
		}
	}
	for _, word := range words2 {
		if _, ok := set2[word]; !ok {
			set2[word] = struct{}{}
		}
	}
	intersection := 0
	for word := range set1 {
		if _, ok := set2[word]; ok {
			intersection++
		}
	}
	union := len(set1) + len(set2) - intersection
	return float64(intersection) / float64(union)
}

// 比较两个字符串的相似程度，相似的时候返回true，不相似的时候返回false
func logCompare(log1, log2 string) bool {
	log1 = logPreprocess(log1)
	log2 = logPreprocess(log2)
	if log1 == "" || log2 == "" {
		return false
	}
	return jaccardDistance(log1, log2) > 0.6

}

func TestMonitor() {
	monitor := NewMonitor()

	tmpAlarm1 := Alarm{
		Id:        15,
		Input:     "test",
		MatchType: "accurate",
		Logs:      []string{"generating core.464", "generating core.464", "generating core.464"},
	}
	//tmpAlarm2 := Alarm{
	//	Id:        2,
	//	Input:     "test",
	//	MatchType: "approximate",
	//	Logs:      []string{"test1", "test2"},
	//}
	monitor.AddMonitor(tmpAlarm1)
	//monitor.AddMonitor(tmpAlarm2)
	time.Sleep(20 * time.Second)
	monitor.RemoveMonitor(1)
	//monitor.RemoveMonitor("2")
}
