package initialize

import (
	"runtime"
	"time"
	"watchAlert/controllers/repo"
	"watchAlert/globals"
	"watchAlert/models"
)

func InitResource() {
	ticker := time.Tick(time.Second * 10)
	layout := "2006-01-02 15:04:05"
	go func() {
		for range ticker {
			curAt := time.Now()
			goNum := runtime.NumGoroutine()
			cleanupOldData(curAt)
			repo.DBCli.Create(&models.ServiceResource{}, models.ServiceResource{
				ID:    uint(curAt.Unix()),
				Time:  curAt.Format(layout),
				Value: goNum,
				Label: "goroutine",
			})
		}
	}()
}

func cleanupOldData(curAt time.Time) {
	cutoffTime := curAt.Add(-6 * time.Hour)
	globals.DBCli.Where("time < ?", cutoffTime).Delete(&models.ServiceResource{})
}
