package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/alemuro/tfiam/pkg/filesmanagement"
)

type Perms struct {
	Version  int                 `json:"version"`
	Data     map[string][]string `json:"data"`
	Resource map[string][]string `json:"resource"`
}

func getPolicies() *Perms {
	resp, err := http.Get("https://amurtra.s3.eu-west-1.amazonaws.com/permissions.json")
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	p := new(Perms)
	json.Unmarshal(body, &p)

	return p
}

func findTerraformVerb(files []string, verb string) []string {
	res := []string{}
	for _, file := range files {
		datas := filesmanagement.FindInFile(file, fmt.Sprintf("^%s", verb))
		for _, data := range datas {
			res = append(res, filesmanagement.ExtractValueRegex(data, fmt.Sprintf("^%s \"([a-zA-Z0-9-_]*)\" \".*$", verb), "$1"))
		}
	}
	return unique(res)
}

func main() {
	policies := getPolicies()

	tffiles := filesmanagement.ApplyFilter(filesmanagement.ListFiles("."), ".*\\.tf$")

	permissions := []string{}

	resources := findTerraformVerb(tffiles, "resource")
	for _, resource := range resources {
		for _, perm := range policies.Resource[resource] {
			permissions = append(permissions, perm)
		}
	}

	datas := findTerraformVerb(tffiles, "data")
	for _, data := range datas {
		for _, perm := range policies.Data[data] {
			permissions = append(permissions, perm)
		}
	}

	policyjson, _ := json.Marshal(generateAWSPolicy(permissions))

	fmt.Println(string(policyjson))
}
