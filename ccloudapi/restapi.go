package ccloudapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

const (
	ccloudAPI       = "https://confluent.cloud/api"
	loginURI        = "/sessions"
	envMetadataURI  = "/env_metadata"
	environmentsURI = "/accounts"
	clustersURI     = "/clusters"
	apiKeysURI      = "/api_keys"
)

var httpClient = &http.Client{Timeout: 30 * time.Second}

// Login performs a login on behalf of the user and creates
// an authentication token that will be used subsequentially.
func Login(username, password string) (*Session, error) {
	loginRequest := LoginRequest{
		Username: username,
		Password: password,
	}
	loginBytes, _ := json.Marshal(loginRequest)
	payload := bytes.NewBuffer(loginBytes)
	_, resp, err := httpRequest(nil, "POST", loginURI, "application/json", payload)
	if err != nil {
		return nil, err
	}
	loginResponse := new(LoginResponse)
	err = json.Unmarshal(resp, &loginResponse)
	if err != nil {
		return nil, err
	}
	return &Session{
		AuthToken: loginResponse.AuthToken,
		User:      loginResponse.User,
	}, nil
}

// GetEnvironmentMetadata Retrieve the last update from the metadata
func GetEnvironmentMetadata(session *Session) ([]*CloudProvider, error) {
	uri := envMetadataURI
	statusCode, resp, err := httpRequest(session, "GET", uri, "application/json", nil)
	if err != nil || statusCode == 404 {
		return nil, err
	}
	envMetadataResponse := new(EnvironmentMetadataResponse)
	err = json.Unmarshal(resp, &envMetadataResponse)
	if err != nil {
		return nil, err
	}
	return envMetadataResponse.CloudProviders, nil
}

// CreateEnvironment creates a new environment
func CreateEnvironment(environment *Environment, session *Session) (*Environment, error) {
	createEnvironmentRequest := CreateEnvironmentRequest{
		Environment: environment,
	}
	createEnvironmentBytes, err := json.Marshal(createEnvironmentRequest)
	if err != nil {
		return nil, err
	}
	payload := bytes.NewBuffer(createEnvironmentBytes)
	_, resp, err := httpRequest(session, "POST", environmentsURI, "application/json", payload)
	if err != nil {
		return nil, err
	}
	createEnvironmentResponse := new(CreateEnvironmentResponse)
	err = json.Unmarshal(resp, &createEnvironmentResponse)
	if err != nil {
		return nil, err
	}
	return createEnvironmentResponse.Environment, nil
}

// UpdateEnvironment updates an existing environment
func UpdateEnvironment(environment *Environment, session *Session) (bool, error) {
	updateEnvironmentRequest := UpdateEnvironmentRequest{
		Environment: environment,
	}
	updateEnvironmentBytes, err := json.Marshal(updateEnvironmentRequest)
	if err != nil {
		return false, err
	}
	payload := bytes.NewBuffer(updateEnvironmentBytes)
	uri := environmentsURI + "/" + environment.ID
	statusCode, _, err := httpRequest(session, "PUT", uri, "application/json", payload)
	if err != nil || statusCode == 404 {
		return false, err
	}
	return true, nil
}

// DeleteEnvironment deletes an existing environment
func DeleteEnvironment(environment *Environment, session *Session) (bool, error) {
	deleteEnvironmentRequest := DeleteEnvironmentRequest{
		Environment: environment,
	}
	deleteEnvironmentBytes, _ := json.Marshal(deleteEnvironmentRequest)
	payload := bytes.NewBuffer(deleteEnvironmentBytes)
	uri := environmentsURI + "/" + environment.ID
	statusCode, _, err := httpRequest(session, "DELETE", uri, "application/json", payload)
	if err != nil || statusCode == 404 {
		return false, err
	}
	return true, nil
}

