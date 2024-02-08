package repo

import (
	"fmt"
	"watchAlert/globals"
)

type GormDBCli struct{}

type InterGormDBCli interface {
	Create(table, value interface{}) error
	Update(value Update) error
	Updates(value Updates) error
	Delete(value Delete) error
}

func NewInterGormDBCli() InterGormDBCli {
	return &GormDBCli{}
}

func (GormDBCli *GormDBCli) Create(table, value interface{}) error {

	tx := globals.DBCli.Begin()
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
	Where  []string
	Update []string
}

func (GormDBCli *GormDBCli) Update(value Update) error {

	//updateData := Update{
	//	Table: dao.People{},
	//	Where: []string{
	//		"Name =?", "xxx", "xxx",
	//	},
	//	Update: []string{
	//		"Description", "xxx",
	//	},
	//}

	tx := globals.DBCli.Begin()
	err := tx.Model(value.Table).
		Where(value.Where[0], value.Where[1:]).
		Update(value.Update[0], value.Update[1:]).Error
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
	Where   []string
	Updates interface{}
}

func (GormDBCli *GormDBCli) Updates(value Updates) error {

	tx := globals.DBCli.Begin()
	err := tx.Model(value.Table).
		Where(value.Where[0], value.Where[1:]).
		Updates(value.Updates).Error
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
	Where []string
}

func (GormDBCli *GormDBCli) Delete(value Delete) error {

	tx := globals.DBCli.Begin()
	err := tx.Where(value.Where[0], value.Where[1:]).
		Delete(&value.Table).Error
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
