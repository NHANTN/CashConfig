package main

import (
	"fmt"
	"log"
	"os"

	"go.uber.org/zap"

	"github.com/cashier-config/server/internal/config"
	"github.com/cashier-config/server/internal/database"
	"github.com/cashier-config/server/internal/handler"
	"github.com/cashier-config/server/internal/model"
	"github.com/cashier-config/server/internal/router"
	"github.com/cashier-config/server/internal/service"
)

func main() {
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Printf("config not found, using defaults: %v", err)
		cfg = config.DefaultConfig()
	}

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	db, err := database.Connect(cfg.Database)
	if err != nil {
		logger.Fatal("database connection failed", zap.Error(err))
	}
	logger.Info("database connected")

	if err := db.AutoMigrate(
		&model.User{},
		&model.Role{},
		&model.Module{},
		&model.Rule{},
		&model.Store{},
		&model.TillList{},
		&model.SyncReport{},
		&model.Var{},
		&model.Group{},
		&model.CsvGenerationLog{},
		&model.OperationLog{},
		&model.SSOState{},
	); err != nil {
		logger.Fatal("auto migrate failed", zap.Error(err))
	}
	logger.Info("database migrated")

	database.Seed(db)

	outputDir := "generated"
	os.MkdirAll(outputDir, 0755)

	authSvc := service.NewAuthService(db, cfg.JWT)
	moduleSvc := service.NewModuleService(db)
	ruleSvc := service.NewRuleService(db)
	storeSvc := service.NewStoreService(db)
	tillListSvc := service.NewTillListService(db)
	varSvc := service.NewVarService(db)
	groupSvc := service.NewGroupService(db)
	csvGenSvc := service.NewCsvGenerateService(db, moduleSvc, ruleSvc, storeSvc, tillListSvc, varSvc, outputDir)
	userSvc := service.NewUserService(db)
	roleSvc := service.NewRoleService(db)

	var ldapSvc *service.LDAPService
	if cfg.LDAP != nil && cfg.LDAP.Enabled {
		ldapSvc = service.NewLDAPService(db, cfg.LDAP)
		logger.Info("LDAP service enabled")
	}

	var ssoSvc *service.SSOService
	if cfg.SSO != nil && cfg.SSO.Enabled {
		var err error
		ssoSvc, err = service.NewSSOService(db, cfg.SSO)
		if err != nil {
			logger.Warn("SSO service initialization failed", zap.Error(err))
		} else {
			logger.Info("SSO service enabled")
		}
	}

	dashH := handler.NewDashboardHandler(db)
	authH := handler.NewAuthHandler(authSvc, ldapSvc, ssoSvc)
	csvGenH := handler.NewCsvGenerateHandler(csvGenSvc)
	moduleH := handler.NewModuleHandler(moduleSvc)
	ruleH := handler.NewRuleHandler(ruleSvc)
	storeH := handler.NewStoreHandler(storeSvc)
	tillListH := handler.NewTillListHandler(tillListSvc)
	varH := handler.NewVarHandler(varSvc)
	groupH := handler.NewGroupHandler(groupSvc)
	scriptH := handler.NewScriptHandler(cfg.WinTillModulesPath)
	userH := handler.NewUserHandler(userSvc)
	roleH := handler.NewRoleHandler(roleSvc)
	logH := handler.NewOperationLogHandler(db)

	r := router.Setup([]handler.Handler{authH, scriptH, dashH, csvGenH, moduleH, ruleH, storeH, tillListH, varH, groupH, userH, roleH, logH}, authSvc, cfg.APIKey)
	r.Run(fmt.Sprintf(":%d", cfg.Server.Port))
}
