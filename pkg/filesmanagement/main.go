package filesmanagement

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

func ApplyFilter(files []string, match string) []string {
	newfiles := []string{}

	for _, file := range files {
		if matchRegex(match, file) {
			newfiles = append(newfiles, file)
		}
	}

	return newfiles
}

func ListFiles(folder string) []string {
	files := []string{}
	err := filepath.Walk(folder,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			files = append(files, path)
			return nil
		})
	if err != nil {
		log.Println(err)
	}

	return files
}

func FindInFile(file string, match string) []string {
	f, err := os.Open(file)
	if err != nil {
		return []string{}
	}
	defer f.Close()

	matchedlines := []string{}

	scanner := bufio.NewScanner(f)

	// https://golang.org/pkg/bufio/#Scanner.Scan
	for scanner.Scan() {
		if matchRegex(match, scanner.Text()) {
			matchedlines = append(matchedlines, scanner.Text())
		}
	}

	return matchedlines
}

func ExtractValueRegex(sentence string, match string, output string) string {
	return regexp.MustCompile(match).ReplaceAllString(sentence, output)
}

func matchRegex(regex string, challenge string) bool {
	m, _ := regexp.Match(regex, []byte(challenge))
	return m
}
