package services

import (
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"io"
	"io/ioutil"
	"watchAlert/controllers/dao"
	"watchAlert/controllers/dto"
	"watchAlert/globals"
)

type RuleGroupService struct{}

type InterRuleGroupService interface {
	SelectRuleGroup() ([]dao.RuleGroupData, error)
	CreateRuleGroup(body io.ReadCloser) error
	UpdateRuleGroup(ruleGroupName string, body io.ReadCloser) error
	DeleteRuleGroup(ruleGroupName string) error
	GetRuleGroup(ruleGroupName string) ([]dao.RuleGroupData, error)
}

func NewInterRuleGroupService() InterRuleGroupService {
	return &RuleGroupService{}
}

func (rgs *RuleGroupService) SelectRuleGroup() ([]dao.RuleGroupData, error) {

	var (
		ruleGroup []dao.RuleGroupData
	)
	err := globals.DBCli.Find(&ruleGroup).Error
	if err != nil {
		return []dao.RuleGroupData{}, err
	}

	return ruleGroup, nil

}

func (rgs *RuleGroupService) CreateRuleGroup(body io.ReadCloser) error {

	var (
		ruleGroup dto.RuleGroup
	)
	jsonByte, _ := ioutil.ReadAll(body)
	_ = json.Unmarshal(jsonByte, &ruleGroup)

	data := NewInterRuleService().SelectPromRules(ruleGroup.Name)

	err := globals.DBCli.Create(&dao.RuleGroupData{
		Model:       gorm.Model{},
		Name:        ruleGroup.Name,
		RuleNumber:  data.Number,
		Description: ruleGroup.Description,
	}).Error

	if err != nil {
		globals.Logger.Sugar().Error("数据写入失败 ->", err)
		return err
	}

	return nil

}

func (rgs *RuleGroupService) UpdateRuleGroup(ruleGroupName string, body io.ReadCloser) error {

	var (
		ruleGroup dao.RuleGroupData
	)
	jsonByte, _ := ioutil.ReadAll(body)
	_ = json.Unmarshal(jsonByte, &ruleGroup)

	err := globals.DBCli.Model(&ruleGroup).
		Where("Name = ?", ruleGroupName).
		Update("Description", ruleGroup.Description).Error
	if err != nil {
		globals.Logger.Sugar().Error("数据更新失败 ->", err)
		return err
	}

	return nil

}

func (rgs *RuleGroupService) DeleteRuleGroup(ruleGroupName string) error {

	var (
		ruleGroup dao.RuleGroupData
	)

	rules := NewInterRuleService().SelectPromRules(ruleGroupName)
	if rules.Number > 0 {
		return fmt.Errorf("当前规则组中规则数不为零, 无法删除！")
	}

	err := globals.DBCli.Where("Name = ?", ruleGroupName).Delete(&ruleGroup).Error
	if err != nil {
		globals.Logger.Sugar().Error("数据删除失败 ->", err)
		return err
	}

	return nil

}

func (rgs *RuleGroupService) GetRuleGroup(ruleGroupName string) ([]dao.RuleGroupData, error) {

	var (
		ruleGroup []dao.RuleGroupData
	)

	err := globals.DBCli.Model(&ruleGroup).Where("Name = ?", ruleGroupName).Find(&ruleGroup).Error

	if err != nil {
		globals.Logger.Sugar().Error("数据查询失败 ->", err)
		return []dao.RuleGroupData{}, err
	}

	return ruleGroup, nil

}
