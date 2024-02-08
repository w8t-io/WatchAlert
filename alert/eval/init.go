package eval

import "watchAlert/alert/queue"

func Initialize() {

	//curAlertsEventCache := queue.NewCurAlertsEventMap()

	NewInterEvalConsumeWork().Run()
	NewInterAlertRuleWork(queue.AlertRuleChannel, queue.QuitAlertRuleChannel).Run()

}
