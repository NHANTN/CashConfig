package handler

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/cashier-config/server/internal/model"
	"github.com/cashier-config/server/internal/service"
)

type TillListHandler struct {
	svc *service.TillListService
}

func NewTillListHandler(svc *service.TillListService) *TillListHandler {
	return &TillListHandler{svc: svc}
}

func (h *TillListHandler) Register(api gin.IRouter, authed gin.IRouter) {
	tills := authed.Group("/till-lists")
	tills.GET("", h.List)
	tills.GET("/reports", h.QueryReports)
	tills.GET("/:id", h.Get)
	tills.POST("", h.Create)
	tills.PUT("/:id", h.Update)
	tills.DELETE("/:id", h.Delete)
	tills.POST("/import", h.ImportCSV)
	tills.GET("/export/csv", h.ExportCSV)
	tills.GET("/:id/reports", h.ListReports)
	tills.GET("/:id/reports/:reportId", h.GetReport)
	tills.POST("/checkin", h.CheckIn)
}

// @Summary      List till lists
// @Description  List till lists with optional filters
// @Tags         TillList
// @Security     ApiKeyAuth
// @Security     BearerAuth
// @Produce      json
// @Param        host_name  query  string  false  "Filter by host name (LIKE)"
// @Param        location   query  string  false  "Filter by location"
// @Param        env        query  string  false  "Filter by environment"
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /till-lists [get]
func (h *TillListHandler) List(c *gin.Context) {
	hostName := c.Query("host_name")
	location := c.Query("location")
	env := c.Query("env")
	list, err := h.svc.List(hostName, location, env)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok", "data": list})
}

// @Summary      Get till list by ID
// @Description  Get a single till list entry
// @Tags         TillList
// @Security     ApiKeyAuth
// @Security     BearerAuth
// @Produce      json
// @Param        id   path  int  true  "Till list ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router       /till-lists/{id} [get]
func (h *TillListHandler) Get(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid id"})
		return
	}
	m, err := h.svc.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok", "data": m})
}

// @Summary      Create till list
// @Description  Create a new till list entry
// @Tags         TillList
// @Security     ApiKeyAuth
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Success      201  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /till-lists [post]
func (h *TillListHandler) Create(c *gin.Context) {
	var m model.TillList
	if err := c.ShouldBindJSON(&m); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}
	if err := h.svc.Create(&m); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"code": 0, "message": "created", "data": m})
}

// @Summary      Update till list
// @Description  Update an existing till list entry
// @Tags         TillList
// @Security     ApiKeyAuth
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id   path  int  true  "Till list ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /till-lists/{id} [put]
func (h *TillListHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid id"})
		return
	}
	var m model.TillList
	if err := c.ShouldBindJSON(&m); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}
	m.ID = id
	if err := h.svc.Update(&m); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "updated", "data": m})
}

// @Summary      Delete till list
// @Description  Delete a till list entry
// @Tags         TillList
// @Security     ApiKeyAuth
// @Security     BearerAuth
// @Produce      json
// @Param        id   path  int  true  "Till list ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /till-lists/{id} [delete]
func (h *TillListHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid id"})
		return
	}
	if err := h.svc.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "deleted"})
}

// @Summary      Import till list CSV
// @Description  Import till list data from CSV file upload
// @Tags         TillList
// @Security     ApiKeyAuth
// @Security     BearerAuth
// @Accept       multipart/form-data
// @Produce      json
// @Param        file  formData  file  true  "CSV file"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /till-lists/import [post]
func (h *TillListHandler) ImportCSV(c *gin.Context) {
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "file required"})
		return
	}
	defer file.Close()
	records, err := csv.NewReader(file).ReadAll()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid csv"})
		return
	}
	count, err := h.svc.ImportCSV(records)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok", "data": gin.H{"imported": count}})
}

// @Summary      Export till list CSV
// @Description  Export till list data as CSV file
// @Tags         TillList
// @Security     ApiKeyAuth
// @Security     BearerAuth
// @Produce      text/csv
// @Param        host_name  query  string  false  "Filter by host name"
// @Param        location   query  string  false  "Filter by location"
// @Param        env        query  string  false  "Filter by environment"
// @Success      200  {file}  string  "CSV file"
// @Failure      500  {object}  map[string]interface{}
// @Router       /till-lists/export/csv [get]
func (h *TillListHandler) ExportCSV(c *gin.Context) {
	hostName := c.Query("host_name")
	location := c.Query("location")
	env := c.Query("env")
	data, filename, err := h.svc.ExportCSV(hostName, location, env)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Data(http.StatusOK, "text/csv; charset=utf-8", data)
}

// @Summary      List reports for a till list
// @Description  Get all sync reports for a specific till list
// @Tags         TillList
// @Security     ApiKeyAuth
// @Security     BearerAuth
// @Produce      json
// @Param        id   path  int  true  "Till list ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /till-lists/{id}/reports [get]
func (h *TillListHandler) ListReports(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid id"})
		return
	}
	reports, err := h.svc.ListReports(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok", "data": reports})
}

// @Summary      Get a single report
// @Description  Get a specific sync report by ID
// @Tags         TillList
// @Security     ApiKeyAuth
// @Security     BearerAuth
// @Produce      json
// @Param        id        path  int  true  "Till list ID"
// @Param        reportId  path  int  true  "Report ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router       /till-lists/{id}/reports/{reportId} [get]
func (h *TillListHandler) GetReport(c *gin.Context) {
	reportID, err := strconv.ParseInt(c.Param("reportId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid report id"})
		return
	}
	r, err := h.svc.GetReport(reportID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "report not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok", "data": r})
}

// @Summary      Query reports by device
// @Description  Query sync reports grouped by device, filtering by host_name or mac_address
// @Tags         TillList
// @Security     ApiKeyAuth
// @Security     BearerAuth
// @Produce      json
// @Param        host_name    query  string  false  "Host name (LIKE)"
// @Param        mac_address  query  string  false  "MAC address (exact)"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Router       /till-lists/reports [get]
func (h *TillListHandler) QueryReports(c *gin.Context) {
	hostName := c.Query("host_name")
	macAddress := c.Query("mac_address")
	result, err := h.svc.QueryReportsByDevice(hostName, macAddress)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok", "data": result})
}

// @Summary      CheckIn
// @Description  Receive check-in data from a till device (raw JSON)
// @Tags         TillList
// @Security     ApiKeyAuth
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /till-lists/checkin [post]
func (h *TillListHandler) CheckIn(c *gin.Context) {
	raw, _ := c.GetRawData()
	if err := h.svc.CheckIn(string(raw)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok"})
}
