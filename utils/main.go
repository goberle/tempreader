package utils

import (
    "os"
)

// IsTrue returns true if the string value is one of the true strings.
func IsTrue(value string) (ret bool) {
    ret = false
    trueValues := [7]string{
        "true",
        "True",
        "TRUE",
        "yes",
        "Yes",
        "YES",
        "1",
    }

    for _, v := range trueValues {
        if value == v {
            ret = true
            break
        }
    }

    return
}

// GetEnv returns the value of the environment variable or provided fallback
// value if the environment variable is not defined.
func GetEnv(key string, fallback string) string {
    if value, ok := os.LookupEnv(key); ok {
        return value
    }

    return fallback
}

// GetEnvBool is the same like GetEnv but for boolean values.
func GetEnvBool(key string, fallback bool) bool {
    if value, ok := os.LookupEnv(key); ok {
        return IsTrue(value)
    }

    return fallback
}
