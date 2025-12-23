package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Project struct {
	Name   string `json:"name"`
	Alias  string `json:"alias,omitempty"`
	Path   string `json:"path"`
	Status string `json:"status"`
}

type Config struct {
	Workspaces []string           `json:"workspaces"`
	Projects   map[string]Project `json:"projects"`
}

func getConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	configDir := filepath.Join(home, ".config", "orbit")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", err
	}
	return filepath.Join(configDir, "orbit.json"), nil
}

func Load() (*Config, error) {
	configPath, err := getConfigPath()
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return &Config{
			Workspaces: []string{},
			Projects:   make(map[string]Project),
		}, nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	if cfg.Projects == nil {
		cfg.Projects = make(map[string]Project)
	}

	return &cfg, nil
}

func Save(cfg *Config) error {
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}

func FindProjectPath(cfg *Config, projectName string) string {

	if project, exists := cfg.Projects[projectName]; exists {
		return project.Path
	}

	for _, workspace := range cfg.Workspaces {
		repoPath := filepath.Join(workspace, "project", "repo", projectName)
		if info, err := os.Stat(repoPath); err == nil && info.IsDir() {
			return filepath.Join(workspace, "project", projectName)
		}
	}

	return ""
}

func GetAllProjects(cfg *Config) []Project {
	projects := []Project{}
	seenPaths := make(map[string]bool)

	for _, project := range cfg.Projects {
		if project.Alias == "" && !seenPaths[project.Path] {
			projects = append(projects, project)
			seenPaths[project.Path] = true
		}
	}

	for _, workspace := range cfg.Workspaces {
		entries, err := os.ReadDir(workspace)
		if err != nil {
			continue
		}

		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}

			projectPath := filepath.Join(workspace, entry.Name())
			if seenPaths[projectPath] {
				continue
			}

			status := "not set"
			if existing, exists := cfg.Projects[entry.Name()]; exists {
				status = existing.Status
			}

			projects = append(projects, Project{
				Name:   entry.Name(),
				Path:   projectPath,
				Status: status,
			})
			seenPaths[projectPath] = true
		}
	}

	return projects
}
