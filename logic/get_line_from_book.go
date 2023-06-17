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
	"strings"

	"github.com/spf13/viper"
)

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

func getRandomLine(file *os.File) (string, error) {
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

func isFileNameTryingToTraverse(fileName string) bool {
	tryingToAttack := filepath.Clean(fileName) != fileName ||
		fileName[0] == filepath.Separator ||
		strings.Contains(fileName, "..")
	return tryingToAttack
}

func getConfiguredBooksDir() string {
	return viper.GetString("BOOKS_DIRECTORY")
}

func ReadRandomLineFromFile(fileName string, lineNumber int) (string, error) {
	// protect against directory traversal attacks,
	// see https://owasp.org/www-community/attacks/Path_Traversal
	if fileName != "" && isFileNameTryingToTraverse(fileName) {
		return "", errors.New("invalid file name - possible directory traversal attempt")
	}

	finalFileName, filePathErr := getBookFilePath(fileName)
	if filePathErr != nil {
		return "", filePathErr
	}

	file, fileOpenErr := os.Open(finalFileName)
	if fileOpenErr != nil {
		if os.IsNotExist(fileOpenErr) {
			return "", errors.Join(fileOpenErr, &BookNotFoundError{BookTitle: fileName})
		}
		return "", fileOpenErr
	}
	defer file.Close()

	if lineNumber <= 0 {
		pickedLine, err := getRandomLine(file)
		if err != nil {
			return "", err
		}
		return pickedLine, nil
	}

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

func getBookFilePath(fileName string) (string, error) {
	booksDir := getConfiguredBooksDir()

	var finalFileName string
	if fileName == "" {
		pickedFile, err := getRandomFile(booksDir)
		if err != nil {
			return "", err
		}
		finalFileName = filepath.Join(booksDir, pickedFile)
	} else {
		finalFileName = filepath.Join(booksDir, fileName)
	}

	finalFileName, absPathErr := filepath.Abs(finalFileName)
	if absPathErr != nil {
		return "", absPathErr
	}
	return finalFileName, nil
}
