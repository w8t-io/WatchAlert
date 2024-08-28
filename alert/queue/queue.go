package queue

var (
	// RecoverWaitMap 存储等待被恢复的告警的 Key
	RecoverWaitMap = make(map[string]int64)
)
