

import (
	"context"
	"log"

	"UAS-backend/app/models"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func RunSeeder() {
	ctx := context.Background()

	// ==============================
	// 1. CREATE ROLE ADMIN
	// ==============================
	adminRoleID := uuid.New()

	_, err := DB.Exec(ctx,
		`INSERT INTO roles (id, name, description, created_at)
		 VALUES ($1, 'admin', 'Administrator system', NOW())
		 ON CONFLICT DO NOTHING`,
		adminRoleID,
	)
	if err != nil {
		log.Println("❌ Failed create admin role:", err)
	} else {
		log.Println("✔️ Admin role ready")
	}

	// ==============================
	// 2. CREATE PERMISSIONS
	// ==============================
	permList := []models.Permission{
		{ID: uuid.New(), Name: "user:read", Resource: "user", Action: "read"},
		{ID: uuid.New(), Name: "user:create", Resource: "user", Action: "create"},
		{ID: uuid.New(), Name: "achievement:read", Resource: "achievement", Action: "read"},
		{ID: uuid.New(), Name: "achievement:verify", Resource: "achievement", Action: "verify"},
	}

	for _, p := range permList {
		_, err := DB.Exec(ctx,
			`INSERT INTO permissions (id, name, resource, action, description)
			 VALUES ($1,$2,$3,$4,'Seed permission')
			 ON CONFLICT DO NOTHING`,
			p.ID, p.Name, p.Resource, p.Action,
		)
		if err != nil {
			log.Println("❌ Failed inserting permission:", p.Name)
		}
	}

	log.Println("✔️ Permissions ready")

	// ==============================
	// 3. ASSIGN ALL PERMISSIONS TO ADMIN ROLE
	// ==============================
	for _, p := range permList {
		_, _ = DB.Exec(ctx,
			`INSERT INTO role_permissions (role_id, permission_id)
			 VALUES ($1,$2)
			 ON CONFLICT DO NOTHING`,
			adminRoleID, p.ID,
		)
	}

	log.Println("✔️ Permissions linked to Admin role")

	// ==============================
	// 4. CREATE ADMIN USER
	// ==============================
	adminUserID := uuid.New()

	hash, _ := bcrypt.GenerateFromPassword([]byte("admin123"), 10)

	_, err = DB.Exec(ctx,
		`INSERT INTO users
		(id, username, email, password_hash, full_name, role_id, is_active, created_at, updated_at)
		 VALUES ($1,'admin','admin@gmail.com',$2,'Administrator',$3,true,NOW(),NOW())
		 ON CONFLICT DO NOTHING`,
		adminUserID, string(hash), adminRoleID,
	)

	if err != nil {
		log.Println("❌ Failed creating admin user:", err)
	} else {
		log.Println("✔️ Admin user ready (username: admin, password: admin123)")
	}
}
