package releases

import (
	"encoding/json"
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/commands/services"
	"github.com/catalyzeio/cli/lib/httpclient"
	"github.com/catalyzeio/cli/models"
)

func CmdUpdate(svcName, releaseName, notes, newReleaseName string, ir IReleases, is services.IServices) error {
	service, err := is.RetrieveByLabel(svcName)
	if err != nil {
		return err
	}
	if service == nil {
		return fmt.Errorf("Could not find a service with the label \"%s\". You can list services with the \"catalyze services\" command.", svcName)
	}
	err = ir.Update(releaseName, service.ID, notes, newReleaseName)
	if err != nil {
		return err
	}
	logrus.Printf("Release '%s' successfully updated", releaseName)
	return nil
}

func (r *SReleases) Update(releaseName, svcID, notes, newReleaseName string) error {
	rls := models.Release{
		ID:    newReleaseName,
		Notes: notes,
	}
	b, err := json.Marshal(rls)
	if err != nil {
		return err
	}
	headers := httpclient.GetHeaders(r.Settings.SessionToken, r.Settings.Version, r.Settings.Pod, r.Settings.UsersID)
	resp, statusCode, err := httpclient.Put(b, fmt.Sprintf("%s%s/environments/%s/services/%s/releases/%s", r.Settings.PaasHost, r.Settings.PaasHostVersion, r.Settings.EnvironmentID, svcID, releaseName), headers)
	if err != nil {
		return err
	}
	return httpclient.ConvertResp(resp, statusCode, nil)
}
