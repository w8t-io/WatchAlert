package repo

import (
	"fmt"
	"gorm.io/gorm"
)

type GormDBCli struct {
	db *gorm.DB
}

type InterGormDBCli interface {
	Create(table, value interface{}) error
	Update(value Update) error
	Updates(value Updates) error
	Delete(value Delete) error
}

func NewInterGormDBCli(db *gorm.DB) InterGormDBCli {
	return &GormDBCli{
		db,
	}
}

func (g GormDBCli) Create(table, value interface{}) error {

	tx := g.db.Begin()
	err := tx.Model(table).Create(value).Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("数据写入失败 -> %s", err)
	}
	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("事务提交失败 -> %s", err)
	}

	return nil

}

type Update struct {
	Table  interface{}
	Where  []interface{}
	Update []string
}

func (g GormDBCli) Update(value Update) error {

	tx := g.db.Begin()
	tx = tx.Model(value.Table)
	for column, val := range value.Where {
		tx = tx.Where(column, val)
	}
	err := tx.Update(value.Update[0], value.Update[1:]).Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("数据更新失败 -> %s", err)
	}
	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("事务提交失败 -> %s", err)
	}

	return nil

}

type Updates struct {
	Table   interface{}
	Where   map[string]interface{}
	Updates interface{}
}

func (g GormDBCli) Updates(value Updates) error {

	tx := g.db.Begin()
	tx = tx.Model(value.Table)
	for column, val := range value.Where {
		tx = tx.Where(column, val)
	}
	err := tx.Updates(value.Updates).Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("数据更新失败 -> %s", err)
	}

	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("事务提交失败 -> %s", err)
	}

	return nil

}

type Delete struct {
	Table interface{}
	Where map[string]interface{}
}

func (g GormDBCli) Delete(value Delete) error {

	tx := g.db.Begin()

	tx = tx.Model(value.Table)
	for column, val := range value.Where {
		tx = tx.Where(column, val)
	}
	err := tx.Delete(value.Table).Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("数据删除失败 -> %s", err)
	}
	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("事务提交失败 -> %s", err)
	}

	return nil

}
