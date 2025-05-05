package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/tidwall/gjson"
)

type APITest struct {
	baseURL       string
	client        *http.Client
	headers       map[string]string
	response      *http.Response
	responseBody  string
	commandOutput string
	store         map[string]any
}

// Global store for variables that can be accessed from other tests
var globalStore = map[string]any{}

func NewAPITest(baseURL string) *APITest {
	return &APITest{
		baseURL: baseURL,
		client:  &http.Client{},
		headers: map[string]string{"Content-Type": "application/json"},
		store:   globalStore,
	}
}

func (a *APITest) replaceVars(text string) string {
	r := regexp.MustCompile(`\${([^}]+)}`)
	return r.ReplaceAllStringFunc(text, func(match string) string {
		// Extract key name without ${ and }
		key := match[2 : len(match)-1]
		if val, ok := a.store[key]; ok {
			return fmt.Sprintf("%v", val)
		}
		return match
	})
}

func (a *APITest) sendRequest(method, endpoint, payload string) error {
	endpoint = a.replaceVars(endpoint)
	payload = a.replaceVars(payload)

	var req *http.Request
	var err error

	if payload != "" {
		req, err = http.NewRequest(method, a.baseURL+endpoint, bytes.NewBufferString(payload))
	} else {
		req, err = http.NewRequest(method, a.baseURL+endpoint, nil)
	}

	if err != nil {
		return err
	}

	for k, v := range a.headers {
		req.Header.Set(k, a.replaceVars(v))
	}

	a.response, err = a.client.Do(req)
	if err != nil {
		return err
	}

	bodyBytes, err := io.ReadAll(a.response.Body)
	if err != nil {
		return err
	}
	a.responseBody = string(bodyBytes)
	a.response.Body.Close()

	return nil
}

func generateFromTag(tag string) (string, error) {
	faker := gofakeit.New(0)
	result, err := faker.Generate(tag)
	if err != nil {
		return "", fmt.Errorf("failed to generate data from tag: %w", err)
	}

	return fmt.Sprintf("%v", result), nil
}

func (a *APITest) generateFakeData(dataSpec string) error {
	pairs := splitByCommaOutsideBrackets(dataSpec)

	for _, pair := range pairs {
		parts := strings.SplitN(strings.TrimSpace(pair), "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid data specification format: %s", pair)
		}

		varName := strings.TrimSpace(parts[0])
		pattern := strings.TrimSpace(parts[1])

		value, err := generateFromTag(pattern)
		if err != nil {
			return fmt.Errorf("error generating fake data for %s: %w", varName, err)
		}

		a.store[varName] = value
	}

	return nil
}

func splitByCommaOutsideBrackets(s string) []string {
	var result []string
	var current strings.Builder
	bracketLevel := 0

	for _, r := range s {
		switch {
		case r == '{' || r == '[':
			bracketLevel++
			current.WriteRune(r)
		case r == '}' || r == ']':
			bracketLevel--
			current.WriteRune(r)
		case r == ',' && bracketLevel == 0:
			result = append(result, current.String())
			current.Reset()
		default:
			current.WriteRune(r)
		}
	}

	if current.Len() > 0 {
		result = append(result, current.String())
	}

	return result
}

func containsSubset(data, subset map[string]any) error {
	for k, expectedVal := range subset {
		actualVal, exists := data[k]
		if !exists {
			return fmt.Errorf("missing key %q", k)
		}

		// Recursive check for nested objects
		if expectedMap, ok := expectedVal.(map[string]any); ok {
			if actualMap, ok := actualVal.(map[string]any); ok {
				if err := containsSubset(actualMap, expectedMap); err != nil {
					return fmt.Errorf("in key %q: %w", k, err)
				}
				continue
			}
		}

		if !reflect.DeepEqual(actualVal, expectedVal) {
			return fmt.Errorf("value mismatch for key %q: expected %v but got %v",
				k, expectedVal, actualVal)
		}
	}
	return nil
}

func (a *APITest) iExecuteCommand(command string) error {
	return a.iExecuteCommandInDirectory(command, "")
}

