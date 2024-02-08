package models

type Member struct {
	UserId       string `json:"userid"`
	UserName     string `json:"username"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	Password     string `json:"password"`
	Role         string `json:"role"`
	CreateBy     string `json:"create_by"`
	CreateAt     int64  `json:"create_at"`
	JoinDuty     string `json:"joinDuty" `
	FeiShuUserId string `json:"fsUserId"`
}
