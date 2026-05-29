package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/cashier-config/server/internal/model"
)

type OperationLogHandler struct {
	db *gorm.DB
}

func NewOperationLogHandler(db *gorm.DB) *OperationLogHandler {
	return &OperationLogHandler{db: db}
}

func (h *OperationLogHandler) Register(api gin.IRouter, authed gin.IRouter) {
	system := authed.Group("/system")
	system.GET("/logs", h.List)
}

// @Summary      List operation logs
// @Tags         System
// @Security     ApiKeyAuth
// @Security     BearerAuth
// @Produce      json
// @Param        target_type  query  string  false  "Filter by target type"
// @Param        action       query  string  false  "Filter by action"
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /system/logs [get]
func (h *OperationLogHandler) List(c *gin.Context) {
	var logs []model.OperationLog
	q := h.db.Order("created_at DESC")
	if targetType := c.Query("target_type"); targetType != "" {
		q = q.Where("target_type = ?", targetType)
	}
	if action := c.Query("action"); action != "" {
		q = q.Where("action = ?", action)
	}
	if err := q.Limit(200).Find(&logs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok", "data": logs})
}
