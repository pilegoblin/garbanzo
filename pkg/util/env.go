package util

import "os"

// Utility function to grab an environment variable or panic. Should only be used on mission critical items.
func GetEnvVarOrPanic(envVarName string) string {
	v, ok := os.LookupEnv(envVarName)
	if !ok {
		panic("unable to find env var: " + envVarName)
	}
	return v

}

// Utility function to grab an environment variable, or return the default if it does not exist
func GetEnvVarOrDefault(envVarName string, defaultValue string) string {
	v, ok := os.LookupEnv(envVarName)
	if !ok {
		return defaultValue
	}
	return v
}
