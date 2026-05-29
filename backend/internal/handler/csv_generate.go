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

// @Summary      Generate CSV files
// @Description  Trigger generation of CSV files. Optionally specify a type: module, rule, store, till, var
// @Tags         CSV
// @Security     ApiKeyAuth
// @Security     BearerAuth
// @Produce      json
// @Param        type  path  string  false  "File type (module/rule/store/till/var)"
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /csv/generate [post]
// @Router       /csv/generate/{type} [post]
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

// @Summary      Download CSV file
// @Description  Download a single CSV file by type
// @Tags         CSV
// @Security     ApiKeyAuth
// @Security     BearerAuth
// @Produce      text/csv
// @Param        type  path  string  true  "File type (module/rule/store/till/var)"
// @Success      200  {file}  string  "CSV file"
// @Failure      404  {object}  map[string]interface{}
// @Router       /csv/download/{type} [get]
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

// @Summary      Download all CSV as ZIP
// @Description  Download all CSV files as a ZIP archive
// @Tags         CSV
// @Security     ApiKeyAuth
// @Security     BearerAuth
// @Produce      application/zip
// @Success      200  {file}  string  "ZIP file"
// @Failure      404  {object}  map[string]interface{}
// @Router       /csv/download/all [get]
func (h *CsvGenerateHandler) DownloadAll(c *gin.Context) {
	data, filename, err := h.svc.GetAllAsZip()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": err.Error()})
		return
	}
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Data(http.StatusOK, "application/zip", data)
}

// @Summary      Get CSV generation history
// @Description  Get history of CSV generation. Optionally filter by type
// @Tags         CSV
// @Security     ApiKeyAuth
// @Security     BearerAuth
// @Produce      json
// @Param        type  path  string  false  "File type (module/rule/store/till/var)"
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /csv/history [get]
// @Router       /csv/history/{type} [get]
func (h *CsvGenerateHandler) History(c *gin.Context) {
	fileType := c.Param("type")
	logs, err := h.svc.GetHistory(fileType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok", "data": logs})
}

// @Summary      CSV diff
// @Description  Compare two versions of a CSV file by type
// @Tags         CSV
// @Security     ApiKeyAuth
// @Security     BearerAuth
// @Produce      json
// @Param        type  path   string  true   "File type (module/rule/store/till/var)"
// @Param        from  query  string  true   "From timestamp"
// @Param        to    query  string  true   "To timestamp"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /csv/diff/{type} [get]
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
