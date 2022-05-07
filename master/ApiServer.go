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

// JobSaveHandler POST
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

// JobDeleteHandler DELETE /job/delete name=job1
func JobDeleteHandler(c *gin.Context) {
	var (
		job    common.Job
		err    error
		oldJob *common.Job
		bytes  []byte
	)
	if err = c.ShouldBind(&job); err != nil {
		goto ERR
	}

	if oldJob, err = G_jobMgr.DeleteJob(job.Name); err != nil {
		goto ERR
	}

	//Return to normal reply
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

// JobListHandler GET list all jobs of crontab
func JobListHandler(c *gin.Context) {
	var (
		jobList []*common.Job
		err     error
		bytes   []byte
	)
	if jobList, err = G_jobMgr.ListJob(); err != nil {
		goto ERR
	}

	//Return to normal reply
	if bytes, err = common.BuildResponse(0, "success", jobList); err == nil {
		c.JSON(http.StatusOK, string(bytes))
	}
	return
ERR:
	//Return exception reply
	if bytes, err = common.BuildResponse(-1, err.Error(), nil); err == nil {
		c.JSON(http.StatusOK, string(bytes))
	}
}

// JobKillHandler POST /job/kill name=job1
func JobKillHandler(c *gin.Context) {
	var (
		job   common.Job
		err   error
		bytes []byte
	)

	if err = c.ShouldBind(&job); err != nil {
		goto ERR
	}

	if err = G_jobMgr.KillJob(job.Name); err != nil {
		goto ERR
	}

	//Return to normal reply
	if bytes, err = common.BuildResponse(0, "success", nil); err == nil {
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
	G_apiServer.router.DELETE("/job/delete", JobDeleteHandler)
	G_apiServer.router.GET("/job/list", JobListHandler)
	G_apiServer.router.POST("/job/kill", JobKillHandler)

	if err = G_apiServer.router.Run(":" + G_config.ApiPort); err != nil {
		return
	}
	return
}
