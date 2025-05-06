package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"regexp"
	"strings"

	"github.com/brianvoe/gofakeit/v7"
)

type APITest struct {
	baseURL       string
	client        *http.Client
	headers       map[string]string
	response      *http.Response
	responseBody  string
	commandOutput string
	store         map[string]any
	debug         bool
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

func generateFromTag(tag string) (string, error) {
	faker := gofakeit.New(0)
	result, err := faker.Generate(tag)
	if err != nil {
		return "", fmt.Errorf("failed to generate data from tag: %w", err)
	}

	return fmt.Sprintf("%v", result), nil
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