// ReadEnvironment reads an existing environment
func ReadEnvironment(id string, session *Session) (*Environment, error) {
	uri := environmentsURI + "/" + id
	statusCode, resp, err := httpRequest(session, "GET", uri, "application/json", nil)
	if err != nil || statusCode == 404 {
		return nil, err
	}
	readEnvironmentResponse := new(ReadEnvironmentResponse)
	err = json.Unmarshal(resp, &readEnvironmentResponse)
	if err != nil {
		return nil, err
	}
	return readEnvironmentResponse.Environment, nil
}

// ListEnvironments list all environments from the account
func ListEnvironments(session *Session) ([]*Environment, error) {
	_, resp, err := httpRequest(session, "GET", environmentsURI, "application/json", nil)
	if err != nil {
		return nil, err
	}
	listEnvironmentResponse := new(ListEnvironmentResponse)
	err = json.Unmarshal(resp, &listEnvironmentResponse)
	if err != nil {
		return nil, err
	}
	return listEnvironmentResponse.Environments, nil
}

// CreateCluster creates a new cluster
func CreateCluster(cluster *Cluster, session *Session) (*Cluster, error) {
	createClusterRequest := CreateClusterRequest{
		Cluster: cluster,
	}
	createClusterBytes, err := json.Marshal(createClusterRequest)
	if err != nil {
		return nil, err
	}
	payload := bytes.NewBuffer(createClusterBytes)
	_, resp, err := httpRequest(session, "POST", clustersURI, "application/json", payload)
	if err != nil {
		return nil, err
	}
	createClusterResponse := new(CreateClusterResponse)
	err = json.Unmarshal(resp, &createClusterResponse)
	if err != nil {
		return nil, err
	}
	return createClusterResponse.Cluster, nil
}

// ReadCluster reads an existing cluster
func ReadCluster(id, environmentID string, session *Session) (*Cluster, error) {
	uri := clustersURI + "/" + id + "?account_id=" + environmentID
	statusCode, resp, err := httpRequest(session, "GET", uri, "application/json", nil)
	if err != nil || statusCode == 404 {
		return nil, err
	}
	readClusterResponse := new(ReadClusterResponse)
	err = json.Unmarshal(resp, &readClusterResponse)
	if err != nil {
		return nil, err
	}
	return readClusterResponse.Cluster, nil
}

// UpdateCluster updates an existing cluster
func UpdateCluster(cluster *Cluster, session *Session) (bool, error) {
	updateClusterRequest := UpdateClusterRequest{
		Cluster: cluster,
	}
	updateClusterBytes, err := json.Marshal(updateClusterRequest)
	if err != nil {
		return false, err
	}
	payload := bytes.NewBuffer(updateClusterBytes)
	uri := clustersURI + "/" + cluster.ID
	statusCode, _, err := httpRequest(session, "PUT", uri, "application/json", payload)
	if err != nil || statusCode == 404 {
		return false, err
	}
	return true, nil
}

// DeleteCluster deletes an existing cluster
func DeleteCluster(cluster *Cluster, session *Session) (bool, error) {
	deleteClusterRequest := DeleteClusterRequest{
		Cluster: cluster,
	}
	deleteClusterBytes, _ := json.Marshal(deleteClusterRequest)
	payload := bytes.NewBuffer(deleteClusterBytes)
	uri := clustersURI + "/" + cluster.ID
	statusCode, _, err := httpRequest(session, "DELETE", uri, "application/json", payload)
	if err != nil || statusCode == 404 {
		return false, err
	}
	return true, nil
}

// ListClusters lists all clusters in the environment
func ListClusters(environmentID string, session *Session) ([]*Cluster, error) {
	uri := clustersURI + "?account_id=" + environmentID
	_, resp, err := httpRequest(session, "GET", uri, "application/json", nil)
	if err != nil {
		return nil, err
	}
	listClusterResponse := new(ListClusterResponse)
	err = json.Unmarshal(resp, &listClusterResponse)
	if err != nil {
		return nil, err
	}
	return listClusterResponse.Clusters, nil
}

