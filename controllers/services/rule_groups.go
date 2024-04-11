package services

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"watchAlert/controllers/repo"
	"watchAlert/models"
	"watchAlert/public/globals"
	"watchAlert/public/utils/cmd"
)

type RuleGroupService struct{}

type InterRuleGroupService interface {
	Create(group models.RuleGroups) error
	Update(group models.RuleGroups) error
	Delete(tid, id string) error
	List(ctx *gin.Context) []models.RuleGroups
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
			Where:   []interface{}{"tenant_id = ? AND id = ?", group.TenantId, group.ID},
			Updates: group,
		})
	if err != nil {
		return err
	}

	return nil

}

func (rgs *RuleGroupService) Delete(tid, id string) error {

	var ruleNum int64
	globals.DBCli.Model(&models.AlertRule{}).Where("tenant_id = ? AND rule_group_id = ?", tid, id).Count(&ruleNum)
	if ruleNum != 0 {
		return fmt.Errorf("无法删除规则组 %s, 因为规则组不为空", id)
	}

	err := repo.DBCli.Delete(repo.Delete{
		Table: &models.RuleGroups{},
		Where: []interface{}{"tenant_id = ? AND id = ?", tid, id},
	})
	if err != nil {
		return err
	}

	return nil

}

func (rgs *RuleGroupService) List(ctx *gin.Context) []models.RuleGroups {

	var resGroup []models.RuleGroups
	db := globals.DBCli
	tid, _ := ctx.Get("TenantID")

	db.Model(&models.RuleGroups{}).Where("tenant_id = ?", tid.(string)).Find(&resGroup)
	for k, v := range resGroup {
		var resRules []models.AlertRule
		db.Model(&models.AlertRule{}).Where("tenant_id = ? AND rule_group_id = ?", tid.(string), v.ID).Find(&resRules)
		resGroup[k].Number = len(resRules)
	}
	return resGroup

}

func (rgs *RuleGroupService) Search() {

}
