/*
Copyright © 2025 Dave Savic
*/

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// syntaxCmd represents the syntax command
var syntaxCmd = &cobra.Command{
	Use:   "syntax",
	Short: "Show the gherkin syntax available",
	Long:  `Show the gherkin syntax available for use in the tests`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Available Gherkin Syntax")
		fmt.Println("=======================")
		fmt.Println("")

		// Request Steps
		fmt.Println("📤 Request Steps")
		fmt.Println("-------------")
		fmt.Println("  Given/When I send a \"GET\" request to \"/api/users\"")
		fmt.Println("  Given/When I send a \"POST\" request to \"/api/users\" with payload:")
		fmt.Println("    ```")
		fmt.Println("    {")
		fmt.Println("      \"name\": \"John Doe\",")
		fmt.Println("      \"email\": \"${email}\"")
		fmt.Println("    }")
		fmt.Println("    ```")
		fmt.Println("")

		// Response Validation Steps
		fmt.Println("✅ Response Validation Steps")
		fmt.Println("------------------------")
		fmt.Println("  Then the response status should be 200")
		fmt.Println("  Then the response property \"data.id\" should be 123")
		fmt.Println("  Then the response property \"data.email\" should be \"test@example.com\"")
		fmt.Println("  Then the response property \"data.created_at\" should not be empty")
		fmt.Println("  Then the response should match JSON:")
		fmt.Println("    ```")
		fmt.Println("    {")
		fmt.Println("      \"success\": true,")
		fmt.Println("      \"data\": {")
		fmt.Println("        \"id\": 123")
		fmt.Println("      }")
		fmt.Println("    }")
		fmt.Println("    ```")
		fmt.Println("  Then the response should contain JSON:")
		fmt.Println("    ```")
		fmt.Println("    {")
		fmt.Println("      \"success\": true")
		fmt.Println("    }")
		fmt.Println("    ```")
		fmt.Println("")

		// State Management Steps
		fmt.Println("🔄 State Management Steps")
		fmt.Println("-----------------------")
		fmt.Println("  Given/When I store the response property \"data.id\" as \"user_id\"")
		fmt.Println("  Given/When I set header \"Authorization\" to \"Bearer ${token}\"")
		fmt.Println("  Given/When I reset all stored variables")
		fmt.Println("")

		// Data Generation Steps
		fmt.Println("🔮 Data Generation Steps")
		fmt.Println("---------------------")
		fmt.Println("  Given/When I generate fake data: \"email={email}, name={firstname}, phone={phone}\"")
		fmt.Println("    # Supported patterns use gofakeit syntax https://github.com/brianvoe/gofakeit")
		fmt.Println("    # Examples: {firstname} {lastname} {email} {phone} etc.")
		fmt.Println("")

		// Command Execution Steps
		fmt.Println("⚙️  Command Execution Steps")
		fmt.Println("-----------------------")
		fmt.Println("  Given/When I execute command \"echo Hello World\"")
		fmt.Println("  Given/When I execute command \"npm install\" in directory \"./frontend\"")
		fmt.Println("  Given/When I execute command \"gradle test\" with timeout 30")
		fmt.Println("")

		// Variable Substitution
		fmt.Println("🔠 Variable Substitution")
		fmt.Println("---------------------")
		fmt.Println("  You can use stored variables in any step with ${variable_name} syntax:")
		fmt.Println("    - In URLs: \"/api/users/${user_id}\"")
		fmt.Println("    - In payloads: \"{ \"user_id\": ${user_id} }\"")
		fmt.Println("    - In headers: \"Bearer ${token}\"")
		fmt.Println("    - In commands: \"curl -X GET ${api_url}\"")
		fmt.Println("")

		// Example Feature
		fmt.Println("📝 Example Feature")
		fmt.Println("--------------")
		fmt.Println("```")
		fmt.Println("Feature: User API Tests")
		fmt.Println("")
		fmt.Println("  Scenario: Create and retrieve a user")
		fmt.Println("    Given I generate fake data: \"email=email, name={firstname} {lastname}\"")
		fmt.Println("    When I send a \"POST\" request to \"/api/users\" with payload:")
		fmt.Println("      \"\"\"")
		fmt.Println("      {")
		fmt.Println("        \"name\": \"${name}\",")
		fmt.Println("        \"email\": \"${email}\"")
		fmt.Println("      }")
		fmt.Println("      \"\"\"")
		fmt.Println("    Then the response status should be 201")
		fmt.Println("    And the response property \"success\" should be true")
		fmt.Println("    And I store the response property \"data.id\" as \"user_id\"")
		fmt.Println("")
		fmt.Println("    When I send a \"GET\" request to \"/api/users/${user_id}\"")
		fmt.Println("    Then the response status should be 200")
		fmt.Println("    And the response property \"data.email\" should be \"${email}\"")
		fmt.Println("```")
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
