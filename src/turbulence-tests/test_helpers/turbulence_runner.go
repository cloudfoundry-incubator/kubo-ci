package test_helpers

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	. "github.com/onsi/ginkgo"
)

type TurbulenceRunner struct {
	username    string
	password    string
	apiEndpoint string
	httpClient  *http.Client
}

type TurbulenceIncident struct {
	ID string `json:"ID"`
	// other fields TBD
}

type TurbulenceIncidents []TurbulenceIncident

func NewTurbulenceRunner() *TurbulenceRunner {

	runner := &TurbulenceRunner{}

	runner.apiEndpoint = os.Getenv("TURBULENCE_API_ENDPOINT")
	if runner.apiEndpoint == "" {
		Fail("TURBULENCE_API_ENDPOINT is not set")
	}

	runner.username = os.Getenv("TURBULENCE_USERNAME")
	if runner.username == "" {
		Fail("TURBULENCE_USERNAME is not set")
	}

	runner.password = os.Getenv("TURBULENCE_PASSWORD")
	if runner.password == "" {
		Fail("TURBULENCE_PASSWORD is not set")
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	runner.httpClient = &http.Client{Transport: tr}

	return runner
}

func PathFromRoot(relativePath string) string {
	_, filename, _, _ := runtime.Caller(0)
	currentDir := filepath.Dir(filename)
	return filepath.Join(currentDir, "..", "..", "..", relativePath)
}

func (runner TurbulenceRunner) ListIncidents() (TurbulenceIncidents, error) {
	getUrl := fmt.Sprintf("https://%s:%s@%s/incidents", runner.username, runner.password, runner.apiEndpoint)

	response, listErr := runner.httpClient.Get(getUrl)
	if listErr != nil {
		return nil, fmt.Errorf("Error listing incidents; [%v]", listErr)
	}

	bytes, bodyErr := ioutil.ReadAll(response.Body)
	if bodyErr != nil {
		return nil, fmt.Errorf("Error parsing incidents list response; [%v]", bodyErr)
	}

	if response.StatusCode >= 300 {
		return nil, fmt.Errorf("Error in listing incidents; [%v]", string(bytes))
	}

	incidents := TurbulenceIncidents{}
	json.Unmarshal(bytes, &incidents)
	return incidents, nil
}

func (runner TurbulenceRunner) GetIncidentById(incidentId string) (TurbulenceIncident, error) {
	getUrl := fmt.Sprintf("https://%s:%s@%s/incidents/%s", runner.username, runner.password, runner.apiEndpoint, incidentId)

	response, getErr := runner.httpClient.Get(getUrl)
	if getErr != nil {
		return TurbulenceIncident{}, fmt.Errorf("Error getting incident; [%v]", getErr)
	}

	bytes, bodyErr := ioutil.ReadAll(response.Body)
	if bodyErr != nil {
		return TurbulenceIncident{}, fmt.Errorf("Error parsing incident response; [%v]", bodyErr)
	}

	if response.StatusCode >= 300 {
		return TurbulenceIncident{}, fmt.Errorf("Error in getting incident; [%v]", string(bytes))
	}

	incident := TurbulenceIncident{}
	json.Unmarshal(bytes, &incident)
	return incident, nil
}

func (runner TurbulenceRunner) ApplyIncident(incidentFile string) (TurbulenceIncident, error) {
	postUrl := fmt.Sprintf("https://%s:%s@%s/incidents", runner.username, runner.password, runner.apiEndpoint)
	fileBytes, fileReadErr := ioutil.ReadFile(incidentFile)
	if fileReadErr != nil {
		return TurbulenceIncident{}, fmt.Errorf("Can't read file [%s]; [%v]", incidentFile, fileReadErr)
	}

	response, postErr := runner.httpClient.Post(postUrl, "application/json", bytes.NewReader(fileBytes))
	if postErr != nil {
		return TurbulenceIncident{}, fmt.Errorf("Error submitting incident; [%v]", postErr)
	}

	bytes, bodyErr := ioutil.ReadAll(response.Body)
	if bodyErr != nil {
		return TurbulenceIncident{}, fmt.Errorf("Error parsing incident submit response; [%v]", bodyErr)
	}

	if response.StatusCode >= 300 {
		return TurbulenceIncident{}, fmt.Errorf("Error in submitting incident; [%v]", string(bytes))
	}

	newIncident := TurbulenceIncident{}
	json.Unmarshal(bytes, &newIncident)
	return newIncident, nil
}
