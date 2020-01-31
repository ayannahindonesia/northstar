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
				Level:    "event_seed",
				Messages: "test log messages one red",
			},
			models.Log{
				Level:    "log_seed",
				Messages: "test log messages two red",
			},
			models.Log{
				Level:    "warning_seed",
				Messages: "test log messages three red",
			},
			models.Log{
				Level:    "error_seed",
				Messages: "test log messages four blue",
			},
			models.Log{
				Level:    "debug_seed",
				Messages: "test log messages five blue",
			},
		}
		for _, log := range logs {
			log.Create()
		}
	}
}

// Unseed removes all seed datas
func Unseed() (err error) {
	seededTables := []string{
		"clients",
		"logs",
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
		}

		if err != nil {
			break
		}
	}

	return err
}
