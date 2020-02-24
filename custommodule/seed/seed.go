package seed

import (
	"northstar/application"
	"northstar/models"
)

// Seed insert dummy datas
func Seed() {
	if application.App.ENV == "development" {
		clients := []models.Client{
			models.Client{
				Name:   "Borrower Service",
				Key:    "borrowerkey",
				Secret: "borrowersecret",
			},
			models.Client{
				Name:   "Lender Service",
				Key:    "lenderkey",
				Secret: "lendersecret",
			},
			models.Client{
				Name:   "Geomapping Service",
				Key:    "geomappingkey",
				Secret: "geomappingsecret",
			},
			models.Client{
				Name:   "Messaging Service",
				Key:    "messagingkey",
				Secret: "messagingsecret",
			},
			models.Client{
				Name:   "Android",
				Key:    "androkey",
				Secret: "androsecret",
			},
			models.Client{
				Name:   "React Frontend",
				Key:    "reactkey",
				Secret: "reactsecret",
			},
		}
		for _, client := range clients {
			client.Create()
		}

		logs := []models.Log{
			models.Log{
				Level:    "info_seed",
				Client:   "Client A",
				Tag:      "c01",
				Note:     "note a",
				UID:      "2",
				Username: "user dummy",
				Messages: "test log messages one red",
			},
			models.Log{
				Level:    "log_seed",
				Client:   "Client B",
				Tag:      "c02",
				Note:     "note b",
				Messages: "test log messages two red",
			},
			models.Log{
				Level:    "warning_seed",
				Client:   "Client C",
				UID:      "4",
				Username: "iamnumber4",
				Messages: "test log messages three red",
			},
			models.Log{
				Level:    "error_seed",
				Client:   "Client C",
				UID:      "4",
				Username: "iamnumber4",
				Messages: "test log messages four blue",
			},
			models.Log{
				Level:    "debug_seed",
				Client:   "Client C",
				UID:      "4",
				Username: "iamnumber4",
				Messages: "test log messages five blue",
			},
		}
		for _, log := range logs {
			log.Create()
		}

		ats := []models.Audittrail{
			models.Audittrail{
				Client:   "Client A _seed_",
				UserID:   "1",
				Username: "iamnumber1",
				Roles:    "[1,2,3]",
				Entity:   "entity a",
				EntityID: "1",
				Action:   "create",
				Original: "",
				New:      `{"messages":"audit trail example 1"}`,
			},
			models.Audittrail{
				Client:   "Client B _seed_",
				UserID:   "2",
				Username: "iamnumber2",
				Roles:    "[1,3]",
				Entity:   "entity h",
				EntityID: "1",
				Action:   "update",
				Original: `{"messages":"original"}`,
				New:      `{"messages":"new"}`,
			},
		}
		for _, at := range ats {
			at.Create()
		}
	}
}

// Unseed removes all seed datas
func Unseed() (err error) {
	seededTables := []string{
		"clients",
		"logs",
		"audittrail",
	}

	for _, s := range seededTables {
		switch s {
		case "logs":
			err = application.App.DB.
				Where("level LIKE ?", "%_seed").
				Unscoped().Delete(&models.Log{}).Error
			break
		case "clients":
			err = application.App.DB.
				Or("key = ?", "borrowerkey").
				Or("key = ?", "lenderkey").
				Or("key = ?", "geomappingkey").
				Or("key = ?", "messagingkey").
				Or("key = ?", "androkey").
				Or("key = ?", "reactkey").
				Unscoped().Delete(&models.Client{}).Error
			break
		case "audittrail":
			err = application.App.DB.
				Where("client LIKE ?", "%_seed_").
				Unscoped().Delete(&models.Audittrail{}).Error
			break
		}

		if err != nil {
			break
		}
	}

	return err
}
