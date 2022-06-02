package view

import (
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

func Loadidnex(e *gin.Engine, v1 *gin.RouterGroup) {
	v1.GET("../account/logout", logout)
	v1.POST("/index", getindex)
	v1.GET("/index", getindex)
}
