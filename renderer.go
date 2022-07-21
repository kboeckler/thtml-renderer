package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	executable, _ := os.Executable()
	programName := filepath.Base(executable)
	programDirPath := filepath.Dir(executable)
	programDir, err := os.Open(programDirPath)
	if err != nil {
		fmt.Printf("%s\n", err)
		fmt.Printf("Usage: %s <thtml-file>\n", programName)
		os.Exit(-1)
	}
	programDirInfo, err := programDir.Stat()
	programDir.Close()
	if err != nil {
		fmt.Printf("%s\n", err)
		fmt.Printf("Usage: %s <thtml-file>\n", programName)
		os.Exit(-1)
	}
	if !programDirInfo.IsDir() {
		panic("dir of executable is not a dir")
	}
	if len(os.Args) < 2 {
		fmt.Printf("Not enough arguments\n")
		fmt.Printf("Usage: %s <thtml-file>\n", programName)
		os.Exit(-1)
	}
	if len(os.Args) > 2 {
		fmt.Printf("Too many arguments\n")
		fmt.Printf("Usage: %s <thtml-file>\n", programName)
		os.Exit(-1)
	}
	thtmlFile := os.Args[1]
	outputHtmlFile, err := getFilenameAsHtmlFile(thtmlFile)
	if err != nil {
		fmt.Printf("%s\n", err)
		fmt.Printf("Usage: %s <thtml-file>\n", programName)
		os.Exit(-1)
	}
	var sf *os.File
	sf, err = os.Open(thtmlFile)
	useProgramDir := false
	if err != nil {
		var errRetry error
		sf, errRetry = os.Open(filepath.Join(programDirPath, filepath.Base(thtmlFile)))
		if errRetry == nil {
			useProgramDir = true
			err = nil
		}
	}
	if err != nil {
		fmt.Printf("%s\n", err)
		fmt.Printf("Usage: %s <thtml-file>\n", programName)
		os.Exit(-1)
	}
	sf.Close()
	targetFilePath := outputHtmlFile
	if useProgramDir {
		targetFilePath = filepath.Join(programDirPath, filepath.Base(outputHtmlFile))
	}
	targetFile, err := os.Create(targetFilePath)
	if err != nil {
		fmt.Printf("%s\n", err)
		fmt.Printf("Usage: %s <thtml-file>\n", programName)
		os.Exit(-1)
	}
	sourceFilePath := thtmlFile
	if useProgramDir {
		sourceFilePath = filepath.Join(programDirPath, filepath.Base(thtmlFile))
	}
	sourceContent, err := ioutil.ReadFile(sourceFilePath)
	targetContent := convertHtml(string(sourceContent))
	_, err = targetFile.WriteString(targetContent)
	if err != nil {
		fmt.Printf("%s\n", err)
		fmt.Printf("Usage: %s <thtml-file>\n", programName)
		os.Exit(-1)
	}
	targetFile.Close()
	fmt.Printf("%s: %d bytes written\n", targetFilePath, len(targetContent))
}

func convertHtml(sourceContent string) string {
	return strings.ReplaceAll(sourceContent, "h1>", "h2>")
}

func getFilenameAsHtmlFile(thtmlFile string) (string, error) {
	lastIndexForEnding := strings.LastIndex(strings.ToLower(thtmlFile), ".thtml")
	if lastIndexForEnding == -1 {
		return "", errors.New("not a .html filename")
	}
	return thtmlFile[:lastIndexForEnding] + ".html", nil
}
