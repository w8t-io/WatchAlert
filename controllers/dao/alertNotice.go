package dao

type AlertNotice struct {
	Uuid         string `json:"uuid"`
	Name         string `json:"name"`
	Env          string `json:"env"`
	NoticeStatus string `json:"noticeStatus"`
	DutyId       string `json:"dutyId"`
	DataSource   string `json:"dataSource"`
	NoticeType   string `json:"noticeType"`
	FeishuChatId string `json:"feishuChatId"`
}
