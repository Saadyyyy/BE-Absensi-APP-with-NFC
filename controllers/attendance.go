package controllers

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"myapp/config"
	"myapp/models"
)

type AttendanceController struct{}

type NFCAttendanceRequest struct {
	NFCUID string `json:"nfc_uid" validate:"required"`
}

type RegisterNFCRequest struct {
	NFCUID    string    `json:"nfc_uid" validate:"required"`
	Name      string    `json:"name" validate:"required"`
	Class     string    `json:"class" validate:"required"`
	StudentID string    `json:"student_id" validate:"required"`
	SchoolID  uuid.UUID `json:"school_id" validate:"required"`
}

// RecordAttendance records attendance using NFC card
func (ac *AttendanceController) RecordAttendance(c echo.Context) error {
	req := new(NFCAttendanceRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	// Find student by NFC UID
	var student models.Student
	result := config.DB.Where("nfc_uid = ? AND is_active = ?", req.NFCUID, true).First(&student)
	if result.Error != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Student not found or card not registered",
		})
	}

	// Get today's date
	today := time.Now().Truncate(24 * time.Hour)

	// Check if attendance already exists for today
	var attendance models.Attendance
	result = config.DB.Where("student_id = ? AND date = ?", student.ID, today).First(&attendance)

	if result.Error != nil {
		// Create new attendance record (check-in)
		now := time.Now()
		attendance = models.Attendance{
			ID:        uuid.New(),
			StudentID: student.ID,
			Date:      today,
			TimeIn:    &now,
			Status:    "present",
		}

		// Check if student is late (assuming school starts at 7:30 AM)
		schoolStartTime := time.Date(now.Year(), now.Month(), now.Day(), 7, 30, 0, 0, now.Location())
		if now.After(schoolStartTime) {
			attendance.Status = "late"
		}

		result = config.DB.Create(&attendance)
		if result.Error != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Failed to record attendance",
			})
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"message":    "Check-in successful",
			"student":    student.Name,
			"class":      student.Class,
			"time_in":    attendance.TimeIn,
			"status":     attendance.Status,
			"attendance": attendance,
		})
	} else {
		// Update existing attendance record (check-out)
		if attendance.TimeOut != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Student already checked out today",
			})
		}

		now := time.Now()
		attendance.TimeOut = &now

		result = config.DB.Save(&attendance)
		if result.Error != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Failed to update attendance",
			})
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"message":    "Check-out successful",
			"student":    student.Name,
			"class":      student.Class,
			"time_out":   attendance.TimeOut,
			"attendance": attendance,
		})
	}
}

// RegisterNFCCard registers a new NFC card for a student
func (ac *AttendanceController) RegisterNFCCard(c echo.Context) error {
	req := new(RegisterNFCRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	// Check if NFC UID already exists
	var existingStudent models.Student
	result := config.DB.Where("nfc_uid = ?", req.NFCUID).First(&existingStudent)
	if result.Error == nil {
		return c.JSON(http.StatusConflict, map[string]string{
			"error": "NFC card already registered",
		})
	}

	// Check if student ID already exists
	result = config.DB.Where("student_id = ?", req.StudentID).First(&existingStudent)
	if result.Error == nil {
		return c.JSON(http.StatusConflict, map[string]string{
			"error": "Student ID already exists",
		})
	}

	// Create new student
	student := models.Student{
		ID:        uuid.New(),
		NFCUID:    req.NFCUID,
		Name:      req.Name,
		Class:     req.Class,
		StudentID: req.StudentID,
		SchoolID:  req.SchoolID,
		IsActive:  true,
	}

	result = config.DB.Create(&student)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to register NFC card",
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "NFC card registered successfully",
		"student": student,
	})
}

// GetAttendanceHistory gets attendance history for a student
func (ac *AttendanceController) GetAttendanceHistory(c echo.Context) error {
	studentID := c.Param("student_id")
	uuid, err := uuid.Parse(studentID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid student ID",
		})
	}

	var attendances []models.Attendance
	result := config.DB.Where("student_id = ?", uuid).Order("date DESC").Limit(30).Find(&attendances)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch attendance history",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"attendances": attendances,
	})
}

// GetTodayAttendance gets today's attendance for all students
func (ac *AttendanceController) GetTodayAttendance(c echo.Context) error {
	today := time.Now().Truncate(24 * time.Hour)

	var attendances []models.Attendance
	result := config.DB.Preload("Student").Where("date = ?", today).Find(&attendances)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch today's attendance",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"date":        today,
		"attendances": attendances,
		"total":       len(attendances),
	})
}