package services

import "prometheus-manager/controllers/repo"

var (
	alertNotice  = repo.RepoGroupEntry.AlertNoticeRepo
	dutySchedule = repo.RepoGroupEntry.DutyScheduleRepo
)
