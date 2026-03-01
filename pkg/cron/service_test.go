package cron

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"
)

func TestSaveStore_FilePermissions(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("file permission bits are not enforced on Windows")
	}

	tmpDir := t.TempDir()
	storePath := filepath.Join(tmpDir, "cron", "jobs.json")

	cs := NewCronService(storePath, nil)

	_, err := cs.AddJob("test", CronSchedule{Kind: "every", EveryMS: int64Ptr(60000)}, "hello", false, "cli", "direct")
	if err != nil {
		t.Fatalf("AddJob failed: %v", err)
	}

	info, err := os.Stat(storePath)
	if err != nil {
		t.Fatalf("Stat failed: %v", err)
	}

	perm := info.Mode().Perm()
	if perm != 0600 {
		t.Errorf("cron store has permission %04o, want 0600", perm)
	}
}

func int64Ptr(v int64) *int64 {
	return &v
}

// --- computeNextRun ---

func TestComputeNextRun_AtFuture(t *testing.T) {
	cs := NewCronService(filepath.Join(t.TempDir(), "cron.json"), nil)
	futureMS := time.Now().Add(time.Hour).UnixMilli()
	schedule := &CronSchedule{Kind: "at", AtMS: &futureMS}
	nowMS := time.Now().UnixMilli()

	result := cs.computeNextRun(schedule, nowMS)
	if result == nil {
		t.Fatal("expected non-nil result for future 'at' schedule")
	}
	if *result != futureMS {
		t.Errorf("result = %d, want %d", *result, futureMS)
	}
}

func TestComputeNextRun_AtPast(t *testing.T) {
	cs := NewCronService(filepath.Join(t.TempDir(), "cron.json"), nil)
	pastMS := time.Now().Add(-time.Hour).UnixMilli()
	schedule := &CronSchedule{Kind: "at", AtMS: &pastMS}
	nowMS := time.Now().UnixMilli()

	result := cs.computeNextRun(schedule, nowMS)
	if result != nil {
		t.Error("expected nil for past 'at' schedule")
	}
}

func TestComputeNextRun_Every(t *testing.T) {
	cs := NewCronService(filepath.Join(t.TempDir(), "cron.json"), nil)
	everyMS := int64(60000) // 1 minute
	schedule := &CronSchedule{Kind: "every", EveryMS: &everyMS}
	nowMS := int64(1_000_000)

	result := cs.computeNextRun(schedule, nowMS)
	if result == nil {
		t.Fatal("expected non-nil result")
	}
	if *result != 1_060_000 {
		t.Errorf("result = %d, want %d", *result, int64(1_060_000))
	}
}

func TestComputeNextRun_EveryZero(t *testing.T) {
	cs := NewCronService(filepath.Join(t.TempDir(), "cron.json"), nil)
	everyMS := int64(0)
	schedule := &CronSchedule{Kind: "every", EveryMS: &everyMS}

	result := cs.computeNextRun(schedule, time.Now().UnixMilli())
	if result != nil {
		t.Error("expected nil for zero interval")
	}
}

func TestComputeNextRun_CronExpr(t *testing.T) {
	cs := NewCronService(filepath.Join(t.TempDir(), "cron.json"), nil)
	schedule := &CronSchedule{Kind: "cron", Expr: "* * * * *"}
	nowMS := time.Now().UnixMilli()

	result := cs.computeNextRun(schedule, nowMS)
	if result == nil {
		t.Fatal("expected non-nil result for valid cron expr")
	}
	if *result <= nowMS {
		t.Error("next run should be in the future")
	}
}

func TestComputeNextRun_CronEmptyExpr(t *testing.T) {
	cs := NewCronService(filepath.Join(t.TempDir(), "cron.json"), nil)
	schedule := &CronSchedule{Kind: "cron", Expr: ""}

	result := cs.computeNextRun(schedule, time.Now().UnixMilli())
	if result != nil {
		t.Error("expected nil for empty cron expression")
	}
}

