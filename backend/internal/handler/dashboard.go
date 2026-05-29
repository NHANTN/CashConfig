package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/cashier-config/server/internal/model"
)

type DashboardHandler struct {
	db *gorm.DB
}

func NewDashboardHandler(db *gorm.DB) *DashboardHandler {
	return &DashboardHandler{db: db}
}

func (h *DashboardHandler) Register(api gin.IRouter, authed gin.IRouter) {
	authed.GET("/dashboard/stats", h.Stats)
}

func (h *DashboardHandler) Stats(c *gin.Context) {
	var moduleCount, ruleCount, storeCount, tillCount, varCount, groupCount, userCount int64
	h.db.Model(&model.Module{}).Count(&moduleCount)
	h.db.Model(&model.Rule{}).Count(&ruleCount)
	h.db.Model(&model.Store{}).Count(&storeCount)
	h.db.Model(&model.TillList{}).Count(&tillCount)
	h.db.Model(&model.Var{}).Count(&varCount)
	h.db.Model(&model.Group{}).Count(&groupCount)
	h.db.Model(&model.User{}).Count(&userCount)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "ok",
		"data": gin.H{
			"module_count":  moduleCount,
			"rule_count":    ruleCount,
			"store_count":   storeCount,
			"till_count":    tillCount,
			"var_count":     varCount,
			"group_count":   groupCount,
			"user_count":    userCount,
		},
	})
}
