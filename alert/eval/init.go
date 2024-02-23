package eval

func Initialize() {

	//curAlertsEventCache := queue.NewCurAlertsEventMap()

	NewInterEvalConsumeWork().Run()
	NewInterAlertRuleWork().Run()

}
