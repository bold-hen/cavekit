package site

import (
	"bufio"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// TaskID pattern: T-001, T-AUTH-001, etc.
var taskIDPattern = regexp.MustCompile(`T-([A-Za-z0-9]+-)*[A-Za-z0-9]+`)

// Task represents a single task from a site file.
type Task struct {
	ID          string
	Title       string
	Spec        string
	Requirement string
	BlockedBy   []string
	Effort      string
	Tier        int
	// Files is an optional list of file globs that this task is expected to
	// touch. When present, the team scheduler uses it for real path-overlap
	// detection instead of falling back to a substring heuristic, and
	// `cavekit team claim` defaults --paths to this list.
	Files []string
}

// Site represents a parsed site file.
type Site struct {
	Path       string
	Name       string
	Tasks      []Task
	TierCounts map[int]int // tier → count of tasks
}

// TotalTasks returns the total number of tasks.
func (s *Site) TotalTasks() int {
	return len(s.Tasks)
}

// TaskByID returns a task by its ID, or nil if not found.
func (s *Site) TaskByID(id string) *Task {
	for i := range s.Tasks {
		if s.Tasks[i].ID == id {
			return &s.Tasks[i]
		}
	}
	return nil
}

// Parse reads and parses a site markdown file.
func Parse(path string) (*Site, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	s := &Site{
		Path:       path,
		TierCounts: make(map[int]int),
	}

	currentTier := -1
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Detect tier headers: "## Tier N" or "## Tier N —"
		if strings.HasPrefix(line, "## Tier ") {
			parts := strings.Fields(line)
			if len(parts) >= 3 {
				if tier, err := strconv.Atoi(parts[2]); err == nil {
					currentTier = tier
				}
			}
			continue
		}

		// Parse table rows: | T-001 | Title | Spec | Req | ... |
		if strings.HasPrefix(line, "| T-") {
			task := parseTableRow(line, currentTier)
			if task != nil {
				s.Tasks = append(s.Tasks, *task)
				s.TierCounts[currentTier]++
			}
		}
	}

	return s, scanner.Err()
}

func parseTableRow(line string, tier int) *Task {
	cells := splitTableRow(line)
	if len(cells) < 4 {
		return nil
	}

	id := strings.TrimSpace(cells[0])
	if !taskIDPattern.MatchString(id) {
		return nil
	}

	task := &Task{
		ID:   id,
		Tier: tier,
	}

	task.Title = strings.TrimSpace(cells[1])

	if len(cells) > 2 {
		task.Spec = strings.TrimSpace(cells[2])
	}
	if len(cells) > 3 {
		task.Requirement = strings.TrimSpace(cells[3])
	}

	// The blockedBy and effort columns vary by tier (tier 0 has no blockedBy column).
	// An optional trailing "Files" column may be present on either tier.
	if tier == 0 {
		// | Task | Title | Spec | Requirement | Effort | Files? |
		if len(cells) > 4 {
			task.Effort = strings.TrimSpace(cells[4])
		}
		if len(cells) > 5 {
			task.Files = parseFiles(cells[5])
		}
	} else {
		// | Task | Title | Spec | Requirement | blockedBy | Effort | Files? |
		if len(cells) > 4 {
			task.BlockedBy = parseBlockedBy(cells[4])
		}
		if len(cells) > 5 {
			task.Effort = strings.TrimSpace(cells[5])
		}
		if len(cells) > 6 {
			task.Files = parseFiles(cells[6])
		}
	}

	return task
}

// parseFiles splits a cell into a list of path globs. Accepts comma or
// semicolon separators so patterns that include commas (e.g. brace expansion)
// can opt into `;`. Empty cells and placeholder dashes are treated as unset.
func parseFiles(cell string) []string {
	cell = strings.TrimSpace(cell)
	if cell == "" || cell == "-" || cell == "—" {
		return nil
	}
	sep := ","
	if strings.Contains(cell, ";") {
		sep = ";"
	}
	parts := strings.Split(cell, sep)
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		out = append(out, p)
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

func splitTableRow(line string) []string {
	// Split on | and trim
	parts := strings.Split(line, "|")
	var cells []string
	for _, p := range parts {
		trimmed := strings.TrimSpace(p)
		if trimmed != "" {
			cells = append(cells, trimmed)
		}
	}
	return cells
}

func parseBlockedBy(cell string) []string {
	cell = strings.TrimSpace(cell)
	if cell == "" {
		return nil
	}
	parts := strings.Split(cell, ",")
	var result []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if taskIDPattern.MatchString(p) {
			result = append(result, p)
		}
	}
	return result
}
