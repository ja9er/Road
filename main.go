package main

import (
	"Road/moudle/sqlmoudle"
	"Road/moudle/view"
	"github.com/gin-gonic/gin"
)

func main() {
	sqlmoudle.InitDB()
	r := gin.Default()
	//r.Static("static", "moudle/static")
	//r.LoadHTMLGlob("moudle/templete/**/*")
	//r.StaticFS("/static", http.Dir("moudle/static/mq-admin"))
	r.Static("/static", "moudle/static/mq-admin")
	r.LoadHTMLGlob("moudle/templete/**/*")
	//r.NoRoute(func(c *gin.Context) {
	//	// 实现内部重定向
	//	c.HTML(http.StatusOK, "test.html", nil)
	//})

	v1 := r.Group("/", view.CheckAuth)
	view.Loadlogin(r)
	view.Loadidnex(r, v1)
	r.Run(":80")
}
