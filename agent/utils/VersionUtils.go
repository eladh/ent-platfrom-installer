package utils

import (
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
	"github.com/hashicorp/go-version"
	"common/utils"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func GetProductsVersions() map[string][]*version.Version {
	var versions = make(map[string][]*version.Version)

	versions["artifactory"] = 	getProductVersions("jfrog-artifactory-pro-zip" ,"artifactory-pro")
	versions["distribution"] = 	getProductVersions("distribution-ubuntu" ,"jfrog-distribution")
	versions["jfmc"] = 	getProductVersions("mc-docker-installers" ,"jfrog-mission-control")
	versions["xray"] = 	getProductVersions("xray-ubuntu" ,"xray")

	return versions
}

func GetHelmVersions() map[string][]*version.Version {
	var versions = make(map[string][]*version.Version)

	res, err := http.Get("https://charts.jfrog.io/")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Find the review items
	doc.Find("a").Each(func(index int, s *goquery.Selection) {
		text := s.Nodes[0].Attr[0].Val
		i:= strings.LastIndex(text, "-")
		if i == -1 {
			return
		}

		artifact:= text[0:i]
		versionMeta , _ := version.NewVersion(text[i+1:strings.LastIndex(text, "tgz")-1])
		versions[artifact] = append(versions[artifact], versionMeta)

	})

	return versions
}


func getProductVersions(name string ,repo string) []*version.Version {
	var versions []*version.Version
	var arr []string

	res, err := http.Get("https://api.bintray.com/search/packages?name=" + name + "&repo=" + repo)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	bodyB, _ := ioutil.ReadAll(res.Body)
	_ = json.Unmarshal([]byte(utils.GetJsonAttribute(".[0].versions" ,string(bodyB))), &arr)

	for _ ,v := range arr {
		versionMeta , _ := version.NewVersion(v)
		versions = append(versions ,versionMeta)
	}

	return versions
}


