package seniorities

import "github.com/masa-finance/tee-worker/v2/pkg/util"

// id represents a LinkedIn seniority level identifier
type Id string

// Seniority level constants
const (
	InTraining         Id = "100"
	EntryLevel         Id = "110"
	Senior             Id = "120"
	Strategic          Id = "130"
	EntryLevelManager  Id = "200"
	ExperiencedManager Id = "210"
	Director           Id = "220"
	VicePresident      Id = "300"
	CXO                Id = "310"
	Partner            Id = "320"
)

var All = util.NewSet(
	InTraining,
	EntryLevel,
	Senior,
	Strategic,
	EntryLevelManager,
	ExperiencedManager,
	Director,
	VicePresident,
	CXO,
	Partner,
)

type SenioritiesConfig struct {
	All                util.Set[Id]
	InTraining         Id
	EntryLevel         Id
	Senior             Id
	Strategic          Id
	EntryLevelManager  Id
	ExperiencedManager Id
	Director           Id
	VicePresident      Id
	CXO                Id
	Partner            Id
}

var Seniorities = SenioritiesConfig{
	All:                *All,
	InTraining:         InTraining,
	EntryLevel:         EntryLevel,
	Senior:             Senior,
	Strategic:          Strategic,
	EntryLevelManager:  EntryLevelManager,
	ExperiencedManager: ExperiencedManager,
	Director:           Director,
	VicePresident:      VicePresident,
	CXO:                CXO,
	Partner:            Partner,
}
