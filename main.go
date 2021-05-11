package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"gopkg.in/yaml.v2"
)

// Config for the application
type Config struct {
	MainRepo string `yaml:"mainrepo"`
	Github   struct {
		User     string `yaml:"user"`
		Apitoken string `yaml:"apitoken"`
	}
}

// PullRequest incomming pull request
type PullRequest interface {
	Comment(string)
	FetchRemote(string) error
}

type gitHubPR struct {
	Number int    `json:"number"`
	Action string `json:"action"`

	PullRequest struct {
		CommentsURL string `json:"comments_url"`
		Head        struct {
			Ref string `json:"ref"`
			SHA string `json:"sha"`

			Repo struct {
				CloneURL string `json:"clone_url"`
			} `json:"repo"`
		} `json:"head"`
	} `json:"pull_request"`
}

func (pr gitHubPR) Comment(body string) {
	jsonStr := []byte(fmt.Sprintf("{\"body\": \"" + body + "\"}"))

	req, err := http.NewRequest("POST", pr.PullRequest.CommentsURL, bytes.NewBuffer(jsonStr))
	if err != nil {
		panic(err)
	}
	req.SetBasicAuth(config.Github.User, config.Github.Apitoken)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		panic(err)
	}
}

func (pr gitHubPR) FetchRemote(path string) error {
	return pullBranch(pr.PullRequest.Head.Repo.CloneURL, pr.PullRequest.Head.SHA, path)
}

func handler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Panic(err)
	}

	pr := gitHubPR{}
	err = json.Unmarshal(body, &pr)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Body:\n%s\n", string(body))

	handlePullRequest(pr)
}

func handlePullRequest(pr PullRequest) {
	pr.Comment("I see you")

	tDir, err := ioutil.TempDir("", "targetDir")
	if err != nil {
		pr.Comment("Internal error")
		panic(err)
	}
	defer os.RemoveAll(tDir)

	err = checkoutBranch(config.MainRepo, "master", tDir)
	if err != nil {
		pr.Comment("Internal error")
		panic(err)
	}
	err = pr.FetchRemote(tDir)
	if err != nil {
		pr.Comment("Unable to merge request")
		panic(err)
	}

	err = pushBranch(tDir)
	if err != nil {
		pr.Comment("Internal error")
		panic(err)
	}
	pr.Comment("Merged your request")
}

var config Config

func main() {
	yamlFile, err := ioutil.ReadFile("config.yml")
	if err != nil {
		fmt.Printf("Unable to open file 'config.yml'")
		return
	}
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		fmt.Printf("Unable to parse yaml file")
		panic(err)
	}

	http.HandleFunc("/github", handler)
	log.Fatal(http.ListenAndServe("192.168.1.210:8380", nil))
}
