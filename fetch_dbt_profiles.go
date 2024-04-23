package main

import (
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type DbtProfile struct {
	Outputs map[string]struct {
		Settings map[string]struct {
			S3Region          string `yaml:"s3_region"`
			S3AccessKeyID     string `yaml:"s3_access_key_id"`
			S3SecretAccessKey string `yaml:"s3_secret_access_key"`
		} `yaml:"settings"`
		DataprocBatch map[string]interface{} `yaml:"dataproc_batch"`
		KeyfileJson   map[string]struct {
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
		PrivateKeyPath             string   `yaml:"private_key_path"`
		DbName                     string   `yaml:"dbname"`
		Database                   string   `yaml:"database"`
		Account                    string   `yaml:"account"`
		Schema                     string   `yaml:"schema"`
		QueryTag                   string   `yaml:"query_tag"`
		Dataset                    string   `yaml:"dataset"`
		Path                       string   `yaml:"path"`
		Role                       string   `yaml:"role"`
		Password                   string   `yaml:"password"`
		User                       string   `yaml:"user"`
		Warehouse                  string   `yaml:"warehouse"`
		Method                     string   `yaml:"method"`
		Host                       string   `yaml:"host"`
		PrivateKey                 string   `yaml:"private_key"`
		Location                   string   `yaml:"location"`
		ConnType                   string   `yaml:"type"`
		Authenticator              string   `yaml:"authenticator"`
		Project                    string   `yaml:"project"`
		SslMode                    string   `yaml:"sslmode"`
		DataprocClusterName        string   `yaml:"dataproc_cluster_name"`
		DataprocRegion             string   `yaml:"dataproc_region"`
		GcsBucket                  string   `yaml:"gcs_bucket"`
		ExecutionProject           string   `yaml:"execution_project"`
		PrivateKeyPassphrase       string   `yaml:"private_key_passphrase"`
		RefreshToken               string   `yaml:"refresh_token"`
		ClientID                   string   `yaml:"client_id"`
		ClientSecret               string   `yaml:"client_secret"`
		TokenURI                   string   `yaml:"token_uri"`
		Token                      string   `yaml:"token"`
		Priority                   string   `yaml:"priority"`
		Keyfile                    string   `yaml:"keyfile"`
		ImpersonateServiceAccount  string   `yaml:"impersonate_service_account"`
		HttpPath                   string   `yaml:"http_path"`
		Extensions                 []string `yaml:"extensions"`
		Scopes                     []string `yaml:"scopes"`
		JobCreationTimeoutSeconds  int      `yaml:"job_creation_timeout_seconds"`
		MaximumBytesBilled         int      `yaml:"maximum_bytes_billed"`
		JobRetryDeadlineSeconds    int      `yaml:"job_retry_deadline_seconds"`
		JobExecutionTimeoutSeconds int      `yaml:"job_execution_timeout_seconds"`
		ConnectTimeout             int      `yaml:"connect_timeout"`
		ConnectRetries             int      `yaml:"connect_retries"`
		Port                       int      `yaml:"port"`
		Threads                    int      `yaml:"threads"`
		ReuseConnections           bool     `yaml:"reuse_connections"`
		RetryAll                   bool     `yaml:"retry_all"`
		RetryOnDatabaseErrors      bool     `yaml:"retry_on_database_errors"`
		ClientSessionKeepAlive     bool     `yaml:"client_session_keep_alive"`
	} `yaml:"outputs"`
	Target string `yaml:"target"`
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
