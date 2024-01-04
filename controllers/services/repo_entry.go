package services

import "watchAlert/controllers/repo"

var (
	alertNotice  = repo.RepoGroupEntry.AlertNoticeRepo
	dutySchedule = repo.RepoGroupEntry.DutyScheduleRepo
)
