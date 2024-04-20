package main

import (
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type DbtProfile struct {
	Target  string `yaml:"target"`
	Outputs map[string]struct {
		ConnType                   string                 `yaml:"type"`
		Account                    string                 `yaml:"account"`
		User                       string                 `yaml:"user"`
		Role                       string                 `yaml:"role"`
		Authenticator              string                 `yaml:"authenticator"`
		Database                   string                 `yaml:"database"`
		Schema                     string                 `yaml:"schema"`
		Project                    string                 `yaml:"project"`
		Dataset                    string                 `yaml:"dataset"`
		Path                       string                 `yaml:"path"`
		Threads                    int                    `yaml:"threads"`
		Password                   string                 `yaml:"password"`
		Port                       int                    `yaml:"port"`
		Warehouse                  string                 `yaml:"warehouse"`
		Method                     string                 `yaml:"method"`
		Host                       string                 `yaml:"host"`
		PrivateKey                 string                 `yaml:"private_key"`
		PrivateKeyPath             string                 `yaml:"private_key_path"`
		PrivateKeyPassphrase       string                 `yaml:"private_key_passphrase"`
		ClientSessionKeepAlive     bool                   `yaml:"client_session_keep_alive"`
		QueryTag                   string                 `yaml:"query_tag"`
		ConnectRetries             int                    `yaml:"connect_retries"`
		ConnectTimeout             int                    `yaml:"connect_timeout"`
		RetryOnDatabaseErrors      bool                   `yaml:"retry_on_database_errors"`
		RetryAll                   bool                   `yaml:"retry_all"`
		ReuseConnections           bool                   `yaml:"reuse_connections"`
		Extensions                 []string               `yaml:"extensions"`
		RefreshToken               string                 `yaml:"refresh_token"`
		ClientID                   string                 `yaml:"client_id"`
		ClientSecret               string                 `yaml:"client_secret"`
		TokenURI                   string                 `yaml:"token_uri"`
		Token                      string                 `yaml:"token"`
		Priority                   string                 `yaml:"priority"`
		Keyfile                    string                 `yaml:"keyfile"`
		JobExecutionTimeoutSeconds int                    `yaml:"job_execution_timeout_seconds"`
		JobCreationTimeoutSeconds  int                    `yaml:"job_creation_timeout_seconds"`
		JobRetryDeadlineSeconds    int                    `yaml:"job_retry_deadline_seconds"`
		Location                   string                 `yaml:"location"`
		MaximumBytesBilled         int                    `yaml:"maximum_bytes_billed"`
		Scopes                     []string               `yaml:"scopes"`
		ImpersonateServiceAccount  string                 `yaml:"impersonate_service_account"`
		ExecutionProject           string                 `yaml:"execution_project"`
		GcsBucket                  string                 `yaml:"gcs_bucket"`
		DataprocRegion             string                 `yaml:"dataproc_region"`
		DataprocClusterName        string                 `yaml:"dataproc_cluster_name"`
		DataprocBatch              map[string]interface{} `yaml:"dataproc_batch"`
		KeyfileJson                map[string]struct {
			Type                    string `yaml:"type"`
			ProjectId               string `yaml:"project_id"`
			PrivateKeyId            string `yaml:"private_key_id"`
			PrivateKey              string `yaml:"private_key"`
			ClientEmail             string `yaml:"client_email"`
			ClientId                string `yaml:"client_id"`
			AuthURI                 string `yaml:"auth_uri"`
			TokenURI                string `yaml:"token_uri"`
			AuthProviderX509CertUrl string `yaml:"auth_provider_x509_cert_url"`
			ClientX509CertUrl       string `yaml:"client_x509_cert_url"`
		} `yaml:"keyfile_json"`
		Settings map[string]struct {
			S3Region          string `yaml:"s3_region"`
			S3AccessKeyID     string `yaml:"s3_access_key_id"`
			S3SecretAccessKey string `yaml:"s3_secret_access_key"`
		} `yaml:"settings"`
	} `yaml:"outputs"`
}

type DbtProfiles map[string]DbtProfile

func FetchDbtProfiles() (DbtProfiles, error) {
	paths := []string{
		filepath.Join(".", "profiles.yml"),
		filepath.Join(os.Getenv("HOME"), ".dbt", "profiles.yml"),
	}
	ps := DbtProfiles{}
	for _, path := range paths {
		pf := DbtProfiles{}
		yf, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		if err = yaml.Unmarshal(yf, pf); err != nil {
			log.Fatalf("Could not read dbt profile, \nlikely unsupported fields or formatting issues\n please open an issue: %v\n", err)
		}
		for k, v := range pf {
			ps[k] = v
		}
	}
	return ps, nil
}
