package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/Akashpg-M/polaris/backend/internal/adapter/repository"
	"github.com/Akashpg-M/polaris/backend/internal/application/orchestration"
	"github.com/Akashpg-M/polaris/backend/internal/core/command"
	"github.com/Akashpg-M/polaris/backend/internal/core/extension"
	taskcore "github.com/Akashpg-M/polaris/backend/internal/core/task"
	"github.com/Akashpg-M/polaris/backend/internal/modules/mobility/routing"
	"github.com/gin-gonic/gin"
)

type OrchestrationAPI struct {
	store   *repository.RegistryStore
	service *orchestration.Service
	metrics *orchestration.Metrics
}

func NewOrchestrationAPI(store *repository.RegistryStore, service *orchestration.Service, metrics *orchestration.Metrics) *OrchestrationAPI {
	return &OrchestrationAPI{store: store, service: service, metrics: metrics}
}

func (a *OrchestrationAPI) Register(r *gin.RouterGroup, registryAPI *RegistryAPI) {
	r.POST("/tasks", registryAPI.Middleware("orchestrate"), a.createTask)
	r.GET("/tasks", registryAPI.Middleware("read"), a.listTasks)
	r.GET("/tasks/:task_id", registryAPI.Middleware("read"), a.getTask)
	r.POST("/tasks/:task_id/cancel", registryAPI.Middleware("orchestrate"), a.cancelTask)
	r.POST("/tasks/:task_id/retry", registryAPI.Middleware("admin_retry"), a.retryTask)
	r.GET("/commands", registryAPI.Middleware("read"), a.listCommands)
	r.GET("/commands/:command_id", registryAPI.Middleware("read"), a.getCommand)
	r.POST("/commands/:command_id/retry", registryAPI.Middleware("admin_retry"), a.retryCommand)
	r.POST("/commands/:command_id/cancel", registryAPI.Middleware("orchestrate"), a.cancelCommand)
	r.GET("/metrics/orchestration", registryAPI.Middleware("read"), func(c *gin.Context) { apiData(c, 200, a.metrics.Snapshot()) })
}

func orchestrationError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, orchestration.ErrInvalidTask), errors.Is(err, orchestration.ErrUnsupportedCommand):
		apiError(c, http.StatusBadRequest, "INVALID_TASK", err.Error())
	case errors.Is(err, orchestration.ErrNoEligibleDevice):
		apiError(c, http.StatusConflict, "NO_ELIGIBLE_DEVICE", err.Error())
	case errors.Is(err, routing.ErrBusy):
		apiError(c, http.StatusTooManyRequests, "ROUTING_BUSY", err.Error())
	case errors.Is(err, routing.ErrTimeout):
		apiError(c, http.StatusGatewayTimeout, "ROUTING_TIMEOUT", err.Error())
	case errors.Is(err, routing.ErrUnavailable):
		apiError(c, http.StatusServiceUnavailable, "ROUTING_UNAVAILABLE", err.Error())
	case errors.Is(err, routing.ErrNoRoute), errors.Is(err, routing.ErrNoRoadNode), errors.Is(err, routing.ErrUnsupportedProfile), errors.Is(err, routing.ErrOutsideRegion):
		apiError(c, http.StatusUnprocessableEntity, err.Error(), err.Error())
	case errors.Is(err, extension.ErrPlanningRequired):
		apiError(c, http.StatusUnprocessableEntity, "PLANNER_UNAVAILABLE", err.Error())
	case errors.Is(err, repository.ErrNotFound):
		apiError(c, http.StatusNotFound, "NOT_FOUND", "Resource was not found")
	case errors.Is(err, repository.ErrForbidden):
		apiError(c, http.StatusForbidden, "FORBIDDEN", "Command identity does not match the authenticated device")
	case errors.Is(err, repository.ErrConflict), errors.Is(err, repository.ErrInvalidTransition):
		apiError(c, http.StatusConflict, "INVALID_STATE_TRANSITION", "Operation is not legal in the current durable state")
	default:
		apiError(c, http.StatusInternalServerError, "ORCHESTRATION_ERROR", "Orchestration operation failed")
	}
}

