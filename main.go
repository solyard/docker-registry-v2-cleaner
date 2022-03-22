package main

import "flag"

func main() {
	readParamsAndRun()
}

func readParamsAndRun() {
	registryURL := flag.String("url", "https://registry.docker.io", "Docker Registry V2 address (HTTPS / HTTP) ")
	registryUsername := flag.String("username", "", "Please set Usermame if basic auth enabled for Docker Registry V2")
	registryPassword := flag.String("password", "", "Please set Password if basic auth enabled for Docker Registry V2")
	imageTagGlob := flag.String("glob", "", "Glob for image filtering by tags (If empty will be match nothing)")
	flag.Parse()

	getAllTags(registryURL, registryUsername, registryPassword, imageTagGlob)
}
