// Package config manages the CLI configuration file (~/.wlt/config.yaml).
package config

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"gopkg.in/yaml.v3"
)

// Config represents the top-level configuration.
type Config struct {
	Active   string             `yaml:"active"`
	Profiles map[string]*Profile `yaml:"profiles"`
}

// Profile represents a single environment profile.
// It holds connection info only (where to reach the backend).
// Authentication (token) and tenant identity (tenant-id) are supplied
// per-call via flags — see internal/cmdutil. Nothing secret is stored here.
type Profile struct {
	BaseURL        string `yaml:"base_url"`
	APIPrefix      string `yaml:"api_prefix"`
	EnterpriseType string `yaml:"enterprise_type"`
}

// Manager manages configuration loading and saving.
type Manager struct {
	path            string
	cfg             *Config
	profileOverride string
	mu              sync.RWMutex
}

// NewManager creates a new config Manager.
// Uses ~/.wlt/config.yaml as default path.
func NewManager() *Manager {
	home, _ := os.UserHomeDir()
	return &Manager{
		path: filepath.Join(home, ".wlt", "config.yaml"),
		cfg:  defaultConfig(),
	}
}

// NewManagerWithPath creates a config Manager with a custom path (for testing).
func NewManagerWithPath(p string) *Manager {
	return &Manager{
		path: p,
		cfg:  defaultConfig(),
	}
}

// SetProfileOverride overrides the active profile (from --profile flag).
func (m *Manager) SetProfileOverride(name string) {
	m.profileOverride = name
}

// Load reads the configuration file. Creates default config if file doesn't exist.
func (m *Manager) Load() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	data, err := os.ReadFile(m.path)
	if err != nil {
		if os.IsNotExist(err) {
			m.cfg = defaultConfig()
			return m.saveLocked()
		}
		return fmt.Errorf("read config: %w", err)
	}
	if err := yaml.Unmarshal(data, m.cfg); err != nil {
		return fmt.Errorf("parse config: %w", err)
	}
	return nil
}

// Save persists the configuration to disk.
func (m *Manager) Save() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.saveLocked()
}

func (m *Manager) saveLocked() error {
	dir := filepath.Dir(m.path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("create config dir: %w", err)
	}
	data, err := yaml.Marshal(m.cfg)
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}
	return os.WriteFile(m.path, data, 0600)
}

// ActiveProfile returns the currently active profile.
// Returns an error if the profile does not exist.
func (m *Manager) ActiveProfile() (*Profile, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	name := m.profileOverride
	if name == "" {
		name = m.cfg.Active
	}
	p, ok := m.cfg.Profiles[name]
	if !ok {
		return nil, fmt.Errorf("profile %q not found", name)
	}
	return p, nil
}

// GetActive returns the active profile name.
func (m *Manager) GetActive() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.cfg.Active
}

// SetActive sets the active profile name.
func (m *Manager) SetActive(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.cfg.Profiles[name]; !ok {
		return fmt.Errorf("profile %q not found", name)
	}
	m.cfg.Active = name
	return nil
}

// SetProfile updates or creates a profile.
func (m *Manager) SetProfile(name string, p *Profile) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.cfg.Profiles[name] = p
	return nil
}

// UpdateProfileField updates a single field in a profile and saves.
func (m *Manager) UpdateProfileField(profileName, field, value string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	p, ok := m.cfg.Profiles[profileName]
	if !ok {
		return fmt.Errorf("profile %q not found", profileName)
	}
	switch field {
	case "base_url":
		p.BaseURL = value
	case "api_prefix":
		p.APIPrefix = value
	case "enterprise_type":
		p.EnterpriseType = value
	default:
		return fmt.Errorf("unknown field: %s", field)
	}
	return m.saveLocked()
}

// GetConfig returns a deep copy of the full config (for display).
func (m *Manager) GetConfig() *Config {
	m.mu.RLock()
	defer m.mu.RUnlock()
	cp := *m.cfg
	cp.Profiles = make(map[string]*Profile, len(m.cfg.Profiles))
	for k, v := range m.cfg.Profiles {
		p := *v
		cp.Profiles[k] = &p
	}
	return &cp
}

// Path returns the config file path.
func (m *Manager) Path() string {
	return m.path
}

func defaultConfig() *Config {
	return &Config{
		Active: "sit",
		Profiles: map[string]*Profile{
			"sit": {
				BaseURL:   "https://erpsit.api.w-lian.com",
				APIPrefix: "/admin-api",
			},
			"prod": {
				BaseURL:   "https://erpapi.w-lian.com",
				APIPrefix: "/admin-api",
			},
		},
	}
}
