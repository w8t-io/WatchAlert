package models

type Dashboard struct {
	TenantId    string `json:"tenantId"`
	ID          string `json:"id" `
	Name        string `json:"name" gorm:"unique"`
	URL         string `json:"url"`
	FolderId    string `json:"folderId"`
	Description string `json:"description"`
}

type DashboardQuery struct {
	TenantId string `json:"tenantId" form:"tenantId"`
	ID       string `json:"id" form:"id"`
	Query    string `json:"query" form:"query"`
}

type DashboardFolders struct {
	TenantId            string `json:"tenantId" form:"tenantId"`
	ID                  string `json:"id" form:"id"`
	Name                string `json:"name"`
	Theme               string `json:"theme" form:"theme"`
	GrafanaHost         string `json:"grafanaHost" form:"grafanaHost"`
	GrafanaFolderId     int    `json:"grafanaFolderId" form:"grafanaFolderId"`
	GrafanaDashboardUid string `json:"grafanaDashboardUid" form:"grafanaDashboardUid" gorm:"-"`
}

type GrafanaDashboardInfo struct {
	Uid      string `json:"uid"`
	Title    string `json:"title"`
	FolderId int    `json:"folderId"`
}

type GrafanaDashboardMeta struct {
	Meta meta `json:"meta"`
}

type meta struct {
	Url string `json:"url"`
}
