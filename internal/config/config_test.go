package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad_CreatesDefault(t *testing.T) {
	dir := t.TempDir()
	m := &Manager{
		path: filepath.Join(dir, ".wlt", "config.yaml"),
		cfg:  defaultConfig(),
	}
	if err := m.Load(); err != nil {
		t.Fatal(err)
	}
	cfg := m.GetConfig()
	if cfg.Active != "sit" {
		t.Errorf("expected active=sit, got %s", cfg.Active)
	}
	if _, ok := cfg.Profiles["sit"]; !ok {
		t.Error("missing sit profile")
	}
	if _, ok := cfg.Profiles["prod"]; !ok {
		t.Error("missing prod profile")
	}
	// Verify file was created
	if _, err := os.Stat(m.path); os.IsNotExist(err) {
		t.Error("config file not created")
	}
}

func TestSaveAndReload(t *testing.T) {
	dir := t.TempDir()
	m := &Manager{
		path: filepath.Join(dir, ".wlt", "config.yaml"),
		cfg:  defaultConfig(),
	}
	if err := m.Load(); err != nil {
		t.Fatal(err)
	}
	// Modify and save
	p, _ := m.ActiveProfile()
	p.EnterpriseType = "test-enterprise"
	if err := m.Save(); err != nil {
		t.Fatal(err)
	}
	// Reload
	m2 := &Manager{
		path: m.path,
		cfg:  defaultConfig(),
	}
	if err := m2.Load(); err != nil {
		t.Fatal(err)
	}
	loaded, _ := m2.ActiveProfile()
	if loaded.EnterpriseType != "test-enterprise" {
		t.Errorf("expected enterprise_type=test-enterprise, got %s", loaded.EnterpriseType)
	}
}

func TestProfileOverride(t *testing.T) {
	dir := t.TempDir()
	m := &Manager{
		path: filepath.Join(dir, ".wlt", "config.yaml"),
		cfg:  defaultConfig(),
	}
	_ = m.Load()
	m.SetProfileOverride("prod")
	p, _ := m.ActiveProfile()
	if p.BaseURL != "https://erpapi.w-lian.com" {
		t.Errorf("expected prod URL, got %s", p.BaseURL)
	}
}

func TestSetActive_Invalid(t *testing.T) {
	dir := t.TempDir()
	m := &Manager{
		path: filepath.Join(dir, ".wlt", "config.yaml"),
		cfg:  defaultConfig(),
	}
	_ = m.Load()
	err := m.SetActive("nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent profile")
	}
}
