package gointelowl

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
)

// AnalyzerConfig represents how an analyzer is configured in IntelOwl.
//
// IntelOwl docs: https://intelowl.readthedocs.io/en/latest/Usage.html#analyzers-customization
type AnalyzerConfig struct {
	BaseConfigurationType
	Type                  string   `json:"type"`
	ExternalService       bool     `json:"external_service"`
	LeaksInfo             bool     `json:"leaks_info"`
	DockerBased           bool     `json:"docker_based"`
	RunHash               bool     `json:"run_hash"`
	RunHashType           string   `json:"run_hash_type"`
	SupportedFiletypes    []string `json:"supported_filetypes"`
	NotSupportedFiletypes []string `json:"not_supported_filetypes"`
	ObservableSupported   []string `json:"observable_supported"`
}

// AnalyzerService handles communication with analyzer related methods of the IntelOwl API.
//
// IntelOwl REST API docs: https://intelowl.readthedocs.io/en/latest/Redoc.html#tag/analyzer
type AnalyzerService struct {
	client *IntelOwlClient
}

// GetConfigs lists down every analyzer configuration in your IntelOwl instance.
//
//	Endpoint: GET /api/get_analyzer_configs
//
// IntelOwl REST API docs: https://intelowl.readthedocs.io/en/latest/Redoc.html#tag/get_analyzer_configs
func (analyzerService *AnalyzerService) GetConfigs(ctx context.Context) (*[]AnalyzerConfig, error) {
	requestUrl := fmt.Sprintf(ANALYZER_CONFIG_URL, analyzerService.client.options.Url)
	contentType := "application/json"
	method := "GET"
	request, err := analyzerService.client.buildRequest(ctx, method, contentType, nil, requestUrl)
	if err != nil {
		return nil, err
	}

	successResp, err := analyzerService.client.newRequest(ctx, request)
	if err != nil {
		return nil, err
	}
	analyzerConfigurationResponse := map[string]AnalyzerConfig{}
	if unmarshalError := json.Unmarshal(successResp.Data, &analyzerConfigurationResponse); unmarshalError != nil {
		return nil, unmarshalError
	}

	analyzerNames := make([]string, 0)
	// *getting all the analyzer key names!
	for analyzerName := range analyzerConfigurationResponse {
		analyzerNames = append(analyzerNames, analyzerName)
	}
	// * sorting them alphabetically
	sort.Strings(analyzerNames)
	analyzerConfigurationList := []AnalyzerConfig{}
	for _, analyzerName := range analyzerNames {
		analyzerConfig := analyzerConfigurationResponse[analyzerName]
		analyzerConfigurationList = append(analyzerConfigurationList, analyzerConfig)
	}
	return &analyzerConfigurationList, nil
}

// HealthCheck checks if the specified analyzer is up and running
//
//	Endpoint: GET /api/analyzer/{NameOfAnalyzer}/healthcheck
//
// IntelOwl REST API docs: https://intelowl.readthedocs.io/en/latest/Redoc.html#tag/analyzer/operation/analyzer_healthcheck_retrieve
func (analyzerService *AnalyzerService) HealthCheck(ctx context.Context, analyzerName string) (bool, error) {
	requestUrl := fmt.Sprintf(ANALYZER_HEALTHCHECK_URL, analyzerService.client.options.Url, analyzerName)
	contentType := "application/json"
	method := "GET"
	request, err := analyzerService.client.buildRequest(ctx, method, contentType, nil, requestUrl)
	if err != nil {
		return false, err
	}
	status := StatusResponse{}
	successResp, err := analyzerService.client.newRequest(ctx, request)
	if err != nil {
		return false, err
	}
	if unmarshalError := json.Unmarshal(successResp.Data, &status); unmarshalError != nil {
		return false, unmarshalError
	}
	return status.Status, nil
}
