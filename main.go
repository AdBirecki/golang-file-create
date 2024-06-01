package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

type PathError struct {
	path string
	err  error
}

func (e *PathError) Error() string {
	return fmt.Sprintf("error accessing path %q: %v", e.path, e.err)
}

const letterRunes = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ!"

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env files", err)
	}
	sourcePath := os.Getenv("SOURCE_PATH")
	generateFiles(sourcePath)

	targetPath := os.Getenv("TARGET_PATH")
	sortFiles(sourcePath, targetPath)
}

func readFileContents(filePath string) (rune, []byte, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return 0, nil, err
	}
	reader := bytes.NewReader(data)
	r, _, err := reader.ReadRune()
	if err != nil {
		return 0, nil, err
	}
	return r, data, nil
}

func createFile(targetPath string, fileName string, data []byte) error {
	err := ensureDirExists(targetPath)
	if err != nil {
		fmt.Fprintln(os.Stdout, []any{"an error occured while ensuring directory created, %q", err}...)
		return err
	}
	fullFilePath := filepath.Join(targetPath, fileName)
	err = os.WriteFile(fullFilePath, data, 0644)
	if err != nil {
		fmt.Fprintln(os.Stdout, []any{"An error creating file %q", fullFilePath}...)
		return err
	}
	fmt.Fprintln(os.Stdout, []any{"File created and written successfully %q", fullFilePath}...)
	return nil
}

func sortFiles(sourcePath, targetPath string) {
	err := ensureDirExists(sourcePath)
	if err != nil {
		fmt.Fprintln(os.Stdout, []any{"an error occured while ensuring directory exists, %q", err}...)
		return
	}

	files, err := os.ReadDir(sourcePath)
	if err != nil {
		fmt.Println("Error readign directory:", err)
	}
	for _, file := range files {
		if !file.IsDir() {
			fullPath := filepath.Join(sourcePath, file.Name())
			dirName, data, err := readFileContents(fullPath)

			if err != nil {
				fmt.Printf("error: %q", err)
				continue
			}

			targetPath := filepath.Join(targetPath, string(dirName))
			err = createFile(targetPath, file.Name(), data)
			if err != nil {
				fmt.Printf("error: %q", err)
				continue
			}
			fmt.Printf("%q, %q ", data, file)

		}
	}

	err = ensureDirExists(targetPath)
	if err != nil {
		fmt.Fprintln(os.Stdout, []any{"an error occured while ensuring directory exists, %q", err}...)
		return
	}

}
func generateFiles(basePath string) {
	err := ensureDirExists(basePath)
	if err != nil {
		fmt.Fprintln(os.Stdout, []any{"an error occured while ensuring directory created, %q", err}...)
		return
	}
	runes := []rune(letterRunes)
	iterator := 0
	for {
		fileContents := randStringRunes(10, runes)
		fmt.Println(fileContents)
		filePath := filepath.Join(basePath, fileContents)
		err = os.WriteFile(filePath, []byte(fileContents), 0755)
		if err != nil {
			fmt.Fprintln(os.Stdout, []any{"An error while writing file %q, content %v, message: %q", filePath, fileContents, err}...)
		}
		if iterator > 10 {
			break
		}
		iterator++
	}
}

func ensureDirExists(path string) error {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(path, 0755)
			if err != nil {
				return &PathError{
					path: path,
					err:  err,
				}
			}
			fmt.Println("Directory created:", path)
		} else {
			return &PathError{path: path, err: err}
		}
	}
	return nil
}

func randStringRunes(n int, runes []rune) string {
	b := make([]rune, n)

	for i := range b {
		ri := randRange(0, len(letterRunes))
		b[i] = runes[ri]
	}
	return string(b)
}

func randRange(min, max int) int {
	return rand.Intn(max-min) + min
}
