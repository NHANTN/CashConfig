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

type ModuleHandler struct {
	svc *service.ModuleService
}

func NewModuleHandler(svc *service.ModuleService) *ModuleHandler {
	return &ModuleHandler{svc: svc}
}

func (h *ModuleHandler) Register(api gin.IRouter, authed gin.IRouter) {
	modules := authed.Group("/modules")
	modules.GET("", h.List)
	modules.GET("/:id", h.Get)
	modules.POST("", h.Create)
	modules.PUT("/:id", h.Update)
	modules.DELETE("/:id", h.Delete)
	modules.POST("/import", h.ImportCSV)
	modules.GET("/export/csv", h.ExportCSV)
}

// @Summary      List modules
// @Tags         Module
// @Security     ApiKeyAuth
// @Security     BearerAuth
// @Produce      json
// @Param        env       query  string  false  "Filter by environment"
// @Param        location  query  string  false  "Filter by location"
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /modules [get]
func (h *ModuleHandler) List(c *gin.Context) {
	env := c.Query("env")
	location := c.Query("location")
	list, err := h.svc.List(env, location)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok", "data": list})
}

// @Summary      Get module by ID
// @Tags         Module
// @Security     ApiKeyAuth
// @Security     BearerAuth
// @Produce      json
// @Param        id   path  int  true  "Module ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router       /modules/{id} [get]
func (h *ModuleHandler) Get(c *gin.Context) {
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

// @Summary      Create module
// @Tags         Module
// @Security     ApiKeyAuth
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Success      201  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /modules [post]
func (h *ModuleHandler) Create(c *gin.Context) {
	var m model.Module
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

// @Summary      Update module
// @Tags         Module
// @Security     ApiKeyAuth
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id   path  int  true  "Module ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /modules/{id} [put]
func (h *ModuleHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid id"})
		return
	}
	var m model.Module
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

// @Summary      Delete module
// @Tags         Module
// @Security     ApiKeyAuth
// @Security     BearerAuth
// @Produce      json
// @Param        id   path  int  true  "Module ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /modules/{id} [delete]
func (h *ModuleHandler) Delete(c *gin.Context) {
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

// @Summary      Import module CSV
// @Tags         Module
// @Security     ApiKeyAuth
// @Security     BearerAuth
// @Accept       multipart/form-data
// @Produce      json
// @Param        file  formData  file  true  "CSV file"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /modules/import [post]
func (h *ModuleHandler) ImportCSV(c *gin.Context) {
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

// @Summary      Export module CSV
// @Tags         Module
// @Security     ApiKeyAuth
// @Security     BearerAuth
// @Produce      text/csv
// @Param        env       query  string  false  "Filter by environment"
// @Param        location  query  string  false  "Filter by location"
// @Success      200  {file}  string  "CSV file"
// @Failure      500  {object}  map[string]interface{}
// @Router       /modules/export/csv [get]
func (h *ModuleHandler) ExportCSV(c *gin.Context) {
	env := c.Query("env")
	location := c.Query("location")
	data, filename, err := h.svc.ExportCSV(env, location)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Data(http.StatusOK, "text/csv; charset=utf-8", data)
}
