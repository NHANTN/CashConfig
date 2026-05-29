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

type StoreHandler struct {
	svc *service.StoreService
}

func NewStoreHandler(svc *service.StoreService) *StoreHandler {
	return &StoreHandler{svc: svc}
}

func (h *StoreHandler) Register(api gin.IRouter, authed gin.IRouter) {
	stores := authed.Group("/stores")
	stores.GET("", h.List)
	stores.GET("/:id", h.Get)
	stores.POST("", h.Create)
	stores.PUT("/:id", h.Update)
	stores.DELETE("/:id", h.Delete)
	stores.POST("/import", h.ImportCSV)
	stores.GET("/export/csv", h.ExportCSV)
}

// @Summary      List stores
// @Tags         Store
// @Security     ApiKeyAuth
// @Security     BearerAuth
// @Produce      json
// @Param        location  query  string  false  "Filter by location"
// @Param        eft       query  string  false  "Filter by EFT type"
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /stores [get]
func (h *StoreHandler) List(c *gin.Context) {
	location := c.Query("location")
	eft := c.Query("eft")
	list, err := h.svc.List(location, eft)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok", "data": list})
}

// @Summary      Get store by ID
// @Tags         Store
// @Security     ApiKeyAuth
// @Security     BearerAuth
// @Produce      json
// @Param        id   path  int  true  "Store ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router       /stores/{id} [get]
func (h *StoreHandler) Get(c *gin.Context) {
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

// @Summary      Create store
// @Tags         Store
// @Security     ApiKeyAuth
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Success      201  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /stores [post]
func (h *StoreHandler) Create(c *gin.Context) {
	var m model.Store
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

// @Summary      Update store
// @Tags         Store
// @Security     ApiKeyAuth
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id   path  int  true  "Store ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /stores/{id} [put]
func (h *StoreHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid id"})
		return
	}
	var m model.Store
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

// @Summary      Delete store
// @Tags         Store
// @Security     ApiKeyAuth
// @Security     BearerAuth
// @Produce      json
// @Param        id   path  int  true  "Store ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /stores/{id} [delete]
func (h *StoreHandler) Delete(c *gin.Context) {
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

// @Summary      Import store CSV
// @Tags         Store
// @Security     ApiKeyAuth
// @Security     BearerAuth
// @Accept       multipart/form-data
// @Produce      json
// @Param        file  formData  file  true  "CSV file"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /stores/import [post]
func (h *StoreHandler) ImportCSV(c *gin.Context) {
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

// @Summary      Export store CSV
// @Tags         Store
// @Security     ApiKeyAuth
// @Security     BearerAuth
// @Produce      text/csv
// @Param        location  query  string  false  "Filter by location"
// @Param        eft       query  string  false  "Filter by EFT type"
// @Success      200  {file}  string  "CSV file"
// @Failure      500  {object}  map[string]interface{}
// @Router       /stores/export/csv [get]
func (h *StoreHandler) ExportCSV(c *gin.Context) {
	location := c.Query("location")
	eft := c.Query("eft")
	data, filename, err := h.svc.ExportCSV(location, eft)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Data(http.StatusOK, "text/csv; charset=utf-8", data)
}
