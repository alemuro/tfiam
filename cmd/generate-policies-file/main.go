package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/alemuro/tfiam/pkg/filesmanagement"
	git "github.com/go-git/go-git/v5"
)

type Perms struct {
	Version  int                 `json:"version"`
	Data     map[string][]string `json:"data"`
	Resource map[string][]string `json:"resource"`
}

func cloneTerraformRepository() {
	git.PlainClone(".terraformrepo", false, &git.CloneOptions{
		URL:      "https://github.com/hashicorp/terraform-provider-aws.git",
		Depth:    1,
		Progress: os.Stdout,
	})
}

func isData(challenge string) bool {
	m, _ := regexp.Match(`^data_.*`, []byte(challenge))
	return m
}

func isResource(challenge string) bool {
	m, _ := regexp.Match(`^resource_.*`, []byte(challenge))
	return m
}

func isTest(challenge string) bool {
	m, _ := regexp.Match(`.*_test.go$`, []byte(challenge))
	return m
}

func getResourceIdentifier(challenge string) string {
	challenge = strings.Replace(challenge, "resource_", "", 1)
	challenge = strings.Replace(challenge, "data_source_", "", 1)
	challenge = strings.Replace(challenge, ".go", "", 1)
	return challenge
}

func findCalls2AWS(t string, resource string) []string {
	filename := fmt.Sprintf(".terraformrepo/aws/%s_%s.go", t, resource)

	perms := []string{}

	f, err := os.Open(filename)
	if err != nil {
		return []string{}
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	lines := filesmanagement.FindInFile(filename, ".*\"github.com/aws/aws-sdk-go/service/([a-zA-Z0-9]*)\"")
	service := ""
	for _, line := range lines {
		service = filesmanagement.ExtractValueRegex(line, ".*\"github.com/aws/aws-sdk-go/service/([a-zA-Z0-9]*)\"", "$1")
	}

	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "conn.") {
			perms = append(perms, fmt.Sprintf("%s:%s", service, regexp.MustCompile(`.*conn.([a-zA-Z0-9]*).*`).ReplaceAllString(scanner.Text(), "$1")))
		}
	}

	return perms
}

func loadFiles() *Perms {
	// files = filesmanagement.ApplyFilter(filesmanagement.ListFiles("."), ".*\\/(resource_|data_).*\\.go$"))
	files, err := ioutil.ReadDir(".terraformrepo/aws/")
	if err != nil {
		log.Fatal(err)
	}

	tfperms := new(Perms)
	tfperms.Version = 1
	tfperms.Data = make(map[string][]string)
	tfperms.Resource = make(map[string][]string)

	for _, f := range files {
		if isData(f.Name()) && !isTest(f.Name()) {
			tfperms.Data[getResourceIdentifier(f.Name())] = findCalls2AWS("data_source", getResourceIdentifier(f.Name()))
		}

		if isResource(f.Name()) && !isTest(f.Name()) {
			tfperms.Resource[getResourceIdentifier(f.Name())] = findCalls2AWS("resource", getResourceIdentifier(f.Name()))
		}
	}

	return tfperms
}

func main() {
	cloneTerraformRepository()

	tfperms := loadFiles()

	output, _ := json.Marshal(tfperms)

	ioutil.WriteFile("./permissions.json", output, 0644)
}
