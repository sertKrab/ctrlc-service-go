package main

import (
	"context"
	"log"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"git.trovefin.com/poc/ctrlc-service-go/internal/config"
	"git.trovefin.com/poc/ctrlc-service-go/internal/domain/errmsg"
	"git.trovefin.com/poc/ctrlc-service-go/internal/domain/permission"
	"git.trovefin.com/poc/ctrlc-service-go/internal/domain/role"
	"git.trovefin.com/poc/ctrlc-service-go/internal/domain/user"
	"git.trovefin.com/poc/ctrlc-service-go/internal/infrastructure/database"
)

func main() {
	_ = godotenv.Load()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("[Seed] config: %v", err)
	}

	db, err := database.NewPostgresDB(cfg)
	if err != nil {
		log.Fatalf("[Seed] db: %v", err)
	}

	ctx := context.Background()
	seedPermissions(ctx, db)
	seedRoles(ctx, db)
	seedAdminUser(ctx, db)
	seedErrorMessages(ctx, db)

	log.Println("[Seed] complete")
}

func seedPermissions(ctx context.Context, db *gorm.DB) {
	perms := []permission.Permission{
		{Resource: "user",  Action: "create", Description: "Create users"},
		{Resource: "user",  Action: "read",   Description: "Read users"},
		{Resource: "user",  Action: "update", Description: "Update users"},
		{Resource: "user",  Action: "delete", Description: "Delete users"},
		{Resource: "role",  Action: "create", Description: "Create roles"},
		{Resource: "role",  Action: "read",   Description: "Read roles"},
		{Resource: "role",  Action: "update", Description: "Update roles"},
		{Resource: "role",  Action: "delete", Description: "Delete roles"},
		{Resource: "audit", Action: "read",   Description: "Read audit logs"},
	}
	db.WithContext(ctx).Clauses(clause.OnConflict{DoNothing: true}).Create(&perms)
	log.Printf("[Seed] permissions: %d rows", len(perms))
}

func seedRoles(ctx context.Context, db *gorm.DB) {
	var allPerms []permission.Permission
	db.WithContext(ctx).Find(&allPerms)

	var userReadPerm []permission.Permission
	for _, p := range allPerms {
		if p.Resource == "user" && p.Action == "read" {
			userReadPerm = append(userReadPerm, p)
		}
	}

	roles := []role.Role{
		{Name: "admin", Description: "Full access", Permissions: allPerms},
		{Name: "user",  Description: "Read-only",   Permissions: userReadPerm},
	}
	for i := range roles {
		existing := role.Role{}
		result := db.WithContext(ctx).Where("name = ?", roles[i].Name).First(&existing)
		if result.Error != nil {
			db.WithContext(ctx).Create(&roles[i])
			log.Printf("[Seed] role created: %s", roles[i].Name)
		} else {
			log.Printf("[Seed] role exists: %s", roles[i].Name)
		}
	}
}

func seedAdminUser(ctx context.Context, db *gorm.DB) {
	var adminRole role.Role
	if err := db.WithContext(ctx).Where("name = ?", "admin").First(&adminRole).Error; err != nil {
		log.Fatalf("[Seed] admin role not found: %v", err)
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte("Admin1234!"), bcrypt.DefaultCost)
	admin := user.User{
		Email:        "admin@localhost",
		PasswordHash: string(hash),
		FirstName:    "System",
		LastName:     "Admin",
		IsActive:     true,
		RoleID:       adminRole.ID,
	}

	existing := user.User{}
	result := db.WithContext(ctx).Where("email = ?", admin.Email).First(&existing)
	if result.Error != nil {
		db.WithContext(ctx).Create(&admin)
		log.Printf("[Seed] admin user created: %s", admin.Email)
	} else {
		log.Printf("[Seed] admin user exists: %s", admin.Email)
	}
}

func seedErrorMessages(ctx context.Context, db *gorm.DB) {
	msgs := []errmsg.ErrorMessage{
		{Code: "AUTH_INVALID_CREDENTIALS",  LocaleTH: "อีเมลหรือรหัสผ่านไม่ถูกต้อง",            LocaleEN: "Invalid email or password"},
		{Code: "AUTH_TOKEN_EXPIRED",        LocaleTH: "Token หมดอายุ กรุณาเข้าสู่ระบบใหม่",      LocaleEN: "Token expired, please login again"},
		{Code: "AUTH_TOKEN_INVALID",        LocaleTH: "Token ไม่ถูกต้อง",                        LocaleEN: "Invalid token"},
		{Code: "AUTH_UNAUTHORIZED",         LocaleTH: "ไม่มีสิทธิ์เข้าถึง",                      LocaleEN: "Unauthorized"},
		{Code: "AUTH_SESSION_REVOKED",      LocaleTH: "Session ถูกยกเลิก กรุณาเข้าสู่ระบบใหม่", LocaleEN: "Session revoked, please login again"},
		{Code: "VALIDATION_REQUIRED",       LocaleTH: "กรุณากรอกข้อมูลให้ครบถ้วน",               LocaleEN: "Required field missing"},
		{Code: "VALIDATION_INVALID_FORMAT", LocaleTH: "รูปแบบข้อมูลไม่ถูกต้อง",                  LocaleEN: "Invalid format"},
		{Code: "USER_NOT_FOUND",            LocaleTH: "ไม่พบผู้ใช้งาน",                           LocaleEN: "User not found"},
		{Code: "USER_ALREADY_EXISTS",       LocaleTH: "มีผู้ใช้งานนี้อยู่แล้ว",                   LocaleEN: "User already exists"},
		{Code: "USER_INACTIVE",             LocaleTH: "บัญชีผู้ใช้ถูกระงับการใช้งาน",             LocaleEN: "User account is inactive"},
		{Code: "ROLE_NOT_FOUND",            LocaleTH: "ไม่พบบทบาทผู้ใช้",                         LocaleEN: "Role not found"},
		{Code: "PERMISSION_DENIED",         LocaleTH: "ไม่มีสิทธิ์ดำเนินการนี้",                  LocaleEN: "Permission denied"},
		{Code: "INTERNAL_ERROR",            LocaleTH: "เกิดข้อผิดพลาดภายในระบบ",                 LocaleEN: "Internal server error"},
	}
	db.WithContext(ctx).Clauses(clause.OnConflict{DoNothing: true}).Create(&msgs)
	log.Printf("[Seed] error_messages: %d rows", len(msgs))
}
