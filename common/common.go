package common

import (
	"crypto/rand"
	"fmt"
	"os"
)

//GetEnv get environment variable passing a default value
func GetEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}

func HandleError(err error, funcName string) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s ----- in func: %s", err.Error(), funcName)
	}
}

// GenerateUUID generate uuid
func GenerateUUID() (uuid string) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	uuid = fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])

	return uuid
}
