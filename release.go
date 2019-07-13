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
	var RepoLocation = "wancom/sandbox"
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
	a, err := HttpPost("https://api.github.com/repos/"+RepoLocation+"/releases", "application/json", bytes.NewBuffer(params))
	if err != nil {
		panic(err)
	}
	println(string(a))
	var ExtList = make(map[string]interface{})

	err = json.Unmarshal(a, &ExtList)
	if err != nil {
		panic(err)
	}
	url := strings.Replace(ExtList["upload_url"].(string), "{?name,label}", "?name=agd.tar.gz", 1)
	url2 := strings.Replace(ExtList["upload_url"].(string), "{?name,label}", "?name=gl.tar.gz", 1)
	file, err := os.Open("agd.tar.gz")
	if err != nil {
		panic(err)
	}

	b, err := HttpPost(url, "application/octet-stream", file)
	println(string(b))
	file2, err := os.Open("gl.tar.gz")
	if err != nil {
		panic(err)
	}
	c, err := HttpPost(url2, "application/octet-stream", file2)
	println(string(c))
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
