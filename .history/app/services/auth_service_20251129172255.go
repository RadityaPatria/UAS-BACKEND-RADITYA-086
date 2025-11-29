package services

import (
    "context"
    "errors"
    "time"

    "UAS-backend/app/models"
    "UAS-backend/app/repositories"
    "UAS-backend/config"

    "github.com/golang-jwt/jwt/v5"
    "golang.org/x/crypto/bcrypt"
)

type AuthService struct {
    UserRepo       *repositories.UserRepository
    RoleRepo       *repositories.RoleRepository
    PermissionRepo *repositories.PermissionRepository
    Cfg            *config.Config
}

// COMMAND: Constructor untuk membuat instance service auth
func NewAuthService(
    userRepo *repositories.UserRepository,
    roleRepo *repositories.RoleRepository,
    permRepo *repositories.PermissionRepository,
    cfg *config.Config,
) *AuthService {
    return &AuthService{
        UserRepo:       userRepo,
        RoleRepo:       roleRepo,
        PermissionRepo: permRepo,
        Cfg:            cfg,
    }
}

// COMMAND: Login utama (validasi kredensial)
func (s *AuthService) Login(ctx context.Context, usernameOrEmail, password string) (*models.User, string, error) {
    // 1. Ambil user berdasarkan username / email
    user, err := s.UserRepo.FindByUsernameOrEmail(ctx, usernameOrEmail)
    if err != nil {
        return nil, "", errors.New("user tidak ditemukan")
    }

    // 2. Cek password
    if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
        return nil, "", errors.New("password salah")
    }

    // 3. Cek status aktif
    if !user.IsActive {
        return nil, "", errors.New("akun tidak aktif")
    }

    // 4. Load permissions
    perms, err := s.PermissionRepo.GetPermissionsByRole(ctx, user.RoleID)
    if err != nil {
        return nil, "", errors.New("gagal load permissions")
    }

    // 5. Generate token
    token, err := s.generateJWT(user, perms)
    if err != nil {
        return nil, "", errors.New("gagal buat token")
    }

    return user, token, nil
}

// COMMAND: Generate JWT dengan role + permission
func (s *AuthService) generateJWT(user *models.User, permissions []string) (string, error) {
    claims := jwt.MapClaims{
        "sub":   user.ID.String(), // user id
        "role":  user.RoleID.String(),
        "perms": permissions,
        "exp":   time.Now().Add(time.Hour * 5).Unix(),
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    signed, err := token.SignedString([]byte(s.Cfg.JWTSecret))
    if err != nil {
        return "", err
    }

    return signed, nil
}

// COMMAND: Ambil profile dari user id (untuk GET /profile)
func (s *AuthService) GetProfile(ctx context.Context, userID string) (*models.User, error) {
    return s.UserRepo.FindByID(ctx, userID)
}
