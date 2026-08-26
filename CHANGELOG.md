# Changelog

All notable changes to `github.com/capy-base/capydbclient` are documented here. The format follows
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/); versions follow
[SemVer](https://semver.org/).

This module is the single Go mirror of the control-plane OpenAPI component schemas
(`backend/internal/httpapi/openapi.json`). Both the `capydb` CLI and the Terraform provider consume
it, so a shape added here must match the spec exactly.

## [Unreleased] - ships as v1.7.0

### Added

- `ProjectExport` and `ExportDownload`, mirroring the new customer export endpoints
  (`POST/GET /v1/projects/{id}/exports`, `GET .../exports/{exportID}/download`): downloadable
  `pg_dump` custom-format archives that expire after 7 days.

## [1.6.0] - 2026-08-20

## [1.6.0] - 2026-08-20

### Added

- The entity and request types the CLI had been re-declaring locally now live here, so the two Go
  clients share one shape instead of each maintaining (and drifting from) its own copy:
  `Principal`, `PreviewDatabase`, `Backup`, `ScheduledBackup`, `UpsertScheduledBackupRequest`,
  `RestorePoint`, `CreateRestorePointRequest`, `CreateRestoreRequest`, `StatusResponse`,
  `StatusComponent`, `RegionStatus`, `ActiveQuerySample`, `SlowQuerySample`,
  `ProjectObservability`, `SQLQueryRequest`, `SQLQueryResult`, `ProjectLogEntry`, `ProjectLogs`,
  `ProjectAuditEvent`, `ProjectExtensionStatus`, `ProjectAlert`, `ProjectIntegration`,
  `ImportPreflightCheck`, `ImportPreflightResult`, `SourceExtension`, `SourceForeignKey`,
  `SourceInspection`, `ImportUpload`, `CreateImportRequest`, `IndexSuggestion`,
  `IndexAdvisorReport`, `UpdateProjectRequest`, `UpdateWebhookEndpointRequest`,
  `CLILoginSessionStartRequest`, `CLILoginSessionStartResponse`, `CLILoginSessionPollResponse`,
  `CLILoginSessionDetailsResponse`, `ProvisionCloudflareDatabaseRequest` and
  `ProvisionCloudflareDatabaseResponse`.
- `Viewer` now carries `Principal` alongside `Organization`, matching the full `GET /v1/me` payload.
  Callers that only read `Organization` are unaffected.
- `Organization.CloudflareAccountID`, set on organizations provisioned and billed through the
  Cloudflare Hyperdrive partner flow.

### Fixed

- `Backup` was missing `verified_at` and `verification_error`, so consumers could not report whether
  a backup had been proven restorable.
- `SourceInspection` was missing `event_triggers`.
- `ProjectAlert.limit_value` / `observed_value` and the `duration_ms` fields on `ActiveQuerySample`
  and `SQLQueryResult` are `int64`, matching the wire contract. The CLI's local copies declared them
  as floats.

### Changed

- Go directive raised to 1.27.0.

## [1.5.0] - 2026-08-05

### Added

- `SchemaColumn.Generated`, distinguishing `stored` from `virtual` generated columns (virtual
  columns exist only on Postgres 18+).
- Test coverage for `Doer` request shaping, retry behaviour, and error decoding.

## [1.4.0] - 2026-07-29

### Added

- `WebhookDelivery`, covering one delivery attempt against a webhook endpoint.

## [1.3.5] - 2026-07-22

### Added

- `DatabaseSchema` and the schema-introspection types (`SchemaNamespace`, `SchemaTable`,
  `SchemaColumn`, `SchemaEnum`, `SchemaExtension`, `SchemaForeignKey`, `SchemaUniqueConstraint`)
  plus `GeneratedTypes`, backing schema dump/diff and server-side type generation.

## [1.3.0] - 2026-07-16

### Added

- `PostgresVersion` on `Project` and `CreateProjectRequest`, for multi-major (16/17/18) placement.

## [1.2.0] and earlier

Initial extraction of the shared `Doer` transport, `NormalizeList`, `APIError`, and the core
`Organization` / `Project` / `Job` / `APIKey` / `WebhookEndpoint` / `ConnectionInfo` types out of the
CLI so the Terraform provider could reuse them.
