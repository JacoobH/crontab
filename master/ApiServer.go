package master

import (
	"github.com/JacoobH/crontab/master/common"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ApiServer struct {
	router *gin.Engine
}

// G_apiServer Singleton
var (
	G_apiServer *ApiServer
)

func JobSaveHandler(c *gin.Context) {
	// POST job={"name":"job1", "command":"echo hello", "cronExpr":"* * * * *"}
	var (
		job    common.Job
		oldJob *common.Job
		err    error
		bytes  []byte
	)
	if err = c.ShouldBind(&job); err != nil {
		goto ERR
	}
	// save to etcd
	if oldJob, err = G_jobMgr.SaveJob(&job); err != nil {
		goto ERR
	}

	//Return to normal reply({"errNo":0, "msg":"", "data":{}})
	if bytes, err = common.BuildResponse(0, "success", oldJob); err == nil {
		c.JSON(http.StatusOK, string(bytes))
	}
	return
ERR:
	//Return exception reply
	if bytes, err = common.BuildResponse(-1, err.Error(), nil); err == nil {
		c.JSON(http.StatusOK, string(bytes))
	}
}

func InitApiServer() (err error) {
	gin.SetMode(gin.ReleaseMode)

	// Assignment singleton
	G_apiServer = &ApiServer{
		router: gin.Default(),
	}

	//配置路由
	G_apiServer.router.POST("/job/save", JobSaveHandler)

	if err = G_apiServer.router.Run(":" + G_config.ApiPort); err != nil {
		return
	}
	return
}
