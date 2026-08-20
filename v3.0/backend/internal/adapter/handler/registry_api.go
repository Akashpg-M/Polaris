package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Akashpg-M/polaris/backend/internal/adapter/repository"
	"github.com/Akashpg-M/polaris/backend/internal/core/auth"
	"github.com/Akashpg-M/polaris/backend/internal/core/registry"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

const (
	operatorPrincipalKey = "operator_principal"
	requestIDKey         = "request_id"
)

func RequestID(c *gin.Context) string {
	if value, ok := c.Get(requestIDKey); ok {
		if id, valid := value.(string); valid && id != "" {
			return id
		}
	}
	v := c.GetHeader("X-Request-ID")
	if v == "" {
		v = auth.NewID()
	}
	c.Set(requestIDKey, v)
	return v
}
func Principal(c *gin.Context) (auth.OperatorPrincipal, bool) {
	v, ok := c.Get(operatorPrincipalKey)
	if !ok {
		return auth.OperatorPrincipal{}, false
	}
	p, ok := v.(auth.OperatorPrincipal)
	return p, ok
}

type RegistryAPI struct {
	store                               *repository.RegistryStore
	redis                               *redis.Client
	staleAfter, offlineAfter, ticketTTL time.Duration
	lifecycleHook                       func(tenant, device, status string)
}

func (a *RegistryAPI) SetLifecycleHook(fn func(tenant, device, status string)) { a.lifecycleHook = fn }

func NewRegistryAPI(store *repository.RegistryStore, redisClient *redis.Client, staleAfter, offlineAfter, ticketTTL time.Duration) *RegistryAPI {
	return &RegistryAPI{store: store, redis: redisClient, staleAfter: staleAfter, offlineAfter: offlineAfter, ticketTTL: ticketTTL}
}

func (a *RegistryAPI) Middleware(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			apiError(c, http.StatusUnauthorized, "UNAUTHENTICATED", "Operator API key is required")
			c.Abort()
			return
		}
		p, err := a.store.ResolveOperator(c.Request.Context(), strings.TrimSpace(strings.TrimPrefix(header, "Bearer ")))
		if err != nil {
			apiError(c, http.StatusUnauthorized, "UNAUTHENTICATED", "Operator API key is invalid")
			c.Abort()
			return
		}
		if !auth.Can(p.Role, permission) {
			apiError(c, http.StatusForbidden, "FORBIDDEN", "Role does not permit this operation")
			c.Abort()
			return
		}
		c.Set(operatorPrincipalKey, p)
		c.Next()
	}
}
func apiData(c *gin.Context, status int, data interface{}) {
	c.JSON(status, gin.H{"data": data, "request_id": RequestID(c)})
}
func apiError(c *gin.Context, status int, code, message string) {
	c.JSON(status, gin.H{"error": gin.H{"code": code, "message": message}, "request_id": RequestID(c)})
}
func registryError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, repository.ErrNotFound):
		apiError(c, 404, "NOT_FOUND", "Resource was not found")
	case errors.Is(err, repository.ErrConflict):
		apiError(c, 409, "CONFLICT", "Resource already exists")
	case errors.Is(err, repository.ErrInvalidTransition):
		apiError(c, 409, "INVALID_LIFECYCLE_TRANSITION", "Lifecycle transition is not allowed")
	default:
		apiError(c, 500, "INTERNAL_ERROR", "Registry operation failed")
	}
}
func tenantFor(c *gin.Context) (string, bool) {
	p, ok := Principal(c)
	if !ok {
		return "", false
	}
	if p.Role == auth.PlatformAdmin {
		tenant := c.GetHeader("X-Tenant-ID")
		if tenant == "" {
			tenant = c.Query("tenant_id")
		}
		return tenant, tenant != ""
	}
	return p.TenantID, p.TenantID != ""
}

