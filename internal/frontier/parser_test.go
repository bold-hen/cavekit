package frontier

import (
	"os"
	"path/filepath"
	"testing"
)

const testFrontier = `---
created: "2026-03-17"
---

# Test Frontier

## Tier 0 — No Dependencies

| Task | Title | Spec | Requirement | Effort |
|------|-------|------|------------|--------|
| T-001 | Go module init | spec-cli.md | R1 | S |
| T-002 | Tmux session create | spec-tmux.md | R1 | M |

---

## Tier 1 — Depends on Tier 0

| Task | Title | Spec | Requirement | blockedBy | Effort |
|------|-------|------|------------|-----------|--------|
| T-009 | PTY-based attach | spec-tmux.md | R3 | T-002 | L |
| T-010 | Status detection | spec-tmux.md | R4 | T-002, T-003 | M |
`

func TestParse(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "build-site.md")
	os.WriteFile(path, []byte(testFrontier), 0644)

	f, err := Parse(path)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}

	if f.TotalTasks() != 4 {
		t.Errorf("TotalTasks() = %d, want 4", f.TotalTasks())
	}

	// Check tier counts
	if f.TierCounts[0] != 2 {
		t.Errorf("Tier 0 count = %d, want 2", f.TierCounts[0])
	}
	if f.TierCounts[1] != 2 {
		t.Errorf("Tier 1 count = %d, want 2", f.TierCounts[1])
	}

	// Check T-001
	t001 := f.TaskByID("T-001")
	if t001 == nil {
		t.Fatal("T-001 not found")
	}
	if t001.Title != "Go module init" {
		t.Errorf("T-001 Title = %q", t001.Title)
	}
	if t001.Spec != "spec-cli.md" {
		t.Errorf("T-001 Spec = %q", t001.Spec)
	}
	if t001.Requirement != "R1" {
		t.Errorf("T-001 Requirement = %q", t001.Requirement)
	}
	if t001.Effort != "S" {
		t.Errorf("T-001 Effort = %q", t001.Effort)
	}
	if t001.Tier != 0 {
		t.Errorf("T-001 Tier = %d", t001.Tier)
	}
	if len(t001.BlockedBy) != 0 {
		t.Errorf("T-001 BlockedBy should be empty, got %v", t001.BlockedBy)
	}

	// Check T-010 (has blockedBy)
	t010 := f.TaskByID("T-010")
	if t010 == nil {
		t.Fatal("T-010 not found")
	}
	if len(t010.BlockedBy) != 2 {
		t.Fatalf("T-010 BlockedBy = %v, want 2 items", t010.BlockedBy)
	}
	if t010.BlockedBy[0] != "T-002" || t010.BlockedBy[1] != "T-003" {
		t.Errorf("T-010 BlockedBy = %v", t010.BlockedBy)
	}
	if t010.Tier != 1 {
		t.Errorf("T-010 Tier = %d", t010.Tier)
	}
}

func TestParse_TaskIDPattern(t *testing.T) {
	// Verify the task ID pattern matches correctly
	tests := []struct {
		id    string
		match bool
	}{
		{"T-001", true},
		{"T-AUTH-001", true},
		{"T-A1-B2-C3", true},
		{"X-001", false},
		{"T001", false},
	}
	for _, tt := range tests {
		got := taskIDPattern.MatchString(tt.id)
		if got != tt.match {
			t.Errorf("taskIDPattern.MatchString(%q) = %v, want %v", tt.id, got, tt.match)
		}
	}
}

func TestTaskByID_NotFound(t *testing.T) {
	f := &Frontier{}
	if f.TaskByID("T-999") != nil {
		t.Error("TaskByID should return nil for non-existent task")
	}
}
