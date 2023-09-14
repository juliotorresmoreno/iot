package handlers

import (
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/juliotorresmoreno/iot/etl/data"
	"github.com/juliotorresmoreno/iot/etl/entity"
	"github.com/juliotorresmoreno/iot/etl/tasks"
	log "github.com/sirupsen/logrus"
)

type JobsHandler struct {
	taskManager *tasks.TaskManager
}

func AttachJobsHandler(g *gin.RouterGroup) {
	h := &JobsHandler{
		taskManager: tasks.DefaultTaskManager,
	}
	g.GET("", h.Find)
	g.PUT("", h.Put)
	g.DELETE("/:uuid", h.Delete)
}

type ImportBody struct {
	Type   string `json:"type"   valid:"required,type(string),in(import)"`
	Source string `json:"source" valid:"required,type(string),in(fake|real)"`
}

func (h *JobsHandler) Put(c *gin.Context) {
	payload := &ImportBody{}
	err := c.Bind(payload)
	if err != nil {
		log.Error(err)
		c.JSON(ErrBadRequest.Status, ErrBadRequest.Body)
		return
	}

	_, err = govalidator.ValidateStruct(payload)
	if err != nil {
		log.Error(err)
		c.JSON(ErrBadRequest.Status, HttpBody{Message: err.Error()})
		return
	}

	var source entity.Source = entity.Fake
	if payload.Source == "real" {
		source = entity.Real
	}
	etl, err := data.MakeETL(source)
	if err != nil {
		log.Error(err)
		c.JSON(ErrBadRequest.Status, HttpBody{Message: err.Error()})
		return
	}
	h.taskManager.Add(etl.Run())

	c.JSON(http.StatusOK, HttpBody{Message: "Task added"})
}

func (h *JobsHandler) Find(c *gin.Context) {
	tasks := h.taskManager.List()

	c.JSON(http.StatusOK, HttpBodyData{
		Data:  tasks,
		Total: len(tasks),
	})
}

func (h *JobsHandler) Delete(c *gin.Context) {
	uuid := c.Param("uuid")
	h.taskManager.Del(uuid)

	c.String(http.StatusNoContent, "")
}
