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

type CreateReleaseReq struct {
	TagName         string `json:"tag_name"`
	TargetCommitish string `json:"target_commitish"`
	Name            string `json:"name"`
	Prerelease      bool   `json:"prerelease"`
}

func main() {
	var RepoLocation = "AizuGeekDojo/EnterLeaveSystem"
	var Branch = os.Args[1]
	var Tag = os.Args[2]

	// CreateRelease Request Data
	var crreqdat = CreateReleaseReq{
		TagName:         Tag,
		TargetCommitish: Branch,
		Name:            Tag,
		Prerelease:      true,
	}
	params, err := json.Marshal(crreqdat)
	if err != nil {
		panic(err)
	}
	relinforaw, err := HttpGet("https://api.github.com/repos/" + RepoLocation + "/releases")
	var relinfo = make(map[string]interface{})
	if err == nil {
		println(string(relinforaw))

		err = json.Unmarshal(relinforaw, &relinfo)
		if err != nil {
			panic(err)
		}

		delrelinfo, err := HttpDelete("https://api.github.com/repos/" + RepoLocation + "/releases/" + relinfo["id"].(string))
		if err != nil {
			panic(err)
		}
		println(string(delrelinfo))
	}
	relinforaw, err = HttpPost("https://api.github.com/repos/"+RepoLocation+"/releases", "application/json", bytes.NewBuffer(params))
	if err != nil {
		panic(err)
	}
	println(string(relinforaw))
	// var ExtList = make(map[string]interface{})

	err = json.Unmarshal(relinforaw, &relinfo)
	if err != nil {
		panic(err)
	}
	agdfile, err := os.Open("agd.tar.gz")
	if err != nil {
		panic(err)
	}

	upagdres, err := HttpPost(strings.Replace(relinfo["upload_url"].(string), "{?name,label}", "?name=agd.tar.gz", 1), "application/octet-stream", agdfile)
	println(string(upagdres))
	glfile, err := os.Open("gl.tar.gz")
	if err != nil {
		panic(err)
	}
	upglres, err := HttpPost(strings.Replace(relinfo["upload_url"].(string), "{?name,label}", "?name=gl.tar.gz", 1), "application/octet-stream", glfile)
	println(string(upglres))
}

func HttpPost(url string, ctype string, dat io.Reader) ([]byte, error) {
	token := os.Getenv("GITHUB_API_KEY")
	buf := &bytes.Buffer{}
	nRead, err := io.Copy(buf, dat)
	if err != nil {
		fmt.Println(err)
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
	content, err := ioutil.ReadAll(resp.Body)
	return content, err
}

func HttpDelete(url string) ([]byte, error) {
	token := os.Getenv("GITHUB_API_KEY")
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
		println(err)
		// return nil, err
	}
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	return content, err
}

func HttpGet(url string) ([]byte, error) {
	token := os.Getenv("GITHUB_API_KEY")
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
	content, err := ioutil.ReadAll(resp.Body)
	return content, err
}
