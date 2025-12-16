package services

import (
	"context"
	"errors"
	"strings"
	"time"

	"UAS-backend/app/models"
	"UAS-backend/app/repositories"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte("SECRET-JWT-UAS-BACKEND")


// LoginHandler godoc
// @Summary      Login user
// @Description  Login menggunakan username atau email dan password
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body  body  object{identifier=string,password=string}  true  "Login payload"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Router       /auth/login [post]
//
// LoginHandler -> login user & generate JWT | FR-001
func LoginHandler(c *fiber.Ctx) error {
	var req struct {
		Identifier string `json:"identifier"`
		Password   string `json:"password"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid body",
		})
	}

	ctx := context.Background()
	resp, err := Login(ctx, req.Identifier, req.Password)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   resp,
	})
}

// Login -> validasi user, role, permission & JWT | FR-001
// (INTERNAL FUNCTION - TIDAK PERLU SWAGGER)
func Login(ctx context.Context, identifier, password string) (interface{}, error) {
	user, err := repositories.GetUserByIdentifier(ctx, identifier)
	if err != nil {
		return nil, errors.New("invalid username or password")
	}

	if !user.IsActive {
		return nil, errors.New("user inactive")
	}

	if bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(password),
	) != nil {
		return nil, errors.New("invalid username or password")
	}

	role, _ := repositories.GetRoleByID(ctx, user.RoleID.String())
	permIDs, _ := repositories.GetPermissionIDsByRoleID(ctx, user.RoleID.String())
	perms, _ := repositories.GetPermissionsByIDs(ctx, permIDs)

	token, err := generateJWT(user, role.Name, perms)
	if err != nil {
		return nil, err
	}

	cleanPerms := []string{}
	for _, p := range perms {
		cleanPerms = append(cleanPerms, p.Name)
	}

	return fiber.Map{
		"token": token,
		"user": fiber.Map{
			"id":          user.ID.String(),
			"username":    user.Username,
			"fullName":    user.FullName,
			"role":        role.Name,
			"permissions": cleanPerms,
		},
	}, nil
}

// RefreshTokenHandler godoc
// @Summary      Refresh JWT token
// @Description  Generate token baru menggunakan token lama (tanpa body)
// @Tags         Auth
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Router       /auth/refresh [post]
//
// RefreshTokenHandler -> refresh JWT tanpa body | FR-002
func RefreshTokenHandler(c *fiber.Ctx) error {
	auth := c.Get("Authorization")
	if auth == "" {
		return c.Status(401).JSON(fiber.Map{
			"error": "missing token",
		})
	}

	tokenStr := strings.Replace(auth, "Bearer ", "", 1)

	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return c.Status(401).JSON(fiber.Map{
			"error": "invalid token",
		})
	}

	claims := token.Claims.(jwt.MapClaims)

	userID := claims["user_id"].(string)
	tokenVersion := int(claims["token_version"].(float64))

	ctx := context.Background()
	user, err := repositories.GetUserByID(ctx, userID)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{
			"error": "user not found",
		})
	}

	if user.TokenVersion != tokenVersion {
		return c.Status(401).JSON(fiber.Map{
			"error": "token revoked",
		})
	}

	user.TokenVersion++
	_ = repositories.UpdateTokenVersion(ctx, user.ID.String(), user.TokenVersion)

	role, _ := repositories.GetRoleByID(ctx, user.RoleID.String())
	permIDs, _ := repositories.GetPermissionIDsByRoleID(ctx, user.RoleID.String())
	perms, _ := repositories.GetPermissionsByIDs(ctx, permIDs)

	newToken, _ := generateJWT(user, role.Name, perms)

	return c.JSON(fiber.Map{
		"status": "success",
		"token":  newToken,
	})
}

// LogoutHandler godoc
// @Summary      Logout user
// @Description  Invalidate JWT token (logout)
// @Tags         Auth
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  map[string]string
// @Router       /auth/logout [post]
//
// LogoutHandler -> invalidate token user | FR-003
func LogoutHandler(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	ctx := context.Background()
	user, _ := repositories.GetUserByID(ctx, userID)

	user.TokenVersion++
	_ = repositories.UpdateTokenVersion(ctx, user.ID.String(), user.TokenVersion)

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "logged out",
	})
}

// GetProfileHandler godoc
// @Summary      Get user profile
// @Description  Ambil informasi user dari JWT
// @Tags         Auth
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Router       /auth/profile [get]
//
// GetProfileHandler -> ambil data user dari JWT | FR-004
func GetProfileHandler(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"user_id":     c.Locals("userID"),
		"role":        c.Locals("role"),
		"permissions": c.Locals("permissions"),
	})
}

// generateJWT -> generate token JWT | INTERNAL (tidak perlu swagger)
func generateJWT(user *models.User, role string, perms []models.Permission) (string, error) {
	ps := []string{}
	for _, p := range perms {
		ps = append(ps, p.Name)
	}

	claims := jwt.MapClaims{
		"user_id":       user.ID.String(),
		"role":          role,
		"permissions":   ps,
		"token_version": user.TokenVersion,
		"exp":           time.Now().Add(24 * time.Hour).Unix(),
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(jwtSecret)
}
