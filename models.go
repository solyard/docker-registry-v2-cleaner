package main

type RegistryCatalog struct {
	Repositories []string `json:"repositories"`
}

type TagsList struct {
	Name string   `json:"name"`
	Tags []string `json:"tags"`
}

type AllTags struct {
	Tags []TagsList
}
