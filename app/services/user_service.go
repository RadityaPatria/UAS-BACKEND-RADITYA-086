package services

import (
	"context"
	"errors" // Diperlukan untuk penanganan error
	 "time"
	"UAS-backend/app/models"
	"UAS-backend/app/repositories"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// ====================================================================
// GET /api/v1/users (Fungsi yang hilang dari stack trace)
// ====================================================================
func GetAllUsers(c *fiber.Ctx) error {
	ctx := context.Background()

	users, err := repositories.GetAllUsers(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   users,
	})
}


// ====================================================================
// GET /api/v1/users/:id (Fungsi yang hilang dari stack trace)
// ====================================================================
func GetUserByID(c *fiber.Ctx) error {
	ctx := context.Background()
	id := c.Params("id")

    // Pastikan ID adalah UUID yang valid sebelum mencari
    if _, err := uuid.Parse(id); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "ID format is invalid. Must be a valid UUID.",
        })
    }

	user, err := repositories.GetUserByID(ctx, id)
	if err != nil {
		// Asumsi repository mengembalikan error spesifik jika tidak ditemukan
        if errors.Is(err, errors.New("record not found")) { // Ganti dengan error repositori spesifik Anda
            return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
                "error": "user not found",
            })
        }
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "user not found",
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   user,
	})
}


// ====================================================================
// POST /api/v1/users (Diperbaiki: Menambah validasi wajib)
// ====================================================================
func CreateUser(c *fiber.Ctx) error {
    var req struct {
        Username      string `json:"username"`
        Email         string `json:"email"`
        Password      string `json:"password"`
        FullName      string `json:"full_name"`
        RoleID        string `json:"role_id"`

        // Student only
        ProgramStudy  string `json:"program_study"`
        AcademicYear  string `json:"academic_year"`
        AdvisorID     string `json:"advisor_id"`

        // Lecturer only
        LecturerID    string `json:"lecturer_id"`
        Department    string `json:"department"`
    }

    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "invalid request body"})
    }

    ctx := context.Background()

    roleUUID, err := uuid.Parse(req.RoleID)
    if err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "invalid role_id"})
    }

    // Hash password
    hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "password hashing failed"})
    }

    // Create User
    user := &models.User{
        ID:           uuid.New(),
        Username:     req.Username,
        Email:        req.Email,
        PasswordHash: string(hash),
        FullName:     req.FullName,
        RoleID:       roleUUID,
        IsActive:     true,
    }

    if err := repositories.CreateUser(ctx, user); err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }

    // ===================================================
    // AUTO CREATE STUDENT
    // ===================================================
    const ROLE_MAHASISWA = "c65ba8ef-8302-430a-9d5e-fa1438fa0eff"

    if user.RoleID.String() == ROLE_MAHASISWA {

        if req.AdvisorID == "" {
            return c.Status(400).JSON(fiber.Map{"error": "advisor_id required for student"})
        }

        advisorUUID, err := uuid.Parse(req.AdvisorID)
        if err != nil {
            return c.Status(400).JSON(fiber.Map{"error": "invalid advisor_id"})
        }

        student := models.Student{
            ID:           uuid.New(),
            UserID:       user.ID,
            StudentID:    req.Username,
            ProgramStudy: req.ProgramStudy,
            AcademicYear: req.AcademicYear,
            AdvisorID:    advisorUUID,
            CreatedAt:    time.Now(),
        }

        if err := repositories.CreateStudent(ctx, &student); err != nil {
            return c.Status(500).JSON(fiber.Map{"error": "failed to create student: " + err.Error()})
        }
    }

    // ===================================================
    // AUTO CREATE LECTURER
    // ===================================================
    const ROLE_DOSEN = "8e9d2e54-aee2-4013-a21d-d3016e8c8295"

    if user.RoleID.String() == ROLE_DOSEN {

        if req.LecturerID == "" {
            return c.Status(400).JSON(fiber.Map{"error": "lecturer_id required for lecturer"})
        }

        if req.Department == "" {
            return c.Status(400).JSON(fiber.Map{"error": "department required for lecturer"})
        }

        lecturer := models.Lecturer{
            ID:         uuid.New(),
            UserID:     user.ID,
            LecturerID: req.LecturerID,
            Department: req.Department,
            CreatedAt:  time.Now(),
        }

        if err := repositories.CreateLecturer(ctx, &lecturer); err != nil {
            return c.Status(500).JSON(fiber.Map{"error": "failed to create lecturer: " + err.Error()})
        }
    }

    return c.JSON(fiber.Map{
        "status": "success",
        "data":   user,
    })
}



// ====================================================================
// PUT /api/v1/users/:id (Diperbaiki: Menggunakan pointer untuk opsional)
// ====================================================================
func UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")

	// Menggunakan pointer untuk field agar bisa membedakan antara tidak dikirim dan string kosong
	var req struct {
		Username *string `json:"username"`
		Email    *string `json:"email"`
		FullName *string `json:"full_name"`
		RoleID   *string `json:"role_id"`
		IsActive *bool   `json:"is_active"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	ctx := context.Background()

	// Find user
	user, err := repositories.GetUserByID(ctx, id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "user not found",
		})
	}
	
	// Update fields hanya jika field tersebut ada (bukan nil)
	if req.Username != nil {
		user.Username = *req.Username
	}
	if req.Email != nil {
		user.Email = *req.Email
	}
	if req.FullName != nil {
		user.FullName = *req.FullName
	}
	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}

	// 1. VALIDASI RoleID (Hanya jika RoleID dikirimkan)
	if req.RoleID != nil {
		newRoleID, err := uuid.Parse(*req.RoleID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "RoleID format is invalid. Must be a valid UUID.",
			})
		}
		user.RoleID = newRoleID
	}

	// Save update
	if err := repositories.UpdateUser(ctx, user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   user,
	})
}


// ====================================================================
// DELETE /api/v1/users/:id (Fungsi yang hilang dari stack trace)
// ====================================================================
func DeleteUser(c *fiber.Ctx) error {
	idString := c.Params("id")
	ctx := context.Background()

	// 1. VALIDASI UUID
	if _, err := uuid.Parse(idString); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID format is invalid. Must be a valid UUID.",
		})
	}
	
	// 2. LANJUTKAN KE REPOSITORY
	err := repositories.DeleteUser(ctx, idString)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// 3. RESPONS SUKSES
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "user deleted",
	})
}


// ====================================================================
// PUT /api/v1/users/:id/role
// ====================================================================
func UpdateUserRole(c *fiber.Ctx) error {
	id := c.Params("id")

	var req struct {
		RoleID string `json:"role_id"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	ctx := context.Background()

	user, err := repositories.GetUserByID(ctx, id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "user not found",
		})
	}

	// 1. VALIDASI RoleID (Mencegah panic)
	newRoleID, err := uuid.Parse(req.RoleID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "RoleID format is invalid. Must be a valid UUID.",
		})
	}

	user.RoleID = newRoleID

	if err := repositories.UpdateUser(ctx, user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "role updated",
	})
}