package test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/daticahealth/cli/lib/httpclient"
	"github.com/daticahealth/cli/models"
)

const (
	BinaryName = "cli"

	Alias     = "ctest"
	EnvID     = "env1"
	EnvName   = "cli-integration-tests"
	Namespace = "pod011234"
	OrgID     = "org1"
	Pod       = "pod1"
	SvcID     = "svc1"
	SvcLabel  = "code-1"

	AliasAlt     = "ctest2"
	EnvIDAlt     = "env2"
	EnvNameAlt   = "cli-integration-tests2"
	NamespaceAlt = "pod0125678"
	OrgIDAlt     = "org2"
	PodAlt       = "pod2"
	SvcIDAlt     = "svc2"
	SvcLabelAlt  = "code-2"
)

func Setup() (*http.ServeMux, *httptest.Server, *url.URL) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	baseURL, _ := url.Parse(server.URL)
	return mux, server, baseURL
}

func Teardown(server *httptest.Server) {
	server.Close()
}

func AssertMethod(t *testing.T, r *http.Request, want string) {
	if want != r.Method {
		t.Errorf("Request method = %v, want %v", r.Method, want)
	}
}

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