func (a *APITest) iExecuteCommandInDirectory(command string, dir string) error {
	command = a.replaceVars(command)
	dir = a.replaceVars(dir)

	if strings.TrimSpace(command) == "" {
		return fmt.Errorf("command is empty")
	}

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", command)
	} else {
		cmd = exec.Command("sh", "-c", command)
	}

	if dir != "" {
		cmd.Dir = dir
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	a.commandOutput = strings.Trim(stdout.String(), "\n")
	if err != nil {
		return fmt.Errorf("command failed: %v\nStdout: %s\nStderr: %s",
			err, a.commandOutput, stderr.String())
	}

	return nil
}

func (a *APITest) iExecuteCommandWithTimeout(command string, timeoutSec int) error {
	done := make(chan error)

	go func() {
		done <- a.iExecuteCommand(command)
	}()

	select {
	case err := <-done:
		return err
	case <-time.After(time.Duration(timeoutSec) * time.Second):
		return fmt.Errorf("command timed out after %d seconds: %s", timeoutSec, command)
	}
}

func (a *APITest) iSendRequestTo(method, endpoint string) error {
	return a.sendRequest(method, endpoint, "")
}

func (a *APITest) iSendRequestToWithPayload(method, endpoint, payload string) error {
	return a.sendRequest(method, endpoint, payload)
}

func (a *APITest) theResponseStatusShouldBe(status int) error {
	if a.response.StatusCode != status {
		return fmt.Errorf("expected status %d but got %d with body %s", status, a.response.StatusCode, a.responseBody)
	}
	return nil
}

func (a *APITest) theResponsePropertyShouldBe(property, expectedValue string) error {
	value := gjson.Get(a.responseBody, property)
	expected := a.replaceVars(expectedValue)

	if !value.Exists() {
		return fmt.Errorf("property %s not found in response", property)
	}

	switch {
	case expected == "true" || expected == "false":
		if value.Bool() != (expected == "true") {
			return fmt.Errorf("expected %s to be %s but got %v", property, expected, value.Bool())
		}
	case strings.HasPrefix(expected, "\"") && strings.HasSuffix(expected, "\""):
		expString := expected[1 : len(expected)-1]
		if value.String() != expString {
			return fmt.Errorf("expected %s to be %s but got %s", property, expString, value.String())
		}
	default:
		if expNum, err := strconv.ParseFloat(expected, 64); err == nil {
			if value.Float() != expNum {
				return fmt.Errorf("expected %s to be %f but got %f", property, expNum, value.Float())
			}
		} else if value.String() != expected {
			return fmt.Errorf("expected %s to be %s but got %s", property, expected, value.String())
		}
	}

	return nil
}

func (a *APITest) iStoreTheResponsePropertyAs(property, variable string) error {
	value := gjson.Get(a.responseBody, property)
	if !value.Exists() {
		return fmt.Errorf("property %s not found in response", property)
	}

	switch value.Type {
	case gjson.String:
		a.store[variable] = value.String()
	case gjson.Number:
		a.store[variable] = value.Float()
	case gjson.True, gjson.False:
		a.store[variable] = value.Bool()
	case gjson.JSON:
		a.store[variable] = value.Raw
	default:
		a.store[variable] = value.Raw
	}

	return nil
}

func (a *APITest) iStoreTheCommandOutputAs(variable string) error {
	if a.commandOutput == "" {
		return fmt.Errorf("command output is empty")
	}
	a.store[variable] = a.commandOutput
	return nil
}

func (a *APITest) iStoreAs(variable, value string) error {
	value = a.replaceVars(value)
	if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
		value = value[1 : len(value)-1]
	}
	if strings.HasPrefix(value, "{") && strings.HasSuffix(value, "}") {
		var jsonObj map[string]any
		if err := json.Unmarshal([]byte(value), &jsonObj); err != nil {
			return fmt.Errorf("invalid JSON format: %w", err)
		}
		a.store[variable] = jsonObj
	} else if strings.HasPrefix(value, "[") && strings.HasSuffix(value, "]") {
		var jsonArray []any
		if err := json.Unmarshal([]byte(value), &jsonArray); err != nil {
			return fmt.Errorf("invalid JSON array format: %w", err)
		}
		a.store[variable] = jsonArray
	} else if num, err := strconv.ParseFloat(value, 64); err == nil {
		a.store[variable] = num
	} else if boolVal, err := strconv.ParseBool(value); err == nil {
		a.store[variable] = boolVal
	} else {
		a.store[variable] = value
	}
	return nil
}

