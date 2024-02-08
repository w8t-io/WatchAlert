package notice

import (
	"fmt"
	"time"
	"watchAlert/controllers/repo"
	"watchAlert/models"
	"watchAlert/utils/templates"
)

type Prometheus struct {
	CardContentMsg string
	repo.DutyScheduleRepo
	f templates.FeiShu
}

func (p *Prometheus) NewTemplate(alert models.AlertCurEvent, notice models.AlertNotice) Prometheus {

	user := p.GetDutyUserInfo(notice.DutyId, time.Now().Format("2006-1-2"))

	switch notice.NoticeType {
	case "FeiShu":
		// 判断是否有安排值班人员
		if len(user.FeiShuUserId) > 1 {
			alert.DutyUser = fmt.Sprintf("<at id=%s></at>", user.FeiShuUserId)
		}
		return Prometheus{CardContentMsg: p.f.Template(alert, notice)}
	}

	return Prometheus{}

}
