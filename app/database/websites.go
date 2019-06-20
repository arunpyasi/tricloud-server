package database

import (
	"fmt"
)

type Website struct {
	Url        string `json:"url"`
	Name       string `json:"name"`
	Active     bool   `json:"active"`
	Timestamp  int64  `json:timestamp`
	Subscriber string `json:subscribers`
}

var WebsiteBucketName = []byte("websites")

func NewWebsite(websiteInfo map[string]string, active bool) (*Website, error) {
	fields := []string{"url", "name", "user"}

	for _, field := range fields {
		value, ok := websiteInfo[field]
		if !ok {
			return nil, fmt.Errorf("field %s not found", value)
		}
	}

	name := websiteInfo["user"] + "_" + websiteInfo["name"]

	return &Website{
		Url:        websiteInfo["url"],
		Name:       name,
		Subscriber: websiteInfo["user"],
		Active:     active,
	}, nil
}

func CreateWebsite(website *Website) error {
	websitebyte, err := Encode(website)
	if err != nil {
		return err
	}

	err = DB.Create([]byte(website.Name), websitebyte, WebsiteBucketName)
	if err != nil {
		return err
	}
	return nil
}

func UpdateWebsite(website *Website) error {
	websitebyte, err := Encode(website)
	if err != nil {
		return err
	}

	err = DB.Update([]byte(website.Name), websitebyte, WebsiteBucketName)
	if err != nil {
		return err
	}
	return nil
}

func GetAllWebsites() ([]*Website, error) {
	websites := []*Website{}

	websitebyte, err := DB.ReadAll(WebsiteBucketName)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	for _, val := range websitebyte {
		website := &Website{}
		err = Decode(val, website)
		websites = append(websites, website)

		if err != nil {
			return nil, err
		}

	}
	return websites, nil
}

func GetWebsite(name string) (*Website, error) {
	website := &Website{}
	websitebyte, err := DB.Read([]byte(name), WebsiteBucketName)

	if err != nil {
		return nil, err
	}

	err = Decode(websitebyte, website)
	if err != nil {
		return nil, err
	}

	return website, nil
}

func DeleteWebsite(name string) error {
	// u, err := GetWebsite(url)
	// if err != nil {
	// 	return err
	// }
	// u.Agents = deleteSliceItem(u.Agents, id)
	// userbyte, err := Encode(u)
	// if err != nil {
	// 	return err
	// }
	// DB.Update([]byte(u.ID), userbyte, UserBucketName)
	return DB.Delete([]byte(name), WebsiteBucketName)
}
