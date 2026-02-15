package utils

import (
	"fmt"
	"time"
)

func ToDuration(str string) time.Duration {
	result, err := time.ParseDuration(str)
	if err != nil {
		fmt.Println(err)
	}

	return result
}
