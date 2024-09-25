package config

import (
	"fmt"
	"os/exec"
	"sort"
	"strconv"
	"strings"

	config "HaseFlrn/ollama_commit/lib/config"
	inputany "HaseFlrn/ollama_commit/lib/inputAny"

	"github.com/charmbracelet/huh"
)

func Config() {
	persistedConf := config.GetConfig()

	newConfig := persistedConf

	// Create a new huh form
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				TitleFunc(func() string {
					return fmt.Sprintf("Chose a model (current: %s)", persistedConf.Model)
				}, persistedConf.Model).
				OptionsFunc(func() []huh.Option[string] { return getOllamaModels(persistedConf.Model) }, "").
				Value(&newConfig.Model),
			inputany.NewInputAny(marshalFloat32, unmarshalFloat32).
				Title(fmt.Sprintf("Chose the temperature (current: %f)", persistedConf.Temperature)).
				Value(&newConfig.Temperature),
			inputany.NewInputAny(marshalInt, unmarshalInt).
				Title(fmt.Sprintf("Chose the port (current: %d)", persistedConf.Ollama_Port)).
				Validate(validatePort).
				Value(&newConfig.Ollama_Port),
		),
	)

	form.Run()

	fmt.Printf("New Config: %+v\n", newConfig)

}

func marshalFloat32(v float32) string {
	return strconv.FormatFloat(float64(v), 'f', -1, 32)
}

func unmarshalFloat32(v string) (float32, error) {
	r, e := strconv.ParseFloat(v, 32)
	return float32(r), e
}

func marshalInt(input int) string {
	return strconv.FormatInt(int64(input), 10)
}

func unmarshalInt(input string) (int, error) {
	return strconv.Atoi(input)
}

func validatePort(port int) error {
	if port < 0 || port > 65535 {
		return fmt.Errorf("invalid port")
	}
	return nil
}

func getOllamaModels(currentModel string) []huh.Option[string] {
	res, err := exec.Command("ollama", "list").Output()
	if err != nil {
		panic(err)
	}
	lines := strings.Split(string(res), "\n")
	models := []string{}
	for _, line := range lines {
		fields := strings.Split(line, " ")
		if len(fields) > 0 {
			models = append(models, fields[0])
		}
	}
	models = models[1 : len(models)-1]

	sort.Slice(models, func(i, j int) bool {
		if models[i] == currentModel {
			return true
		}
		if models[j] == currentModel {
			return false
		}
		return models[i] < models[j]
	})

	if len(models) == 0 {
		// TODO: Check if color-coding is doable
		fmt.Println("No models found. Please create/pull a model first.")
		fmt.Println("Browse models at https://ollama.com/library to check available models.")
		fmt.Println("For how to create a model, visit https://github.com/ollama/ollama?tab=readme-ov-file#customize-a-model")
	}

	return huh.NewOptions(models...)
}
