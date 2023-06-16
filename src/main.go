package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/aws/aws-lambda-go/lambda"
)

type Target struct {
	Env   string `json:"env"`
	DB    string `json:"db"`
	Table string `json:"table"`
}

type RDSEnv struct {
	User     string
	Password string
	Host     string
	DB       string
}

type RedshiftEnv struct {
	Host     string
	User     string
	Password string
	Schema   string
	DB       string
	Table    string
	S3Bucket string
}

func main() {
	if isDebug() {
		t := Target{
			Env:   "dev",
			DB:    "my_db",
			Table: "my_table",
		}
		if err := handleRequest(t); err != nil {
			fmt.Printf("%v\n", err)
		}
	} else {
		lambda.Start(handleRequest)
	}
}

func isDebug() bool {
	return os.Getenv("DEBUG") == "1"
}

func handleRequest(t Target) error {
	setEmbulkHome()
	if err := setEnv(t); err != nil {
		return fmt.Errorf("failed to setEnv: %w", err)
	}

	embulkPath := filepath.Join(embulkHome(), "bin/embulk")
	cmd := exec.Command("java", "-jar", embulkPath, "run", configPath())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to cmd.Run: %w", err)
	}
	return nil
}

func setEmbulkHome() error {
	if embulkHome := os.Getenv("EMBULK_HOME"); embulkHome != "" {
		return nil
	}

	if isDebug() {
		home := os.Getenv("HOME")
		if home == "" {
			return errors.New("both EMBULK_HOME and HOME are emtpy")
		}
		os.Setenv("EMBULK_HOME", filepath.Join(home, ".embulk"))
	} else {
		os.Setenv("EMBULK_HOME", "/embulk")
	}
	return nil
}

func embulkHome() string {
	return os.Getenv("EMBULK_HOME")
}

func setEnv(t Target) error {
	rdsEnv, err := rdsEnv(t)
	if err != nil {
		return fmt.Errorf("failed to rdsEnv: %w", err)
	}
	os.Setenv("RDS_HOST", rdsEnv.Host)
	os.Setenv("RDS_USER", rdsEnv.User)
	os.Setenv("RDS_PASSWORD", rdsEnv.Password)
	os.Setenv("RDS_DB", rdsEnv.DB)
	os.Setenv("RDS_TABLE", t.Table)

	redshiftEnv, err := redshiftEnv(t)
	if err != nil {
		return fmt.Errorf("failed to redshiftEnv: %w", err)
	}
	os.Setenv("REDSHIFT_HOST", redshiftEnv.Host)
	os.Setenv("REDSHIFT_USER", redshiftEnv.User)
	os.Setenv("REDSHIFT_PASSWORD", redshiftEnv.Password)
	os.Setenv("REDSHIFT_SCHEMA", redshiftEnv.Schema)
	os.Setenv("REDSHIFT_DB", redshiftEnv.DB)
	os.Setenv("REDSHIFT_TABLE", redshiftEnv.Table)
	os.Setenv("REDSHIFT_S3_BUCKET", redshiftEnv.S3Bucket)
	return nil
}

func rdsEnv(t Target) (*RDSEnv, error) {
	// ... Write code to get RDS environments from Parameter Store or somewhere else ...

	env := &RDSEnv{
		User:     "your_user",
		Password: "your_password",
		Host:     "your_host",
		DB:       "your_db",
	}
	return env, nil
}

func redshiftEnv(t Target) (*RedshiftEnv, error) {
	// ... Write code to get Redshift environments from Parameter Store or somewhere else ...

	env := &RedshiftEnv{
		Host:     "your_host",
		User:     "your_user",
		Password: "your_password",
		Schema:   "your_schema",
		DB:       "your_db",
		Table:    "your_table",
		S3Bucket: "your_s3_bucket",
	}
	return env, nil
}

func configPath() string {
	if isDebug() {
		p, err := filepath.Abs("../config/rds_to_redshift_dev_local.yml.liquid")
		if err != nil {
			return ""
		}
		return p
	} else {
		return filepath.Join(embulkHome(), "config/rds_to_redshift.yml.liquid")
	}
}
