package repo

import (
	"watchAlert/models"
	"watchAlert/public/globals"
)

type UserRepo struct{}

func (ur UserRepo) SearchUser(r models.MemberQuery) ([]models.Member, error) {
	var data []models.Member
	var db = globals.DBCli.Model(&models.Member{})
	db.Where("user_name LIKE ? OR email Like ? OR phone LIKE ?", "%"+r.Query+"%", "%"+r.Query+"%", "%"+r.Query+"%")
	err := db.Find(&data).Error
	if err != nil {
		return nil, err
	}
	return data, nil
}
