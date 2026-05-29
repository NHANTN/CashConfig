package handler

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

type ScriptHandler struct {
	modulesPath string
}

func NewScriptHandler(modulesPath string) *ScriptHandler {
	return &ScriptHandler{modulesPath: modulesPath}
}

var scriptExts = map[string]bool{
	".ps1": true, ".bat": true, ".cmd": true, ".exe": true,
}

func (h *ScriptHandler) Register(api gin.IRouter, authed gin.IRouter) {
	api.GET("/script-files", h.ListFiles)
}

// @Summary      List script files
// @Description  List available PowerShell/BAT/CMD script files in the win-till-modules directory
// @Tags         Script
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /script-files [get]
func (h *ScriptHandler) ListFiles(c *gin.Context) {
	absPath, err := filepath.Abs(h.modulesPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}
	var files []string
	filepath.Walk(absPath, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		ext := strings.ToLower(filepath.Ext(info.Name()))
		if !scriptExts[ext] {
			return nil
		}
		rel, _ := filepath.Rel(absPath, path)
		files = append(files, filepath.ToSlash(rel))
		return nil
	})
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok", "data": files})
}
