package main

import (
	"context"
	"errors"
	"fmt"
	"log"

	"git.trovefin.com/poc/ctrlc-service-go/internal/config"
	"git.trovefin.com/poc/ctrlc-service-go/internal/domain/errmsg"
	"git.trovefin.com/poc/ctrlc-service-go/internal/domain/permission"
	"git.trovefin.com/poc/ctrlc-service-go/internal/domain/role"
	"git.trovefin.com/poc/ctrlc-service-go/internal/domain/user"
	"git.trovefin.com/poc/ctrlc-service-go/internal/infrastructure/database"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type deviceFixture struct {
	Email  string
	Device string
	Status string
}

const (
	deviceStatusActive   = "active"
	deviceStatusBlocked  = "blocked"
	deviceStatusInactive = "inactive"
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
	if err := seedDevices(ctx, db); err != nil {
		log.Fatalf("[Seed] devices: %v", err)
	}

	log.Println("[Seed] complete")
}

func seedPermissions(ctx context.Context, db *gorm.DB) {
	perms := []permission.Permission{
		{Resource: "user", Action: "create", Description: "Create users"},
		{Resource: "user", Action: "read", Description: "Read users"},
		{Resource: "user", Action: "update", Description: "Update users"},
		{Resource: "user", Action: "delete", Description: "Delete users"},
		{Resource: "role", Action: "create", Description: "Create roles"},
		{Resource: "role", Action: "read", Description: "Read roles"},
		{Resource: "role", Action: "update", Description: "Update roles"},
		{Resource: "role", Action: "delete", Description: "Delete roles"},
		{Resource: "audit", Action: "read", Description: "Read audit logs"},
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
		{Name: "user", Description: "Read-only", Permissions: userReadPerm},
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
		{Code: "AUTH_INVALID_CREDENTIALS", LocaleTH: "อีเมลหรือรหัสผ่านไม่ถูกต้อง", LocaleEN: "Invalid email or password"},
		{Code: "AUTH_TOKEN_EXPIRED", LocaleTH: "Token หมดอายุ กรุณาเข้าสู่ระบบใหม่", LocaleEN: "Token expired, please login again"},
		{Code: "AUTH_TOKEN_INVALID", LocaleTH: "Token ไม่ถูกต้อง", LocaleEN: "Invalid token"},
		{Code: "AUTH_UNAUTHORIZED", LocaleTH: "ไม่มีสิทธิ์เข้าถึง", LocaleEN: "Unauthorized"},
		{Code: "AUTH_SESSION_REVOKED", LocaleTH: "Session ถูกยกเลิก กรุณาเข้าสู่ระบบใหม่", LocaleEN: "Session revoked, please login again"},
		{Code: "VALIDATION_REQUIRED", LocaleTH: "กรุณากรอกข้อมูลให้ครบถ้วน", LocaleEN: "Required field missing"},
		{Code: "VALIDATION_INVALID_FORMAT", LocaleTH: "รูปแบบข้อมูลไม่ถูกต้อง", LocaleEN: "Invalid format"},
		{Code: "USER_NOT_FOUND", LocaleTH: "ไม่พบผู้ใช้งาน", LocaleEN: "User not found"},
		{Code: "USER_ALREADY_EXISTS", LocaleTH: "มีผู้ใช้งานนี้อยู่แล้ว", LocaleEN: "User already exists"},
		{Code: "USER_INACTIVE", LocaleTH: "บัญชีผู้ใช้ถูกระงับการใช้งาน", LocaleEN: "User account is inactive"},
		{Code: "ROLE_NOT_FOUND", LocaleTH: "ไม่พบบทบาทผู้ใช้", LocaleEN: "Role not found"},
		{Code: "PERMISSION_DENIED", LocaleTH: "ไม่มีสิทธิ์ดำเนินการนี้", LocaleEN: "Permission denied"},
		{Code: "INTERNAL_ERROR", LocaleTH: "เกิดข้อผิดพลาดภายในระบบ", LocaleEN: "Internal server error"},
		{Code: "DEVICE_NOT_FOUND", LocaleTH: "ไม่พบอุปกรณ์นี้ในระบบ", LocaleEN: "Device not found"},
		{Code: "DEVICE_BLOCKED", LocaleTH: "อุปกรณ์นี้ถูกระงับการใช้งาน", LocaleEN: "Device is blocked"},
	}
	db.WithContext(ctx).Clauses(clause.OnConflict{DoNothing: true}).Create(&msgs)
	log.Printf("[Seed] error_messages: %d rows", len(msgs))
}

func seedDevices(ctx context.Context, db *gorm.DB) error {
	entries := []deviceFixture{
		{Email: "admin@localhost", Device: "a1b2c3d4-e5f6-4a7b-8c9d-0e1f2a3b4c5d", Status: deviceStatusActive},
		{Email: "blocked@localhost", Device: "b2c3d4e5-f6a7-4b8c-9d0e-1f2a3b4c5d6e", Status: deviceStatusBlocked},
		{Email: "inactive@localhost", Device: "c3d4e5f6-a7b8-4c9d-0e1f-2a3b4c5d6e7f", Status: deviceStatusInactive},
	}
	if ok, err := hasTable(ctx, db, "devices"); err != nil {
		return fmt.Errorf("check devices table: %w", err)
	} else if !ok {
		log.Printf("[Seed] devices table not present; skipping auth device fixtures")
		return nil
	}

	pinHash, err := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("generate test PIN hash: %w", err)
	}

	var adminRole role.Role
	if err := db.WithContext(ctx).Where("name = ?", "admin").First(&adminRole).Error; err != nil {
		return fmt.Errorf("admin role not found: %w", err)
	}

	for _, e := range entries {
		var u user.User
		if err := db.WithContext(ctx).Where("email = ?", e.Email).First(&u).Error; err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("find fixture user %s: %w", e.Email, err)
			}
			u = user.User{
				Email:        e.Email,
				PasswordHash: string(pinHash),
				FirstName:    "Test",
				LastName:     e.Status,
				IsActive:     true,
				RoleID:       adminRole.ID,
			}
			if err := db.WithContext(ctx).Create(&u).Error; err != nil {
				return fmt.Errorf("create fixture user %s for device %s: %w", e.Email, e.Device, err)
			}
			log.Printf("[Seed] test user created: %s", e.Email)
		}

		result := db.WithContext(ctx).Exec(`
INSERT INTO devices (device_uid, user_id, pin_hash, status)
VALUES (?, ?, ?, ?)
ON CONFLICT (device_uid) DO UPDATE
SET user_id = EXCLUDED.user_id,
    pin_hash = EXCLUDED.pin_hash,
    status = EXCLUDED.status,
    updated_at = NOW()`,
			e.Device, u.ID, string(pinHash), e.Status,
		)
		if result.Error != nil {
			return fmt.Errorf("upsert fixture device %s: %w", e.Device, result.Error)
		}
		log.Printf("[Seed] device upserted: %s (status=%s)", e.Device, e.Status)
	}
	return verifySeedDevices(ctx, db, entries)
}

func verifySeedDevices(ctx context.Context, db *gorm.DB, entries []deviceFixture) error {
	for _, e := range entries {
		var status string
		var email string
		err := db.WithContext(ctx).Raw(`
SELECT d.status, u.email
FROM devices d
JOIN users u ON u.id = d.user_id
WHERE d.device_uid = ?`, e.Device).Row().Scan(&status, &email)
		if err != nil {
			return fmt.Errorf("verify fixture device %s: %w", e.Device, err)
		}
		if status != e.Status {
			return fmt.Errorf("verify fixture device %s: status=%s, want %s", e.Device, status, e.Status)
		}
		if email != e.Email {
			return fmt.Errorf("verify fixture device %s: owner=%s, want %s", e.Device, email, e.Email)
		}
	}
	log.Printf("[Seed] verified auth device fixtures: %d rows", len(entries))
	return nil
}

func hasTable(ctx context.Context, db *gorm.DB, table string) (bool, error) {
	var exists bool
	err := db.WithContext(ctx).Raw(`SELECT to_regclass(?) IS NOT NULL`, "public."+table).Row().Scan(&exists)
	return exists, err
}
