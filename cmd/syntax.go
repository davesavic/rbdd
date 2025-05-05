/*
Copyright © 2025 Dave Savic
*/

package cmd

import (
	"github.com/spf13/cobra"
)

// syntaxCmd represents the syntax command
var syntaxCmd = &cobra.Command{
	Use:   "syntax",
	Short: "Show the gherkin syntax available",
	Long:  `Show the gherkin syntax available for use in the tests`,
	Run: func(cmd *cobra.Command, args []string) {
		output := `
--- Making requests ---
Gherkin Syntax: I send a "METHOD" request to "ENDPOINT"
Description: This step sends a request to the specified endpoint using the specified HTTP method (GET, POST, PUT, DELETE, etc.).
Example: Given/When I send a "GET" request to "/api/users"

Gherkin Syntax: I send a "METHOD" request to "ENDPOINT" with payload:
Description: This step sends a request to the specified endpoint using the specified HTTP method with the provided payload.
Example: 
Given/When I send a "POST" request to "/api/users" with payload:
  """
  {
	"name": "John Doe",
	"email": "john@example.com"
  }
  """

--- Validating responses ---
Gherkin Syntax: the response status should be STATUS_CODE
Description: This step checks if the response status code matches the expected status code.
Example: Then the response status should be 200

Gherkin Syntax: the response property "JSON_PATH" should be VALUE
Description: This step checks if the response property at the specified JSON path matches the expected value.
Example: Then the response property "data.user.email" should be "john@example.com"

Gherkin Syntax: the response property "JSON_PATH" should not be empty
Description: This step checks if the response property at the specified JSON path is not empty.
Example: Then the response property "data.user.id" should not be empty

Gherkin Syntax: the response should match JSON:
Description: This step checks if the entire response matches the expected JSON structure.
Example:
Then the response should match JSON:
  """
  {
	"success": true,
	"data": {
	  "id": 123,
	  "name": "John Doe"
	}
  }
  """

Gherkin Syntax: the response should contain JSON:
Description: This step checks if the response contains the specified JSON structure.
Example:
Then the response should contain JSON:
  """
  {
	"success": true
  }
  """

--- Managing state ---
Gherkin Syntax: I store the response property "JSON_PATH" as "VARIABLE_NAME"
Description: This step stores the value of the response property at the specified JSON path into a variable.
Example: And I store the response property "data.token" as "auth_token"

Gherkin Syntax: I store the command output as "VARIABLE_NAME"
Description: This step stores the output of a command into a variable.
Example: And I store the command output as "db_result"

Gherkin Syntax: I store "VALUE" as "VARIABLE_NAME"
Description: This step stores a specified value into a variable.
Example: And I store "Bearer ${auth_token}" as "authorization"

Gherkin Syntax: I set header "HEADER_NAME" to "HEADER_VALUE"
Description: This step sets a specified header to a specified value.
Example: And I set header "Authorization" to "${authorization}"

Gherkin Syntax: I reset all variables
Description: This step resets all stored variables to their initial state.
Example: And I reset all variables

Gherkin Syntax: I reset variables "VARIABLE_LIST"
Description: This step resets specified variables to their initial state.
Example: And I reset variables "user_id, auth_token"

--- Data generation ---
Gherkin Syntax: I generate fake data: "PATTERN"
Description: This step generates fake data based on the specified pattern using the gofakeit library.
Example: Given I generate fake data: "email={email}, name={firstname} {lastname}, phone={phone}"

--- Command execution ---
Gherkin Syntax: I execute command "COMMAND"
Description: This step executes a specified command in the shell.
Example: When I execute command "echo 'Hello World'"

Gherkin Syntax: I execute command "COMMAND" in directory "DIRECTORY"
Description: This step executes a specified command in the shell within a specified directory.
Example: When I execute command "npm install" in directory "./frontend"

Gherkin Syntax: I execute command "COMMAND" with timeout SECONDS
Description: This step executes a specified command in the shell with a specified timeout in seconds.
Example: When I execute command "gradle build" with timeout 30

Gherkin Syntax: the command output should match "PATTERN"
Description: This step checks if the command output matches the specified pattern.
Example: Then the command output should match "Success"

Gherkin Syntax: the command output should contain "TEXT"
Description: This step checks if the command output contains the specified text.
Example: Then the command output should contain "Build completed"
		`
		cmd.Println(output)
	},
}

func init() {
	rootCmd.AddCommand(syntaxCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// syntaxCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// syntaxCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
