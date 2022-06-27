package view

import (
	"Road/moudle/searchinfo/scan"
	"Road/moudle/sqlmoudle"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Querysearch struct {
	Queryinfo string `json:"query"`
}

var (
	Upgrader = websocket.Upgrader{
		//允许跨域
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func Handle(ctx *gin.Context) {
	wsConn, _ := Upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	for {
		data2 := sqlmoudle.Querytaskinfo()
		buf, _ := json.Marshal(data2)
		wsConn.WriteMessage(websocket.TextMessage, buf)
		time.Sleep(30 * time.Second)
	}
}

func gettaskinfo(resp *gin.Context) {
	res := sqlmoudle.Querytaskinfo()

	resp.JSON(http.StatusOK, res)
}
func setsearchtask(resp *gin.Context) {
	var json Querysearch
	resp.BindJSON(&json)
	Task_Id := strconv.FormatInt(time.Now().Unix(), 10)
	log.Println(json.Queryinfo)
	go func() {
		sqlmoudle.Inserttask(sqlmoudle.Taskjob{
			Task_Id:   Task_Id,
			Fofaquery: json.Queryinfo,
			Progress:  1,
		})
		scan.FofaGet(Task_Id, json.Queryinfo)
		sqlmoudle.Upadtetask(sqlmoudle.Taskjob{
			Task_Id:  Task_Id,
			Progress: 2,
		})
	}()
	resp.JSONP(http.StatusOK, true)
}
func deletetask(resp *gin.Context) {
	var json sqlmoudle.Taskjob
	resp.BindJSON(&json)
	flag1 := sqlmoudle.Deletetask(json)
	flag2 := sqlmoudle.Deletebannerfromtask(json.Task_Id)
	resp.JSONP(http.StatusOK, flag1 && flag2)
}

func Loadtask(e *gin.Engine, v1 *gin.RouterGroup) {
	v1.POST("/admin/task/send", setsearchtask)
	v1.POST("/admin/task/delete", deletetask)
	v1.GET("/admin/task/info", gettaskinfo)
	v1.Any("/admin/task/websocket", Handle)
}
