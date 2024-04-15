package models

import (
	"bytes"
	"fmt"
	"watchAlert/public/globals"
)

type AuditLog struct {
	TenantId   string `json:"tenantId" form:"tenantId"`
	ID         string `json:"id" form:"id"`
	Username   string `json:"username" form:"username"`
	IPAddress  string `json:"ipAddress" form:"ipAddress"`
	Method     string `json:"method"`
	Path       string `json:"path" form:"path"`
	CreatedAt  int64  `json:"createdAt" form:"createdAt"`
	StatusCode int    `json:"statusCode"`
	Body       string `json:"body"`
	AuditType  string `json:"auditType"`
}

type AuditLogQuery struct {
	TenantId string `json:"tenantId" form:"tenantId"`
	Query    string `json:"query" form:"query"`
	Scope    string `json:"scope" form:"scope"`
	Page
}

type Page struct {
	Total     string `json:"total" form:"total"`
	PageIndex string `json:"pageIndex" form:"pageIndex"`
	PageSize  string `json:"pageSize" form:"pageSize"`
}

type AuditLogResponse struct {
	List       []AuditLog
	PageIndex  int64
	PageSize   int64
	TotalCount int64
}

func (a AuditLog) String() string {
	buf := bytes.NewBuffer(nil)
	buf.WriteString("{")
	buf.WriteString(fmt.Sprintf("id: %s", a.ID))
	buf.WriteString(fmt.Sprintf("username: %s", a.Username))
	buf.WriteString(fmt.Sprintf("ip_address: %s", a.IPAddress))
	buf.WriteString(fmt.Sprintf("method: %s", a.Method))
	buf.WriteString(fmt.Sprintf("path: %s", a.Path))
	buf.WriteString(fmt.Sprintf("createdAt: %d", a.CreatedAt))
	buf.WriteString(fmt.Sprintf("statusCode: %d", a.StatusCode))
	buf.WriteString("}")
	return buf.String()
}

func (a AuditLog) SaveAuditLog() error {
	err := globals.DBCli.Model(&AuditLog{}).Create(a).Error
	if err != nil {
		return err
	}
	return nil
}
