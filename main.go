package main

import (
	"Road/moudle/sqlmoudle"
	"Road/moudle/view"
	"github.com/gin-gonic/gin"
	"io/ioutil"
)

func main() {
	sqlmoudle.InitDB()
	r := gin.Default()
	//r.Static("static", "moudle/static")
	//r.LoadHTMLGlob("moudle/templete/**/*")
	//r.StaticFS("/static", http.Dir("moudle/static/mq-admin"))
	r.Static("/static", "moudle/static/mq-admin")
	//r.LoadHTMLGlob("moudle/templete/**/*")
	r.NoRoute(func(resp *gin.Context) {
		// 实现内部重定向
		resp.Request.URL.Path = "/"
		content, err := ioutil.ReadFile("moudle/static/mq-admin/pages/login/login.html")
		if err != nil {
			resp.Writer.WriteHeader(404)
			resp.Writer.WriteString("Not Found")
			return
		}
		resp.Writer.WriteHeader(200)
		resp.Writer.Header().Add("Accept", "text/html")
		resp.Writer.Write(content)
		resp.Writer.Flush()
	})
	v1 := r.Group("/")
	view.Loadlogin(r)
	view.Loadidnex(r, v1)
	r.Run(":80")
}
