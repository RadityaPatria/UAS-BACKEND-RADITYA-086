package services

import (
	"context"

	"UAS-backend/app/models"
	"UAS-backend/app/repositories"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)


// ====================================================================
// GET /api/v1/users
// ====================================================================
func GetAllUsers(c *fiber.Ctx) error {
	ctx := context.Background()

	users, err := repositories.GetAllUsers(ctx)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   users,
	})
}


// ====================================================================
// GET /api/v1/users/:id
// ====================================================================
func GetUserByID(c *fiber.Ctx) error {
	ctx := context.Background()
	id := c.Params("id")

	user, err := repositories.GetUserByID(ctx, id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "user not found",
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   user,
	})
}


// ====================================================================
// POST /api/v1/users
// ====================================================================
func CreateUser(c *fiber.Ctx) error {
	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
		FullName string `json:"full_name"`
		RoleID   string `json:"role_id"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	ctx := context.Background()

	// Hash Password
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "failed to hash password",
		})
	}

	// Create user struct
	user := &models.User{
		ID:           uuid.New(),
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hash),
		FullName:     req.FullName,
		RoleID:       uuid.MustParse(req.RoleID),
		IsActive:     true,
	}

	// Insert to repository
	if err := repositories.CreateUser(ctx, user); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   user,
	})
}


// ====================================================================
// PUT /api/v1/users/:id
// ====================================================================
func UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")

	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		FullName string `json:"full_name"`
		RoleID   string `json:"role_id"`
		IsActive bool   `json:"is_active"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	ctx := context.Background()

	// Find user
	user, err := repositories.GetUserByID(ctx, id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "user not found",
		})
	}

	// Update fields
	user.Username = req.Username
	user.Email = req.Email
	user.FullName = req.FullName
	user.RoleID = uuid.MustParse(req.RoleID)
	user.IsActive = req.IsActive

	// Save update
	if err := repositories.UpdateUser(ctx, user); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   user,
	})
}


// ====================================================================
// DELETE /api/v1/users/:id
// ====================================================================
func DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	ctx := context.Background()

	err := repositories.DeleteUser(ctx, id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

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
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	ctx := context.Background()

	user, err := repositories.GetUserByID(ctx, id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "user not found",
		})
	}

	user.RoleID = uuid.MustParse(req.RoleID)

	if err := repositories.UpdateUser(ctx, user); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "role updated",
	})
}
