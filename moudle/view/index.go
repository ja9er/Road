package view

import (
	"Road/moudle/sqlmoudle"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

func postindex(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "loginsuccess",
	})
}

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
	res := sqlmoudle.Queryuser("%")

	resp.JSON(http.StatusOK, res)
}
func Loadidnex(e *gin.Engine, v1 *gin.RouterGroup) {
	v1.GET("../account/logout", logout)
	v1.Any("/index", getindex)
	v1.GET("/admin/user", getuser)
	v1.Any("/admin/user/userinfo", getuserinfo)
}