func (a *RegistryAPI) Register(r *gin.RouterGroup) {
	r.POST("/tenants", a.Middleware("mutate"), a.createTenant)
	r.GET("/tenants/:tenant_id", a.Middleware("read"), a.getTenant)
	r.PATCH("/tenants/:tenant_id", a.Middleware("mutate"), a.patchTenant)
	r.POST("/projects", a.Middleware("mutate"), a.createProject)
	r.GET("/projects", a.Middleware("read"), a.listProjects)
	r.GET("/projects/:project_id", a.Middleware("read"), a.getProject)
	r.PATCH("/projects/:project_id", a.Middleware("mutate"), a.patchProject)
	r.POST("/devices", a.Middleware("mutate"), a.createDevice)
	r.GET("/devices", a.Middleware("read"), a.listDevices)
	r.GET("/devices/:device_id", a.Middleware("read"), a.getDevice)
	r.PATCH("/devices/:device_id", a.Middleware("mutate"), a.patchDevice)
	r.POST("/devices/:device_id/activate", a.Middleware("mutate"), a.lifecycle("ACTIVE"))
	r.POST("/devices/:device_id/suspend", a.Middleware("mutate"), a.lifecycle("SUSPENDED"))
	r.POST("/devices/:device_id/decommission", a.Middleware("mutate"), a.lifecycle("DECOMMISSIONED"))
	r.GET("/capabilities", a.Middleware("read"), a.allCapabilities)
	r.GET("/devices/:device_id/capabilities", a.Middleware("read"), a.deviceCapabilities)
	r.PUT("/devices/:device_id/capabilities/:capability_id", a.Middleware("mutate"), a.putCapability)
	r.DELETE("/devices/:device_id/capabilities/:capability_id", a.Middleware("mutate"), a.removeCapability)
	r.POST("/devices/:device_id/credentials", a.Middleware("mutate"), a.issueCredential)
	r.GET("/devices/:device_id/credentials", a.Middleware("read"), a.listCredentials)
	r.POST("/devices/:device_id/credentials/:credential_id/revoke", a.Middleware("mutate"), a.revokeCredential)
	r.POST("/devices/:device_id/credentials/rotate", a.Middleware("mutate"), a.rotateCredential)
	r.POST("/devices/:device_id/connection-ticket", a.Middleware("mutate"), a.connectionTicket)
	r.GET("/devices/:device_id/twin", a.Middleware("read"), a.getTwin)
	r.GET("/twins", a.Middleware("read"), a.listTwins)
	r.GET("/audit-events", a.Middleware("audit"), a.listAudit)
	r.POST("/dashboard-ticket", a.Middleware("read"), a.dashboardTicket)
}