func (a *OrchestrationAPI) createTask(c *gin.Context) {
	tenant, ok := tenantFor(c)
	if !ok {
		apiError(c, 400, "TENANT_REQUIRED", "Tenant scope is required")
		return
	}
	principal, _ := Principal(c)
	var in struct {
		ProjectID    *string               `json:"project_id"`
		TaskType     string                `json:"task_type"`
		Priority     string                `json:"priority"`
		Requirements taskcore.Requirements `json:"requirements"`
		Target       json.RawMessage       `json:"target"`
		ExpiresAt    *time.Time            `json:"expires_at"`
	}
	if c.ShouldBindJSON(&in) != nil {
		apiError(c, 400, "INVALID_REQUEST", "A valid task document is required")
		return
	}
	if in.Priority == "" {
		in.Priority = "NORMAL"
	}
	expiresAt := time.Now().UTC().Add(5 * time.Minute)
	if in.ExpiresAt != nil {
		expiresAt = in.ExpiresAt.UTC()
	}
	result, err := a.service.CreateTask(c, tenant, principal, RequestID(c), orchestration.CreateTaskInput{ProjectID: in.ProjectID, TaskType: in.TaskType, Priority: in.Priority, Requirements: in.Requirements, Target: in.Target, ExpiresAt: expiresAt, CorrelationID: c.GetHeader("X-Correlation-ID")})
	if err != nil {
		orchestrationError(c, err)
		return
	}
	apiData(c, http.StatusCreated, result)
}

func (a *OrchestrationAPI) listTasks(c *gin.Context) {
	tenant, ok := tenantFor(c)
	if !ok {
		apiError(c, 400, "TENANT_REQUIRED", "Tenant scope is required")
		return
	}
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	result, err := a.store.ListTasks(c, tenant, limit, c.Query("cursor"), c.Query("status"), c.Query("device_id"))
	if err != nil {
		orchestrationError(c, err)
		return
	}
	apiData(c, 200, result)
}

func (a *OrchestrationAPI) getTask(c *gin.Context) {
	tenant, ok := tenantFor(c)
	if !ok {
		apiError(c, 400, "TENANT_REQUIRED", "Tenant scope is required")
		return
	}
	v, err := a.store.GetTask(c, tenant, c.Param("task_id"))
	if err != nil {
		orchestrationError(c, err)
		return
	}
	commands, err := a.store.ListCommands(c, tenant, 100, "", "", v.TaskID, "")
	if err != nil {
		orchestrationError(c, err)
		return
	}
	apiData(c, 200, gin.H{"task": v, "commands": commands})
}

func (a *OrchestrationAPI) cancelTask(c *gin.Context) {
	tenant, _ := tenantFor(c)
	principal, _ := Principal(c)
	if err := a.store.CancelTask(c, tenant, c.Param("task_id"), principal.APIKeyID, RequestID(c)); err != nil {
		orchestrationError(c, err)
		return
	}
	apiData(c, 200, gin.H{"task_id": c.Param("task_id"), "status": taskcore.Cancelled})
}

func (a *OrchestrationAPI) retryTask(c *gin.Context) {
	tenant, _ := tenantFor(c)
	principal, _ := Principal(c)
	var in struct {
		TTLSeconds int `json:"ttl_seconds"`
	}
	_ = c.ShouldBindJSON(&in)
	result, err := a.service.RetryTask(c, tenant, c.Param("task_id"), principal, RequestID(c), time.Duration(in.TTLSeconds)*time.Second)
	if err != nil {
		orchestrationError(c, err)
		return
	}
	apiData(c, 200, result)
}

func (a *OrchestrationAPI) listCommands(c *gin.Context) {
	tenant, ok := tenantFor(c)
	if !ok {
		apiError(c, 400, "TENANT_REQUIRED", "Tenant scope is required")
		return
	}
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	result, err := a.store.ListCommands(c, tenant, limit, c.Query("cursor"), c.Query("status"), c.Query("task_id"), c.Query("device_id"))
	if err != nil {
		orchestrationError(c, err)
		return
	}
	apiData(c, 200, result)
}

func (a *OrchestrationAPI) getCommand(c *gin.Context) {
	tenant, ok := tenantFor(c)
	if !ok {
		apiError(c, 400, "TENANT_REQUIRED", "Tenant scope is required")
		return
	}
	v, err := a.store.GetCommand(c, tenant, c.Param("command_id"))
	if err != nil {
		orchestrationError(c, err)
		return
	}
	apiData(c, 200, v)
}

func (a *OrchestrationAPI) retryCommand(c *gin.Context) {
	tenant, _ := tenantFor(c)
	principal, _ := Principal(c)
	if err := a.store.RetryCommand(c, tenant, c.Param("command_id"), principal.APIKeyID, RequestID(c), true); err != nil {
		orchestrationError(c, err)
		return
	}
	apiData(c, 200, gin.H{"command_id": c.Param("command_id"), "status": command.Pending})
}

func (a *OrchestrationAPI) cancelCommand(c *gin.Context) {
	tenant, _ := tenantFor(c)
	principal, _ := Principal(c)
	v, err := a.store.GetCommand(c, tenant, c.Param("command_id"))
	if err == nil {
		err = a.store.CancelTask(c, tenant, v.TaskID, principal.APIKeyID, RequestID(c))
	}
	if err != nil {
		orchestrationError(c, err)
		return
	}
	apiData(c, 200, gin.H{"command_id": v.CommandID, "status": command.Cancelled})
}