func TestComputeNextRun_UnknownKind(t *testing.T) {
	cs := NewCronService(filepath.Join(t.TempDir(), "cron.json"), nil)
	schedule := &CronSchedule{Kind: "unknown"}

	result := cs.computeNextRun(schedule, time.Now().UnixMilli())
	if result != nil {
		t.Error("expected nil for unknown schedule kind")
	}
}

// --- AddJob ---

func TestAddJob(t *testing.T) {
	storePath := filepath.Join(t.TempDir(), "cron.json")
	cs := NewCronService(storePath, nil)

	job, err := cs.AddJob("test-job", CronSchedule{Kind: "every", EveryMS: int64Ptr(60000)}, "hello", false, "cli", "direct")
	if err != nil {
		t.Fatalf("AddJob failed: %v", err)
	}
	if job.Name != "test-job" {
		t.Errorf("Name = %q, want %q", job.Name, "test-job")
	}
	if job.ID == "" {
		t.Error("ID should not be empty")
	}
	if !job.Enabled {
		t.Error("job should be enabled")
	}
	if job.DeleteAfterRun {
		t.Error("'every' jobs should have DeleteAfterRun=false")
	}

	// Verify persisted
	if _, err := os.Stat(storePath); err != nil {
		t.Errorf("store file should exist: %v", err)
	}
}

func TestAddJob_AtKind_DeleteAfterRun(t *testing.T) {
	cs := NewCronService(filepath.Join(t.TempDir(), "cron.json"), nil)
	futureMS := time.Now().Add(time.Hour).UnixMilli()

	job, err := cs.AddJob("one-time", CronSchedule{Kind: "at", AtMS: &futureMS}, "task", false, "", "")
	if err != nil {
		t.Fatal(err)
	}
	if !job.DeleteAfterRun {
		t.Error("'at' jobs should have DeleteAfterRun=true")
	}
}

// --- RemoveJob ---

func TestRemoveJob(t *testing.T) {
	cs := NewCronService(filepath.Join(t.TempDir(), "cron.json"), nil)

	job, err := cs.AddJob("to-remove", CronSchedule{Kind: "every", EveryMS: int64Ptr(60000)}, "msg", false, "", "")
	if err != nil {
		t.Fatal(err)
	}

	removed := cs.RemoveJob(job.ID)
	if !removed {
		t.Error("RemoveJob should return true")
	}

	jobs := cs.ListJobs(true)
	if len(jobs) != 0 {
		t.Errorf("expected 0 jobs, got %d", len(jobs))
	}
}

func TestRemoveJob_NotFound(t *testing.T) {
	cs := NewCronService(filepath.Join(t.TempDir(), "cron.json"), nil)

	removed := cs.RemoveJob("nonexistent")
	if removed {
		t.Error("RemoveJob should return false for nonexistent ID")
	}
}

// --- generateID ---

func TestGenerateID_Unique(t *testing.T) {
	seen := make(map[string]bool)
	for i := 0; i < 100; i++ {
		id := generateID()
		if id == "" {
			t.Fatal("generateID returned empty string")
		}
		if seen[id] {
			t.Fatalf("duplicate ID generated: %s", id)
		}
		seen[id] = true
	}
}

func TestGenerateID_Length(t *testing.T) {
	id := generateID()
	// hex.EncodeToString(8 bytes) = 16 hex chars
	if len(id) != 16 {
		t.Errorf("expected 16 chars, got %d (%q)", len(id), id)
	}
}

// --- ListJobs ---

func TestListJobs_IncludeDisabled(t *testing.T) {
	cs := NewCronService(filepath.Join(t.TempDir(), "cron.json"), nil)

	_, _ = cs.AddJob("enabled", CronSchedule{Kind: "every", EveryMS: int64Ptr(60000)}, "msg", false, "", "")
	job2, _ := cs.AddJob("to-disable", CronSchedule{Kind: "every", EveryMS: int64Ptr(60000)}, "msg", false, "", "")
	cs.EnableJob(job2.ID, false)

	all := cs.ListJobs(true)
	enabled := cs.ListJobs(false)

	if len(all) != 2 {
		t.Errorf("expected 2 total jobs, got %d", len(all))
	}
	if len(enabled) != 1 {
		t.Errorf("expected 1 enabled job, got %d", len(enabled))
	}
}

