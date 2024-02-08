package repo

type RepoGroup struct {
	AlertNoticeRepo
	DutyScheduleRepo
}

var RepoGroupEntry = new(RepoGroup)

var (
	DBCli = NewInterGormDBCli()
)
