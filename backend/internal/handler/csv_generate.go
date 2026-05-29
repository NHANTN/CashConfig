package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/cashier-config/server/internal/service"
)

type CsvGenerateHandler struct {
	svc *service.CsvGenerateService
}

func NewCsvGenerateHandler(svc *service.CsvGenerateService) *CsvGenerateHandler {
	return &CsvGenerateHandler{svc: svc}
}

func (h *CsvGenerateHandler) Register(api gin.IRouter, authed gin.IRouter) {
	csv := authed.Group("/csv")
	csv.POST("/generate", h.Generate)
	csv.POST("/generate/:type", h.Generate)
	csv.GET("/download/:type", h.Download)
	csv.GET("/download/all", h.DownloadAll)
	csv.GET("/history", h.History)
	csv.GET("/history/:type", h.History)
	csv.GET("/diff/:type", h.Diff)
}

func (h *CsvGenerateHandler) Generate(c *gin.Context) {
	fileType := c.Param("type")
	files, err := h.svc.Generate(fileType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}
	operator, _ := c.Get("username")
	timestamp, err := h.svc.SaveGeneration(files, fmt.Sprint(operator))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "ok",
		"data": gin.H{
			"timestamp": timestamp,
			"files":     len(files),
			"types":     fileType,
		},
	})
}

func (h *CsvGenerateHandler) Download(c *gin.Context) {
	fileType := c.Param("type")
	data, filename, err := h.svc.GetFileData(fileType)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": err.Error()})
		return
	}
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Data(http.StatusOK, "text/csv; charset=utf-8", data)
}

func (h *CsvGenerateHandler) DownloadAll(c *gin.Context) {
	data, filename, err := h.svc.GetAllAsZip()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": err.Error()})
		return
	}
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Data(http.StatusOK, "application/zip", data)
}

func (h *CsvGenerateHandler) History(c *gin.Context) {
	fileType := c.Param("type")
	logs, err := h.svc.GetHistory(fileType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok", "data": logs})
}

func (h *CsvGenerateHandler) Diff(c *gin.Context) {
	fileType := c.Param("type")
	from := c.Query("from")
	to := c.Query("to")
	if from == "" || to == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "from and to params required"})
		return
	}
	result, err := h.svc.Diff(fileType, from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok", "data": result})
}
