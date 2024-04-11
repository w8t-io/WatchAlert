package notice

import (
	"fmt"
	"time"
	"watchAlert/controllers/repo"
	"watchAlert/models"
	templates2 "watchAlert/public/utils/templates"
)

type Template struct {
	CardContentMsg string
	repo.DutyScheduleRepo
	f templates2.FeiShu
	d templates2.DingDing
}

func (p *Template) NewTemplate(alert models.AlertCurEvent, notice models.AlertNotice) Template {

	user := p.GetDutyUserInfo(notice.DutyId, time.Now().Format("2006-1-2"))

	switch notice.NoticeType {
	case "FeiShu":
		// 判断是否有安排值班人员
		if len(user.DutyUserId) > 1 {
			alert.DutyUser = fmt.Sprintf("<at id=%s></at>", user.DutyUserId)
		}
		return Template{CardContentMsg: p.f.Template(alert, notice)}
	case "DingDing":
		if len(user.DutyUserId) > 1 {
			alert.DutyUser = fmt.Sprintf("%s", user.DutyUserId)
		}
		return Template{CardContentMsg: p.d.Template(alert, notice)}
	}

	return Template{}

}
