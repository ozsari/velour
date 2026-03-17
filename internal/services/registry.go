package services

import "github.com/ozsari/velour/internal/models"

// Registry holds all app definitions. Category files append via init().
var Registry []models.ServiceDefinition

// testingApps limits which apps are shown during testing. Set to nil to show all.
var testingApps map[string]bool = nil

func GetRegistry() []models.ServiceDefinition {
	result := make([]models.ServiceDefinition, 0, len(Registry))
	for _, s := range Registry {
		if testingApps != nil && !testingApps[s.ID] {
			continue
		}
		if len(s.InstallTypes) == 0 {
			s.InstallTypes = []models.InstallType{models.InstallDocker}
		}
		result = append(result, s)
	}
	return result
}

func FindByID(id string) *models.ServiceDefinition {
	for _, s := range Registry {
		if s.ID == id {
			return &s
		}
	}
	return nil
}