// CreateAPIKey creates a new api key
func CreateAPIKey(environmentID, clusterID string, session *Session) (*APIKey, error) {
	cluster := Cluster{ID: clusterID}
	createAPIKeyRequest := CreateAPIKeyRequest{
		APIKeyRequest: APIKeyRequestBody{
			EnvironmentID: environmentID,
			Clusters:      []*Cluster{&cluster},
		},
	}
	createAPIKeyBytes, err := json.Marshal(createAPIKeyRequest)
	if err != nil {
		return nil, err
	}
	payload := bytes.NewBuffer(createAPIKeyBytes)
	_, resp, err := httpRequest(session, "POST", apiKeysURI, "application/json", payload)
	if err != nil {
		return nil, err
	}
	createAPIKeyResponse := new(CreateAPIKeyResponse)
	err = json.Unmarshal(resp, &createAPIKeyResponse)
	if err != nil {
		return nil, err
	}
	return createAPIKeyResponse.APIKey, nil
}

// ReadAPIKey reads an existing api key
func ReadAPIKey(environmentID, clusterID, key string, session *Session) (*APIKey, error) {
	uri := apiKeysURI + "?account_id=" + environmentID + "&cluster_id=" + clusterID + "&key=" + key
	statusCode, resp, err := httpRequest(session, "GET", uri, "application/json", nil)
	if err != nil || statusCode == 404 {
		return nil, err
	}
	readAPIKeyResponse := new(ReadAPIKeyResponse)
	err = json.Unmarshal(resp, &readAPIKeyResponse)
	if err != nil {
		return nil, err
	}
	apiKeys := readAPIKeyResponse.APIKeys
	if len(apiKeys) > 0 {
		return apiKeys[0], nil
	}
	return nil, nil
}

// DeleteAPIKey deletes an existing api key
func DeleteAPIKey(environmentID, clusterID, id string, session *Session) (bool, error) {
	cluster := Cluster{ID: clusterID}
	apiKeyID, _ := strconv.Atoi(id)
	deleteAPIKeyRequest := DeleteAPIKeyRequest{
		APIKeyRequest: APIKeyRequestBody{
			ID:            apiKeyID,
			EnvironmentID: environmentID,
			Clusters:      []*Cluster{&cluster},
		},
	}
	deleteAPIKeyBytes, _ := json.Marshal(deleteAPIKeyRequest)
	payload := bytes.NewBuffer(deleteAPIKeyBytes)
	uri := apiKeysURI + "/" + id
	statusCode, _, err := httpRequest(session, "DELETE", uri, "application/json", payload)
	if err != nil || statusCode == 404 {
		return false, err
	}
	return true, nil
}

// generic function to handle HTTP requests throughout the code
func httpRequest(session *Session, method, uri, contentType string, payload io.Reader) (int, []byte, error) {
	url := fmt.Sprintf("%s%s", ccloudAPI, uri)
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return 0, nil, err
	}
	req.Header.Set("Content-Type", contentType)
	if session != nil {
		req.Header.Set("Cookie", fmt.Sprintf("%s%s", "auth_token=", session.AuthToken))
	}
	resp, err := httpClient.Do(req)
	statusCode := resp.StatusCode
	if err != nil {
		return statusCode, nil, err
	}

	if resp != nil {
		defer resp.Body.Close()
	}
	if statusCode < 200 || statusCode > 299 {
		return statusCode, nil, createError(resp)
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	return statusCode, bytes, err
}

// Creates a easier-to-read error
func createError(resp *http.Response) error {
	decoder := json.NewDecoder(resp.Body)
	errorResp := new(ErrorResponse)
	err := decoder.Decode(&errorResp)
	if err == nil {
		return fmt.Errorf("%s: %s", resp.Status, errorResp.Error.Message)
	}
	return fmt.Errorf("%s", resp.Status)
}
