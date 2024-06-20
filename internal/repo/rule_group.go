package repo

import (
	"fmt"
	"gorm.io/gorm"
	"watchAlert/internal/models"
	"watchAlert/pkg/utils/cmd"
)

type (
	RuleGroupRepo struct {
		entryRepo
	}

	InterRuleGroupRepo interface {
		List(req models.RuleGroupQuery) (models.RuleGroupResponse, error)
		Create(req models.RuleGroups) error
		Update(req models.RuleGroups) error
		Delete(req models.RuleGroupQuery) error
	}
)

func newRuleGroupInterface(db *gorm.DB, g InterGormDBCli) InterRuleGroupRepo {
	return &RuleGroupRepo{
		entryRepo{
			g:  g,
			db: db,
		},
	}
}

func (r RuleGroupRepo) List(req models.RuleGroupQuery) (models.RuleGroupResponse, error) {
	var (
		data  []models.RuleGroups
		db    = r.db.Model(&models.RuleGroups{})
		count int64
	)

	pageIndexInt := req.Page.Index
	pageSizeInt := req.Page.Size

	db.Where("tenant_id = ?", req.TenantId)

	if req.Query != "" {
		db.Where("id LIKE ? OR name LIKE ? OR description LIKE ?",
			"%"+req.Query+"%", "%"+req.Query+"%", "%"+req.Query+"%")
	}

	db.Count(&count)

	db.Limit(int(pageSizeInt)).Offset(int((pageIndexInt - 1) * pageSizeInt))

	err := db.Find(&data).Error
	if err != nil {
		return models.RuleGroupResponse{}, err
	}

	for k, v := range data {
		var resRules []models.AlertRule
		r.db.Model(&models.AlertRule{}).Where("tenant_id = ? AND rule_group_id = ?", req.TenantId, v.ID).Find(&resRules)
		data[k].Number = len(resRules)
	}
	return models.RuleGroupResponse{
		List: data,
		Page: models.Page{
			Index: pageIndexInt,
			Size:  pageSizeInt,
			Total: count,
		},
	}, nil
}

func (r RuleGroupRepo) Create(req models.RuleGroups) error {
	var resGroup models.RuleGroups
	r.db.Model(&models.RuleGroups{}).Where("name = ?", req.Name).First(&resGroup)
	if resGroup.Name != "" {
		return fmt.Errorf("规则组名称已存在")
	}

	nr := req
	nr.ID = "rg-" + cmd.RandId()
	err := r.g.Create(models.RuleGroups{}, nr)
	if err != nil {
		return err
	}

	return nil
}

func (r RuleGroupRepo) Update(req models.RuleGroups) error {
	u := Updates{
		Table: &models.RuleGroups{},
		Where: map[string]interface{}{
			"tenant_id = ?": req.TenantId,
			"id = ?":        req.ID,
		},
		Updates: req,
	}

	err := r.g.Updates(u)
	if err != nil {
		return err
	}

	return nil
}

func (r RuleGroupRepo) Delete(req models.RuleGroupQuery) error {
	var ruleNum int64
	r.db.Model(&models.AlertRule{}).Where("tenant_id = ? AND rule_group_id = ?", req.TenantId, req.ID).
		Count(&ruleNum)
	if ruleNum != 0 {
		return fmt.Errorf("无法删除规则组 %s, 因为规则组不为空", req.ID)
	}

	d := Delete{
		Table: models.RuleGroups{},
		Where: map[string]interface{}{
			"tenant_id = ?": req.TenantId,
			"id = ?":        req.ID,
		},
	}

	err := r.g.Delete(d)
	if err != nil {
		return err
	}

	return nil
}
