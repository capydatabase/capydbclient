package capydbclient

import (
	"encoding/json"
	"testing"
	"time"
)

// decode unmarshals a spec-shaped payload, failing the test on error.
func decode[T any](t *testing.T, payload string) T {
	t.Helper()
	var value T
	if err := json.Unmarshal([]byte(payload), &value); err != nil {
		t.Fatalf("decode %T: %v", value, err)
	}
	return value
}

// The payloads below are written against the OpenAPI component schemas, so a
// mistyped JSON tag shows up as a zero value here rather than in a consumer.

func TestDecodeProjectChannelAndConnections(t *testing.T) {
	project := decode[Project](t, `{"postgres_version":"20","postgres_channel":"beta","postgres_warning":"not for production"}`)
	if project.PostgresChannel != "beta" || project.PostgresWarning != "not for production" {
		t.Fatalf("project channel fields = %q/%q", project.PostgresChannel, project.PostgresWarning)
	}

	connections := decode[ConnectionInfo](t, `{"username":"owner","direct_url":"d","pooled_url":"p","app":{"username":"app_user","direct_url":"ad","pooled_url":"ap"}}`)
	if connections.App == nil || connections.App.Username != "app_user" || connections.App.DirectURL != "ad" || connections.App.PooledURL != "ap" {
		t.Fatalf("connections.app = %+v", connections.App)
	}
	if without := decode[ConnectionInfo](t, `{"username":"owner"}`); without.App != nil {
		t.Fatalf("connections without app decoded app = %+v", without.App)
	}

	status := decode[AppRoleStatus](t, `{"available":true,"enabled":true,"username":"app_user","created_at":"2026-09-30T10:00:00Z"}`)
	if !status.Available || !status.Enabled || status.CreatedAt == nil || status.RotatedAt != nil {
		t.Fatalf("app role status = %+v", status)
	}
}

func TestDecodeRegionsAndVersions(t *testing.T) {
	regions := decode[RegionsResponse](t, `{"regions":["eu-north-1"],"region_details":[{"id":"eu-north-1","display_name":"EU North 1","location":"Helsinki, Finland"}]}`)
	if len(regions.RegionDetails) != 1 || regions.RegionDetails[0].DisplayName != "EU North 1" || regions.RegionDetails[0].Location == "" {
		t.Fatalf("regions = %+v", regions)
	}

	versions := decode[PostgresVersionsResponse](t, `{"versions":[{"version":"17","channel":"stable","default":true,"production_ready":true},{"version":"19","channel":"current","default":false,"production_ready":true},{"version":"20","channel":"beta","default":false,"production_ready":false}]}`)
	if len(versions.Versions) != 3 || !versions.Versions[0].Default || !versions.Versions[1].ProductionReady || versions.Versions[2].ProductionReady || versions.Versions[2].Channel != "beta" {
		t.Fatalf("versions = %+v", versions)
	}
}

func TestDecodeRestoreAndUpgrade(t *testing.T) {
	restore := decode[CreateRestoreResponse](t, `{"job":{"id":"job_1","state":"queued"},"pitr":{"requested_restore_time":"2026-09-30T10:00:00Z","restore_time":"2026-09-30T09:58:00Z","restore_time_clamped":true}}`)
	if restore.Job.ID != "job_1" || restore.PITR == nil || !restore.PITR.RestoreTimeClamped {
		t.Fatalf("restore = %+v", restore)
	}
	if !restore.PITR.RestoreTime.Before(restore.PITR.RequestedRestoreTime) {
		t.Fatalf("restore times = %+v", restore.PITR)
	}
	if backup := decode[CreateRestoreResponse](t, `{"job":{"id":"job_2"}}`); backup.PITR != nil {
		t.Fatalf("backup restore decoded pitr = %+v", backup.PITR)
	}

	upgrade := decode[MajorUpgradeStatusResponse](t, `{"upgrade":{"from_major":"17","to_major":"18","state":"rollback_available","rollback_available_until":"2026-10-03T10:00:00Z","created_at":"2026-09-30T10:00:00Z","updated_at":"2026-09-30T10:00:00Z"}}`)
	if upgrade.Upgrade == nil || upgrade.Upgrade.ToMajor != "18" || upgrade.Upgrade.RollbackAvailableUntil == nil {
		t.Fatalf("upgrade = %+v", upgrade.Upgrade)
	}
	if none := decode[MajorUpgradeStatusResponse](t, `{"upgrade":null}`); none.Upgrade != nil {
		t.Fatalf("null upgrade decoded = %+v", none.Upgrade)
	}
}

