package capydbclient

import (
	"encoding/json"
	"time"
)

// This file holds the canonical Go representations of the control-plane entities
// and request bodies that the CLI and the Terraform provider both consume. They
// mirror the component schemas in backend/internal/httpapi/openapi.json (the
// same source the TypeScript SDK is generated from) so the two Go clients share
// one shape instead of each re-declaring - and drifting - their own copies.
//
// Response entities carry the full schema field set; a consumer that only reads
// a subset simply ignores the rest. Date-time fields use time.Time (required)
// or *time.Time (optional), matching the backend model.

// Organization is a billing/identity tenant.
type Organization struct {
	BillingCustomerID     string     `json:"billing_customer_id,omitempty"`
	BillingEmail          string     `json:"billing_email,omitempty"`
	BillingName           string     `json:"billing_name,omitempty"`
	BillingPeriodEnd      *time.Time `json:"billing_period_end,omitempty"`
	BillingPlan           string     `json:"billing_plan"`
	BillingProductID      string     `json:"billing_product_id,omitempty"`
	BillingProvider       string     `json:"billing_provider"`
	BillingStatus         string     `json:"billing_status"`
	BillingSubscriptionID string     `json:"billing_subscription_id,omitempty"`
	ClerkOrganizationID   string     `json:"clerk_organization_id,omitempty"`
	ClerkOrganizationSlug string     `json:"clerk_organization_slug,omitempty"`
	// CloudflareAccountID links the organization to the Cloudflare account that
	// provisioned and is billed for it (partner-provisioned organizations only).
	CloudflareAccountID  string     `json:"cloudflare_account_id,omitempty"`
	CreatedAt            time.Time  `json:"created_at"`
	ID                   string     `json:"id"`
	Name                 string     `json:"name"`
	Slug                 string     `json:"slug"`
	SuspendedAt          *time.Time `json:"suspended_at,omitempty"`
	SuspendedReason      string     `json:"suspended_reason,omitempty"`
	UpdatedAt            time.Time  `json:"updated_at"`
	VercelInstallationID string     `json:"vercel_installation_id,omitempty"`
}

// Viewer is the GET /v1/me payload: the resolved caller identity plus the
// organization it is bound to (nil for the platform admin token, which is not
// scoped to one).
type Viewer struct {
	Organization *Organization `json:"organization"`
	Principal    Principal     `json:"principal"`
}

