package renderTemplates

import (
	"context"
	"encoding/json"
	"fmt"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkcontact "github.com/larksuite/oapi-sdk-go/v3/service/contact/v3"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	"watchAlert/controllers/dto"
	"watchAlert/globals"
)

func PushFeiShu(chatId string, cardContentJson []string) error {

	for _, v := range cardContentJson {
		req := larkim.NewCreateMessageReqBuilder().
			ReceiveIdType(`chat_id`).
			Body(larkim.NewCreateMessageReqBodyBuilder().
				ReceiveId(chatId).
				MsgType(`interactive`).
				Content(v).
				Build()).
			Build()

		resp, err := globals.FeiShuCli.Im.Message.Create(context.Background(), req, larkcore.WithTenantAccessToken(globals.Config.FeiShu.Token))
		// 处理错误
		if err != nil {
			globals.Logger.Sugar().Error("消息卡片发送失败 ->", err)
			return fmt.Errorf("消息卡片发送失败 -> %s", err)
		}

		// 服务端错误处理
		if !resp.Success() {
			globals.Logger.Sugar().Error(resp.Code, resp.Msg, resp.RequestId())
			return fmt.Errorf("响应错误 -> %s", err)
		}

		globals.Logger.Sugar().Info("消息卡片发送成功 ->", string(resp.RawBody))
	}

	return nil
}

func GetFeiShuUserInfo(userID string) dto.FeiShuUserInfo {

	// 创建请求对象
	req := larkcontact.NewGetUserReqBuilder().
		UserId(userID).
		UserIdType(`user_id`).
		DepartmentIdType(`open_department_id`).
		Build()

	// 发起请求
	// 如开启了SDK的Token管理功能，就无需在请求时调用larkcore.WithTenantAccessToken("-xxx")来手动设置租户Token了
	resp, err := globals.FeiShuCli.Contact.User.Get(context.Background(), req, larkcore.WithTenantAccessToken(globals.Config.FeiShu.Token))

	// 处理错误
	if err != nil {
		globals.Logger.Sugar().Error("获取飞书用户信息失败 ->", err)
		return dto.FeiShuUserInfo{}
	}

	var feiShuUserInfo dto.FeiShuUserInfo
	respJson, _ := json.Marshal(resp)
	_ = json.Unmarshal(respJson, &feiShuUserInfo)

	return feiShuUserInfo

}
