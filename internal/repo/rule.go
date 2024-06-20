package repo

import (
	"gorm.io/gorm"
	"watchAlert/internal/models"
	"watchAlert/pkg/utils/cmd"
)

type (
	RuleRepo struct {
		entryRepo
	}

	InterRuleRepo interface {
		GetQuota(id string) bool
		Search(r models.AlertRuleQuery) (models.AlertRule, error)
		List(r models.AlertRuleQuery) (models.RuleResponse, error)
		Create(r models.AlertRule) error
		Update(r models.AlertRule) error
		Delete(r models.AlertRuleQuery) error
		GetRuleIsExist(ruleId string) bool
	}
)

func newRuleInterface(db *gorm.DB, g InterGormDBCli) InterRuleRepo {
	return &RuleRepo{
		entryRepo{
			g:  g,
			db: db,
		},
	}
}

func (rr RuleRepo) GetQuota(id string) bool {
	var (
		db     = rr.db.Model(&models.Tenant{})
		data   models.Tenant
		Number int64
	)

	db.Where("id = ?", id)
	db.Find(&data)

	rr.db.Model(&models.AlertRule{}).Where("tenant_id = ?", id).Count(&Number)

	if Number < data.RuleNumber {
		return true
	}

	return false
}

func (rr RuleRepo) Search(r models.AlertRuleQuery) (models.AlertRule, error) {
	var data models.AlertRule

	db := rr.db.Model(&models.AlertRule{})
	db.Where("tenant_id = ? AND rule_group_id = ? AND rule_id = ?", r.TenantId, r.RuleGroupId, r.RuleId)
	err := db.First(&data).Error
	if err != nil {
		return data, err
	}

	return data, nil
}

func (rr RuleRepo) List(r models.AlertRuleQuery) (models.RuleResponse, error) {
	var (
		data  []models.AlertRule
		count int64
	)

	db := rr.db.Model(&models.AlertRule{})
	db.Where("tenant_id = ?", r.TenantId)
	db.Where("rule_group_id = ?", r.RuleGroupId)

	if r.Query != "" {
		db.Where("rule_id LIKE ? OR rule_name LIKE ? OR description LIKE ?",
			"%"+r.Query+"%", "%"+r.Query+"%", "%"+r.Query+"%")
	}

	if r.Status != "all" {
		var s bool
		switch r.Status {
		case "enabled":
			s = true
		case "disabled":
			s = false
		}
		db.Where("enabled = ?", s)
	}

	db.Count(&count)

	db.Limit(int(r.Page.Size)).Offset(int((r.Page.Index - 1) * r.Page.Size))

	err := db.Find(&data).Error

	if err != nil {
		return models.RuleResponse{}, err
	}

	return models.RuleResponse{
		List: data,
		Page: models.Page{
			Total: count,
			Index: r.Page.Index,
			Size:  r.Page.Size,
		},
	}, nil
}

func (rr RuleRepo) Create(r models.AlertRule) error {
	nr := r
	nr.RuleId = "a-" + cmd.RandId()

	err := rr.g.Create(models.AlertRule{}, nr)
	if err != nil {
		return err
	}

	return nil
}

func (rr RuleRepo) Update(r models.AlertRule) error {
	u := Updates{
		Table: &models.AlertRule{},
		Where: map[string]interface{}{
			"tenant_id = ?": r.TenantId,
			"rule_id = ?":   r.RuleId,
		},
		Updates: r,
	}

	err := rr.g.Updates(u)
	if err != nil {
		return err
	}

	return nil
}

func (rr RuleRepo) Delete(r models.AlertRuleQuery) error {
	var alertRule models.AlertRule
	d := Delete{
		Table: alertRule,
		Where: map[string]interface{}{
			"tenant_id = ?": r.TenantId,
			"rule_id = ?":   r.RuleId,
		},
	}

	err := rr.g.Delete(d)
	if err != nil {
		return err
	}

	return nil
}

func (rr RuleRepo) GetRuleIsExist(ruleId string) bool {
	var ruleNum int64
	rr.DB().Model(&models.AlertRule{}).
		Where("rule_id = ? AND enabled = ?", ruleId, "1").
		Count(&ruleNum)
	if ruleNum > 0 {
		return true
	}

	return false
}
