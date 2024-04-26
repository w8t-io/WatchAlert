package repo

import (
	"fmt"
	"gorm.io/gorm"
	"watchAlert/internal/models"
)

type (
	NoticeRepo struct {
		entryRepo
	}

	InterNoticeRepo interface {
		Get(r models.NoticeQuery) (models.AlertNotice, error)
		GetQuota(id string) bool
		Search(req models.NoticeQuery) ([]models.AlertNotice, error)
		List(req models.NoticeQuery) ([]models.AlertNotice, error)
		Create(r models.AlertNotice) error
		Update(r models.AlertNotice) error
		Delete(r models.NoticeQuery) error
	}
)

func newNoticeInterface(db *gorm.DB, g InterGormDBCli) InterNoticeRepo {
	return &NoticeRepo{
		entryRepo{
			g:  g,
			db: db,
		},
	}
}

func (nr NoticeRepo) GetQuota(id string) bool {
	var (
		db     = nr.db.Model(&models.Tenant{})
		data   models.Tenant
		Number int64
	)

	db.Where("id = ?", id)
	db.Find(&data)

	nr.db.Model(&models.AlertNotice{}).Where("tenant_id = ?", id).Count(&Number)

	if Number < data.NoticeNumber {
		return true
	}

	return false
}

func (nr NoticeRepo) Get(r models.NoticeQuery) (models.AlertNotice, error) {
	var alertNoticeData models.AlertNotice
	db := nr.db.Model(&models.AlertNotice{}).Where("tenant_id = ? AND uuid = ?", r.TenantId, r.Uuid)
	err := db.First(&alertNoticeData).Error
	if err != nil {
		return alertNoticeData, err
	}

	return alertNoticeData, nil
}

func (nr NoticeRepo) Search(req models.NoticeQuery) ([]models.AlertNotice, error) {
	var data []models.AlertNotice
	var db = nr.db.Model(&models.AlertNotice{})
	db.Where("name LIKE ? OR env LIKE ? OR notice_type LIKE ?", "%"+req.Query+"%", "%"+req.Query+"%", "%"+req.Query+"%")
	err := db.Find(&data).Error
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (nr NoticeRepo) List(req models.NoticeQuery) ([]models.AlertNotice, error) {
	var alertNoticeObject []models.AlertNotice
	db := nr.db.Model(&models.AlertNotice{})

	db.Where("tenant_id = ?", req.TenantId)
	err := db.Find(&alertNoticeObject).Error
	if err != nil {
		return nil, err
	}

	return alertNoticeObject, nil
}

func (nr NoticeRepo) Create(r models.AlertNotice) error {
	err := nr.g.Create(models.AlertNotice{}, r)
	if err != nil {
		return err
	}
	return nil
}

func (nr NoticeRepo) Update(r models.AlertNotice) error {
	u := Updates{
		Table: models.AlertNotice{},
		Where: map[string]interface{}{
			"tenant_id = ?": r.TenantId,
			"uuid = ?":      r.Uuid,
		},
		Updates: r,
	}
	err := nr.g.Updates(u)
	if err != nil {
		return err
	}
	return nil
}

func (nr NoticeRepo) Delete(r models.NoticeQuery) error {

	var ruleNum1, ruleNum2 int64
	db := nr.db.Model(&models.AlertRule{})
	db.Where("notice_id = ?", r.Uuid).Count(&ruleNum1)
	db.Where("notice_group LIKE ?", "%"+r.Uuid+"%").Count(&ruleNum2)
	if ruleNum1 != 0 || ruleNum2 != 0 {
		return fmt.Errorf("无法删除通知对象 %s, 因为已有告警规则绑定", r.Uuid)
	}

	d := Delete{
		Table: models.AlertNotice{},
		Where: map[string]interface{}{
			"tenant_id = ?": r.TenantId,
			"uuid = ?":      r.Uuid,
		},
	}
	err := nr.g.Delete(d)
	if err != nil {
		return err
	}
	return nil
}