func TestDecodeLintLogsAndObservability(t *testing.T) {
	body := decode[struct {
		Lint LintReport `json:"lint"`
	}](t, `{"lint":{"findings":[{"rule":"missing_primary_key","severity":"warning","object":"public.events","message":"no primary key","fix":"ALTER TABLE ..."}],"skipped":["unused_index: statistics too young"]}}`)
	if len(body.Lint.Findings) != 1 || body.Lint.Findings[0].Fix == "" || len(body.Lint.Skipped) != 1 {
		t.Fatalf("lint = %+v", body.Lint)
	}

	errs := decode[ProjectLogErrors](t, `{"groups":[{"key":"22008:date/time field value out of range: ?","sqlstate":"22008","class":"data_exception","retryable":false,"severity":"error","count":25,"variants":1,"first_seen":"2026-11-08T13:15:18Z","last_seen":"2026-11-08T14:02:13Z","message":"date/time field value out of range: \"1992-10-00\"","statement":"update \"games\" g set","hint":"Perhaps you need a different \"DateStyle\" setting.","context":"unnamed portal parameter $2256 = '...'","user":"usr_ludography","database":"db_ludography"}],"since":"2026-11-08T13:00:00Z","scanned":500,"occurrences":25,"truncated":true}`)
	if len(errs.Groups) != 1 || errs.Groups[0].SQLState != "22008" || errs.Groups[0].Count != 25 || errs.Groups[0].Retryable || !errs.Truncated || errs.Scanned != 500 {
		t.Fatalf("ProjectLogErrors decoded wrong: %+v", errs)
	}
	if errs.Groups[0].Context != "unnamed portal parameter $2256 = '...'" || errs.Groups[0].Class != "data_exception" {
		t.Fatalf("ProjectLogErrorGroup continuation fields decoded wrong: %+v", errs.Groups[0])
	}

	search := decode[ProjectLogSearch](t, `{"entries":[{"timestamp":"2026-09-30T10:00:00Z","severity":"error","message":"relation does not exist","cursor":"c1","sqlstate":"42P01","pid":4242,"user":"owner","database":"app"}],"next_cursor":"n1","truncated":true}`)
	entry := search.Entries[0]
	if entry.SQLState != "42P01" || entry.PID != 4242 || entry.User != "owner" || entry.Database != "app" || search.NextCursor != "n1" || !search.Truncated {
		t.Fatalf("search = %+v", search)
	}

	observability := decode[ProjectObservability](t, `{"wake":{"wakes":3,"timed_wakes":2,"p50_ms":148.5,"p95_ms":310,"max_ms":320,"window_hours":168}}`)
	wake := observability.Wake
	if wake == nil || wake.Wakes != 3 || wake.TimedWakes != 2 || wake.P50Ms == nil || *wake.P50Ms != 148.5 || wake.MaxMs == nil || *wake.MaxMs != 320 || wake.WindowHours != 168 {
		t.Fatalf("wake = %+v", wake)
	}
	if observability.LockWaits != nil || observability.AutovacuumPriority != nil {
		t.Fatalf("a payload without Postgres 19 readings decoded some: %+v", observability)
	}

	pg19 := decode[ProjectObservability](t, `{"slow_queries":[{"calls":8,"mean_time_ms":1,"query":"SELECT $1","rows":8,"total_time_ms":8,"generic_plan_calls":3,"custom_plan_calls":5}],"lock_waits":{"types":[{"lock_type":"relation","waits":2,"wait_time_ms":1500.5,"fastpath_exceeded":0,"stats_reset":"2026-11-12T15:03:11Z"}]},"autovacuum_priority":{"tables":[{"schema":"public","table":"events","toast":true,"score":1.2,"xid_score":0.1,"mxid_score":0,"vacuum_score":1.2,"vacuum_insert_score":0,"analyze_score":0.4,"do_vacuum":true,"do_analyze":false,"for_wraparound":false}]}}`)
	slow := pg19.SlowQueries[0]
	if slow.GenericPlanCalls == nil || *slow.GenericPlanCalls != 3 || slow.CustomPlanCalls == nil || *slow.CustomPlanCalls != 5 {
		t.Fatalf("plan calls decoded wrong: %+v", slow)
	}
	if pg19.LockWaits == nil || len(pg19.LockWaits.Types) != 1 || pg19.LockWaits.Types[0].WaitTimeMs != 1500.5 || pg19.LockWaits.Types[0].StatsReset == nil {
		t.Fatalf("lock waits decoded wrong: %+v", pg19.LockWaits)
	}
	if pg19.AutovacuumPriority == nil || len(pg19.AutovacuumPriority.Tables) != 1 {
		t.Fatalf("autovacuum priority decoded wrong: %+v", pg19.AutovacuumPriority)
	}
	if table := pg19.AutovacuumPriority.Tables[0]; !table.Toast || table.Score != 1.2 || !table.DoVacuum || table.Table != "events" {
		t.Fatalf("autovacuum table decoded wrong: %+v", table)
	}

	unmeasured := decode[ProjectWakeLatency](t, `{"wakes":0,"timed_wakes":0,"p50_ms":null,"p95_ms":null,"max_ms":null,"window_hours":168}`)
	if unmeasured.P50Ms != nil || unmeasured.P95Ms != nil || unmeasured.MaxMs != nil {
		t.Fatalf("unmeasured wake = %+v", unmeasured)
	}
}