// --- Start / Stop ---

func TestStartStop(t *testing.T) {
	storePath := filepath.Join(t.TempDir(), "cron.json")
	cs := NewCronService(storePath, nil)

	if err := cs.Start(); err != nil {
		t.Fatalf("Start failed: %v", err)
	}

	cs.mu.RLock()
	running := cs.running
	cs.mu.RUnlock()
	if !running {
		t.Error("expected running=true after Start")
	}

	// Start again should be no-op
	if err := cs.Start(); err != nil {
		t.Fatalf("second Start failed: %v", err)
	}

	cs.Stop()

	cs.mu.RLock()
	running = cs.running
	cs.mu.RUnlock()
	if running {
		t.Error("expected running=false after Stop")
	}

	// Stop again should be no-op
	cs.Stop()
}

// --- UpdateJob ---

func TestUpdateJob_Success(t *testing.T) {
	storePath := filepath.Join(t.TempDir(), "cron.json")
	cs := NewCronService(storePath, nil)

	job, err := cs.AddJob("original", CronSchedule{Kind: "every", EveryMS: int64Ptr(60000)}, "msg", false, "", "")
	if err != nil {
		t.Fatal(err)
	}

	job.Name = "updated"
	if err := cs.UpdateJob(job); err != nil {
		t.Fatalf("UpdateJob failed: %v", err)
	}

	jobs := cs.ListJobs(true)
	if len(jobs) != 1 {
		t.Fatalf("expected 1 job, got %d", len(jobs))
	}
	if jobs[0].Name != "updated" {
		t.Errorf("Name = %q, want 'updated'", jobs[0].Name)
	}
}

func TestUpdateJob_NotFound(t *testing.T) {
	cs := NewCronService(filepath.Join(t.TempDir(), "cron.json"), nil)

	fakeJob := &CronJob{ID: "nonexistent"}
	err := cs.UpdateJob(fakeJob)
	if err == nil {
		t.Error("expected error for nonexistent job")
	}
}

// --- Status ---

func TestStatus(t *testing.T) {
	storePath := filepath.Join(t.TempDir(), "cron.json")
	cs := NewCronService(storePath, nil)

	status := cs.Status()
	if status["enabled"] != false {
		t.Error("expected enabled=false before Start")
	}
	if status["jobs"] != 0 {
		t.Errorf("expected 0 jobs, got %v", status["jobs"])
	}

	job, _ := cs.AddJob("test", CronSchedule{Kind: "every", EveryMS: int64Ptr(60000)}, "msg", false, "", "")
	_, _ = cs.AddJob("test2", CronSchedule{Kind: "every", EveryMS: int64Ptr(60000)}, "msg", false, "", "")
	cs.EnableJob(job.ID, false)

	status = cs.Status()
	if status["jobs"] != 2 {
		t.Errorf("expected 2 total jobs, got %v", status["jobs"])
	}
}

func TestStatus_Running(t *testing.T) {
	storePath := filepath.Join(t.TempDir(), "cron.json")
	cs := NewCronService(storePath, nil)

	if err := cs.Start(); err != nil {
		t.Fatal(err)
	}
	defer cs.Stop()

	status := cs.Status()
	if status["enabled"] != true {
		t.Error("expected enabled=true after Start")
	}
}

// --- Load ---