func (a *RegistryAPI) createTenant(c *gin.Context) {
	p, _ := Principal(c)
	if p.Role != auth.PlatformAdmin {
		apiError(c, 403, "FORBIDDEN", "Only platform admins create tenants")
		return
	}
	var in struct {
		TenantID    string          `json:"tenant_id"`
		DisplayName string          `json:"display_name"`
		Metadata    json.RawMessage `json:"metadata"`
	}
	if c.ShouldBindJSON(&in) != nil || in.TenantID == "" || in.DisplayName == "" {
		apiError(c, 400, "INVALID_REQUEST", "tenant_id and display_name are required")
		return
	}
	t := registry.Tenant{TenantID: in.TenantID, DisplayName: in.DisplayName, Status: "ACTIVE", Metadata: jsonOrEmptyRaw(in.Metadata)}
	if err := a.store.CreateTenant(c, t, p.APIKeyID, RequestID(c)); err != nil {
		registryError(c, err)
		return
	}
	apiData(c, 201, t)
}
func (a *RegistryAPI) getTenant(c *gin.Context) {
	p, _ := Principal(c)
	id := c.Param("tenant_id")
	if p.Role != auth.PlatformAdmin && p.TenantID != id {
		_ = a.store.Audit(c, p.TenantID, p.APIKeyID, "CROSS_TENANT_ACCESS_DENIED", "tenant", id, RequestID(c), "DENIED")
		apiError(c, 404, "NOT_FOUND", "Resource was not found")
		return
	}
	v, err := a.store.GetTenant(c, id)
	if err != nil {
		registryError(c, err)
		return
	}
	apiData(c, 200, v)
}
func (a *RegistryAPI) patchTenant(c *gin.Context) {
	p, _ := Principal(c)
	if p.Role != auth.PlatformAdmin {
		apiError(c, 403, "FORBIDDEN", "Only platform admins change tenant lifecycle")
		return
	}
	var in struct {
		Status string `json:"status"`
	}
	if c.ShouldBindJSON(&in) != nil {
		apiError(c, 400, "INVALID_REQUEST", "status is required")
		return
	}
	if err := a.store.SetTenantStatus(c, c.Param("tenant_id"), in.Status, p.APIKeyID, RequestID(c)); err != nil {
		registryError(c, err)
		return
	}
	if a.lifecycleHook != nil {
		a.lifecycleHook(c.Param("tenant_id"), "", in.Status)
	}
	apiData(c, 200, gin.H{"tenant_id": c.Param("tenant_id"), "status": in.Status})
}
func jsonOrEmptyRaw(v json.RawMessage) []byte {
	if len(v) == 0 {
		return []byte(`{}`)
	}
	return v
}
func (a *RegistryAPI) createProject(c *gin.Context) {
	tenant, ok := tenantFor(c)
	if !ok {
		apiError(c, 400, "TENANT_REQUIRED", "Tenant scope is required")
		return
	}
	p, _ := Principal(c)
	var in struct {
		Name        string          `json:"name"`
		Description *string         `json:"description"`
		Metadata    json.RawMessage `json:"metadata"`
	}
	if c.ShouldBindJSON(&in) != nil || in.Name == "" {
		apiError(c, 400, "INVALID_REQUEST", "name is required")
		return
	}
	v := registry.Project{ProjectID: auth.NewID(), TenantID: tenant, Name: in.Name, Description: in.Description, Status: "ACTIVE", Metadata: jsonOrEmptyRaw(in.Metadata)}
	if err := a.store.CreateProject(c, v, p.APIKeyID, RequestID(c)); err != nil {
		registryError(c, err)
		return
	}
	apiData(c, 201, v)
}
func (a *RegistryAPI) listProjects(c *gin.Context) {
	tenant, ok := tenantFor(c)
	if !ok {
		apiError(c, 400, "TENANT_REQUIRED", "Tenant scope is required")
		return
	}
	v, err := a.store.ListProjects(c, tenant)
	if err != nil {
		registryError(c, err)
		return
	}
	apiData(c, 200, v)
}
func (a *RegistryAPI) getProject(c *gin.Context) {
	tenant, ok := tenantFor(c)
	if !ok {
		apiError(c, 400, "TENANT_REQUIRED", "Tenant scope is required")
		return
	}
	v, err := a.store.GetProject(c, tenant, c.Param("project_id"))
	if err != nil {
		registryError(c, err)
		return
	}
	apiData(c, 200, v)
}
func (a *RegistryAPI) patchProject(c *gin.Context) {
	tenant, ok := tenantFor(c)
	if !ok {
		apiError(c, 400, "TENANT_REQUIRED", "Tenant scope is required")
		return
	}
	p, _ := Principal(c)
	var in struct {
		Name        *string         `json:"name"`
		Description *string         `json:"description"`
		Status      *string         `json:"status"`
		Metadata    json.RawMessage `json:"metadata"`
	}
	if c.ShouldBindJSON(&in) != nil {
		apiError(c, 400, "INVALID_REQUEST", "A valid project update is required")
		return
	}
	if err := a.store.UpdateProject(c, tenant, c.Param("project_id"), in.Name, in.Description, in.Status, in.Metadata, p.APIKeyID, RequestID(c)); err != nil {
		registryError(c, err)
		return
	}
	v, err := a.store.GetProject(c, tenant, c.Param("project_id"))
	if err != nil {
		registryError(c, err)
		return
	}
	apiData(c, 200, v)
}
func (a *RegistryAPI) createDevice(c *gin.Context) {
	tenant, ok := tenantFor(c)
	if !ok {
		apiError(c, 400, "TENANT_REQUIRED", "Tenant scope is required")
		return
	}
	p, _ := Principal(c)
	var in struct {
		DeviceID     string          `json:"device_id"`
		ProjectID    *string         `json:"project_id"`
		DeviceTypeID string          `json:"device_type_id"`
		DisplayName  string          `json:"display_name"`
		Metadata     json.RawMessage `json:"metadata"`
	}
	if c.ShouldBindJSON(&in) != nil || in.DeviceID == "" || in.DeviceTypeID == "" || in.DisplayName == "" {
		apiError(c, 400, "INVALID_REQUEST", "device_id, device_type_id and display_name are required")
		return
	}
	v := registry.Device{TenantID: tenant, DeviceID: in.DeviceID, ProjectID: in.ProjectID, DeviceTypeID: in.DeviceTypeID, DisplayName: in.DisplayName, LifecycleStatus: "REGISTERED", Metadata: jsonOrEmptyRaw(in.Metadata)}
	if err := a.store.CreateDevice(c, v, p.APIKeyID, RequestID(c)); err != nil {
		registryError(c, err)
		return
	}
	apiData(c, 201, v)
}
func (a *RegistryAPI) listDevices(c *gin.Context) {
	tenant, ok := tenantFor(c)
	if !ok {
		apiError(c, 400, "TENANT_REQUIRED", "Tenant scope is required")
		return
	}
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	v, err := a.store.ListDevices(c, tenant, limit, c.Query("cursor"), c.Query("project_id"), c.Query("device_type"), c.Query("lifecycle_status"), c.Query("capability"))
	if err != nil {
		registryError(c, err)
		return
	}
	apiData(c, 200, v)
}
func (a *RegistryAPI) getDevice(c *gin.Context) {
	tenant, ok := tenantFor(c)
	if !ok {
		return
	}
	v, err := a.store.GetDevice(c, tenant, c.Param("device_id"))
	if err != nil {
		registryError(c, err)
		return
	}
	apiData(c, 200, v)
}
func (a *RegistryAPI) patchDevice(c *gin.Context) {
	tenant, ok := tenantFor(c)
	if !ok {
		apiError(c, 400, "TENANT_REQUIRED", "Tenant scope is required")
		return
	}
	p, _ := Principal(c)
	var in struct {
		DisplayName     *string         `json:"display_name"`
		FirmwareVersion *string         `json:"firmware_version"`
		SoftwareVersion *string         `json:"software_version"`
		ModelVersion    *string         `json:"model_version"`
		Metadata        json.RawMessage `json:"metadata"`
	}
	if c.ShouldBindJSON(&in) != nil {
		apiError(c, 400, "INVALID_REQUEST", "A valid device update is required")
		return
	}
	if err := a.store.UpdateDevice(c, tenant, c.Param("device_id"), in.DisplayName, in.FirmwareVersion, in.SoftwareVersion, in.ModelVersion, in.Metadata, p.APIKeyID, RequestID(c)); err != nil {
		registryError(c, err)
		return
	}
	v, err := a.store.GetDevice(c, tenant, c.Param("device_id"))
	if err != nil {
		registryError(c, err)
		return
	}
	apiData(c, 200, v)
}
func (a *RegistryAPI) lifecycle(next string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tenant, ok := tenantFor(c)
		if !ok {
			return
		}
		p, _ := Principal(c)
		if err := a.store.SetLifecycle(c, tenant, c.Param("device_id"), next, p.APIKeyID, RequestID(c)); err != nil {
			registryError(c, err)
			return
		}
		if a.lifecycleHook != nil {
			a.lifecycleHook(tenant, c.Param("device_id"), next)
		}
		apiData(c, 200, gin.H{"device_id": c.Param("device_id"), "lifecycle_status": next})
	}
}
func (a *RegistryAPI) allCapabilities(c *gin.Context) {
	v, err := a.store.AllCapabilities(c)
	if err != nil {
		registryError(c, err)
		return
	}
	apiData(c, 200, v)
}
func (a *RegistryAPI) deviceCapabilities(c *gin.Context) {
	tenant, _ := tenantFor(c)
	v, err := a.store.ListCapabilities(c, tenant, c.Param("device_id"))
	if err != nil {
		registryError(c, err)
		return
	}
	apiData(c, 200, v)
}
func (a *RegistryAPI) putCapability(c *gin.Context) {
	tenant, _ := tenantFor(c)
	p, _ := Principal(c)
	var in struct {
		Configuration json.RawMessage `json:"configuration"`
	}
	_ = c.ShouldBindJSON(&in)
	if err := a.store.PutCapability(c, tenant, c.Param("device_id"), c.Param("capability_id"), jsonOrEmptyRaw(in.Configuration), p.APIKeyID, RequestID(c)); err != nil {
		registryError(c, err)
		return
	}
	apiData(c, 200, gin.H{"enabled": true})
}
func (a *RegistryAPI) removeCapability(c *gin.Context) {
	tenant, _ := tenantFor(c)
	p, _ := Principal(c)
	if err := a.store.RemoveCapability(c, tenant, c.Param("device_id"), c.Param("capability_id"), p.APIKeyID, RequestID(c)); err != nil {
		registryError(c, err)
		return
	}
	apiData(c, 200, gin.H{"enabled": false})
}
func (a *RegistryAPI) issueCredential(c *gin.Context) {
	tenant, _ := tenantFor(c)
	p, _ := Principal(c)
	meta, raw, err := a.store.IssueCredential(c, tenant, c.Param("device_id"), p.APIKeyID, RequestID(c), nil)
	if err != nil {
		registryError(c, err)
		return
	}
	apiData(c, 201, gin.H{"credential": meta, "secret": raw})
}
func (a *RegistryAPI) listCredentials(c *gin.Context) {
	tenant, _ := tenantFor(c)
	v, err := a.store.ListCredentials(c, tenant, c.Param("device_id"))
	if err != nil {
		registryError(c, err)
		return
	}
	apiData(c, 200, v)
}
func (a *RegistryAPI) revokeCredential(c *gin.Context) {
	tenant, _ := tenantFor(c)
	p, _ := Principal(c)
	if err := a.store.RevokeCredential(c, tenant, c.Param("device_id"), c.Param("credential_id"), p.APIKeyID, RequestID(c)); err != nil {
		registryError(c, err)
		return
	}
	apiData(c, 200, gin.H{"status": "REVOKED"})
}
func (a *RegistryAPI) rotateCredential(c *gin.Context) {
	tenant, _ := tenantFor(c)
	p, _ := Principal(c)
	var in struct {
		CredentialID string `json:"credential_id"`
	}
	if c.ShouldBindJSON(&in) != nil || in.CredentialID == "" {
		apiError(c, 400, "INVALID_REQUEST", "credential_id is required")
		return
	}
	meta, raw, err := a.store.RotateCredential(c, tenant, c.Param("device_id"), in.CredentialID, p.APIKeyID, RequestID(c))
	if err != nil {
		registryError(c, err)
		return
	}
	apiData(c, 201, gin.H{"credential": meta, "secret": raw})
}
func (a *RegistryAPI) connectionTicket(c *gin.Context) {
	tenant, _ := tenantFor(c)
	credentials, err := a.store.ListCredentials(c, tenant, c.Param("device_id"))
	if err != nil || len(credentials) == 0 {
		registryError(c, repository.ErrNotFound)
		return
	}
	active := ""
	for _, v := range credentials {
		if v.Status == "ACTIVE" {
			active = v.CredentialID
			break
		}
	}
	if active == "" {
		registryError(c, repository.ErrNotFound)
		return
	}
	ticket, err := a.store.CreateTicket(c, tenant, c.Param("device_id"), active, a.ticketTTL)
	if err != nil {
		registryError(c, err)
		return
	}
	apiData(c, 201, gin.H{"ticket": ticket, "expires_in_seconds": int(a.ticketTTL.Seconds())})
}

