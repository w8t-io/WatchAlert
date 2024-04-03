package services

import "watchAlert/controllers/repo"

var (
	alertNotice  = repo.RepoGroupEntry.NoticeRepo
	dutySchedule = repo.RepoGroupEntry.DutyScheduleRepo
)
