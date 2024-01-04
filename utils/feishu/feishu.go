package feishu

import (
	"context"
	"encoding/json"
	"fmt"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	"watchAlert/controllers/dto"
	"watchAlert/globals"
)

// GetFeiShuChatsID 获取机器人所在的群列表
func GetFeiShuChatsID() dto.FeiShuChats {

	// 创建请求对象
	req := larkim.NewListChatReqBuilder().
		SortType(`ByCreateTimeAsc`).
		PageSize(20).
		Build()

	// 发起请求
	// 如开启了SDK的Token管理功能，就无需在请求时调用larkcore.WithTenantAccessToken("-xxx")来手动设置租户Token了
	resp, err := globals.FeiShuCli.Im.Chat.List(context.Background(), req, larkcore.WithTenantAccessToken(globals.Config.FeiShu.Token))

	// 处理错误
	if err != nil {
		fmt.Println(err)
		return dto.FeiShuChats{}
	}

	// 服务端错误处理
	if !resp.Success() {
		fmt.Println(resp.Code, resp.Msg, resp.RequestId())
		return dto.FeiShuChats{}
	}

	var feiShuChats dto.FeiShuChats

	respJson, _ := json.Marshal(resp.Data)
	_ = json.Unmarshal(respJson, &feiShuChats)
	return feiShuChats

}

func CheckFeiShuChatId(chatId string) bool {

	data := GetFeiShuChatsID()
	for _, v := range data.Items {
		if v.ChatId == chatId {
			return true
		}
	}

	return false

}
