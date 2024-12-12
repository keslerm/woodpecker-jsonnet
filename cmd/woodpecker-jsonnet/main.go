package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/google/go-jsonnet"
	"github.com/itchyny/json2yaml"
)

type Payload struct {
	Repo struct {
		ID            int    `json:"id"`
		UID           string `json:"uid"`
		UserID        int    `json:"user_id"`
		Namespace     string `json:"namespace"`
		Name          string `json:"name"`
		Slug          string `json:"slug"`
		Scm           string `json:"scm"`
		GitHTTPURL    string `json:"git_http_url"`
		GitSSHURL     string `json:"git_ssh_url"`
		Link          string `json:"link"`
		DefaultBranch string `json:"default_branch"`
		Private       bool   `json:"private"`
		Visibility    string `json:"visibility"`
		Active        bool   `json:"active"`
		Config        string `json:"config"`
		Trusted       bool   `json:"trusted"`
		Protected     bool   `json:"protected"`
		IgnoreForks   bool   `json:"ignore_forks"`
		IgnorePulls   bool   `json:"ignore_pulls"`
		CancelPulls   bool   `json:"cancel_pulls"`
		Timeout       int    `json:"timeout"`
		Counter       int    `json:"counter"`
		Synced        int    `json:"synced"`
		Created       int    `json:"created"`
		Updated       int    `json:"updated"`
		Version       int    `json:"version"`
	} `json:"repo"`
	Pipeline struct {
		Author       string   `json:"author"`
		AuthorAvatar string   `json:"author_avatar"`
		AuthorEmail  string   `json:"author_email"`
		Branch       string   `json:"branch"`
		ChangedFiles []string `json:"changed_files"`
		Commit       string   `json:"commit"`
		CreatedAt    int      `json:"created_at"`
		DeployTo     string   `json:"deploy_to"`
		EnqueuedAt   int      `json:"enqueued_at"`
		Error        string   `json:"error"`
		Event        string   `json:"event"`
		FinishedAt   int      `json:"finished_at"`
		ID           int      `json:"id"`
		LinkURL      string   `json:"link_url"`
		Message      string   `json:"message"`
		Number       int      `json:"number"`
		Parent       int      `json:"parent"`
		Ref          string   `json:"ref"`
		Refspec      string   `json:"refspec"`
		CloneURL     string   `json:"clone_url"`
		ReviewedAt   int      `json:"reviewed_at"`
		ReviewedBy   string   `json:"reviewed_by"`
		Sender       string   `json:"sender"`
		Signed       bool     `json:"signed"`
		StartedAt    int      `json:"started_at"`
		Status       string   `json:"status"`
		Timestamp    int      `json:"timestamp"`
		Title        string   `json:"title"`
		UpdatedAt    int      `json:"updated_at"`
		Verified     bool     `json:"verified"`
	} `json:"pipeline"`
	Configs []Config `json:"configs"`
}

type Output struct {
	Configs []Config `json:"configs"`
}

type Config struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func main() {
	fmt.Println("Starting")
	server()
}

func server() {
	vm := jsonnet.MakeVM()

	http.HandleFunc("/ci", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		payload := &Payload{}
		err := json.NewDecoder(r.Body).Decode(payload)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		fmt.Println("got payload:", payload)

		result := &Output{}

		for _, c := range payload.Configs {

			if strings.HasSuffix(c.Name, "jsonnet") {
				jsonStr, err := vm.EvaluateAnonymousSnippet(c.Name, c.Data)
				if err != nil {
					log.Fatal(err)
				}

				var output strings.Builder
				input := strings.NewReader(jsonStr)
				if err := json2yaml.Convert(&output, input); err != nil {
					log.Fatalln(err)
				}

				result.Configs = append(result.Configs, Config{
					Name: c.Name,
					Data: output.String(),
				})
			} else {
				result.Configs = append(result.Configs, c)
			}
		}

		str, err := json.Marshal(result)
		if err != nil {
			log.Fatal(err)
		}

		w.Write(str)
	})

	if err := http.ListenAndServe(":8080", nil); err != http.ErrServerClosed {
		panic(err)
	}
}
