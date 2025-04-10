package utils

import (
    "fmt"
    "log"
    "os"
    "path/filepath"
    "time"
)

var Logger *log.Logger

func InitLogger() error {
    // Create logs directory if it doesn't exist
    if err := os.MkdirAll("logs", 0755); err != nil {
        return fmt.Errorf("failed to create logs directory: %v", err)
    }

    // Create log file with timestamp
    logFile := filepath.Join("logs", fmt.Sprintf("api_%s.log", time.Now().Format("2006-01-02")))
    f, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
    if err != nil {
        return fmt.Errorf("failed to open log file: %v", err)
    }

    Logger = log.New(f, "", log.LstdFlags|log.Lshortfile)
    return nil
}