package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

func main () {
	checkPrerequisites()
	result := isGitRepo()

	if !result {
		fmt.Println("This is not a git repository. Please run this command in a git repository.")
		os.Exit(1)
	}

	diff := getGitDiff()

	// Build the prompt
	prompt := buildPrompt(diff)

	// Build the commit message
	commitMessage := askOllama(prompt)

	fmt.Printf("Commit Message:\n\n%v\n", commitMessage)

	if !promptConfirmation() {
		fmt.Println("Exiting...")
		os.Exit(0)
	}

	// Commit the changes to the repo 
	commitChanges(commitMessage)
}

func checkPrerequisites() {
	if !commandExists("git") {
		fmt.Println("Git is not installed. Please install Git and try again.")
		os.Exit(1)
	}
	// Add more checks here (ollama,...)
}

func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func isGitRepo() bool {
	_, err := exec.Command("git", "rev-parse", "--is-inside-work-tree", "&>/dev/null").Output()
	return err != nil 
}

func getGitDiff() string {
	r, err := exec.Command("git", "diff", "--cached").Output()

	if err != nil {
		log.Fatal(err)
	}

	if len(r) == 0 {
		fmt.Print("No changes to commit. Make sure you add the changes to the staging area before running this command.\n")
		os.Exit(0)
	}

	return string(r)
}

func promptConfirmation() bool {
	fmt.Println("\n--------------------------------")
	fmt.Print("Do you want to continue? (y/n): ")
	var response string
	fmt.Scanln(&response)
	response = strings.ToLower(response)
	response = strings.TrimSpace(response)
	return response == "y" || response == "yes"
}

func buildPrompt(diff string) string {
	base := `You are a seasoned developer, writing your commit message. 
Your answer must convey the following commit message format exactly:

<type>([scope]): <description>

[optional body]

With <type> being one of the following:
[ 'build', 'chore', 'ci', 'docs', 'feat', 'fix', 'perf', 'refactor', 'revert', 'style', 'test' ]

Whith scope being an field that is used to specify the module that the commit is related to.

With <description> being a short and concise description of the changes made.

With [optional body] being a more detailed description of the changes made. If it is used, it must be separated from the description by a empty line.

You must create a concise and descriptive commit message based on the following diff:

%v

Only answer with the commit message as plain text and do not describe your reasoning.
If there are a lot of changes, you are allowed to use bullet points as description body.

- DO NOT LIE.
- DO NOT HALLUCINATE. 
- DO NOT ADD ANY NUMBERS OR SPECIAL CHARACTERS.
- DO NOT ADD A FULL STOP AT THE END.
- RESPOND ONLY WITH THE COMMIT MESSAGE.`
	prompt := fmt.Sprintf(base, diff)

	return prompt
}

type OllamaResponse struct {
	Model string `json:"model"`
	Created_At string `json:"created_at"`
	Response string `json:"response"`
	Done bool `json:"done"`
	Done_Reason string `json:"done_reason"`
	Context []int `json:"context"`
	Total_Duration int `json:"total_duration"`
	Load_Duration int `json:"load_duration"`
	Prompt_Eval_Count int `json:"prompt_eval_count"`
	Prompt_Eval_Duration int `json:"prompt_eval_duration"`
	Eval_Count int `json:"eval_count"`
	Eval_Duration int `json:"eval_duration"` 
}

func askOllama(prompt string) string {
	_, err := http.Get("http://localhost:11434/")
	if err != nil {
		fmt.Println("Ollama is not running. Please start Ollama and try again.")
		os.Exit(1)
	}


	// TODO: Make model configurable
	requestData := map[string]interface{}{
		"model": "llama3",
		"prompt": prompt,
		"stream": false,
		"options": map[string]float32{
			"temperature": 0.1,
		},
	}

	jsonData, err := json.Marshal(requestData)
	if err != nil {
		fmt.Println("Error while marshalling JSON data.")
		os.Exit(1)
	}

	r, _ := http.NewRequest("POST","http://localhost:11434/api/generate", bytes.NewBuffer(jsonData) )
	r.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(r)
	if err != nil {
		fmt.Println("Error while sending request to Ollama.")
		os.Exit(1)
	}

	defer resp.Body.Close()

	body,_:= io.ReadAll(resp.Body)

	var ollamaResponse OllamaResponse
	err = json.Unmarshal(body, &ollamaResponse)
	if err != nil {
		fmt.Printf("Error while unmarshalling JSON data. %v\n", err)
		os.Exit(1)
	}
	
	return ollamaResponse.Response
}

func commitChanges(commitMessage string) {
	// Commit the changes to the repo
	_, err := exec.Command("git", "commit", "-m", commitMessage).Output()
	if err != nil {
		log.Fatal(err)
	}
}