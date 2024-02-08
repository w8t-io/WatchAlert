package models

type UserRole struct {
	ID              string            `json:"id"`
	Name            string            `json:"name"`
	Description     string            `json:"description"`
	PermissionsJson []UserPermissions `json:"permissions" gorm:"-"`
	Permissions     string            `json:"-" gorm:"permissions"`
	CreateAt        int64             `json:"create_at"`
}
