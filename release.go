package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const (
	githubAPIURL = "https://api.github.com/repos/"
	githubAPIKEY = "GITHUB_TOKEN"
)

// CreateReleaseReq is CreateReleaseReq's struct
type CreateReleaseReq struct {
	TagName         string `json:"tag_name"`
	TargetCommitish string `json:"target_commitish"`
	Name            string `json:"name"`
	Prerelease      bool   `json:"prerelease"`
}

func main() {
	if len(os.Args) < 3 {
		panic("args must be need")
	}
	var (
		repoLocation = os.Args[1]
		branch       = os.Args[2]
		tag          = os.Args[3]
	)

	// CreateRelease Request Data
	var crreqdat = CreateReleaseReq{
		TagName:         tag,
		TargetCommitish: branch,
		Name:            tag,
		Prerelease:      branch != "master",
	}

	fmt.Println(repoLocation, branch, tag)
	fmt.Println(crreqdat)

	params, err := json.Marshal(crreqdat)
	if err != nil {
		panic(err)
	}

	relinforaw, err := httpGet(githubAPIURL + repoLocation + "/releases")
	var relinfos = make([]map[string]interface{}, 0)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(relinforaw))

	err = json.Unmarshal(relinforaw, &relinfos)
	if err != nil {
		panic(err)
	}

	for _, relinfo := range relinfos {
		if relinfo["id"] != nil && relinfo["draft"] == true {
			_, err = httpDelete(githubAPIURL + repoLocation + "/releases/" + fmt.Sprintf("%.0f", relinfo["id"].(float64)))
			if err != nil {
				panic(err)
			}
		}
	}

	relinforaw, err = httpPost(githubAPIURL+repoLocation+"/releases", "application/json", bytes.NewBuffer(params))
	if err != nil {
		panic(err)
	}

	var relinfo = make(map[string]interface{})

	fmt.Println(string(relinforaw))

	err = json.Unmarshal(relinforaw, &relinfo)
	if err != nil {
		panic(err)
	}

	fmt.Println(relinfo)

	agdfile, err := os.Open("agd.tar.gz")
	if err != nil {
		panic(err)
	}

	_, err = httpPost(strings.Replace(relinfo["upload_url"].(string), "{?name,label}", "?name=agd.tar.gz", 1), "application/octet-stream", agdfile)
	if err != nil {
		panic(err)
	}

	glfile, err := os.Open("gl.tar.gz")
	if err != nil {
		panic(err)
	}
	_, err = httpPost(strings.Replace(relinfo["upload_url"].(string), "{?name,label}", "?name=gl.tar.gz", 1), "application/octet-stream", glfile)
	if err != nil {
		panic(err)
	}
}

func httpPost(url string, ctype string, dat io.Reader) ([]byte, error) {
	token := os.Getenv(githubAPIKEY)
	buf := &bytes.Buffer{}
	nRead, err := io.Copy(buf, dat)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(
		"POST",
		url,
		buf,
	)
	if err != nil {
		return nil, err
	}

	// Content-Type 設定
	req.Header.Set("Content-Type", ctype)
	req.Header.Set("Content-Length", strconv.FormatInt(nRead, 10))
	req.Header.Set("Authorization", "token "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

func httpDelete(url string) ([]byte, error) {
	token := os.Getenv(githubAPIKEY)
	req, err := http.NewRequest(
		"DELETE",
		url,
		nil,
	)
	if err != nil {
		return nil, err
	}

	// Content-Type 設定
	req.Header.Set("Authorization", "token "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

func httpGet(url string) ([]byte, error) {
	token := os.Getenv(githubAPIKEY)

	req, err := http.NewRequest(
		"GET",
		url,
		nil,
	)
	if err != nil {
		return nil, err
	}

	// Content-Type 設定
	req.Header.Set("Authorization", "token "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}
