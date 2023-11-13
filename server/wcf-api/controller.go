package wcf

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"

	"github.com/rehiy/wechat-rest-api/config"
	"github.com/rehiy/wechat-rest-api/wcf-sdk"
)

var wc *wcf.Client

func initWCF() {

	var err error

	wc, err = wcf.NewWCF(config.Wcf.Address)
	if err != nil {
		panic(err)
	}

}

func isLogin(c *gin.Context) {

	c.Set("Payload", wc.IsLogin())

}

func getSelfWxid(c *gin.Context) {

	c.Set("Payload", wc.GetSelfWxid())

}

func getUserInfo(c *gin.Context) {

	c.Set("Payload", wc.GetUserInfo())

}

func getMsgTypes(c *gin.Context) {

	c.Set("Payload", wc.GetMsgTypes())

}

func getContacts(c *gin.Context) {

	c.Set("Payload", wc.GetContacts())

}

func getFriends(c *gin.Context) {

	c.Set("Payload", wc.GetFriends())

}

func getDbNames(c *gin.Context) {

	c.Set("Payload", wc.GetDbNames())

}

func getDbTables(c *gin.Context) {

	db := c.Param("db")
	c.Set("Payload", wc.GetDbTables(db))

}

func refreshPyq(c *gin.Context) {

	id := cast.ToUint64(c.Param("id"))
	c.Set("Payload", wc.RefreshPyq(id))

}

func getChatRoomMembers(c *gin.Context) {

	roomid := c.Param("roomid")
	c.Set("Payload", wc.GetChatRoomMembers(roomid))

}
