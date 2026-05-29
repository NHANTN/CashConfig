package database

import (
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/cashier-config/server/internal/model"
)

func Seed(db *gorm.DB) {
	// Create default roles
	var roleCount int64
	db.Model(&model.Role{}).Count(&roleCount)
	if roleCount == 0 {
		adminRole := model.Role{
			Name:        "超级管理员",
			Code:        "super_admin",
			Permissions: `["系统全部功能"]`,
			Status:      1,
		}
		if err := db.Create(&adminRole).Error; err != nil {
			log.Printf("seed: failed to create admin role: %v", err)
			return
		}
		log.Println("seed: created super_admin role")

		opsRole := model.Role{
			Name:        "运营管理员",
			Code:        "ops_admin",
			Permissions: `["门店管理","设备管理","支付配置","收银台布局","查看报表"]`,
			Status:      1,
		}
		if err := db.Create(&opsRole).Error; err != nil {
			log.Printf("seed: failed to create ops role: %v", err)
			return
		}
		log.Println("seed: created ops_admin role")
	}

	// Create default admin user
	var userCount int64
	db.Model(&model.User{}).Count(&userCount)
	if userCount == 0 {
		hash, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("seed: failed to hash password: %v", err)
			return
		}
		var adminRole model.Role
		if err := db.Where("code = ?", "super_admin").First(&adminRole).Error; err != nil {
			log.Printf("seed: admin role not found: %v", err)
			return
		}
		admin := model.User{
			Username:     "admin",
			PasswordHash: string(hash),
			Name:         "超级管理员",
			RoleID:       adminRole.ID,
			Status:       1,
		}
		if err := db.Create(&admin).Error; err != nil {
			log.Printf("seed: failed to create admin user: %v", err)
			return
		}
		log.Println("seed: created admin user (admin / admin123)")
	}
}
