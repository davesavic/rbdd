package app

import (
	"os"

	"github.com/cucumber/godog"
)

func InitializeTestSuite(ctx *godog.TestSuiteContext) {
	api := NewAPITest(os.Getenv("API_BASE_URL"))
	InitializeScenario(api, ctx.ScenarioContext())
}

func InitializeScenario(api *APITest, ctx *godog.ScenarioContext) {
	// Request steps
	ctx.Step(`^I send a "([^"]*)" request to "([^"]*)"$`, api.iSendRequestTo)
	ctx.Step(`^I send a "([^"]*)" request to "([^"]*)" with payload:$`, api.iSendRequestToWithPayload)

	// Response validation steps
	ctx.Step(`^the response status should be (\d+)$`, api.theResponseStatusShouldBe)
	ctx.Step(`^the response property "([^"]*)" should be (.*?)$`, api.theResponsePropertyShouldBe)
	ctx.Step(`^the response property "([^"]*)" should not be empty$`, api.theResponsePropertyShouldNotBeEmpty)
	ctx.Step(`^the response should match JSON:$`, api.theResponseShouldMatchJSON)
	ctx.Step(`^the response should contain JSON:$`, api.theResponseShouldContainJSON)

	// State management steps
	ctx.Step(`^I store the response property "([^"]*)" as "([^"]*)"$`, api.iStoreTheResponsePropertyAs)
	ctx.Step(`^I store the command output as "([^"]*)"$`, api.iStoreTheCommandOutputAs)
	ctx.Step(`^I store "([^"]*)" as "([^"]*)"$`, api.iStoreAs)
	ctx.Step(`^I set header "([^"]*)" to "([^"]*)"$`, api.iSetHeaderTo)
	ctx.Step(`^I reset all variables$`, api.iResetAllVariables)
	ctx.Step(`^I reset variables "([^"]*)"$`, api.iResetVariables)

	// Data generation steps
	ctx.Step(`^I generate fake data: "([^"]*)"$`, api.iGenerateFakeData)

	// Command execution steps
	ctx.Step(`^I execute command "([^"]*)"$`, api.iExecuteCommand)
	ctx.Step(`^I execute command "([^"]*)" in directory "([^"]*)"$`, api.iExecuteCommandInDirectory)
	ctx.Step(`^I execute command "([^"]*)" with timeout (\d+)$`, api.iExecuteCommandWithTimeout)
	ctx.Step(`^the command output should match "([^"]*)"$`, api.theCommandOutputShouldMatch)
	ctx.Step(`^the command output should contain "([^"]*)"$`, api.theCommandOutputShouldContain)

	// Debugging steps
	ctx.Step(`^I start debugging$`, api.iStartDebugging)
	ctx.Step(`^I stop debugging$`, api.iStopDebugging)
}
