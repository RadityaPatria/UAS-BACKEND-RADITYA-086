package services

import (
    "context"
    "time" // ← WAJIB untuk time.Now()

    "UAS-backend/app/models"       // ← WAJIB untuk models.Student
    "UAS-backend/app/repositories"

    "github.com/gofiber/fiber/v2"
    "github.com/google/uuid"
)
 
// GET /api/v1/students
func GetAllStudents(c *fiber.Ctx) error {
    ctx := context.Background()

    list, err := repositories.GetAllStudents(ctx)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }

    return c.JSON(fiber.Map{"status": "success", "data": list})
}

// GET /api/v1/students/:id
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

    return c.JSON(fiber.Map{"status": "success", "data": student})
}

// GET /api/v1/students/:id/achievements
func GetStudentAchievements(c *fiber.Ctx) error {
    ctx := context.Background()
    id := c.Params("id")

    sid, err := uuid.Parse(id)
    if err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
    }

    list, err := repositories.GetAchievementReferencesByStudentIDs(ctx, []string{sid.String()})
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }

    return c.JSON(fiber.Map{"status": "success", "data": list})
}

// PUT /api/v1/students/:id/advisor
func UpdateStudentAdvisor(c *fiber.Ctx) error {
    ctx := context.Background()
    id := c.Params("id")

    var body struct {
        AdvisorID string `json:"advisor_id"`
    }
    if err := c.BodyParser(&body); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
    }

    sid, _ := uuid.Parse(id)
    aid, _ := uuid.Parse(body.AdvisorID)

    err := repositories.UpdateAdvisor(ctx, sid, aid)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }

    return c.JSON(fiber.Map{"status": "success", "message": "advisor updated"})
}

// POST /api/v1/students
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
        CreatedAt:    time.Now(), // ← FIX
    }

    if err := repositories.CreateStudent(ctx, &student); err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }

    return c.JSON(fiber.Map{
        "status": "success",
        "data":   student,
    })
}
