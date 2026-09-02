# Changelog

All notable changes to `github.com/capydatabase/capydbclient` are documented here. The format follows
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/); versions follow
[SemVer](https://semver.org/).

This module is the single Go mirror of the control-plane OpenAPI component schemas
(`backend/internal/httpapi/openapi.json`). Both the `capydb` CLI and the Terraform provider consume
it, so a shape added here must match the spec exactly.

## [Unreleased]

## [1.9.0] - 2026-09-02

### Added

- `SQLQueryRequest.ReadOnly` (`read_only`): asks the server to run the statement inside a
  `READ ONLY` transaction, refusing every write with SQLSTATE 25006.
- `ProjectApproval`: the response shape of `POST /v1/projects/{projectID}/approvals`, the mint
  call for single-use destructive-action approval tokens.

### Changed

- `CreateRestoreRequest`: `ConfirmProjectOverwrite` is replaced by `ApprovalToken`, carrying a
  single-use `project.restore_overwrite` approval token when `target_kind` is `"project"`.

## [1.8.0] - 2026-08-31

### Added

- `SQLQueryRequest.AllowUnqualifiedWrites`.
- `IndexSuggestion.EstimatedCostReductionPct` and `IndexSuggestion.QueryID`, plus
  `IndexAdvisorReport.CostEstimatesAvailable` - the advisor now reports what an index would buy,
  not only what it would cost to store.

## [1.7.0] - 2026-08-26

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

- Go directive raised to 1.27.1.

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

[Unreleased]: https://github.com/capy-base/capydbclient/compare/v1.9.0...HEAD
[1.9.0]: https://github.com/capy-base/capydbclient/compare/v1.8.0...v1.9.0
[1.8.0]: https://github.com/capy-base/capydbclient/compare/v1.7.0...v1.8.0
[1.7.0]: https://github.com/capy-base/capydbclient/compare/v1.6.0...v1.7.0
[1.6.0]: https://github.com/capy-base/capydbclient/compare/v1.6.0...v1.6.0
[1.6.0]: https://github.com/capy-base/capydbclient/compare/v1.5.0...v1.6.0
[1.5.0]: https://github.com/capy-base/capydbclient/compare/v1.4.0...v1.5.0
[1.4.0]: https://github.com/capy-base/capydbclient/compare/v1.3.5...v1.4.0
[1.3.5]: https://github.com/capy-base/capydbclient/compare/v1.3.0...v1.3.5
[1.3.0]: https://github.com/capy-base/capydbclient/compare/v1.2.0...v1.3.0
[1.2.0]: https://github.com/capy-base/capydbclient/releases/tag/v1.2.0
