package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	log "github.com/gookit/slog"
	"github.com/ryanuber/go-glob"
)

func httpDeleteBySHA(url, username, password, name, sha string) int {
	client := &http.Client{}

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/v2/%s/manifests/%v", url, name, strings.TrimSpace(sha)), nil)
	if err != nil {
		log.Fatal(err.Error())
	}
	req.SetBasicAuth(username, password)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err.Error())
	}

	return resp.StatusCode
}

func httpGetHeader(url, username, password, tag, imageName string) string {
	client := &http.Client{}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/v2/%s/manifests/%v", url, imageName, tag), nil)
	if err != nil {
		log.Fatal(err.Error())
	}

	req.SetBasicAuth(username, password)
	req.Header.Set("Accept", "application/vnd.docker.distribution.manifest.v2+json")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err.Error())
	}
	return resp.Header.Get("Docker-Content-Digest")
}

func httpGetBody(url, username, password string) []byte {
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err.Error())
	}

	req.SetBasicAuth(username, password)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err.Error())
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err.Error())
	}

	return body
}

func getImagesCatalog(url, username, password string) *RegistryCatalog {
	registryCatalogBody := httpGetBody(fmt.Sprintf("%s/v2/_catalog", url), username, password)
	var catalog *RegistryCatalog
	json.Unmarshal(registryCatalogBody, &catalog)

	return catalog
}

func getAllTags(url, username, password, imageTagGlob *string) {
	catalog := getImagesCatalog(*url, *username, *password)
	alltags := AllTags{}
	var sortedTagsArray []string
	sortedTags := AllTags{}

	for _, v := range catalog.Repositories {
		tagsListBody := httpGetBody(fmt.Sprintf("%s/v2/%s/tags/list", *url, v), *username, *password)
		var tagslist *TagsList
		json.Unmarshal(tagsListBody, &tagslist)
		alltags.Tags = append(alltags.Tags, *tagslist)
	}

	for _, taglist := range alltags.Tags {
		for _, tag := range taglist.Tags {
			if glob.Glob(*imageTagGlob, tag) {
				sortedTagsArray = append(sortedTagsArray, tag)
			}
		}
		sortedTags.Tags = append(sortedTags.Tags, TagsList{Name: taglist.Name, Tags: sortedTagsArray})
		sortedTagsArray = nil
	}
	removeImageBySHA(sortedTags, url, username, password)
}

func removeImageBySHA(sortedTags AllTags, url, username, password *string) {
	counter := 0

	for _, taglist := range sortedTags.Tags {
		for _, tag := range taglist.Tags {
			sha := httpGetHeader(*url, *username, *password, tag, taglist.Name)
			statusCode := httpDeleteBySHA(*url, *username, *password, taglist.Name, sha)

			if statusCode != 202 {
				log.Error("Error while send delete request for SHA: %s Status code: %s", sha, statusCode)
				log.Warnf("Address: %s/v2/%s/manifests/%s", *url, taglist.Name, sha)
			} else {
				counter += 1
				log.Infof("Removed tag with SHA: %v, for image %v. Removed images count: %v", sha, taglist.Name, counter)
			}
		}
	}
}
