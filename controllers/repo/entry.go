package repo

type RepoGroup struct {
	NoticeRepo
	DutyScheduleRepo
}

var RepoGroupEntry = new(RepoGroup)

var (
	DBCli = NewInterGormDBCli()
)
