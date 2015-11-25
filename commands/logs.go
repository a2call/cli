package commands

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/catalyzeio/catalyze/config"
	"github.com/catalyzeio/catalyze/helpers"
	"github.com/catalyzeio/catalyze/models"
)

// Logs is a way to stream logs from Kibana to your local terminal. This is
// useful because Kibana is hard to look at because it splits every single
// log statement into a separate block that spans multiple lines so it's
// not very cohesive. This is intended to be similar to the `heroku logs`
// command.
func Logs(queryString string, tail bool, hours int, minutes int, seconds int, settings *models.Settings) {
	if settings.Username == "" || settings.Password == "" {
		// sometimes this will be filled in from env variables
		// if it is, just use that and don't prompt them
		settings.Username = ""
		settings.Password = ""
		fmt.Println("Please enter your logging dashboard credentials")
	}
	// if we remove the session token, the CLI will prompt for the
	// username/password normally. It will also set the username/password
	// on the settings object.
	sessionToken := settings.SessionToken
	settings.SessionToken = ""

	helpers.SignIn(settings)

	env := helpers.RetrieveEnvironment("pod", settings)
	domain := env.Data.DNSName
	if domain == "" {
		domain = fmt.Sprintf("%s.catalyze.io", env.Data.Namespace)
	}
	appLogsIdentifier := "source"
	appLogsValue := "app"
	if strings.HasPrefix(domain, "pod01") {
		appLogsIdentifier = "syslog_program"
		appLogsValue = "supervisord"
	}

	urlString := fmt.Sprintf("https://%s/__es", domain)

	offset := time.Duration(hours)*time.Hour + time.Duration(minutes)*time.Minute + time.Duration(seconds)*time.Second
	timestamp := time.Now().In(time.UTC).Add(-1 * offset)

	from := 0

	var tr = &http.Transport{
		TLSClientConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},
	}
	client := &http.Client{
		Transport: tr,
	}

	settings.SessionToken = sessionToken
	config.SaveSettings(settings)

	fmt.Println("        @timestamp       -        message")
	for {
		queryBytes := generateQuery(queryString, appLogsIdentifier, appLogsValue, timestamp, from)
		reader := bytes.NewReader(queryBytes)

		req, _ := http.NewRequest("GET", fmt.Sprintf("%s/_search", urlString), reader)
		req.SetBasicAuth(settings.Username, settings.Password)

		resp, err := client.Do(req)
		if err != nil {
			fmt.Println(err.Error())
		}
		respBody, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			fmt.Println(fmt.Errorf("%d %s", resp.StatusCode, string(respBody)).Error())
			os.Exit(1)
		}
		var logs models.Logs
		json.Unmarshal(respBody, &logs)
		for _, lh := range *logs.Hits.Hits {
			fmt.Printf("%s - %s\n", lh.Fields["@timestamp"][0], lh.Fields["message"][0])
		}
		if !tail {
			break
		}
		time.Sleep(2 * time.Second)
		from += len(*logs.Hits.Hits)
	}
}

func generateQuery(queryString, appLogsIdentifier, appLogsValue string, timestamp time.Time, from int) []byte {
	query := `{
	"fields": ["@timestamp", "message", "` + appLogsIdentifier + `"],
	"query": {
		"wildcard": {
			"message": "` + queryString + `"
		}
	},
	"filter": {
		"bool": {
			"must": [
				{"term": {"` + appLogsIdentifier + `": "` + appLogsValue + `"}},
				{"range": {"@timestamp": {"gt": "` + fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:%02dZ", timestamp.Year(), timestamp.Month(), timestamp.Day(), timestamp.Hour(), timestamp.Minute(), timestamp.Second()) + `"}}}
			]
		}
	},
	"sort": {
		"@timestamp": {
			"order": "asc"
		},
		"message": {
			"order": "asc"
		}
	},
	"from": ` + fmt.Sprintf("%d", from) + `,
	"size": 50
	}`
	var buf bytes.Buffer
	json.Compact(&buf, []byte(query))
	return buf.Bytes()
}