func connectivity(lastSeen string, stale, offline time.Duration) (string, *time.Time) {
	if lastSeen == "" {
		return "NEVER_CONNECTED", nil
	}
	ms, err := strconv.ParseInt(lastSeen, 10, 64)
	if err != nil {
		return "NEVER_CONNECTED", nil
	}
	t := time.UnixMilli(ms).UTC()
	age := time.Since(t)
	if age > offline {
		return "OFFLINE", &t
	}
	if age > stale {
		return "STALE", &t
	}
	return "ONLINE", &t
}
func (a *RegistryAPI) twin(ctx context.Context, tenant, id string) (gin.H, error) {
	device, err := a.store.GetDevice(ctx, tenant, id)
	if err != nil {
		return nil, err
	}
	caps, err := a.store.ListCapabilities(ctx, tenant, id)
	if err != nil {
		return nil, err
	}
	key := "polaris:twin:" + tenant + ":" + id
	state, err := a.redis.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	status, last := connectivity(state["last_seen_at"], a.staleAfter, a.offlineAfter)
	var reported interface{}
	if raw := state["reported_state"]; raw != "" {
		_ = json.Unmarshal([]byte(raw), &reported)
	}
	components := map[string]interface{}{}
	for field, raw := range state {
		if strings.HasPrefix(field, "component:") {
			var component interface{}
			if json.Unmarshal([]byte(raw), &component) == nil {
				components[strings.TrimPrefix(field, "component:")] = component
			}
		}
	}
	return gin.H{"tenant_id": tenant, "device_id": id, "device": device, "capabilities": caps, "reported_state": reported, "components": components, "desired_state": nil, "connectivity": gin.H{"status": status, "last_seen_at": last}}, nil
}
func (a *RegistryAPI) getTwin(c *gin.Context) {
	tenant, _ := tenantFor(c)
	v, err := a.twin(c, tenant, c.Param("device_id"))
	if err != nil {
		registryError(c, err)
		return
	}
	apiData(c, 200, v)
}
func (a *RegistryAPI) listTwins(c *gin.Context) {
	tenant, _ := tenantFor(c)
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	devices, err := a.store.ListDevices(c, tenant, limit, c.Query("cursor"), c.Query("project_id"), c.Query("device_type"), c.Query("lifecycle_status"), c.Query("capability"))
	if err != nil {
		registryError(c, err)
		return
	}
	out := []gin.H{}
	for _, d := range devices {
		v, e := a.twin(c, tenant, d.DeviceID)
		if e == nil {
			if filter := c.Query("connectivity_status"); filter == "" || v["connectivity"].(gin.H)["status"] == filter {
				out = append(out, v)
			}
		}
	}
	apiData(c, 200, out)
}
func (a *RegistryAPI) listAudit(c *gin.Context) {
	tenant, _ := tenantFor(c)
	v, err := a.store.ListAudit(c, tenant)
	if err != nil {
		registryError(c, err)
		return
	}
	apiData(c, 200, v)
}
func (a *RegistryAPI) dashboardTicket(c *gin.Context) {
	p, _ := Principal(c)
	tenant, _ := tenantFor(c)
	ticket, err := a.store.CreateOperatorTicket(c, p, tenant, a.ticketTTL)
	if err != nil {
		registryError(c, err)
		return
	}
	apiData(c, 201, gin.H{"ticket": ticket, "expires_in_seconds": int(a.ticketTTL.Seconds())})
}