func (a *APITest) iSetHeaderTo(header, value string) error {
	a.headers[header] = value
	return nil
}

func (a *APITest) iResetAllVariables() error {
	a.store = map[string]any{}
	return nil
}

func (a *APITest) iResetVariables(variables string) error {
	vars := strings.SplitSeq(variables, ",")
	for variable := range vars {
		varName := strings.TrimSpace(variable)
		delete(a.store, varName)
	}
	return nil
}

func (a *APITest) theResponsePropertyShouldNotBeEmpty(property string) error {
	value := gjson.Get(a.responseBody, property)
	if !value.Exists() || value.String() == "" {
		return fmt.Errorf("property %s is empty or not found", property)
	}
	return nil
}

func (a *APITest) theResponseShouldMatchJSON(expected string) error {
	templated := a.replaceVars(expected)

	var expectedObj any
	var actualObj any

	if err := json.Unmarshal([]byte(templated), &expectedObj); err != nil {
		validJSON, jsonErr := a.makeValidJSON(expected)
		if jsonErr != nil {
			return fmt.Errorf("invalid expected JSON: %w", err)
		}
		if err := json.Unmarshal([]byte(validJSON), &expectedObj); err != nil {
			return fmt.Errorf("still invalid JSON after fixing: %w", err)
		}
	}

	if err := json.Unmarshal([]byte(a.responseBody), &actualObj); err != nil {
		return fmt.Errorf("invalid response JSON: %w", err)
	}

	if !reflect.DeepEqual(expectedObj, actualObj) {
		return fmt.Errorf("JSON mismatch\nExpected: %v\nActual: %v",
			expectedObj, actualObj)
	}

	return nil
}

func (a *APITest) theResponseShouldContainJSON(expected string) error {
	templated := a.replaceVars(expected)

	var expectedMap map[string]any
	var actualMap map[string]any

	if err := json.Unmarshal([]byte(templated), &expectedMap); err != nil {
		validJSON, jsonErr := a.makeValidJSON(expected)
		if jsonErr != nil {
			return fmt.Errorf("invalid expected JSON: %w", err)
		}
		if err := json.Unmarshal([]byte(validJSON), &expectedMap); err != nil {
			return fmt.Errorf("still invalid JSON after fixing: %w", err)
		}
	}

	if err := json.Unmarshal([]byte(a.responseBody), &actualMap); err != nil {
		return fmt.Errorf("invalid response JSON: %w", err)
	}

	if err := containsSubset(actualMap, expectedMap); err != nil {
		return fmt.Errorf("JSON subset mismatch: %w", err)
	}

	return nil
}

func (a *APITest) makeValidJSON(template string) (string, error) {
	var jsonObj map[string]any

	placeholderMap := make(map[string]string)
	uniqueMarker := "_PLACEHOLDER_"
	counter := 0

	re := regexp.MustCompile(`\${([^}]+)}`)
	tempJSON := re.ReplaceAllStringFunc(template, func(match string) string {
		key := match[2 : len(match)-1]
		placeholder := fmt.Sprintf("%s%d%s", uniqueMarker, counter, uniqueMarker)
		placeholderMap[placeholder] = key
		counter++
		return fmt.Sprintf("\"%s\"", placeholder)
	})

	if err := json.Unmarshal([]byte(tempJSON), &jsonObj); err != nil {
		return "", fmt.Errorf("malformed JSON template: %w", err)
	}

	var replaceInValue func(any) any
	replaceInValue = func(val any) any {
		switch v := val.(type) {
		case string:
			for placeholder, varName := range placeholderMap {
				if strings.Contains(v, placeholder) {
					if storeVal, ok := a.store[varName]; ok {
						return storeVal
					}
				}
			}
			return v
		case map[string]any:
			result := make(map[string]any)
			for k, mapVal := range v {
				result[k] = replaceInValue(mapVal)
			}
			return result
		case []any:
			result := make([]any, len(v))
			for i, arrVal := range v {
				result[i] = replaceInValue(arrVal)
			}
			return result
		default:
			return v
		}
	}

	processedObj := replaceInValue(jsonObj)

	result, err := json.Marshal(processedObj)
	if err != nil {
		return "", fmt.Errorf("failed to create valid JSON: %w", err)
	}

	return string(result), nil
}
