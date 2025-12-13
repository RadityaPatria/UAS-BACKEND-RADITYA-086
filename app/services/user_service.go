package services

import (
	"context"
	"errors"
	"time"

	"UAS-backend/app/models"
	"UAS-backend/app/repositories"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// GetAllUsers -> ambil semua user | FR-009
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

// GetUserByID -> ambil user berdasarkan ID | FR-002
func GetUserByID(c *fiber.Ctx) error {
	ctx := context.Background()
	id := c.Params("id")

	if _, err := uuid.Parse(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid uuid",
		})
	}

	user, err := repositories.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, errors.New("record not found")) {
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

// CreateUser -> tambah user + auto student / lecturer | FR-009
func CreateUser(c *fiber.Ctx) error {
	var req struct {
		Username     string `json:"username"`
		Email        string `json:"email"`
		Password     string `json:"password"`
		FullName     string `json:"full_name"`
		RoleID       string `json:"role_id"`

		// Student
		ProgramStudy string `json:"program_study"`
		AcademicYear string `json:"academic_year"`
		AdvisorID    string `json:"advisor_id"`

		// Lecturer
		LecturerID string `json:"lecturer_id"`
		Department string `json:"department"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request body"})
	}

	ctx := context.Background()

	roleUUID, err := uuid.Parse(req.RoleID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid role_id"})
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "password hashing failed"})
	}

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

	// Auto create student
	const ROLE_MAHASISWA = "c65ba8ef-8302-430a-9d5e-fa1438fa0eff"
	if user.RoleID.String() == ROLE_MAHASISWA {

		if req.AdvisorID == "" {
			return c.Status(400).JSON(fiber.Map{"error": "advisor_id required"})
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
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
	}

	// Auto create lecturer
	const ROLE_DOSEN = "8e9d2e54-aee2-4013-a21d-d3016e8c8295"
	if user.RoleID.String() == ROLE_DOSEN {

		if req.LecturerID == "" || req.Department == "" {
			return c.Status(400).JSON(fiber.Map{"error": "lecturer_id & department required"})
		}

		lecturer := models.Lecturer{
			ID:         uuid.New(),
			UserID:     user.ID,
			LecturerID: req.LecturerID,
			Department: req.Department,
			CreatedAt:  time.Now(),
		}

		if err := repositories.CreateLecturer(ctx, &lecturer); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   user,
	})
}

// UpdateUser -> update data user | FR-009
func UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")

	var req struct {
		Username *string `json:"username"`
		Email    *string `json:"email"`
		FullName *string `json:"full_name"`
		RoleID   *string `json:"role_id"`
		IsActive *bool   `json:"is_active"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request body"})
	}

	ctx := context.Background()

	user, err := repositories.GetUserByID(ctx, id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "user not found"})
	}

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
	if req.RoleID != nil {
		roleUUID, err := uuid.Parse(*req.RoleID)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "invalid role_id"})
		}
		user.RoleID = roleUUID
	}

	if err := repositories.UpdateUser(ctx, user); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   user,
	})
}

// DeleteUser -> hapus user | FR-009
func DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	ctx := context.Background()

	if _, err := uuid.Parse(id); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid uuid"})
	}

	if err := repositories.DeleteUser(ctx, id); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "user deleted",
	})
}

// UpdateUserRole -> update role user | FR-009
func UpdateUserRole(c *fiber.Ctx) error {
	id := c.Params("id")

	var req struct {
		RoleID string `json:"role_id"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request body"})
	}

	ctx := context.Background()

	user, err := repositories.GetUserByID(ctx, id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "user not found"})
	}

	roleUUID, err := uuid.Parse(req.RoleID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid role_id"})
	}

	user.RoleID = roleUUID

	if err := repositories.UpdateUser(ctx, user); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "role updated",
	})
}
