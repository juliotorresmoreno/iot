package handlers

import (
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/juliotorresmoreno/iot/etl/db"
	"github.com/juliotorresmoreno/iot/etl/entity"
	"github.com/juliotorresmoreno/iot/etl/kafka"
	log "github.com/sirupsen/logrus"
)

var manager = db.DefaultManager

type JobsHandler struct{}

func AttachJobsHandler(g *gin.RouterGroup) {
	h := &JobsHandler{}
	g.PUT("/import", h.Import)
}

type ImportBody struct {
	Type string `json:"type" valid:"required,type(string),in(fake|real)"`
}

func (h *JobsHandler) Import(c *gin.Context) {
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

	kafkaCli := kafka.NewKafkaClient("jobs")
	var source entity.Source
	if payload.Type == "real" {
		source = entity.Real
	} else if payload.Type == "fake" {
		source = entity.Fake
	}
	dbCli, err := db.MakeManager(source)
	if err != nil {
		log.Error(err)
		c.JSON(ErrInternalServerError.Status, ErrInternalServerError.Body)
		return
	}

	kafkaCli.Pub(entity.Job{
		Type: "add",
		Name: "",
		Data: map[string]string{"source": payload.Type},
	})

	log.Println(kafkaCli, dbCli)

	c.JSON(http.StatusOK, HttpBody{Message: "Task added"})
}
