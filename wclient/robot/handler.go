package robot

import (
	"sort"
	"strings"

	"github.com/opentdp/wechat-rest/dbase/keyword"
	"github.com/opentdp/wechat-rest/dbase/profile"
	"github.com/opentdp/wechat-rest/dbase/setting"
	"github.com/opentdp/wechat-rest/wcferry"
)

var handlers = []*Handler{}
var handlerPre = []*Handler{}
var handlerMap = map[string]*Handler{}

type HandlerFunc func(*wcferry.WxMsg) string

type Handler struct {
	Level    int32       // 0:不限制 7:群管理 9:创始人
	Order    int32       // 排序，越小越靠前
	ChatAble bool        // 是否允许在私聊使用
	RoomAble bool        // 是否允许在群聊使用
	Command  string      // 指令
	Describe string      // 指令的描述信息
	PreCheck HandlerFunc // 前置检查，可拦截文本聊天内容
	Callback HandlerFunc // 指令回调，返回回复内容
}

func GetHandlers() []*Handler {

	return handlers

}

func ResetHandlers() {

	hlst := []*Handler{}
	hpre := []*Handler{}
	hmap := map[string]*Handler{}

	// 获取指令列表
	hlst = append(hlst, aiHandler()...)
	hlst = append(hlst, apiHandler()...)
	hlst = append(hlst, badHandler()...)
	hlst = append(hlst, banHandler()...)
	hlst = append(hlst, helpHandler()...)
	hlst = append(hlst, revokeHandler()...)
	hlst = append(hlst, roomHandler()...)
	hlst = append(hlst, topHandler()...)
	hlst = append(hlst, wgetHandler()...)

	// 指令列表排序
	sort.Slice(hlst, func(i, j int) bool {
		return hlst[i].Order < hlst[j].Order
	})

	// 获取指令映射
	for _, v := range hlst {
		hmap[v.Command] = v
		if v.PreCheck != nil {
			hpre = append(hpre, v)
		}
	}

	// 获取别名数据
	kws, err := keyword.FetchAll(&keyword.FetchAllParam{Group: "handler"})
	if err == nil && len(kws) > 0 {
		for _, v := range kws {
			if hmap[v.Target] != nil {
				hmap[v.Phrase+"@"+v.Roomid] = hmap[v.Target]
			}
		}
	}

	// 更新全局数据
	handlers, handlerPre, handlerMap = hlst, hpre, hmap

}

func ApplyHandlers(msg *wcferry.WxMsg) string {

	// 前置钩子
	for _, v := range handlerPre {
		if txt := v.PreCheck(msg); txt != "" {
			return txt
		}
	}

	// 清理空白
	content := strings.TrimSpace(msg.Content)
	content = strings.ReplaceAll(content, " ", " ")
	if content == "" {
		return ""
	}

	// 解析指令
	params := strings.SplitN(content, " ", 2)
	handler := handlerMap[params[0]] // 默认
	if handler == nil {
		handler = handlerMap[params[0]+"@"+prid(msg)] // 群聊
		if handler == nil {
			handler = handlerMap[params[0]+"@-"] // 全局
			if handler == nil {
				if content[0] == '/' {
					return setting.InvalidHandler
				}
				return ""
			}
		}
	}

	// 验证权限
	if handler.Level > 0 {
		up, _ := profile.Fetch(&profile.FetchParam{Wxid: msg.Sender, Roomid: prid(msg)})
		if up.Level < handler.Level {
			return setting.InvalidHandler
		}
	}

	// 验证场景
	if (msg.IsGroup && !handler.RoomAble) || (!msg.IsGroup && !handler.ChatAble) {
		return setting.InvalidHandler
	}

	// 重写消息
	if len(params) > 1 {
		msg.Content = strings.TrimSpace(params[1])
	} else {
		msg.Content = ""
	}

	// 执行指令
	return handler.Callback(msg)

}
