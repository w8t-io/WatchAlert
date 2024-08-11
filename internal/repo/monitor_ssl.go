package repo

import (
	"fmt"
	"gorm.io/gorm"
	"watchAlert/internal/global"
	"watchAlert/internal/models"
)

type (
	monitorSSLRepo struct {
		entryRepo
	}

	InterMonitorSSLRepo interface {
		Get(r models.MonitorSSLRuleQuery) (models.MonitorSSLRule, error)
		List(req models.MonitorSSLRuleQuery) ([]models.MonitorSSLRule, error)
		Create(r models.MonitorSSLRule) error
		Update(r models.MonitorSSLRule) error
		Delete(r models.MonitorSSLRuleQuery) error
	}
)

func newMonitorSSLInterface(db *gorm.DB, g InterGormDBCli) InterMonitorSSLRepo {
	return &monitorSSLRepo{
		entryRepo{
			g:  g,
			db: db,
		},
	}
}

func (m monitorSSLRepo) Create(r models.MonitorSSLRule) error {
	err := m.g.Create(models.MonitorSSLRule{}, r)
	if err != nil {
		return err
	}
	return nil
}

func (m monitorSSLRepo) Update(r models.MonitorSSLRule) error {
	u := Updates{
		Table: models.MonitorSSLRule{},
		Where: map[string]interface{}{
			"tenant_id = ?": r.TenantId,
			"id = ?":        r.ID,
		},
		Updates: r,
	}
	fmt.Println("--->", r)
	err := m.g.Updates(u)
	if err != nil {
		global.Logger.Sugar().Error(err.Error())
		return err
	}
	return nil
}

func (m monitorSSLRepo) Delete(r models.MonitorSSLRuleQuery) error {
	d := Delete{
		Table: models.MonitorSSLRule{},
		Where: map[string]interface{}{
			"tenant_id = ?": r.TenantId,
			"id = ?":        r.ID,
		},
	}
	err := m.g.Delete(d)
	if err != nil {
		return err
	}
	return nil
}

func (m monitorSSLRepo) List(req models.MonitorSSLRuleQuery) ([]models.MonitorSSLRule, error) {
	var Objects []models.MonitorSSLRule
	db := m.db.Model(&models.MonitorSSLRule{})
	db.Where("tenant_id = ?", req.TenantId)

	if req.Query != "" {
		db.Where("id LIKE ? OR name LIKE ? OR domain LIKE ?", "%"+req.Query+"%", "%"+req.Query+"%", "%"+req.Query+"%")
	}

	err := db.Find(&Objects).Error
	if err != nil {
		return nil, err
	}

	return Objects, nil
}

func (m monitorSSLRepo) Get(r models.MonitorSSLRuleQuery) (models.MonitorSSLRule, error) {
	var Object models.MonitorSSLRule
	db := m.db.Model(&models.MonitorSSLRule{}).Where("tenant_id = ? AND id = ?", r.TenantId, r.ID)
	err := db.First(&Object).Error
	if err != nil {
		return Object, err
	}

	return Object, nil
}
