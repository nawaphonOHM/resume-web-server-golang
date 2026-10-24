// Package config provides application configuration loading and validation.
package config

import (
	"bufio"
	"errors"
	"fmt"
	"os"
)

var defaultDotEnvFiles = []string{".env", "configs/.env"}

func loadDotEnvFiles(filenames ...string) error {
	if len(filenames) == 0 {
		return loadDefaultDotEnv()
	}
	return loadExplicitDotEnvFiles(filenames)
}

func loadDefaultDotEnv() error {
	for _, name := range defaultDotEnvFiles {
		err := loadExplicitEnvFile(name)
		if !errors.Is(err, os.ErrNotExist) {
			return err
		}
	}
	return nil
}

func loadExplicitDotEnvFiles(filenames []string) error {
	for _, f := range filenames {
		if err := loadExplicitEnvFile(f); err != nil {
			return err
		}
	}
	return nil
}

func loadExplicitEnvFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("open env file: %w", err)
	}
	defer closeFile(file)
	return scanEnvLines(file)
}

func closeFile(file *os.File) {
	if err := file.Close(); err != nil {
		return
	}
}

func scanEnvLines(file *os.File) error {
	scanner := bufio.NewScanner(file)
	if err := parseScannerLines(scanner); err != nil {
		return err
	}
	return checkScannerError(scanner)
}

func parseScannerLines(scanner *bufio.Scanner) error {
	for scanner.Scan() {
		if err := parseEnvLine(scanner.Text()); err != nil {
			return err
		}
	}
	return nil
}

func checkScannerError(scanner *bufio.Scanner) error {
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scan env file: %w", err)
	}
	return nil
}
