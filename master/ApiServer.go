package master

import (
	"fmt"
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
	)
	if err = c.ShouldBind(&job); err != nil {
		goto ERR
	}
	// save to etcd
	if oldJob, err = G_jobMgr.SaveJob(&job); err != nil {
		goto ERR
	}

	//Return to normal reply({"errNo":0, "msg":"", "data":{}})
	c.JSON(http.StatusOK, common.BuildResponse(0, "success", oldJob))
	return
ERR:
	//Return exception reply
	c.JSON(http.StatusOK, common.BuildResponse(-1, err.Error(), nil))
}

// JobDeleteHandler DELETE /job/delete name=job1
func JobDeleteHandler(c *gin.Context) {
	var (
		job    common.Job
		err    error
		oldJob *common.Job
	)
	if err = c.ShouldBind(&job); err != nil {
		fmt.Println(job)
		goto ERR
	}
	fmt.Println(job)
	if oldJob, err = G_jobMgr.DeleteJob(job.Name); err != nil {
		goto ERR
	}

	//Return to normal reply
	c.JSON(http.StatusOK, common.BuildResponse(0, "success", oldJob))
	return
ERR:
	//Return exception reply
	c.JSON(http.StatusOK, common.BuildResponse(-1, err.Error(), nil))
}

// JobListHandler GET list all jobs of crontab
func JobListHandler(c *gin.Context) {
	var (
		jobList []*common.Job
		err     error
	)
	if jobList, err = G_jobMgr.ListJob(); err != nil {
		goto ERR
	}

	//Return to normal reply
	//c.JSON(http.StatusOK, gin.H{"errno": 0, "msg": "success", "data": jobList})
	c.JSON(http.StatusOK, common.BuildResponse(0, "success", jobList))
	return
ERR:
	//Return exception reply
	c.JSON(http.StatusOK, common.BuildResponse(-1, err.Error(), nil))
}

// JobKillHandler POST /job/kill name=job1
func JobKillHandler(c *gin.Context) {
	var (
		job common.Job
		err error
	)

	if err = c.ShouldBind(&job); err != nil {
		goto ERR
	}

	if err = G_jobMgr.KillJob(job.Name); err != nil {
		goto ERR
	}

	//Return to normal reply
	c.JSON(http.StatusOK, common.BuildResponse(0, "success", nil))
	return
ERR:
	//Return exception reply
	c.JSON(http.StatusOK, common.BuildResponse(-1, err.Error(), nil))
}

func InitApiServer() (err error) {
	var (
		jobGroup *gin.RouterGroup
	)
	gin.SetMode(gin.ReleaseMode)

	// Assignment singleton
	G_apiServer = &ApiServer{
		router: gin.Default(),
	}

	//Configure the routing
	jobGroup = G_apiServer.router.Group("/job")
	{
		jobGroup.POST("/save", JobSaveHandler)
		jobGroup.DELETE("/delete", JobDeleteHandler)
		jobGroup.GET("/list", JobListHandler)
		jobGroup.POST("/kill", JobKillHandler)
	}

	// set static file directory
	G_apiServer.router.LoadHTMLGlob(G_config.Webroot)
	G_apiServer.router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	if err = G_apiServer.router.Run(":" + G_config.ApiPort); err != nil {
		return
	}
	return
}
