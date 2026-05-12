package utils

import (
	"bufio"
	"os"
)

func ReadFile(filePath string) ([]string, error) {

	lines := make([]string, 0)
	file, err := os.Open(filePath)

	if err != nil {

		return nil, err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		linha := scanner.Text()
		lines = append(lines, linha)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}
