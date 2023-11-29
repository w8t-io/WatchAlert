package rules

import (
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"io"
	"io/ioutil"
	"prometheus-manager/globals"
	"prometheus-manager/models/dao"
)

func SelectRuleGroup() ([]dao.RuleGroupData, error) {

	var (
		ruleGroup []dao.RuleGroupData
	)
	err := globals.DBCli.Find(&ruleGroup).Error
	if err != nil {
		return []dao.RuleGroupData{}, err
	}

	return ruleGroup, nil

}

func CreateRuleGroup(body io.ReadCloser) error {

	var (
		ruleGroup RuleGroup
	)
	jsonByte, _ := ioutil.ReadAll(body)
	_ = json.Unmarshal(jsonByte, &ruleGroup)

	data := SelectPromRules(ruleGroup.Name)

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

func UpdateRuleGroup(ruleGroupName string, body io.ReadCloser) error {

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

func DeleteRuleGroup(ruleGroupName string) error {

	var (
		ruleGroup dao.RuleGroupData
	)

	rules := SelectPromRules(ruleGroupName)
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

func GetRuleGroup(ruleGroupName string) ([]dao.RuleGroupData, error) {

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