func TestDecodeStatesPreferencesAndKeys(t *testing.T) {
	store := decode[KVStore](t, `{"state":"stopped","stopped_reason":"org_suspended"}`)
	if store.State != "stopped" || store.StoppedReason != "org_suspended" {
		t.Fatalf("kv store = %+v", store)
	}

	key := decode[APIKey](t, `{"id":"key_1","manager":true}`)
	if !key.Manager {
		t.Fatalf("api key manager = %v", key.Manager)
	}

	preferences := decode[NotificationPreferences](t, `{"organization_id":"org_1","alert_emails_enabled":false,"alert_email_recipients":["ops@example.com"],"billing_email_recipients":[],"updated_at":null}`)
	if preferences.AlertEmailsEnabled || len(preferences.AlertEmailRecipients) != 1 || preferences.UpdatedAt != nil {
		t.Fatalf("preferences = %+v", preferences)
	}

	encoded, err := json.Marshal(PutNotificationPreferencesRequest{AlertEmailsEnabled: true, AlertEmailRecipients: []string{}, BillingEmailRecipients: []string{"billing@example.com"}})
	if err != nil {
		t.Fatal(err)
	}
	if want := `{"alert_email_recipients":[],"alert_emails_enabled":true,"billing_email_recipients":["billing@example.com"]}`; string(encoded) != want {
		t.Fatalf("put body = %s, want %s", encoded, want)
	}

	history := decode[StatusHistoryResponse](t, `{"days":2,"from":"2026-09-29","to":"2026-09-30","generated_at":"2026-09-30T10:00:00Z","regions":[{"region":"eu-north-1","uptime_percent":99.5,"days":[{"date":"2026-09-29","status":null,"uptime_percent":null},{"date":"2026-09-30","status":"degraded","uptime_percent":99}]}],"incidents":[{"id":"inc_1","title":"t","impact":"minor","status":"resolved","regions":[],"started_at":"2026-09-30T08:00:00Z","resolved_at":"2026-09-30T09:00:00Z","created_at":"2026-09-30T08:00:00Z","updated_at":"2026-09-30T09:00:00Z","updates":[{"id":"u1","status":"resolved","message":"fixed","created_at":"2026-09-30T09:00:00Z"}]}]}`)
	region := history.Regions[0]
	if region.Days[0].Status != nil || region.Days[1].Status == nil || *region.Days[1].Status != "degraded" || history.Incidents[0].ResolvedAt == nil {
		t.Fatalf("history = %+v", history)
	}
	if !history.GeneratedAt.Equal(time.Date(2026, 9, 30, 10, 0, 0, 0, time.UTC)) {
		t.Fatalf("history generated_at = %v", history.GeneratedAt)
	}
}
