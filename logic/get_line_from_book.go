// logic contains the business logic of the application.
//
// Since this application's actual focus is to demonstrate various API
// consumption issues, the business logic is kind of random.
// We're just returning lines from books.
package logic

import (
	"bufio"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
)

const bookAssetsDirName = "books"

func getRandomFile(dirPath string) (string, error) {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return "", err
	}

	if len(files) == 0 {
		return "", fmt.Errorf("no files found in directory")
	}

	//nolint:gosec // we're not using this for anything security related
	randIndex := rand.Intn(len(files))
	return files[randIndex].Name(), nil
}

func getRandomLine(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lines := []string{}
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if len(lines) == 0 {
		return "", fmt.Errorf("no lines found in file")
	}

	//nolint:gosec // we're not using this for anything security related
	randIndex := rand.Intn(len(lines))
	return lines[randIndex], nil
}

type LineNotFoundError struct {
	LineNumber int
}

func (e *LineNotFoundError) Error() string {
	return fmt.Sprintf("line number %d not found in file", e.LineNumber)
}

type BookNotFoundError struct {
	BookTitle string
}

func (e *BookNotFoundError) Error() string {
	return fmt.Sprintf("book %s not found in library", e.BookTitle)
}

func ReadRandomLineFromFile(fileName string, lineNumber int) (string, error) {
	var finalFileName string
	if fileName == "" {
		pickedFile, err := getRandomFile(bookAssetsDirName)
		if err != nil {
			return "", err
		}
		finalFileName = filepath.Join(bookAssetsDirName, pickedFile)
	} else {
		finalFileName = filepath.Join(bookAssetsDirName, fileName)
	}

	if lineNumber <= 0 {
		pickedLine, err := getRandomLine(fileName)
		if err != nil {
			return "", err
		}
		return pickedLine, nil
	}

	file, err := os.Open(finalFileName)
	if err != nil {
		if os.IsNotExist(err) {
			return "", errors.Join(err, &BookNotFoundError{BookTitle: fileName})
		}
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	currentLine := 0
	for scanner.Scan() {
		currentLine++
		if currentLine == lineNumber {
			return scanner.Text(), nil
		}
	}

	return "", fmt.Errorf("line number %d not found in file", lineNumber)
}
