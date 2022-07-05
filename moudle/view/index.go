package view

import (
	"Road/moudle/sqlmoudle"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

func getindex(resp *gin.Context) {
	//c.HTML(http.StatusOK, "index.html", nil)
	content, err := ioutil.ReadFile("moudle/static/mq-admin/index.html")
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
func logout(c *gin.Context) {
	c.SetCookie("HMACCOUNT", "", 0, "/", "", false, true)
	c.SetCookie("name", "", 0, "/", "", false, true)
	c.Redirect(http.StatusFound, "../login")
}

func getuser(resp *gin.Context) {
	resp.Request.URL.Path = "/"
	content, err := ioutil.ReadFile("moudle/static/mq-admin/pages/welcome.html")
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

func getuserinfo(resp *gin.Context) {
	userinfo := sqlmoudle.Queryuser("%")
	resp.JSON(http.StatusOK, userinfo)
}

func gettaskdata(resp *gin.Context) {
	flag := resp.DefaultQuery("flag", "")
	if flag != "" {
		if flag == "0,1" {
			userinfo := sqlmoudle.Queryifno()
			resp.JSON(http.StatusOK, userinfo)
			return
		}
		int, _ := strconv.Atoi(flag)
		userinfo := sqlmoudle.QueryPOCmatch(int)
		resp.JSON(http.StatusOK, userinfo)
		return
	}
	userinfo := sqlmoudle.Queryifno()
	resp.JSON(http.StatusOK, userinfo)
}

func getdata(resp *gin.Context) {
	resp.Request.URL.Path = "/"
	content, err := ioutil.ReadFile("moudle/static/mq-admin/pages/target/target.html")
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

func gettaskmanager(resp *gin.Context) {
	resp.Request.URL.Path = "/"
	content, err := ioutil.ReadFile("moudle/static/mq-admin/pages/task/taskmanager.html")
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

func getconnectcount(resp *gin.Context) {
	flag := resp.DefaultQuery("type", "")
	if flag == "get" {
		data2 := sqlmoudle.Queryconnnectcount()

		resp.Writer.WriteHeader(200)
		resp.Writer.Header().Add("Accept", "text/html")
		buf, _ := json.Marshal(data2)
		resp.Writer.Write(buf)
		resp.Writer.Flush()
		return
	}
	wsConn2, _ := Upgrader.Upgrade(resp.Writer, resp.Request, nil)
	for {
		data2 := sqlmoudle.Queryconnnectcount()
		buf, _ := json.Marshal(data2)
		wsConn2.WriteMessage(websocket.TextMessage, buf)
		time.Sleep(10 * time.Second)
	}
}

func Loadidnex(e *gin.Engine, v1 *gin.RouterGroup) {
	v1.GET("../account/logout", logout)
	v1.Any("/index", getindex)
	v1.Any("/index/console/websocket", getconnectcount)
	v1.GET("/index/console/info", getconnectcount)

	v1.GET("/admin/user", getuser)
	v1.GET("/admin/user/userinfo", getuserinfo)
	v1.GET("/admin/data", getdata)
	v1.GET("/admin/data/info", gettaskdata)

	v1.GET("/admin/task", gettaskmanager)
}
