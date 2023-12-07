package initialize

import (
	"fmt"
	lark "github.com/larksuite/oapi-sdk-go/v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"prometheus-manager/globals"
	"prometheus-manager/models/dao"
	"prometheus-manager/pkg/cache"
)

func InitClient() {

	feiShuClient()
	cacheClient()
	sqlClient()

}

func feiShuClient() {

	globals.FeiShuCli = lark.NewClient(globals.Config.FeiShu.AppID, globals.Config.FeiShu.AppSecret, lark.WithEnableTokenCache(true))

}

func cacheClient() {

	globals.CacheCli = cache.NewMemoryCache()

}

func sqlClient() {

	// 初始化本地 test.db 数据库文件
	//db, err := gorm.Open(sqlite.Open("data/sql.db"), &gorm.Config{})

	sql := globals.Config.MySQL
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local&timeout=%s", sql.User, sql.Pass, sql.Host, sql.Port, sql.DBName, sql.Timeout)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		fmt.Println("failed to connect database")
	}

	// 检查 Product 结构是否变化，变化则进行迁移
	err = db.AutoMigrate(
		&dao.RuleGroupData{},
		&dao.People{},
		&dao.PeopleGroup{},
		&dao.JoinsPeopleGroup{},
		&dao.DutySystem{},
	)
	if err != nil {
		return
	}

	globals.DBCli = db

}
