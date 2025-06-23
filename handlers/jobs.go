package handlers

import (
	"job-api/config"
	"job-api/models"
	"job-api/utils"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateJobRequest struct {
	Title       string `json:"title" validate:"required,min=1,max=100"`
	Description string `json:"description" validate:"required,min=20,max=2000"`
	Location    string `json:"location"`
}

type UpdateJobRequest struct {
	Title       string `json:"title" validate:"required,min=1,max=100"`
	Description string `json:"description" validate:"required,min=20,max=2000"`
	Location    string `json:"location"`
}

func CreateJob(c *gin.Context) {
	var req CreateJobRequest
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
	createdBy := userID.(uuid.UUID)

	job := models.Job{
		Title:       req.Title,
		Description: req.Description,
		Location:    req.Location,
		CreatedBy:   createdBy,
	}

	if err := config.DB.Create(&job).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.BaseResponse{
			Success: false,
			Message: "Failed to create job",
			Object:  nil,
			Errors:  []string{err.Error()},
		})
		return
	}

	// Load creator information
	config.DB.Preload("Creator").First(&job, job.ID)

	c.JSON(http.StatusCreated, models.BaseResponse{
		Success: true,
		Message: "Job created successfully",
		Object:  job,
	})
}

func UpdateJob(c *gin.Context) {
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

	var req UpdateJobRequest
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

	job.Title = req.Title
	job.Description = req.Description
	job.Location = req.Location

	if err := config.DB.Save(&job).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.BaseResponse{
			Success: false,
			Message: "Failed to update job",
			Object:  nil,
			Errors:  []string{err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, models.BaseResponse{
		Success: true,
		Message: "Job updated successfully",
		Object:  job,
	})
}

func DeleteJob(c *gin.Context) {
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

	userID, _ := c.Get("user_id")
	currentUserID := userID.(uuid.UUID)

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

	if err := config.DB.Delete(&job).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.BaseResponse{
			Success: false,
			Message: "Failed to delete job",
			Object:  nil,
			Errors:  []string{err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, models.BaseResponse{
		Success: true,
		Message: "Job deleted successfully",
		Object:  nil,
	})
}

func BrowseJobs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	title := c.Query("title")
	location := c.Query("location")
	companyName := c.Query("company_name")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	query := config.DB.Model(&models.Job{}).Preload("Creator")

	if title != "" {
		query = query.Where("LOWER(title) LIKE ?", "%"+strings.ToLower(title)+"%")
	}
	if location != "" {
		query = query.Where("LOWER(location) LIKE ?", "%"+strings.ToLower(location)+"%")
	}
	if companyName != "" {
		query = query.Joins("JOIN users ON jobs.created_by = users.id").
			Where("LOWER(users.name) LIKE ?", "%"+strings.ToLower(companyName)+"%")
	}

	var total int64
	query.Count(&total)

	var jobs []models.Job
	if err := query.Offset(offset).Limit(pageSize).Find(&jobs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.BaseResponse{
			Success: false,
			Message: "Failed to fetch jobs",
			Object:  nil,
			Errors:  []string{err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, models.PaginatedResponse{
		Success:    true,
		Message:    "Jobs retrieved successfully",
		Object:     jobs,
		PageNumber: page,
		PageSize:   pageSize,
		TotalSize:  total,
	})
}

func GetJobDetails(c *gin.Context) {
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

	var job models.Job
	if err := config.DB.Preload("Creator").First(&job, jobUUID).Error; err != nil {
		c.JSON(http.StatusNotFound, models.BaseResponse{
			Success: false,
			Message: "Job not found",
			Object:  nil,
		})
		return
	}

	c.JSON(http.StatusOK, models.BaseResponse{
		Success: true,
		Message: "Job details retrieved successfully",
		Object:  job,
	})
}

func GetMyJobs(c *gin.Context) {
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

	var total int64
	config.DB.Model(&models.Job{}).Where("created_by = ?", currentUserID).Count(&total)

	var jobs []models.Job
	if err := config.DB.Where("created_by = ?", currentUserID).
		Offset(offset).Limit(pageSize).Find(&jobs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.BaseResponse{
			Success: false,
			Message: "Failed to fetch jobs",
			Object:  nil,
			Errors:  []string{err.Error()},
		})
		return
	}

	// Add application count for each job
	for i := range jobs {
		var count int64
		config.DB.Model(&models.Application{}).Where("job_id = ?", jobs[i].ID).Count(&count)
		// You can add this to a custom response struct if needed
	}

	c.JSON(http.StatusOK, models.PaginatedResponse{
		Success:    true,
		Message:    "Jobs retrieved successfully",
		Object:     jobs,
		PageNumber: page,
		PageSize:   pageSize,
		TotalSize:  total,
	})
}
