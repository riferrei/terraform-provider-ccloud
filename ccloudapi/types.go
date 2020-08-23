package ccloudapi

// User type
type User struct {
	ID             int    `json:"id"`
	Email          string `json:"email"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	OrganizationID int    `json:"organization_id"`
}

// Session type
type Session struct {
	AuthToken string
	User      User
}

// LoginRequest type
type LoginRequest struct {
	Username string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse type
type LoginResponse struct {
	AuthToken string `json:"token"`
	User      User   `json:"user"`
}

// Environment type
type Environment struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	OrganizationID int    `json:"organization_id"`
}

// CreateEnvironmentRequest type
type CreateEnvironmentRequest struct {
	Environment *Environment `json:"account"`
}

// CreateEnvironmentResponse type
type CreateEnvironmentResponse struct {
	Environment *Environment `json:"account"`
}

// UpdateEnvironmentRequest type
type UpdateEnvironmentRequest struct {
	Environment *Environment `json:"account"`
}

// DeleteEnvironmentRequest type
type DeleteEnvironmentRequest struct {
	Environment *Environment `json:"account"`
}

// ReadEnvironmentResponse type
type ReadEnvironmentResponse struct {
	Environment *Environment `json:"account"`
}

// ListEnvironmentResponse type
type ListEnvironmentResponse struct {
	Environments []*Environment `json:"accounts"`
}

// Cluster type
type Cluster struct {
	ID             string `json:"id,omitempty"`
	EnvironmentID  string `json:"account_id,omitempty"`
	Name           string `json:"name,omitempty"`
	CloudProvider  string `json:"service_provider,omitempty"`
	CloudRegion    string `json:"region,omitempty"`
	NetworkIngress int    `json:"network_ingress,omitempty"`
	NetworkEgress  int    `json:"network_egress,omitempty"`
	Storage        int    `json:"storage,omitempty"`
	Durability     string `json:"durability,omitempty"`

	OrganizationID  int    `json:"organization_id,omitempty"`
	ClusterEndpoint string `json:"endpoint,omitempty"`
	APIEndpoint     string `json:"api_endpoint,omitempty"`
}

// CreateClusterRequest type
type CreateClusterRequest struct {
	Cluster *Cluster `json:"config"`
}

// CreateClusterResponse type
type CreateClusterResponse struct {
	Cluster *Cluster `json:"cluster"`
}

// ReadClusterResponse type
type ReadClusterResponse struct {
	Cluster *Cluster `json:"cluster"`
}

// UpdateClusterRequest type
type UpdateClusterRequest struct {
	Cluster *Cluster `json:"cluster"`
}

// DeleteClusterRequest type
type DeleteClusterRequest struct {
	Cluster *Cluster `json:"cluster"`
}

// ListClusterResponse type
type ListClusterResponse struct {
	Clusters []*Cluster `json:"clusters"`
}

// APIKey type
type APIKey struct {
	ID     int    `json:"id,omitempty"`
	Key    string `json:"key"`
	Secret string `json:"secret"`
}

// APIKeyRequestBody type
type APIKeyRequestBody struct {
	ID            int        `json:"id,omitempty"`
	EnvironmentID string     `json:"account_id"`
	Clusters      []*Cluster `json:"logical_clusters"`
}

// CreateAPIKeyRequest type
type CreateAPIKeyRequest struct {
	APIKeyRequest APIKeyRequestBody `json:"api_key"`
}

// CreateAPIKeyResponse type
type CreateAPIKeyResponse struct {
	APIKey *APIKey `json:"api_key"`
}

// DeleteAPIKeyRequest type
type DeleteAPIKeyRequest struct {
	APIKeyRequest APIKeyRequestBody `json:"api_key"`
}

// ReadAPIKeyResponse type
type ReadAPIKeyResponse struct {
	APIKeys []*APIKey `json:"api_keys"`
}

// InternalError type
type InternalError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// ErrorResponse type
type ErrorResponse struct {
	Error InternalError `json:"error"`
}

// CloudRegion type
type CloudRegion struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// CloudProvider type
type CloudProvider struct {
	ID      string         `json:"id"`
	Regions []*CloudRegion `json:"regions"`
}

// EnvironmentMetadataResponse type
type EnvironmentMetadataResponse struct {
	CloudProviders []*CloudProvider `json:"clouds"`
}
