package site

import (
	"os"
	"path/filepath"
	"testing"
)

const testSite = `---
created: "2026-03-17"
---

# Test Site

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
	os.WriteFile(path, []byte(testSite), 0644)

	s, err := Parse(path)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}

	if s.TotalTasks() != 4 {
		t.Errorf("TotalTasks() = %d, want 4", s.TotalTasks())
	}

	// Check tier counts
	if s.TierCounts[0] != 2 {
		t.Errorf("Tier 0 count = %d, want 2", s.TierCounts[0])
	}
	if s.TierCounts[1] != 2 {
		t.Errorf("Tier 1 count = %d, want 2", s.TierCounts[1])
	}

	// Check T-001
	t001 := s.TaskByID("T-001")
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
	t010 := s.TaskByID("T-010")
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
	s := &Site{}
	if s.TaskByID("T-999") != nil {
		t.Error("TaskByID should return nil for non-existent task")
	}
}

const testSiteWithFiles = `# Site with Files

## Tier 0 — No Dependencies

| Task | Title | Spec | Requirement | Effort | Files |
|------|-------|------|------------|--------|-------|
| T-001 | Auth module | spec.md | R1 | S | src/auth/**, tests/auth/** |
| T-002 | No paths | spec.md | R1 | M | - |

## Tier 1 — Depends on Tier 0

| Task | Title | Spec | Requirement | blockedBy | Effort | Files |
|------|-------|------|------------|-----------|--------|-------|
| T-003 | DB layer | spec.md | R2 | T-001 | M | internal/db/**; migrations/** |
`

func TestParse_FilesColumn(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "build-site.md")
	if err := os.WriteFile(path, []byte(testSiteWithFiles), 0644); err != nil {
		t.Fatalf("write: %v", err)
	}
	s, err := Parse(path)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}

	t001 := s.TaskByID("T-001")
	if t001 == nil {
		t.Fatal("T-001 not found")
	}
	if len(t001.Files) != 2 || t001.Files[0] != "src/auth/**" || t001.Files[1] != "tests/auth/**" {
		t.Errorf("T-001 Files = %v", t001.Files)
	}

	t002 := s.TaskByID("T-002")
	if t002 == nil {
		t.Fatal("T-002 not found")
	}
	if len(t002.Files) != 0 {
		t.Errorf("T-002 Files should be empty, got %v", t002.Files)
	}

	t003 := s.TaskByID("T-003")
	if t003 == nil {
		t.Fatal("T-003 not found")
	}
	if len(t003.Files) != 2 || t003.Files[0] != "internal/db/**" || t003.Files[1] != "migrations/**" {
		t.Errorf("T-003 Files = %v", t003.Files)
	}
}
