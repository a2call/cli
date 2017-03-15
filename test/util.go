package test

import (
	"github.com/daticahealth/cli/lib/httpclient"
	"github.com/daticahealth/cli/models"
)

func GetSettings(baseURL string) *models.Settings {
	return &models.Settings{
		HTTPManager: httpclient.NewTLSHTTPManager(false),
		PaasHost:    baseURL,
		Environments: map[string]models.AssociatedEnvV2{
			Alias: models.AssociatedEnvV2{
				Name:          EnvName,
				EnvironmentID: EnvID,
				Pod:           Pod,
				OrgID:         OrgID,
			},
		},
		EnvironmentID: EnvID,
		Pods: &[]models.Pod{
			models.Pod{
				Name: Pod,
			},
			models.Pod{
				Name: PodAlt,
			},
		},
	}
}
