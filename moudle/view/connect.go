package view

import (
	"Road/moudle/sqlmoudle"
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

func getconnect(resp *gin.Context) {
	resp.Request.URL.Path = "/"
	content, err := ioutil.ReadFile("moudle/static/mq-admin/pages/connect/connect_equipment.html")
	if err != nil {
		resp.Writer.WriteHeader(404)
		resp.Writer.WriteString("Not Found")
		return
	}
	resp.Writer.WriteHeader(200)
	resp.Writer.Header().Add("Accept", "text/html")
	resp.Writer.Write(content)
	resp.Writer.Flush()
}
func getconnectinfo(resp *gin.Context) {
	linkinfo := sqlmoudle.Queryconnectinfo()
	resp.JSON(http.StatusOK, linkinfo)
}

func getconnect_taskinfo(resp *gin.Context) {
	flag := resp.DefaultQuery("uuid", "")
	if flag != "" {
		linkinfo := sqlmoudle.Queryconnecttaskinfo(flag)
		resp.JSON(http.StatusOK, linkinfo)
	}

}
func insertequipmenttask(resp *gin.Context) {
	var json sqlmoudle.Linktask
	resp.BindJSON(&json)
	Task_Id := strconv.FormatInt(time.Now().Unix(), 12)
	flag := sqlmoudle.Insertconnecttask(sqlmoudle.Linktask{
		Task_Id:       Task_Id,
		UUID:          json.UUID,
		Update_time:   sql.NullTime{Time: time.Now(), Valid: true},
		Update_Order:  json.Update_Order,
		Update_Result: "",
		Update_Status: 0,
	})
	resp.JSONP(http.StatusOK, flag)
}

func deleteequipmenttask(resp *gin.Context) {
	var json sqlmoudle.Linktask
	resp.BindJSON(&json)
	//Task_Id := strconv.FormatInt(time.Now().Unix(), 12)
	flag := sqlmoudle.Deleteconnecttask(json)
	resp.JSONP(http.StatusOK, flag)
}

func deleteequipment(resp *gin.Context) {
	var json sqlmoudle.Linkinfo
	resp.BindJSON(&json)
	//Task_Id := strconv.FormatInt(time.Now().Unix(), 12)
	flag := sqlmoudle.DeleteconnectEquipment(json)
	resp.JSONP(http.StatusOK, flag)
}

func Handleconnectequipment(ctx *gin.Context) {
	wsConn, _ := Upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	lastlength := len(sqlmoudle.Queryconnectinfo())
	for {
		temp := sqlmoudle.Queryconnectinfo()
		templength := len(temp)
		if templength-lastlength > 0 {
			wsConn.WriteMessage(websocket.TextMessage, []byte("{\"result\":\"true\",\"IP\":\""+temp[templength-1].Ip_addr+"\"}"))
		} else {
			wsConn.WriteMessage(websocket.TextMessage, []byte("{\"result\":\"false\",\"IP\":\"\"}"))
		}
		lastlength = templength
		time.Sleep(5 * time.Second)
	}
}

func Loadconnect(e *gin.Engine, v1 *gin.RouterGroup) {
	v1.Any("/admin/connect/websocket", Handleconnectequipment)
	v1.GET("/admin/connect", getconnect)
	v1.GET("/admin/connect/info", getconnectinfo)
	v1.GET("/admin/connect/task/info", getconnect_taskinfo)
	v1.POST("/admin/connect/task/init", insertequipmenttask)
	v1.POST("/admin/connect/task/delete", deleteequipmenttask)
	v1.POST("/admin/connect/equipment/delete", deleteequipment)
}
