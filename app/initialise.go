package app

import (
	"os"

	"github.com/cucumber/godog"
)

func InitializeScenario(ctx *godog.ScenarioContext) {
	api := NewAPITest(os.Getenv("API_BASE_URL"))

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
	ctx.Step(`^I store "([^"]*)" as "([^"]*)"$`, api.iStoreAs)
	ctx.Step(`^I set header "([^"]*)" to "([^"]*)"$`, api.iSetHeaderTo)
	ctx.Step(`^I reset all variables$`, api.iResetAllVariables)
	ctx.Step(`^I reset variables "([^"]*)"$`, api.iResetVariables)

	// Data generation steps
	ctx.Step(`^I generate fake data: "([^"]*)"$`, api.generateFakeData)

	// Command execution steps
	ctx.Step(`^I execute command "([^"]*)"$`, api.iExecuteCommand)
	ctx.Step(`^I execute command "([^"]*)" in directory "([^"]*)"$`, api.iExecuteCommandInDirectory)
	ctx.Step(`^I execute command "([^"]*)" with timeout (\d+)$`, api.iExecuteCommandWithTimeout)
}
