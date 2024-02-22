package services

import (
	"fmt"
	"watchAlert/controllers/repo"
	"watchAlert/globals"
	"watchAlert/models"
	"watchAlert/utils/cmd"
)

type RuleGroupService struct{}

type InterRuleGroupService interface {
	Create(group models.RuleGroups) error
	Update(group models.RuleGroups) error
	Delete(id string) error
	List() []models.RuleGroups
}

func NewInterRuleGroupService() InterRuleGroupService {
	return &RuleGroupService{}
}

func (rgs *RuleGroupService) Create(group models.RuleGroups) error {

	var resGroup models.RuleGroups
	globals.DBCli.Model(&models.RuleGroups{}).Where("name = ?", group.Name).First(&resGroup)
	if resGroup.Name != "" {
		return fmt.Errorf("规则组名称已存在")
	}

	group.ID = "rg-" + cmd.RandId()
	err := repo.DBCli.Create(models.RuleGroups{}, group)
	if err != nil {
		return err
	}

	return nil

}

func (rgs *RuleGroupService) Update(group models.RuleGroups) error {

	err := repo.DBCli.Updates(
		repo.Updates{
			Table:   &models.RuleGroups{},
			Where:   []string{"id = ?", group.ID},
			Updates: group,
		})
	if err != nil {
		return err
	}

	return nil

}

func (rgs *RuleGroupService) Delete(id string) error {

	err := repo.DBCli.Delete(repo.Delete{
		Table: &models.RuleGroups{},
		Where: []string{"id = ?", id},
	})
	if err != nil {
		return err
	}

	return nil

}

func (rgs *RuleGroupService) List() []models.RuleGroups {

	var resGroup []models.RuleGroups

	globals.DBCli.Model(&models.RuleGroups{}).Find(&resGroup)
	for k, v := range resGroup {
		var resRules []models.AlertRule
		globals.DBCli.Model(&models.AlertRule{}).Where("rule_group_id = ?", v.ID).Find(&resRules)
		resGroup[k].Number = len(resRules)
	}
	return resGroup

}

func (rgs *RuleGroupService) Search() {

}
