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
	checkProgramDirValid(programDirPath)
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
	sourceFilePath, targetFilePath, err := getFilePaths(err, thtmlFile, programDirPath, programName, outputHtmlFile)
	if err != nil {
		fmt.Printf("%s\n", err)
		fmt.Printf("Usage: %s <thtml-file>\n", programName)
		os.Exit(-1)
	}
	sourceContent, err := ioutil.ReadFile(sourceFilePath)
	targetContent := convertHtml(string(sourceContent))
	targetFile, err := os.Create(targetFilePath)
	if err != nil {
		fmt.Printf("%s\n", err)
		fmt.Printf("Usage: %s <thtml-file>\n", programName)
		os.Exit(-1)
	}
	_, err = targetFile.WriteString(targetContent)
	if err != nil {
		fmt.Printf("%s\n", err)
		fmt.Printf("Usage: %s <thtml-file>\n", programName)
		os.Exit(-1)
	}
	targetFile.Close()
	fmt.Printf("%s: %d bytes written\n", targetFilePath, len(targetContent))
}

func checkProgramDirValid(programDirPath string) {
	programDir, err := os.Open(programDirPath)
	if err != nil {
		panic("dir of executable cannot be opened to read")
	}
	programDirInfo, err := programDir.Stat()
	defer programDir.Close()
	if err != nil {
		panic("dir of executable cannot be read")
	}
	if !programDirInfo.IsDir() {
		panic("dir of executable is not a dir")
	}
}

func getFilenameAsHtmlFile(thtmlFile string) (string, error) {
	lastIndexForEnding := strings.LastIndex(strings.ToLower(thtmlFile), ".thtml")
	if lastIndexForEnding == -1 {
		return "", errors.New("not a .html filename")
	}
	return thtmlFile[:lastIndexForEnding] + ".html", nil
}

func getFilePaths(err error, thtmlFile string, programDirPath string, programName string, outputHtmlFile string) (string, string, error) {
	useProgramDir := false
	var sf *os.File
	sf, err = os.Open(thtmlFile)
	if err != nil {
		var errRetry error
		sf, errRetry = os.Open(filepath.Join(programDirPath, filepath.Base(thtmlFile)))
		if errRetry == nil {
			useProgramDir = true
			err = nil
		}
	}
	if err != nil {
		return "", "", err
	}
	sf.Close()
	sourceFilePath := thtmlFile
	if useProgramDir {
		sourceFilePath = filepath.Join(programDirPath, filepath.Base(thtmlFile))
	}
	targetFilePath := outputHtmlFile
	if useProgramDir {
		targetFilePath = filepath.Join(programDirPath, filepath.Base(outputHtmlFile))
	}
	return sourceFilePath, targetFilePath, nil
}

func convertHtml(sourceContent string) string {
	return strings.ReplaceAll(sourceContent, "h1>", "h2>")
}
