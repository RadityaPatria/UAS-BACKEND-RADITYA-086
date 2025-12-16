package services

import (
	"context"
	"time"

	"UAS-backend/app/models"
	"UAS-backend/app/repositories"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)


// @Summary      Get all students
// @Description  Ambil semua data mahasiswa
// @Tags         Students
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]string
// @Router       /students [get]
//
// GetAllStudents -> ambil semua mahasiswa | FR-010
func GetAllStudents(c *fiber.Ctx) error {
	ctx := context.Background()

	list, err := repositories.GetAllStudents(ctx)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   list,
	})
}


// @Summary      Get student by ID
// @Description  Ambil detail mahasiswa berdasarkan ID
// @Tags         Students
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      string  true  "Student ID (UUID)"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /students/{id} [get]
//
// GetStudentByID -> ambil detail mahasiswa | FR-010
func GetStudentByID(c *fiber.Ctx) error {
	ctx := context.Background()
	id := c.Params("id")

	sid, err := uuid.Parse(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}

	student, err := repositories.GetStudentByID(ctx, sid)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "student not found"})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   student,
	})
}


// @Summary      Get student achievements
// @Description  Ambil daftar prestasi mahasiswa
// @Tags         Students
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      string  true  "Student ID (UUID)"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /students/{id}/achievements [get]
//
// GetStudentAchievements -> ambil prestasi mahasiswa | FR-011
func GetStudentAchievements(c *fiber.Ctx) error {
	ctx := context.Background()
	id := c.Params("id")

	sid, err := uuid.Parse(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}

	list, err := repositories.GetAchievementReferencesByStudentIDs(
		ctx,
		[]string{sid.String()},
	)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   list,
	})
}


// @Summary      Update student advisor
// @Description  Update dosen wali mahasiswa
// @Tags         Students
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id    path  string  true  "Student ID (UUID)"
// @Param        body  body  map[string]string  true  "advisor_id"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /students/{id}/advisor [put]
//
// UpdateStudentAdvisor -> update dosen wali mahasiswa | FR-012
func UpdateStudentAdvisor(c *fiber.Ctx) error {
	ctx := context.Background()
	id := c.Params("id")

	var body struct {
		AdvisorID string `json:"advisor_id"`
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
	}

	sid, err := uuid.Parse(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid student id"})
	}

	aid, err := uuid.Parse(body.AdvisorID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid advisor id"})
	}

	if err := repositories.UpdateAdvisor(ctx, sid, aid); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "advisor updated",
	})
}

// @Summary      Create student
// @Description  Tambah mahasiswa baru
// @Tags         Students
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body  body  map[string]string  true  "user_id, student_id, program_study, academic_year, advisor_id"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /students [post]
//
// CreateStudent -> tambah mahasiswa baru | FR-010
func CreateStudent(c *fiber.Ctx) error {
	ctx := context.Background()

	var req struct {
		UserID       string `json:"user_id"`
		StudentID    string `json:"student_id"`
		ProgramStudy string `json:"program_study"`
		AcademicYear string `json:"academic_year"`
		AdvisorID    string `json:"advisor_id"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
	}

	student := models.Student{
		ID:           uuid.New(),
		UserID:       uuid.MustParse(req.UserID),
		StudentID:    req.StudentID,
		ProgramStudy: req.ProgramStudy,
		AcademicYear: req.AcademicYear,
		AdvisorID:    uuid.MustParse(req.AdvisorID),
		CreatedAt:    time.Now(),
	}

	if err := repositories.CreateStudent(ctx, &student); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   student,
	})
}
