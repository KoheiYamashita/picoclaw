package skills

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSkillsInfoValidate(t *testing.T) {
	testcases := []struct {
		name        string
		skillName   string
		description string
		wantErr     bool
		errContains []string
	}{
		{
			name:        "valid-skill",
			skillName:   "valid-skill",
			description: "a valid skill description",
			wantErr:     false,
		},
		{
			name:        "empty-name",
			skillName:   "",
			description: "description without name",
			wantErr:     true,
			errContains: []string{"name is required"},
		},
		{
			name:        "empty-description",
			skillName:   "skill-without-description",
			description: "",
			wantErr:     true,
			errContains: []string{"description is required"},
		},
		{
			name:        "empty-both",
			skillName:   "",
			description: "",
			wantErr:     true,
			errContains: []string{"name is required", "description is required"},
		},
		{
			name:        "name-with-spaces",
			skillName:   "skill with spaces",
			description: "invalid name with spaces",
			wantErr:     true,
			errContains: []string{"name must be alphanumeric with hyphens"},
		},
		{
			name:        "name-with-underscore",
			skillName:   "skill_underscore",
			description: "invalid name with underscore",
			wantErr:     true,
			errContains: []string{"name must be alphanumeric with hyphens"},
		},
		{
			name:        "name-too-long",
			skillName:   strings.Repeat("a", MaxNameLength+1),
			description: "a valid description",
			wantErr:     true,
			errContains: []string{"name exceeds"},
		},
		{
			name:        "description-too-long",
			skillName:   "valid-name",
			description: strings.Repeat("x", MaxDescriptionLength+1),
			wantErr:     true,
			errContains: []string{"description exceeds"},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			info := SkillInfo{
				Name:        tc.skillName,
				Description: tc.description,
			}
			err := info.validate()
			if tc.wantErr {
				assert.Error(t, err)
				for _, msg := range tc.errContains {
					assert.ErrorContains(t, err, msg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// --- parseSimpleYAML ---

func TestParseSimpleYAML_SingleKV(t *testing.T) {
	sl := &SkillsLoader{}
	result := sl.parseSimpleYAML("name: my-skill")
	if result["name"] != "my-skill" {
		t.Errorf("name = %q, want %q", result["name"], "my-skill")
	}
}

func TestParseSimpleYAML_MultipleKV(t *testing.T) {
	sl := &SkillsLoader{}
	result := sl.parseSimpleYAML("name: my-skill\ndescription: does stuff")
	if result["name"] != "my-skill" {
		t.Errorf("name = %q", result["name"])
	}
	if result["description"] != "does stuff" {
		t.Errorf("description = %q", result["description"])
	}
}

func TestParseSimpleYAML_ColonInValue(t *testing.T) {
	sl := &SkillsLoader{}
	result := sl.parseSimpleYAML("description: does: many things")
	if result["description"] != "does: many things" {
		t.Errorf("description = %q, want %q", result["description"], "does: many things")
	}
}

func TestParseSimpleYAML_QuoteRemoval(t *testing.T) {
	sl := &SkillsLoader{}
	tests := []struct {
		input string
		want  string
	}{
		{`name: "quoted"`, "quoted"},
		{`name: 'single'`, "single"},
	}
	for _, tt := range tests {
		result := sl.parseSimpleYAML(tt.input)
		if result["name"] != tt.want {
			t.Errorf("parseSimpleYAML(%q)[name] = %q, want %q", tt.input, result["name"], tt.want)
		}
	}
}

func TestParseSimpleYAML_EmptyLines(t *testing.T) {
	sl := &SkillsLoader{}
	result := sl.parseSimpleYAML("name: test\n\n\ndescription: value")
	if len(result) != 2 {
		t.Errorf("expected 2 entries, got %d", len(result))
	}
}

func TestParseSimpleYAML_Comments(t *testing.T) {
	sl := &SkillsLoader{}
	result := sl.parseSimpleYAML("# comment\nname: test")
	if _, ok := result["# comment"]; ok {
		t.Error("comment should be skipped")
	}
	if result["name"] != "test" {
		t.Errorf("name = %q", result["name"])
	}
}

func TestParseSimpleYAML_EmptyInput(t *testing.T) {
	sl := &SkillsLoader{}
	result := sl.parseSimpleYAML("")
	if len(result) != 0 {
		t.Errorf("expected empty map, got %d entries", len(result))
	}
}

// --- extractFrontmatter ---

func TestExtractFrontmatter_Valid(t *testing.T) {
	sl := &SkillsLoader{}
	content := "---\nname: test\n---\nbody"
	result := sl.extractFrontmatter(content)
	if result != "name: test" {
		t.Errorf("frontmatter = %q, want %q", result, "name: test")
	}
}

func TestExtractFrontmatter_NoFrontmatter(t *testing.T) {
	sl := &SkillsLoader{}
	result := sl.extractFrontmatter("no frontmatter here")
	if result != "" {
		t.Errorf("expected empty, got %q", result)
	}
}

func TestExtractFrontmatter_Empty(t *testing.T) {
	sl := &SkillsLoader{}
	result := sl.extractFrontmatter("")
	if result != "" {
		t.Errorf("expected empty, got %q", result)
	}
}

// --- stripFrontmatter ---

func TestStripFrontmatter_WithFrontmatter(t *testing.T) {
	sl := &SkillsLoader{}
	content := "---\nname: test\n---\nbody content"
	result := sl.stripFrontmatter(content)
	if strings.Contains(result, "---") {
		t.Errorf("frontmatter should be stripped, got %q", result)
	}
	if !strings.Contains(result, "body content") {
		t.Errorf("body should be preserved, got %q", result)
	}
}

func TestStripFrontmatter_MultiLineFrontmatter(t *testing.T) {
	sl := &SkillsLoader{}
	content := "---\nname: test\ndescription: value\n---\nbody content"
	result := sl.stripFrontmatter(content)
	if strings.Contains(result, "---") {
		t.Errorf("multi-line frontmatter should be stripped, got %q", result)
	}
	if !strings.Contains(result, "body content") {
		t.Errorf("body should be preserved, got %q", result)
	}
}

func TestStripFrontmatter_WithoutFrontmatter(t *testing.T) {
	sl := &SkillsLoader{}
	content := "just body content"
	result := sl.stripFrontmatter(content)
	if result != content {
		t.Errorf("content should be unchanged, got %q", result)
	}
}

// --- escapeXML ---

func TestEscapeXML_Ampersand(t *testing.T) {
	got := escapeXML("a & b")
	if got != "a &amp; b" {
		t.Errorf("got %q", got)
	}
}

func TestEscapeXML_LessThan(t *testing.T) {
	got := escapeXML("a < b")
	if got != "a &lt; b" {
		t.Errorf("got %q", got)
	}
}

func TestEscapeXML_GreaterThan(t *testing.T) {
	got := escapeXML("a > b")
	if got != "a &gt; b" {
		t.Errorf("got %q", got)
	}
}

func TestEscapeXML_NoSpecialChars(t *testing.T) {
	input := "hello world"
	got := escapeXML(input)
	if got != input {
		t.Errorf("expected no change, got %q", got)
	}
}

// --- LoadSkill / ListSkills integration ---

func TestLoadSkill_FromWorkspace(t *testing.T) {
	tmpDir := t.TempDir()
	skillDir := filepath.Join(tmpDir, "skills", "test-skill")
	if err := os.MkdirAll(skillDir, 0755); err != nil {
		t.Fatal(err)
	}
	skillContent := "---\n{\"name\":\"test-skill\",\"description\":\"A test skill\"}\n---\nSkill body content"
	if err := os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte(skillContent), 0644); err != nil {
		t.Fatal(err)
	}

	sl := NewSkillsLoader(tmpDir, "", "")
	content, ok := sl.LoadSkill("test-skill")
	if !ok {
		t.Fatal("expected skill to be found")
	}
	if !strings.Contains(content, "Skill body content") {
		t.Errorf("expected body content, got %q", content)
	}
	if strings.Contains(content, "---") {
		t.Error("frontmatter should be stripped")
	}
}

func TestLoadSkill_NotFound(t *testing.T) {
	sl := NewSkillsLoader(t.TempDir(), "", "")
	_, ok := sl.LoadSkill("nonexistent")
	if ok {
		t.Error("expected skill not to be found")
	}
}

// --- ListSkills ---

func TestListSkills_WorkspaceSource(t *testing.T) {
	tmpDir := t.TempDir()
	skillDir := filepath.Join(tmpDir, "skills", "my-skill")
	if err := os.MkdirAll(skillDir, 0755); err != nil {
		t.Fatal(err)
	}
	content := "---\n{\"name\":\"my-skill\",\"description\":\"Test skill\"}\n---\nBody"
	if err := os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	sl := NewSkillsLoader(tmpDir, "", "")
	skills := sl.ListSkills()
	if len(skills) != 1 {
		t.Fatalf("expected 1 skill, got %d", len(skills))
	}
	if skills[0].Name != "my-skill" {
		t.Errorf("name = %q", skills[0].Name)
	}
	if skills[0].Source != "workspace" {
		t.Errorf("source = %q, want workspace", skills[0].Source)
	}
}

func TestListSkills_GlobalSource(t *testing.T) {
	dataDir := t.TempDir()
	globalDir := t.TempDir()
	skillDir := filepath.Join(globalDir, "global-skill")
	if err := os.MkdirAll(skillDir, 0755); err != nil {
		t.Fatal(err)
	}
	content := "---\n{\"name\":\"global-skill\",\"description\":\"A global skill\"}\n---\nBody"
	if err := os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	sl := NewSkillsLoader(dataDir, globalDir, "")
	skills := sl.ListSkills()
	if len(skills) != 1 {
		t.Fatalf("expected 1 skill, got %d", len(skills))
	}
	if skills[0].Source != "global" {
		t.Errorf("source = %q, want global", skills[0].Source)
	}
}

func TestListSkills_BuiltinSource(t *testing.T) {
	dataDir := t.TempDir()
	builtinDir := t.TempDir()
	skillDir := filepath.Join(builtinDir, "builtin-skill")
	if err := os.MkdirAll(skillDir, 0755); err != nil {
		t.Fatal(err)
	}
	content := "---\n{\"name\":\"builtin-skill\",\"description\":\"A builtin skill\"}\n---\nBody"
	if err := os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	sl := NewSkillsLoader(dataDir, "", builtinDir)
	skills := sl.ListSkills()
	if len(skills) != 1 {
		t.Fatalf("expected 1 skill, got %d", len(skills))
	}
	if skills[0].Source != "builtin" {
		t.Errorf("source = %q, want builtin", skills[0].Source)
	}
}

func TestListSkills_WorkspaceOverridesGlobal(t *testing.T) {
	dataDir := t.TempDir()
	globalDir := t.TempDir()

	// Create same-named skill in both workspace and global
	wsSkillDir := filepath.Join(dataDir, "skills", "shared-skill")
	if err := os.MkdirAll(wsSkillDir, 0755); err != nil {
		t.Fatal(err)
	}
	wsContent := "---\n{\"name\":\"shared-skill\",\"description\":\"Workspace version\"}\n---\nWorkspace"
	if err := os.WriteFile(filepath.Join(wsSkillDir, "SKILL.md"), []byte(wsContent), 0644); err != nil {
		t.Fatal(err)
	}

	glSkillDir := filepath.Join(globalDir, "shared-skill")
	if err := os.MkdirAll(glSkillDir, 0755); err != nil {
		t.Fatal(err)
	}
	glContent := "---\n{\"name\":\"shared-skill\",\"description\":\"Global version\"}\n---\nGlobal"
	if err := os.WriteFile(filepath.Join(glSkillDir, "SKILL.md"), []byte(glContent), 0644); err != nil {
		t.Fatal(err)
	}

	sl := NewSkillsLoader(dataDir, globalDir, "")
	skills := sl.ListSkills()
	if len(skills) != 1 {
		t.Fatalf("expected 1 skill (workspace wins), got %d", len(skills))
	}
	if skills[0].Source != "workspace" {
		t.Errorf("source = %q, want workspace", skills[0].Source)
	}
	if skills[0].Description != "Workspace version" {
		t.Errorf("description = %q, want Workspace version", skills[0].Description)
	}
}

func TestListSkills_InvalidSkillExcluded(t *testing.T) {
	tmpDir := t.TempDir()
	// Skill with invalid name (contains underscore)
	skillDir := filepath.Join(tmpDir, "skills", "bad_name")
	if err := os.MkdirAll(skillDir, 0755); err != nil {
		t.Fatal(err)
	}
	// No frontmatter, so name comes from dir name which is invalid
	if err := os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte("No frontmatter"), 0644); err != nil {
		t.Fatal(err)
	}

	sl := NewSkillsLoader(tmpDir, "", "")
	skills := sl.ListSkills()
	if len(skills) != 0 {
		t.Errorf("expected 0 skills (invalid excluded), got %d", len(skills))
	}
}

func TestListSkills_Empty(t *testing.T) {
	sl := NewSkillsLoader(t.TempDir(), "", "")
	skills := sl.ListSkills()
	if len(skills) != 0 {
		t.Errorf("expected 0 skills, got %d", len(skills))
	}
}

// --- getSkillMetadata ---

func TestGetSkillMetadata_JSONFrontmatter(t *testing.T) {
	tmpDir := t.TempDir()
	skillPath := filepath.Join(tmpDir, "SKILL.md")
	content := "---\n{\"name\":\"json-skill\",\"description\":\"JSON desc\"}\n---\nBody"
	if err := os.WriteFile(skillPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	sl := &SkillsLoader{}
	meta := sl.getSkillMetadata(skillPath)
	if meta == nil {
		t.Fatal("expected non-nil metadata")
	}
	if meta.Name != "json-skill" {
		t.Errorf("name = %q", meta.Name)
	}
	if meta.Description != "JSON desc" {
		t.Errorf("description = %q", meta.Description)
	}
}

func TestGetSkillMetadata_YAMLFrontmatter(t *testing.T) {
	tmpDir := t.TempDir()
	skillPath := filepath.Join(tmpDir, "SKILL.md")
	content := "---\nname: yaml-skill\ndescription: YAML desc\n---\nBody"
	if err := os.WriteFile(skillPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	sl := &SkillsLoader{}
	meta := sl.getSkillMetadata(skillPath)
	if meta == nil {
		t.Fatal("expected non-nil metadata")
	}
	if meta.Name != "yaml-skill" {
		t.Errorf("name = %q", meta.Name)
	}
	if meta.Description != "YAML desc" {
		t.Errorf("description = %q", meta.Description)
	}
}

func TestGetSkillMetadata_NoFrontmatter(t *testing.T) {
	tmpDir := t.TempDir()
	skillDir := filepath.Join(tmpDir, "my-fallback-skill")
	if err := os.MkdirAll(skillDir, 0755); err != nil {
		t.Fatal(err)
	}
	skillPath := filepath.Join(skillDir, "SKILL.md")
	if err := os.WriteFile(skillPath, []byte("No frontmatter body"), 0644); err != nil {
		t.Fatal(err)
	}

	sl := &SkillsLoader{}
	meta := sl.getSkillMetadata(skillPath)
	if meta == nil {
		t.Fatal("expected non-nil metadata")
	}
	if meta.Name != "my-fallback-skill" {
		t.Errorf("name = %q, want directory name", meta.Name)
	}
}

// --- BuildSkillsSummary ---

func TestBuildSkillsSummary_WithSkills(t *testing.T) {
	tmpDir := t.TempDir()
	skillDir := filepath.Join(tmpDir, "skills", "test-skill")
	if err := os.MkdirAll(skillDir, 0755); err != nil {
		t.Fatal(err)
	}
	content := "---\n{\"name\":\"test-skill\",\"description\":\"A <test> & skill\"}\n---\nBody"
	if err := os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	sl := NewSkillsLoader(tmpDir, "", "")
	summary := sl.BuildSkillsSummary()
	if summary == "" {
		t.Fatal("expected non-empty summary")
	}
	if !strings.Contains(summary, "<skills>") {
		t.Error("missing <skills> tag")
	}
	if !strings.Contains(summary, "test-skill") {
		t.Error("missing skill name")
	}
	// Check XML escaping
	if !strings.Contains(summary, "&lt;test&gt;") {
		t.Error("XML escaping not applied for <test>")
	}
	if !strings.Contains(summary, "&amp;") {
		t.Error("XML escaping not applied for &")
	}
}

func TestBuildSkillsSummary_NoSkills(t *testing.T) {
	sl := NewSkillsLoader(t.TempDir(), "", "")
	summary := sl.BuildSkillsSummary()
	if summary != "" {
		t.Errorf("expected empty summary, got %q", summary)
	}
}

// --- LoadSkillsForContext ---

func TestLoadSkillsForContext_MultipleSkills(t *testing.T) {
	tmpDir := t.TempDir()
	for _, name := range []string{"skill-a", "skill-b"} {
		dir := filepath.Join(tmpDir, "skills", name)
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatal(err)
		}
		content := "---\n{\"name\":\"" + name + "\",\"description\":\"desc\"}\n---\nBody of " + name
		if err := os.WriteFile(filepath.Join(dir, "SKILL.md"), []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	sl := NewSkillsLoader(tmpDir, "", "")
	result := sl.LoadSkillsForContext([]string{"skill-a", "skill-b"})
	if !strings.Contains(result, "skill-a") {
		t.Error("missing skill-a")
	}
	if !strings.Contains(result, "skill-b") {
		t.Error("missing skill-b")
	}
	if !strings.Contains(result, "\n\n---\n\n") {
		t.Error("missing separator between skills")
	}
}

func TestLoadSkillsForContext_NonexistentSkillExcluded(t *testing.T) {
	tmpDir := t.TempDir()
	dir := filepath.Join(tmpDir, "skills", "real-skill")
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatal(err)
	}
	content := "---\n{\"name\":\"real-skill\",\"description\":\"desc\"}\n---\nBody"
	if err := os.WriteFile(filepath.Join(dir, "SKILL.md"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	sl := NewSkillsLoader(tmpDir, "", "")
	result := sl.LoadSkillsForContext([]string{"real-skill", "nonexistent"})
	if !strings.Contains(result, "real-skill") {
		t.Error("missing real-skill")
	}
	if strings.Contains(result, "nonexistent") {
		t.Error("nonexistent skill should not appear")
	}
}

func TestLoadSkillsForContext_Empty(t *testing.T) {
	sl := NewSkillsLoader(t.TempDir(), "", "")
	result := sl.LoadSkillsForContext(nil)
	if result != "" {
		t.Errorf("expected empty, got %q", result)
	}
}
