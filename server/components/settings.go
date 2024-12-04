package components

import (
	"errors"
	"lommix/wichtelbot/server/store"
	"os"
	"strconv"
	"time"
)

type HttpSettings struct {
	Port int
}

type HttpsSettings struct {
	Port        int
	SslCertPath string
	SslKeyPath  string
}

type Settings struct {
	Http           HttpSettings
	Https          *HttpsSettings
	ExpirySettings store.ExpirySettings
}

// Helper function to wrap an error with context.
func wrapError(envVar string, err error) error {
	return errors.New("error loading " + envVar + ": " + err.Error())
}

func LoadSettingsFromEnv() (*Settings, error) {
	httpPort, err := getEnvAsInt("HTTP_PORT", 8000)
	if err != nil {
		return nil, wrapError("HTTP_PORT", err)
	}

	httpSettings := HttpSettings{
		Port: httpPort,
	}

	var httpsSettings *HttpsSettings
	useHttps, err := getEnvAsBool("USE_HTTPS", false)
	if err != nil {
		return nil, wrapError("HTTPS_PORT", err)
	}

	if useHttps {
		httpsPort, err := getEnvAsInt("HTTPS_PORT", 8080)
		if err != nil {
			return nil, wrapError("HTTPS_PORT", err)
		}

		sslCertPath := getEnv("SSL_CERT_PATH", "")
		if sslCertPath == "" {
			return nil, errors.New("SSL_CERT_PATH must be set if HTTPS_PORT is set")
		}

		sslKeyPath := getEnv("SSL_KEY_PATH", "")
		if sslKeyPath == "" {
			return nil, errors.New("SSL_KEY_PATH must be set if HTTPS_PORT is set")
		}

		httpsSettings = &HttpsSettings{
			Port:        httpsPort,
			SslCertPath: sslCertPath,
			SslKeyPath:  sslKeyPath,
		}
	}

	createdTimeout_h, err := getEnvAsInt("CREATED_TIMEOUT_DURATION_H", 24)
	if err != nil {
		return nil, wrapError("CREATED_TIMEOUT_DURATION_H", err)
	}

	joiningTimeout_h, err := getEnvAsInt("JOINING_TIMEOUT_DURATION_H", 72)
	if err != nil {
		return nil, wrapError("JOINING_TIMEOUT_DURATION_H", err)
	}

	playedTimeout_h, err := getEnvAsInt("PLAYED_TIMEOUT_DURATION_H", 72)
	if err != nil {
		return nil, wrapError("PLAYED_TIMEOUT_DURATION_H", err)
	}
	expirySettings := store.ExpirySettings{
		CreatedTimeoutDuration: time.Duration(createdTimeout_h) * time.Hour,
		JoiningTimeoutDuration: time.Duration(joiningTimeout_h) * time.Hour,
		PlayedTimeoutDuration:  time.Duration(playedTimeout_h) * time.Hour,
	}

	settings := &Settings{
		Http:           httpSettings,
		Https:          httpsSettings,
		ExpirySettings: expirySettings,
	}

	return settings, nil
}

// Helper function to get a string environment variable with a fallback default value.
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// Helper function to get an integer environment variable with a fallback default value.
func getEnvAsInt(key string, defaultValue int) (int, error) {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue, nil
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return 0, err
	}
	return value, nil
}

// Helper function to get a boolean environment variable with a fallback default value.
func getEnvAsBool(key string, defaultValue bool) (bool, error) {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue, nil
	}
	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		return false, err
	}
	return value, nil
}
