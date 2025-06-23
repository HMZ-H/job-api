package handlers

import (
	"job-api/config"
	"job-api/models"
	"job-api/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ApplyJobRequest struct {
	ResumeLink  string `json:"resume_link" validate:"required,url"`
	CoverLetter string `json:"cover_letter" validate:"max=200"`
}

type UpdateApplicationStatusRequest struct {
	Status models.ApplicationStatus `json:"status" validate:"required,oneof=Applied Reviewed Interview Rejected Hired"`
}

func ApplyForJob(c *gin.Context) {
	jobID := c.Param("id")
	jobUUID, err := uuid.Parse(jobID)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.BaseResponse{
			Success: false,
			Message: "Invalid job ID",
			Object:  nil,
		})
		return
	}

	var req ApplyJobRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.BaseResponse{
			Success: false,
			Message: "Invalid request data",
			Object:  nil,
			Errors:  []string{err.Error()},
		})
		return
	}

	if err := utils.ValidateStruct(req); err != nil {
		c.JSON(http.StatusBadRequest, models.BaseResponse{
			Success: false,
			Message: "Validation failed",
			Object:  nil,
			Errors:  []string{err.Error()},
		})
		return
	}

	userID, _ := c.Get("user_id")
	applicantID := userID.(uuid.UUID)

	// Check if job exists
	var job models.Job
	if err := config.DB.First(&job, jobUUID).Error; err != nil {
		c.JSON(http.StatusNotFound, models.BaseResponse{
			Success: false,
			Message: "Job not found",
			Object:  nil,
		})
		return
	}

	// Check if user already applied
	var existingApplication models.Application
	if err := config.DB.Where("applicant_id = ? AND job_id = ?", applicantID, jobUUID).
		First(&existingApplication).Error; err == nil {
		c.JSON(http.StatusConflict, models.BaseResponse{
			Success: false,
			Message: "You have already applied to this job",
			Object:  nil,
			Errors:  []string{"Duplicate application"},
		})
		return
	}

	application := models.Application{
		ApplicantID: applicantID,
		JobID:       jobUUID,
		ResumeLink:  req.ResumeLink,
		CoverLetter: req.CoverLetter,
		Status:      models.StatusApplied,
	}

	if err := config.DB.Create(&application).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.BaseResponse{
			Success: false,
			Message: "Failed to submit application",
			Object:  nil,
			Errors:  []string{err.Error()},
		})
		return
	}

	// Load relationships
	config.DB.Preload("Applicant").Preload("Job").First(&application, application.ID)

	c.JSON(http.StatusCreated, models.BaseResponse{
		Success: true,
		Message: "Application submitted successfully",
		Object:  application,
	})
}

func GetMyApplications(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	userID, _ := c.Get("user_id")
	applicantID := userID.(uuid.UUID)

	var total int64
	config.DB.Model(&models.Application{}).Where("applicant_id = ?", applicantID).Count(&total)

	var applications []models.Application
	if err := config.DB.Where("applicant_id = ?", applicantID).
		Preload("Job").Preload("Job.Creator").
		Offset(offset).Limit(pageSize).Find(&applications).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.BaseResponse{
			Success: false,
			Message: "Failed to fetch applications",
			Object:  nil,
			Errors:  []string{err.Error()},
		})
		return
	}

	// Transform response to include required fields
	type ApplicationResponse struct {
		ID          uuid.UUID `json:"id"`
		JobTitle    string    `json:"job_title"`
		CompanyName string    `json:"company_name"`
		Status      string    `json:"status"`
		AppliedAt   string    `json:"applied_at"`
	}

	var response []ApplicationResponse
	for _, app := range applications {
		response = append(response, ApplicationResponse{
			ID:          app.ID,
			JobTitle:    app.Job.Title,
			CompanyName: app.Job.Creator.Name,
			Status:      string(app.Status),
			AppliedAt:   app.AppliedAt.Format("2006-01-02 15:04:05"),
		})
	}

	c.JSON(http.StatusOK, models.PaginatedResponse{
		Success:    true,
		Message:    "Applications retrieved successfully",
		Object:     response,
		PageNumber: page,
		PageSize:   pageSize,
		TotalSize:  total,
	})
}

