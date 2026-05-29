package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/cashier-config/server/internal/model"
	"github.com/cashier-config/server/internal/service"
)

type GroupHandler struct {
	svc *service.GroupService
}

func NewGroupHandler(svc *service.GroupService) *GroupHandler {
	return &GroupHandler{svc: svc}
}

func (h *GroupHandler) Register(api gin.IRouter, authed gin.IRouter) {
	authed.GET("/groups/all", h.ListAll)
	groups := authed.Group("/groups")
	groups.GET("", h.List)
	groups.POST("", h.Create)
	groups.GET("/:id", h.Get)
	groups.PUT("/:id", h.Update)
	groups.DELETE("/:id", h.Delete)
}

// @Summary      List groups
// @Tags         Group
// @Security     ApiKeyAuth
// @Security     BearerAuth
// @Produce      json
// @Param        name  query  string  false  "Filter by name (LIKE)"
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /groups [get]
func (h *GroupHandler) List(c *gin.Context) {
	name := c.Query("name")
	list, err := h.svc.List(name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok", "data": list})
}

// @Summary      List all groups (no filter)
// @Tags         Group
// @Security     ApiKeyAuth
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /groups/all [get]
func (h *GroupHandler) ListAll(c *gin.Context) {
	list, err := h.svc.ListAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok", "data": list})
}

// @Summary      Get group by ID
// @Tags         Group
// @Security     ApiKeyAuth
// @Security     BearerAuth
// @Produce      json
// @Param        id   path  int  true  "Group ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router       /groups/{id} [get]
func (h *GroupHandler) Get(c *gin.Context) {
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

// @Summary      Create group
// @Tags         Group
// @Security     ApiKeyAuth
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Success      201  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /groups [post]
func (h *GroupHandler) Create(c *gin.Context) {
	var m model.Group
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

// @Summary      Update group
// @Tags         Group
// @Security     ApiKeyAuth
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id   path  int  true  "Group ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /groups/{id} [put]
func (h *GroupHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid id"})
		return
	}
	var m model.Group
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

// @Summary      Delete group
// @Tags         Group
// @Security     ApiKeyAuth
// @Security     BearerAuth
// @Produce      json
// @Param        id   path  int  true  "Group ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /groups/{id} [delete]
func (h *GroupHandler) Delete(c *gin.Context) {
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