func TestLoad(t *testing.T) {
	storePath := filepath.Join(t.TempDir(), "cron.json")
	cs := NewCronService(storePath, nil)

	_, err := cs.AddJob("persist-test", CronSchedule{Kind: "every", EveryMS: int64Ptr(60000)}, "msg", false, "", "")
	if err != nil {
		t.Fatal(err)
	}

	// Create a new service pointing at the same store file
	cs2 := NewCronService(storePath, nil)
	if err := cs2.Load(); err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	jobs := cs2.ListJobs(true)
	if len(jobs) != 1 {
		t.Fatalf("expected 1 job after Load, got %d", len(jobs))
	}
	if jobs[0].Name != "persist-test" {
		t.Errorf("Name = %q", jobs[0].Name)
	}
	if !jobs[0].Enabled {
		t.Error("Enabled should be true after round-trip")
	}
	if jobs[0].Schedule.Kind != "every" {
		t.Errorf("Schedule.Kind = %q, want every", jobs[0].Schedule.Kind)
	}
	if jobs[0].ID == "" {
		t.Error("ID should survive round-trip")
	}
}

// --- getNextWakeMS ---

func TestGetNextWakeMS_NoJobs(t *testing.T) {
	cs := NewCronService(filepath.Join(t.TempDir(), "cron.json"), nil)
	result := cs.getNextWakeMS()
	if result != nil {
		t.Error("expected nil for no jobs")
	}
}

func TestGetNextWakeMS_OneEnabledJob(t *testing.T) {
	cs := NewCronService(filepath.Join(t.TempDir(), "cron.json"), nil)
	_, _ = cs.AddJob("job1", CronSchedule{Kind: "every", EveryMS: int64Ptr(60000)}, "msg", false, "", "")

	result := cs.getNextWakeMS()
	if result == nil {
		t.Fatal("expected non-nil for enabled job with next run")
	}
}

func TestGetNextWakeMS_MultipleJobs_ReturnsEarliest(t *testing.T) {
	cs := NewCronService(filepath.Join(t.TempDir(), "cron.json"), nil)
	// Add two jobs with different intervals
	_, _ = cs.AddJob("long", CronSchedule{Kind: "every", EveryMS: int64Ptr(600000)}, "msg", false, "", "")
	_, _ = cs.AddJob("short", CronSchedule{Kind: "every", EveryMS: int64Ptr(10000)}, "msg", false, "", "")

	result := cs.getNextWakeMS()
	if result == nil {
		t.Fatal("expected non-nil")
	}

	// The shorter interval job should have the earliest next wake
	jobs := cs.ListJobs(true)
	var shortNext *int64
	for _, j := range jobs {
		if j.Name == "short" {
			shortNext = j.State.NextRunAtMS
		}
	}
	if shortNext == nil {
		t.Fatal("short job should have NextRunAtMS")
	}
	if *result != *shortNext {
		t.Errorf("nextWake = %d, want %d (short job)", *result, *shortNext)
	}
}

// --- EnableJob ---

func TestEnableJob_Enable(t *testing.T) {
	cs := NewCronService(filepath.Join(t.TempDir(), "cron.json"), nil)
	job, _ := cs.AddJob("test", CronSchedule{Kind: "every", EveryMS: int64Ptr(60000)}, "msg", false, "", "")

	// Disable then re-enable
	cs.EnableJob(job.ID, false)
	result := cs.EnableJob(job.ID, true)
	if result == nil {
		t.Fatal("expected non-nil result")
	}
	if !result.Enabled {
		t.Error("expected enabled=true")
	}
	if result.State.NextRunAtMS == nil {
		t.Error("expected NextRunAtMS to be set after enable")
	}
}

func TestEnableJob_Disable(t *testing.T) {
	cs := NewCronService(filepath.Join(t.TempDir(), "cron.json"), nil)
	job, _ := cs.AddJob("test", CronSchedule{Kind: "every", EveryMS: int64Ptr(60000)}, "msg", false, "", "")

	result := cs.EnableJob(job.ID, false)
	if result == nil {
		t.Fatal("expected non-nil result")
	}
	if result.Enabled {
		t.Error("expected enabled=false")
	}
	if result.State.NextRunAtMS != nil {
		t.Error("expected NextRunAtMS=nil after disable")
	}
}

func TestEnableJob_NotFound(t *testing.T) {
	cs := NewCronService(filepath.Join(t.TempDir(), "cron.json"), nil)
	result := cs.EnableJob("nonexistent", true)
	if result != nil {
		t.Error("expected nil for nonexistent job")
	}
}