func GetJobApplications(c *gin.Context) {
	jobID := c.Param("id")
	jobUUID, err := uuid.Parse(jobID)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.BaseResponse{
			Success: false,
			Message: "Invalid job ID",
			Object:  nil,
		})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	userID, _ := c.Get("user_id")
	currentUserID := userID.(uuid.UUID)

	// Check if job exists and belongs to current user
	var job models.Job
	if err := config.DB.First(&job, jobUUID).Error; err != nil {
		c.JSON(http.StatusNotFound, models.BaseResponse{
			Success: false,
			Message: "Job not found",
			Object:  nil,
		})
		return
	}

	if job.CreatedBy != currentUserID {
		c.JSON(http.StatusForbidden, models.BaseResponse{
			Success: false,
			Message: "Unauthorized access",
			Object:  nil,
		})
		return
	}

	var total int64
	config.DB.Model(&models.Application{}).Where("job_id = ?", jobUUID).Count(&total)

	var applications []models.Application
	if err := config.DB.Where("job_id = ?", jobUUID).
		Preload("Applicant").
		Offset(offset).Limit(pageSize).Find(&applications).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.BaseResponse{
			Success: false,
			Message: "Failed to fetch applications",
			Object:  nil,
			Errors:  []string{err.Error()},
		})
		return
	}

	// Transform response to include required fields
	type ApplicationResponse struct {
		ID            uuid.UUID `json:"id"`
		ApplicantName string    `json:"applicant_name"`
		ResumeLink    string    `json:"resume_link"`
		CoverLetter   string    `json:"cover_letter"`
		Status        string    `json:"status"`
		AppliedAt     string    `json:"applied_at"`
	}

	var response []ApplicationResponse
	for _, app := range applications {
		response = append(response, ApplicationResponse{
			ID:            app.ID,
			ApplicantName: app.Applicant.Name,
			ResumeLink:    app.ResumeLink,
			CoverLetter:   app.CoverLetter,
			Status:        string(app.Status),
			AppliedAt:     app.AppliedAt.Format("2006-01-02 15:04:05"),
		})
	}

	c.JSON(http.StatusOK, models.PaginatedResponse{
		Success:    true,
		Message:    "Applications retrieved successfully",
		Object:     response,
		PageNumber: page,
		PageSize:   pageSize,
		TotalSize:  total,
	})
}

func UpdateApplicationStatus(c *gin.Context) {
	applicationID := c.Param("id")
	appUUID, err := uuid.Parse(applicationID)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.BaseResponse{
			Success: false,
			Message: "Invalid application ID",
			Object:  nil,
		})
		return
	}

	var req UpdateApplicationStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.BaseResponse{
			Success: false,
			Message: "Invalid request data",
			Object:  nil,
			Errors:  []string{err.Error()},
		})
		return
	}

	if err := utils.ValidateStruct(req); err != nil {
		c.JSON(http.StatusBadRequest, models.BaseResponse{
			Success: false,
			Message: "Validation failed",
			Object:  nil,
			Errors:  []string{err.Error()},
		})
		return
	}

	userID, _ := c.Get("user_id")
	currentUserID := userID.(uuid.UUID)

	var application models.Application
	if err := config.DB.Preload("Job").First(&application, appUUID).Error; err != nil {
		c.JSON(http.StatusNotFound, models.BaseResponse{
			Success: false,
			Message: "Application not found",
			Object:  nil,
		})
		return
	}

	// Check if current user owns the job
	if application.Job.CreatedBy != currentUserID {
		c.JSON(http.StatusForbidden, models.BaseResponse{
			Success: false,
			Message: "Unauthorized",
			Object:  nil,
		})
		return
	}

	application.Status = req.Status
	if err := config.DB.Save(&application).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.BaseResponse{
			Success: false,
			Message: "Failed to update application status",
			Object:  nil,
			Errors:  []string{err.Error()},
		})
		return
	}

	// Load relationships for response
	config.DB.Preload("Applicant").Preload("Job").First(&application, application.ID)

	c.JSON(http.StatusOK, models.BaseResponse{
		Success: true,
		Message: "Application status updated successfully",
		Object:  application,
	})
}
