package models

type UserRole struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Permissions []UserPermissions `json:"permissions" gorm:"permissions;serializer:json"`
	CreateAt    int64             `json:"create_at"`
}

type UserRoleQuery struct {
	ID          string `json:"id" form:"id"`
	Name        string `json:"name" form:"name"`
	Description string `json:"description" form:"description"`
}