// Project is a managed Postgres database.
type Project struct {
	CreatedAt              time.Time `json:"created_at"`
	DatabaseName           string    `json:"database_name"`
	DirectPort             int       `json:"direct_port"`
	Environment            string    `json:"environment"`
	ID                     string    `json:"id"`
	IdleTransactionTimeout string    `json:"idle_transaction_timeout"`
	LastError              string    `json:"last_error,omitempty"`
	LatestJobID            string    `json:"latest_job_id,omitempty"`
	MaxConnections         int       `json:"max_connections"`
	Name                   string    `json:"name"`
	OrganizationID         string    `json:"organization_id"`
	Plan                   string    `json:"plan"`
	PooledPort             int       `json:"pooled_port"`
	// PostgresVersion is the major version of the project's database. Empty
	// while the database is still provisioning.
	PostgresVersion   string `json:"postgres_version,omitempty"`
	PrimaryInstanceID string `json:"primary_instance_id,omitempty"`
	PublicHost        string `json:"public_host,omitempty"`
	Region            string `json:"region"`
	RoleName          string `json:"role_name"`
	// RuntimeStatus overlays the scale-to-zero lifecycle in customer
	// vocabulary: "active", "paused", or "resuming". Empty while the database
	// is provisioning or being deleted (see State).
	RuntimeStatus     string    `json:"runtime_status,omitempty"`
	Slug              string    `json:"slug"`
	SSLMode           string    `json:"ssl_mode,omitempty"`
	State             string    `json:"state"`
	StatementTimeout  string    `json:"statement_timeout"`
	StorageLimitBytes int64     `json:"storage_limit_bytes"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// Job is an asynchronous control-plane operation. Poll until State is
// "completed" or "failed".
type Job struct {
	Attempts            int        `json:"attempts"`
	ClaimedAt           *time.Time `json:"claimed_at,omitempty"`
	ClaimedBy           string     `json:"claimed_by,omitempty"`
	CompletedAt         *time.Time `json:"completed_at,omitempty"`
	CreatedAt           time.Time  `json:"created_at"`
	Error               string     `json:"error,omitempty"`
	HostID              string     `json:"host_id,omitempty"`
	ID                  string     `json:"id"`
	InstanceID          string     `json:"instance_id,omitempty"`
	LastExitCode        int        `json:"last_exit_code,omitempty"`
	LastStderr          string     `json:"last_stderr,omitempty"`
	LastStdout          string     `json:"last_stdout,omitempty"`
	LockedResource      string     `json:"locked_resource,omitempty"`
	MaxAttempts         int        `json:"max_attempts"`
	OrganizationID      string     `json:"organization_id"`
	PreviewDatabaseID   string     `json:"preview_database_id,omitempty"`
	ProjectID           string     `json:"project_id,omitempty"`
	RetryClassification string     `json:"retry_classification,omitempty"`
	// Result is the structured result payload of a completed job. The API
	// returns it only for job types whose result is part of the product
	// contract (currently project.import_follow_status).
	Result    json.RawMessage `json:"result,omitempty"`
	StartedAt *time.Time      `json:"started_at,omitempty"`
	State     string          `json:"state"`
	Type      string          `json:"type"`
	UpdatedAt time.Time       `json:"updated_at"`
}

// ConnectionInfo is a project or preview database's connection endpoints.
type ConnectionInfo struct {
	DirectURL string `json:"direct_url,omitempty"`
	PooledURL string `json:"pooled_url,omitempty"`
	Username  string `json:"username"`
}

// APIKey is an organization or project-scoped API key. Plaintext secrets are
// returned only on creation, never on list endpoints.
type APIKey struct {
	CreatedAt       time.Time  `json:"created_at"`
	CreatedByUserID string     `json:"created_by_user_id,omitempty"`
	DeviceName      string     `json:"device_name,omitempty"`
	ExpiresAt       *time.Time `json:"expires_at,omitempty"`
	ID              string     `json:"id"`
	IsActive        bool       `json:"is_active"`
	KeyPrefix       string     `json:"key_prefix"`
	LastUsedAt      *time.Time `json:"last_used_at,omitempty"`
	Name            string     `json:"name"`
	OrganizationID  string     `json:"organization_id"`
	ProjectID       string     `json:"project_id,omitempty"`
	RevokedAt       *time.Time `json:"revoked_at,omitempty"`
	Scopes          []string   `json:"scopes"`
	Source          string     `json:"source"`
}

// WebhookEndpoint is an outbound webhook receiver.
type WebhookEndpoint struct {
	CreatedAt      time.Time `json:"created_at"`
	Description    string    `json:"description,omitempty"`
	EventTypes     []string  `json:"event_types"`
	ID             string    `json:"id"`
	IsActive       bool      `json:"is_active"`
	OrganizationID string    `json:"organization_id"`
	UpdatedAt      time.Time `json:"updated_at"`
	URL            string    `json:"url"`
}

// WebhookDelivery is one delivery attempt record for a webhook endpoint. It is
// returned by the delivery listing and by the test/redeliver actions (both of
// which enqueue a fresh pending delivery).
type WebhookDelivery struct {
	Attempts       int        `json:"attempts"`
	CreatedAt      time.Time  `json:"created_at"`
	DeliveredAt    *time.Time `json:"delivered_at,omitempty"`
	EndpointID     string     `json:"endpoint_id"`
	EventType      string     `json:"event_type"`
	ID             string     `json:"id"`
	LastError      string     `json:"last_error,omitempty"`
	MaxAttempts    int        `json:"max_attempts"`
	NextAttemptAt  time.Time  `json:"next_attempt_at"`
	OrganizationID string     `json:"organization_id"`
	// Payload is the JSON event envelope that was (or will be) delivered.
	Payload any `json:"payload"`
	// ResponseStatus is the HTTP status of the last attempt; zero when no
	// attempt has completed yet.
	ResponseStatus int `json:"response_status,omitempty"`
	// State is one of pending, delivered, or failed.
	State     string    `json:"state"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateProjectRequest is the POST /v1/projects body.
type CreateProjectRequest struct {
	Environment    string `json:"environment,omitempty"`
	Name           string `json:"name"`
	OrganizationID string `json:"organization_id,omitempty"`
	// PostgresVersion picks the database's major version ("16", "17", "18").
	// Omit for the platform default. Immutable after creation; previews and
	// restores inherit it.
	PostgresVersion string `json:"postgres_version,omitempty"`
	Region          string `json:"region,omitempty"`
	Slug            string `json:"slug,omitempty"`
}

// CreatePreviewRequest is the create-preview-database body.
type CreatePreviewRequest struct {
	Mode     string `json:"mode,omitempty"`
	Name     string `json:"name,omitempty"`
	TTLHours int    `json:"ttl_hours,omitempty"`
}

// CreateAPIKeyRequest is the create-API-key body.
type CreateAPIKeyRequest struct {
	DeviceName string     `json:"device_name,omitempty"`
	ExpiresAt  *time.Time `json:"expires_at,omitempty"`
	Name       string     `json:"name"`
	ProjectID  string     `json:"project_id,omitempty"`
	Scopes     []string   `json:"scopes"`
}

// CreateWebhookEndpointRequest is the create-webhook-endpoint body.
type CreateWebhookEndpointRequest struct {
	Description string   `json:"description,omitempty"`
	EventTypes  []string `json:"event_types,omitempty"`
	URL         string   `json:"url"`
}

// DatabaseSchema is the canonical schema document from GET
// /v1/projects/{id}/schema (and its preview-database sibling): one
// introspection pass covering schemas, tables, views, columns, constraints,
// enums and extensions. It is the input contract for type generation and
// schema diffing.
type DatabaseSchema struct {
	DatabaseName    string            `json:"database_name"`
	Extensions      []SchemaExtension `json:"extensions"`
	PostgresVersion string            `json:"postgres_version"`
	Schemas         []SchemaNamespace `json:"schemas"`
}

// SchemaExtension is an installed Postgres extension.
type SchemaExtension struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// SchemaNamespace is one Postgres schema (namespace) with its relations and
// enum types.
type SchemaNamespace struct {
	Enums  []SchemaEnum  `json:"enums"`
	Name   string        `json:"name"`
	Tables []SchemaTable `json:"tables"`
}

// SchemaEnum is a user-defined enum type.
type SchemaEnum struct {
	Comment string   `json:"comment,omitempty"`
	Name    string   `json:"name"`
	Values  []string `json:"values"`
}

// SchemaTable is a relation; Kind is one of table, partitioned_table, view,
// materialized_view, foreign_table.
type SchemaTable struct {
	Columns           []SchemaColumn           `json:"columns"`
	Comment           string                   `json:"comment,omitempty"`
	ForeignKeys       []SchemaForeignKey       `json:"foreign_keys"`
	Kind              string                   `json:"kind"`
	Name              string                   `json:"name"`
	PrimaryKey        []string                 `json:"primary_key"`
	UniqueConstraints []SchemaUniqueConstraint `json:"unique_constraints"`
}

// SchemaColumn is one column of a relation. For array columns UDTName and
// UDTSchema describe the element type.
type SchemaColumn struct {
	// ArrayDims is the declared array dimension count (>= 1 when IsArray).
	ArrayDims int    `json:"array_dims,omitempty"`
	Comment   string `json:"comment,omitempty"`
	DataType  string `json:"data_type"`
	Default   string `json:"default,omitempty"`
	// Generated is "stored" or "virtual" for generated columns, empty otherwise.
	// Virtual generated columns exist only on Postgres 18+; on 16/17 every
	// generated column is stored.
	Generated   string `json:"generated,omitempty"`
	Identity    string `json:"identity,omitempty"`
	IsArray     bool   `json:"is_array"`
	IsEnum      bool   `json:"is_enum"`
	IsGenerated bool   `json:"is_generated"`
	IsNullable  bool   `json:"is_nullable"`
	Name        string `json:"name"`
	Position    int    `json:"position"`
	UDTName     string `json:"udt_name"`
	UDTSchema   string `json:"udt_schema"`
}

// SchemaForeignKey is a foreign-key constraint.
type SchemaForeignKey struct {
	Columns           []string `json:"columns"`
	Name              string   `json:"name"`
	OnDelete          string   `json:"on_delete"`
	OnUpdate          string   `json:"on_update"`
	ReferencedColumns []string `json:"referenced_columns"`
	ReferencedSchema  string   `json:"referenced_schema"`
	ReferencedTable   string   `json:"referenced_table"`
}

// SchemaUniqueConstraint is a unique constraint.
type SchemaUniqueConstraint struct {
	Columns []string `json:"columns"`
	Name    string   `json:"name"`
}

// GeneratedTypes is one generated source file from GET
// /v1/projects/{id}/schema/types: TypeScript interfaces, Zod schemas or a
// Drizzle schema rendered server-side from the live database schema.
type GeneratedTypes struct {
	Content  string `json:"content"`
	Filename string `json:"filename"`
	Language string `json:"language"`
	Style    string `json:"style,omitempty"`
}

// Principal is the resolved caller identity returned by GET /v1/me: how the
// request authenticated, the organization it is bound to (if any), and the
// scopes it carries.
type Principal struct {
	// AuthSource is how the caller authenticated (api_key, clerk_session,
	// admin_token, ...).
	AuthSource            string `json:"auth_source"`
	ClerkOrganizationID   string `json:"clerk_organization_id,omitempty"`
	ClerkOrganizationRole string `json:"clerk_organization_role,omitempty"`
	ClerkOrganizationSlug string `json:"clerk_organization_slug,omitempty"`
	// IsAdmin is true for platform operators.
	IsAdmin        bool     `json:"is_admin"`
	OrganizationID string   `json:"organization_id,omitempty"`
	Scopes         []string `json:"scopes"`
	UserID         string   `json:"user_id,omitempty"`
}

// PreviewDatabase is a short-lived branch of a project database, cloned from
// its parent's storage and reaped when TTLExpiresAt passes.
type PreviewDatabase struct {
	CreatedAt         time.Time  `json:"created_at"`
	DatabaseName      string     `json:"database_name"`
	DirectPort        int        `json:"direct_port"`
	ID                string     `json:"id"`
	LastError         string     `json:"last_error,omitempty"`
	Mode              string     `json:"mode"`
	Name              string     `json:"name"`
	PooledPort        int        `json:"pooled_port"`
	ProjectID         string     `json:"project_id"`
	PublicHost        string     `json:"public_host,omitempty"`
	RoleName          string     `json:"role_name"`
	SourceBackupKey   string     `json:"source_backup_key,omitempty"`
	SourceDatabase    string     `json:"source_database,omitempty"`
	SourceKind        string     `json:"source_kind,omitempty"`
	SourceRestoreTime *time.Time `json:"source_restore_time,omitempty"`
	SSLMode           string     `json:"ssl_mode,omitempty"`
	State             string     `json:"state"`
	TTLExpiresAt      time.Time  `json:"ttl_expires_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

// Backup is one completed (or in-flight) logical backup of a project database.
// VerificationState records whether the backup has been restored into a
// throwaway database and proven readable; VerifiedAt and VerificationError
// carry the outcome of that check.
type Backup struct {
	BackupKey         string     `json:"backup_key"`
	CreatedAt         time.Time  `json:"created_at"`
	DatabaseName      string     `json:"database_name"`
	ID                string     `json:"id"`
	Label             string     `json:"label,omitempty"`
	ProjectID         string     `json:"project_id"`
	SizeBytes         int64      `json:"size_bytes"`
	State             string     `json:"state"`
	VerificationError string     `json:"verification_error,omitempty"`
	VerificationState string     `json:"verification_state"`
	VerifiedAt        *time.Time `json:"verified_at,omitempty"`
}

// ProjectExport is a downloadable logical export (pg_dump custom-format
// archive) of the project database. Completed exports stay downloadable until
// ExpiresAt, then the artifact is deleted and the export flips to "expired".
type ProjectExport struct {
	CreatedAt    time.Time `json:"created_at"`
	DatabaseName string    `json:"database_name"`
	ExpiresAt    time.Time `json:"expires_at"`
	ID           string    `json:"id"`
	ProjectID    string    `json:"project_id"`
	SizeBytes    int64     `json:"size_bytes"`
	State        string    `json:"state"`
}

// ExportDownload is a short-lived presigned download URL for a completed
// export. Request a fresh one per download.
type ExportDownload struct {
	DownloadURL string    `json:"download_url"`
	ExpiresAt   time.Time `json:"expires_at"`
}

// ScheduledBackup is a project's recurring daily backup configuration. A
// project has at most one (the default schedule).
type ScheduledBackup struct {
	CreatedAt      time.Time  `json:"created_at"`
	CronHour       int        `json:"cron_hour"`
	CronMinute     int        `json:"cron_minute"`
	ID             string     `json:"id"`
	IsActive       bool       `json:"is_active"`
	Label          string     `json:"label"`
	LastJobID      string     `json:"last_job_id,omitempty"`
	LastRunAt      *time.Time `json:"last_run_at,omitempty"`
	OrganizationID string     `json:"organization_id"`
	ProjectID      string     `json:"project_id"`
	RetentionDays  int        `json:"retention_days"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// UpsertScheduledBackupRequest configures the project's single default
// schedule. A RetentionDays of 0 lets the server apply its default.
type UpsertScheduledBackupRequest struct {
	CronHour      int    `json:"cron_hour"`
	CronMinute    int    `json:"cron_minute"`
	IsActive      bool   `json:"is_active"`
	Label         string `json:"label,omitempty"`
	RetentionDays int    `json:"retention_days,omitempty"`
}

// RestorePoint is a labelled position a project can be restored to - either a
// storage checkpoint (Kind "snapshot") or a backup (Kind "backup").
type RestorePoint struct {
	BackupID           string     `json:"backup_id,omitempty"`
	BackupKey          string     `json:"backup_key,omitempty"`
	CreatedAt          time.Time  `json:"created_at"`
	CreatedByActorID   string     `json:"created_by_actor_id,omitempty"`
	CreatedByActorKind string     `json:"created_by_actor_kind"`
	ID                 string     `json:"id"`
	Kind               string     `json:"kind"`
	Label              string     `json:"label"`
	Note               string     `json:"note,omitempty"`
	OrganizationID     string     `json:"organization_id"`
	PITRTime           *time.Time `json:"pitr_time,omitempty"`
	ProjectID          string     `json:"project_id"`
	State              string     `json:"state"`
	UpdatedAt          time.Time  `json:"updated_at"`
}

// CreateRestorePointRequest labels a restore point. PITRTime is an RFC 3339
// instant inside the project's PITR window.
type CreateRestorePointRequest struct {
	BackupKey string `json:"backup_key,omitempty"`
	Kind      string `json:"kind,omitempty"`
	Label     string `json:"label"`
	Note      string `json:"note,omitempty"`
	PITRTime  string `json:"pitr_time,omitempty"`
}

// CreateRestoreRequest starts a restore. TargetKind selects whether the result
// lands in a new preview database or overwrites the project itself; the latter
// requires a single-use project.restore_overwrite ApprovalToken (minted via
// POST /v1/projects/{projectID}/approvals) and a non-production project.
type CreateRestoreRequest struct {
	AllowUnverifiedBackup bool   `json:"allow_unverified_backup,omitempty"`
	ApprovalToken         string `json:"approval_token,omitempty"`
	BackupKey             string `json:"backup_key,omitempty"`
	PreviewID             string `json:"preview_id,omitempty"`
	PreviewName           string `json:"preview_name,omitempty"`
	Recreate              bool   `json:"recreate,omitempty"`
	RestorePointID        string `json:"restore_point_id,omitempty"`
	RestoreTime           string `json:"restore_time,omitempty"`
	TargetKind            string `json:"target_kind,omitempty"`
	TTLHours              int    `json:"ttl_hours,omitempty"`
}

// ProjectApproval is a freshly minted single-use approval for one destructive
// project action ("project.delete" on production projects, or
// "project.restore_overwrite"). Token is the raw ap_ value, returned exactly
// once at mint - the server keeps only its hash. Tokens expire ~10 minutes
// after mint; presenting a consumed token to the destructive endpoint returns
// the job it authorized instead of repeating the action.
type ProjectApproval struct {
	Action    string `json:"action"`
	ExpiresAt string `json:"expires_at"`
	ID        string `json:"id"`
	Token     string `json:"token"`
}

// StatusComponent is one component's line in the public status payload.
type StatusComponent struct {
	Message string `json:"message,omitempty"`
	Name    string `json:"name"`
	Status  string `json:"status"`
}

// RegionStatus is one region's rolled-up health in the public status payload
// (worst node in the region wins).
type RegionStatus struct {
	Region string `json:"region"`
	Status string `json:"status"`
}

// StatusResponse is the unauthenticated GET /status payload. Regions is absent
// on deployments older than the per-region rollup, and GeneratedAt is the
// enriched timestamp that supersedes UpdatedAt.
type StatusResponse struct {
	Components  []StatusComponent `json:"components"`
	GeneratedAt time.Time         `json:"generated_at"`
	Regions     []RegionStatus    `json:"regions,omitempty"`
	Service     string            `json:"service"`
	Status      string            `json:"status"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

// ActiveQuerySample is one in-flight backend on the project database.
type ActiveQuerySample struct {
	DurationMs    int64  `json:"duration_ms"`
	PID           int    `json:"pid"`
	Query         string `json:"query"`
	State         string `json:"state"`
	Username      string `json:"username"`
	WaitEvent     string `json:"wait_event,omitempty"`
	WaitEventType string `json:"wait_event_type,omitempty"`
}

// SlowQuerySample is one normalized statement from pg_stat_statements.
type SlowQuerySample struct {
	Calls       int64   `json:"calls"`
	MeanTimeMs  float64 `json:"mean_time_ms"`
	Query       string  `json:"query"`
	Rows        int64   `json:"rows"`
	TotalTimeMs float64 `json:"total_time_ms"`
}

// ProjectObservability is the live resource and query picture for one project
// database. SlowQueries is empty unless PgStatStatements is true.
type ProjectObservability struct {
	ActiveQueries          []ActiveQuerySample `json:"active_queries"`
	Alerts                 []string            `json:"alerts"`
	ConnectionCount        int                 `json:"connection_count"`
	ConnectionLimit        int                 `json:"connection_limit"`
	ConnectionUsagePercent float64             `json:"connection_usage_percent"`
	DatabaseSizeBytes      int64               `json:"database_size_bytes"`
	PgStatStatements       bool                `json:"pg_stat_statements"`
	SlowQueries            []SlowQuerySample   `json:"slow_queries"`
	StorageLimitBytes      int64               `json:"storage_limit_bytes"`
	StorageUsagePercent    float64             `json:"storage_usage_percent"`
}

// SQLQueryRequest runs one statement against the project database. MaxRows
// caps the returned rows (server default 200, maximum 1000).
type SQLQueryRequest struct {
	// AllowUnqualifiedWrites permits statements the control plane's
	// destructive-statement guard otherwise refuses: an UPDATE or DELETE with
	// no top-level WHERE, and TRUNCATE. It defaults to false, so a caller is
	// guarded unless it opts out - intended for an interactive console where a
	// person has already expressed intent, not for automated callers.
	AllowUnqualifiedWrites bool   `json:"allow_unqualified_writes,omitempty"`
	MaxRows                int    `json:"max_rows,omitempty"`
	Query                  string `json:"query"`
	// ReadOnly runs the statement inside a READ ONLY transaction, so the
	// server itself refuses every write (DML, DDL, TRUNCATE, SELECT INTO,
	// sequence advancement) with SQLSTATE 25006, returned as a 400 carrying
	// the refusal. Executor-proven; recommended for exploratory callers.
	ReadOnly bool `json:"read_only,omitempty"`
}

// SQLQueryResult is one statement's result set. Truncated is true when the row
// cap cut the result short.
type SQLQueryResult struct {
	Columns    []string         `json:"columns"`
	DurationMs int64            `json:"duration_ms"`
	RowCount   int              `json:"row_count"`
	Rows       []map[string]any `json:"rows"`
	Truncated  bool             `json:"truncated"`
}

// ProjectLogEntry is one database log line: journal timestamp, severity parsed
// from the Postgres log format, message, and the cursor that resumes a tail
// strictly after this entry.
type ProjectLogEntry struct {
	Cursor    string    `json:"cursor"`
	Message   string    `json:"message"`
	Severity  string    `json:"severity"`
	Timestamp time.Time `json:"timestamp"`
}

// ProjectLogs is one log fetch: entries ascending by time plus the cursor to
// resume a tail from (empty when the fetch returned nothing).
type ProjectLogs struct {
	Entries    []ProjectLogEntry `json:"entries"`
	NextCursor string            `json:"next_cursor,omitempty"`
}

// ProjectAuditEvent is one recorded action against a project or organization.
// Metadata is action-specific JSON.
type ProjectAuditEvent struct {
	Action         string          `json:"action"`
	ActorID        string          `json:"actor_id,omitempty"`
	ActorKind      string          `json:"actor_kind"`
	CreatedAt      time.Time       `json:"created_at"`
	ID             string          `json:"id"`
	Metadata       json.RawMessage `json:"metadata"`
	OrganizationID string          `json:"organization_id"`
	ProjectID      string          `json:"project_id,omitempty"`
}

// ProjectExtensionStatus describes a Postgres extension and whether it is
// enabled for the project database.
type ProjectExtensionStatus struct {
	// AvailableVersion is what the platform now provides; InstalledVersion is
	// what this database actually has. CREATE EXTENSION pins the version present
	// at the time, so the two diverge after a platform package upgrade until the
	// customer applies an update.
	AvailableVersion string `json:"available_version,omitempty"`
	// Category groups the extension in listings (core, ai, search, ...).
	Category         string `json:"category,omitempty"`
	DefaultVersion   string `json:"default_version,omitempty"`
	Description      string `json:"description"`
	Enabled          bool   `json:"enabled"`
	InstalledVersion string `json:"installed_version,omitempty"`
	Name             string `json:"name"`
	// RequiresRestart marks extensions that load a shared library, so enabling
	// or disabling them RESTARTS the database and drops open connections for a
	// few seconds (pg_cron, pgaudit, pg_qualstats).
	RequiresRestart bool `json:"requires_restart,omitempty"`
	Trusted         bool `json:"trusted"`
	// UpdateAvailable is false for extensions CapyDB manages itself - those are
	// kept current automatically and are not the customer's to bump.
	UpdateAvailable bool `json:"update_available,omitempty"`
}

// ProjectAlert is a triggered resource alert (storage, connections, ...) for a
// project database.
type ProjectAlert struct {
	AcknowledgedAt *time.Time `json:"acknowledged_at,omitempty"`
	ID             string     `json:"id"`
	Kind           string     `json:"kind"`
	LastNotifiedAt *time.Time `json:"last_notified_at,omitempty"`
	LimitValue     int64      `json:"limit_value"`
	ObservedValue  int64      `json:"observed_value"`
	ProjectID      string     `json:"project_id"`
	ResolvedAt     *time.Time `json:"resolved_at,omitempty"`
	Severity       string     `json:"severity"`
	TriggeredAt    time.Time  `json:"triggered_at"`
}

// ProjectIntegration is a project's link to an external provider (Vercel,
// Netlify, Cloudflare, Clerk, ...). Config is provider-specific JSON and never
// carries secrets - HasCredentials reports whether a stored secret exists.
type ProjectIntegration struct {
	Config         json.RawMessage `json:"config"`
	CreatedAt      time.Time       `json:"created_at"`
	ExternalID     string          `json:"external_id,omitempty"`
	HasCredentials bool            `json:"has_credentials"`
	ID             string          `json:"id"`
	OrganizationID string          `json:"organization_id"`
	ProjectID      string          `json:"project_id"`
	Provider       string          `json:"provider"`
	State          string          `json:"state"`
	UpdatedAt      time.Time       `json:"updated_at"`
}

// ImportPreflightCheck is one named check in an import preflight report.
type ImportPreflightCheck struct {
	Detail string `json:"detail,omitempty"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

// SourceExtension is an extension installed on an import source database.
type SourceExtension struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// SourceForeignKey is a foreign key on an import source that points into a
// provider-managed schema (Supabase auth.users and the like).
type SourceForeignKey struct {
	Constraint  string `json:"constraint"`
	SourceTable string `json:"source_table"`
	TargetTable string `json:"target_table"`
}

// SourceInspection is what the preflight learned about the source database,
// including the provider-coupling facts that decide whether an import can be
// lifted cleanly.
type SourceInspection struct {
	AppSchemas               []string           `json:"app_schemas,omitempty"`
	DatabaseSizeBytes        int64              `json:"database_size_bytes"`
	EventTriggers            []string           `json:"event_triggers,omitempty"`
	Extensions               []SourceExtension  `json:"extensions"`
	ProviderAuthFKs          []SourceForeignKey `json:"provider_auth_fks,omitempty"`
	ProviderAuthPolicyTables []string           `json:"provider_auth_policy_tables,omitempty"`
	ProviderSchemas          []string           `json:"provider_schemas,omitempty"`
	ServerVersion            string             `json:"server_version"`
}

// ImportPreflightResult is the verdict on whether a source database can be
// imported into the project, with the evidence behind it.
type ImportPreflightResult struct {
	Checks            []ImportPreflightCheck `json:"checks"`
	OK                bool                   `json:"ok"`
	Source            SourceInspection       `json:"source"`
	StorageLimitBytes int64                  `json:"storage_limit_bytes"`
	TargetVersion     string                 `json:"target_version"`
}

// ImportUpload is a presigned object-storage PUT slot for a pg_dump file.
type ImportUpload struct {
	ExpiresAt time.Time `json:"expires_at"`
	ObjectKey string    `json:"object_key"`
	UploadURL string    `json:"upload_url"`
}

// CreateImportRequest starts an import from exactly one of a live source
// connection URL or a previously uploaded dump file's object key. Confirm is
// required by the API - an import writes over the project's live database.
type CreateImportRequest struct {
	Confirm   bool   `json:"confirm"`
	Recreate  bool   `json:"recreate,omitempty"`
	SourceURL string `json:"source_url,omitempty"`
	UploadKey string `json:"upload_key,omitempty"`
}

// IndexSuggestion is one candidate index the advisor derived from the
// project's real query predicates.
type IndexSuggestion struct {
	DDL string `json:"ddl"`
	// EstimatedCostReductionPct is how much cheaper the planner expects the
	// statement behind this candidate to become, 0..100, measured by planning
	// that statement with and without the index present hypothetically. It is
	// what the index BUYS, where EstimatedSizeBytes is what it COSTS.
	//
	// Nil when it could not be measured; nil is not zero, and zero means the
	// index was measured and would not help.
	EstimatedCostReductionPct *float64 `json:"estimated_cost_reduction_pct,omitempty"`
	// EstimatedSizeBytes is measured by building the index hypothetically -
	// nothing is written to the database. Zero when the estimate is unavailable.
	EstimatedSizeBytes int64  `json:"estimated_size_bytes,omitempty"`
	IndexMethod        string `json:"index_method,omitempty"`
	// QueryID identifies the statement the reduction was measured against,
	// matching the query identifier in query insights.
	QueryID *int64 `json:"query_id,omitempty"`
	Table   string `json:"table,omitempty"`
}

// IndexAdvisorReport is the index advisor's answer for one project.
type IndexAdvisorReport struct {
	Available bool `json:"available"`
	MinFilter int  `json:"min_filter"`
	// MissingExtensions lists what to enable before the advisor can run, or
	// before size estimates appear.
	MinSelectivity    int      `json:"min_selectivity"`
	MissingExtensions []string `json:"missing_extensions"`
	Reason            string   `json:"reason,omitempty"`
	// CostEstimatesAvailable reports whether any suggestion could be measured
	// against the statement it is meant to help. When false, the suggestions
	// are unranked candidates rather than measured recommendations.
	CostEstimatesAvailable bool              `json:"cost_estimates_available"`
	SizeEstimatesAvailable bool              `json:"size_estimates_available"`
	Suggestions            []IndexSuggestion `json:"suggestions"`
}

// UpdateProjectRequest is the PATCH /v1/projects/{id} body. Omitted fields are
// left unchanged.
type UpdateProjectRequest struct {
	Environment *string `json:"environment,omitempty"`
}

// UpdateWebhookEndpointRequest is the PATCH webhook-endpoint body. Omitted
// fields are left unchanged.
type UpdateWebhookEndpointRequest struct {
	Description *string   `json:"description,omitempty"`
	EventTypes  *[]string `json:"event_types,omitempty"`
	IsActive    *bool     `json:"is_active,omitempty"`
	URL         *string   `json:"url,omitempty"`
}

// CLILoginSessionStartRequest opens a browser login session for a CLI or agent.
type CLILoginSessionStartRequest struct {
	DeviceName string     `json:"device_name,omitempty"`
	ExpiresAt  *time.Time `json:"expires_at,omitempty"`
	Name       string     `json:"name,omitempty"`
	Scopes     []string   `json:"scopes,omitempty"`
	Source     string     `json:"source,omitempty"`
}

// CLILoginSessionStartResponse carries the session id plus the poll token that
// gates delivery of the minted key. The poll token is the sole secret in the
// flow - send it in a header, never a query string.
type CLILoginSessionStartResponse struct {
	ExpiresAt time.Time `json:"expires_at"`
	PollToken string    `json:"poll_token"`
	SessionID string    `json:"session_id"`
	State     string    `json:"state"`
}

// CLILoginSessionPollResponse is the poll result. PlaintextAPIKey is populated
// exactly once, on the poll that observes the authorization.
type CLILoginSessionPollResponse struct {
	AuthorizedAt     *time.Time `json:"authorized_at,omitempty"`
	ExpiresAt        time.Time  `json:"expires_at"`
	OrganizationID   string     `json:"organization_id,omitempty"`
	OrganizationName string     `json:"organization_name,omitempty"`
	OrganizationSlug string     `json:"organization_slug,omitempty"`
	PlaintextAPIKey  string     `json:"plaintext_api_key,omitempty"`
	SessionID        string     `json:"session_id"`
	State            string     `json:"state"`
}

// CLILoginSessionDetailsResponse describes a pending session to the browser
// page that is about to authorize it.
type CLILoginSessionDetailsResponse struct {
	DeviceName   string     `json:"device_name,omitempty"`
	ExpiresAt    time.Time  `json:"expires_at"`
	KeyExpiresAt *time.Time `json:"key_expires_at,omitempty"`
	KeyName      string     `json:"key_name"`
	Scopes       []string   `json:"scopes"`
	SessionID    string     `json:"session_id"`
	// Source is "cli" or "agent".
	Source string `json:"source"`
	// State is "pending" or "completed".
	State string `json:"state"`
}

// ProvisionCloudflareDatabaseRequest carries a Cloudflare-issued authorization
// (account id, issuance timestamp, signature) plus the shape of the database to
// create. The signature is the only authentication - this call carries no
// CapyDB credential.
type ProvisionCloudflareDatabaseRequest struct {
	AccountID       string `json:"account_id"`
	AccountName     string `json:"account_name,omitempty"`
	BillingPlan     string `json:"billing_plan,omitempty"`
	Name            string `json:"name"`
	PostgresVersion string `json:"postgres_version,omitempty"`
	Region          string `json:"region,omitempty"`
	Signature       string `json:"signature"`
	Timestamp       int64  `json:"timestamp"`
}

// ProvisionCloudflareDatabaseResponse is the created project plus the
// organization it landed in and the provision job to poll.
type ProvisionCloudflareDatabaseResponse struct {
	Job          Job          `json:"job"`
	Organization Organization `json:"organization"`
	Project      Project      `json:"project"`
}
