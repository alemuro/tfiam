package main

import (
	"encoding/json"
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

func getFile() *Perms {
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

func unique(stringSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range stringSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func getResources(files []string) []string {
	res := []string{}
	for _, file := range files {
		resources := filesmanagement.FindInFile(file, "^resource")
		for _, resource := range resources {
			res = append(res, filesmanagement.ExtractValueRegex(resource, "^resource \"([a-zA-Z0-9-_]*)\" \".*$", "$1"))
		}
	}
	return unique(res)
}

func getDatas(files []string) []string {
	res := []string{}
	for _, file := range files {
		datas := filesmanagement.FindInFile(file, "^data")
		for _, data := range datas {
			res = append(res, filesmanagement.ExtractValueRegex(data, "^data \"([a-zA-Z0-9-_]*)\" \".*$", "$1"))
		}
	}
	return unique(res)
}

func main() {
	log.Println("Hello world")
	policies := getFile()

	tffiles := filesmanagement.ApplyFilter(filesmanagement.ListFiles("."), ".*\\.tf$")

	permissions := []string{}

	resources := getResources(tffiles)
	for _, resource := range resources {
		for _, perm := range policies.Resource[resource] {
			permissions = append(permissions, perm)
		}
	}

	datas := getDatas(tffiles)
	for _, data := range datas {
		for _, perm := range policies.Data[data] {
			permissions = append(permissions, perm)
		}
	}

	log.Println(permissions)
}
