package initialize

import (
	"fmt"
	"github.com/go-redis/redis"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"watchAlert/globals"
	models2 "watchAlert/models"
)

func InitClient() {

	sqlClient()
	redisClient()
	InitPermissionsSQL()
	InitUserRolesSQL()

}

func redisClient() {

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", globals.Config.Redis.Host, globals.Config.Redis.Port),
		Password: globals.Config.Redis.Pass,
		DB:       0, // 使用默认的数据库
	})

	// 尝试连接到 Redis 服务器
	_, err := client.Ping().Result()
	if err != nil {
		log.Printf("redis Connection Failed %s", err)
		return
	}

	globals.RedisCli = client

}

func sqlClient() {

	// 初始化本地 test.db 数据库文件
	//db, err := gorm.Open(sqlite.Open("data/sql.db"), &gorm.Config{})

	sql := globals.Config.MySQL
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4,utf8&parseTime=True&loc=Local&timeout=%s", sql.User, sql.Pass, sql.Host, sql.Port, sql.DBName, sql.Timeout)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		fmt.Println("failed to connect database")
	}

	// 检查 Product 结构是否变化，变化则进行迁移
	err = db.AutoMigrate(
		&models2.DutySchedule{},
		&models2.DutyManagement{},
		&models2.AlertNotice{},
		&models2.AlertDataSource{},
		&models2.AlertRule{},
		&models2.AlertCurEvent{},
		&models2.AlertHisEvent{},
		&models2.AlertSilences{},
		&models2.Member{},
		&models2.UserRole{},
		&models2.UserPermissions{},
		&models2.NoticeTemplateExample{},
		&models2.RuleGroups{},
		&models2.RuleTemplateGroup{},
		&models2.RuleTemplate{},
	)
	if err != nil {
		return
	}

	globals.DBCli = db

}
