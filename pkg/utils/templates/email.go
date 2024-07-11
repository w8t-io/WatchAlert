package templates

import "watchAlert/internal/models"

func emailTemplate(alert models.AlertCurEvent, notice models.AlertNotice) string {
	return ParserTemplate("Event", alert, notice.Template)
}
