package view

import (
	"Road/moudle/common"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

func getpocinfo(resp *gin.Context) {
	res := common.Pluginfos
	resp.JSON(http.StatusOK, res)
}
func showpoc(resp *gin.Context) {
	resp.Request.URL.Path = "/"
	content, err := ioutil.ReadFile("moudle/static/mq-admin/pages/showpoc/showpoc.html")
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
func Loadshowpoc(e *gin.Engine, v1 *gin.RouterGroup) {
	v1.GET("/admin/poc/info", getpocinfo)
	v1.GET("/admin/pocmanager", showpoc)

}
