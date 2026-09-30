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
	project := decode[Project](t, `{"postgres_version":"19","postgres_channel":"beta","postgres_warning":"not for production"}`)
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

	versions := decode[PostgresVersionsResponse](t, `{"versions":[{"version":"17","channel":"stable","default":true,"production_ready":true},{"version":"19","channel":"beta","default":false,"production_ready":false}]}`)
	if len(versions.Versions) != 2 || !versions.Versions[0].Default || versions.Versions[1].ProductionReady || versions.Versions[1].Channel != "beta" {
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
