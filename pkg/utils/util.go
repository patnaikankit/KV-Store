package utils

import (
	"fmt"
	"os"
	"time"
)

func AddLog(op, status, key string) {
	newLog := fmt.Sprintf("%s - %s operation - Key: %s - Status: %s\n", time.Now().Format(time.RFC3339),
		op,
		key,
		status,
	)

	f, err := os.OpenFile("logs.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening log file:", err)
		return
	}
	defer f.Close()

	_, err = f.WriteString(newLog)
	if err != nil {
		fmt.Println("Error writing log:", err)
	}
}
