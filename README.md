# rbdd

A lightweight CLI tool for testing RESTful APIs using behavior-driven development with Gherkin syntax.

## Installation
```bash
go install github.com/davesavic/rbdd@latest
```
or you can use the precompiled binaries for your platform.

Download the latest release from the [releases page](https://github.com/davesavic/rbdd/releases)

## Quick start
1. Create a feature file with Gherkin syntax. For example, `features/users.feature`:
```gherkin
# features/users.feature
Feature: User API

  Scenario: Create a new user
    Given I generate fake data: "email={email}, name={firstname} {lastname}"
    When I send a "POST" request to "/api/users" with payload:
      """
      {
        "name": "${name}",
        "email": "${email}"
      }
      """
    Then the response status should be 201
    And I store the response property "id" as "user_id"

    When I send a "GET" request to "/api/users/${user_id}"
    Then the response status should be 200
    And the response property "email" should be "${email}"
```

2. Run the feature file using rbdd:
```bash
API_BASE_URL=http://localhost:8080/api rbdd -d features
```

## Usage
```bash
➜  rbdd --help
A simple command line tool to test your backend api using cucumber tests written in gherkin syntax. Easily generate fake data using faker 
and test your api with a few simple commands. Run your tests in parallel and get the results in a simple format.

Usage:
  rbdd [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  run         Run the cucumber tests
  syntax      Show the gherkin syntax available

Flags:
      --config string   config file (default is $HOME/.rbdd.env)
  -h, --help            help for rbdd
  -t, --toggle          Help message for toggle

Use "rbdd [command] --help" for more information about a command.
```

## Features
### Requests
```gherkin
When I send a "GET" request to "/users"
When I send a "POST" request to "/users" with payload:
  """
  {
    "name": "John"
  }
  """
```

### Response validation
```gherkin
Then the response status should be 200
Then the response property "data.id" should be 123
Then the response property "data.email" should not be empty
Then the response should match JSON:
  """
  {
    "success": true
  }
  """
Then the response should contain JSON:
  """
  {
    "data": {
      "active": true
    }
  }
  """
``` 

### State management
```gherkin
Given I set header "Authorization" to "Bearer token123"
When I store the response property "id" as "user_id"
When I globally store the response property "access_token" as "token"
When I reset all scoped variables
When I reset all global variables
```

### Data generation
```gherkin
Given I generate fake data: "email=email, name=name, id=uuid"
```

### Command execution
```gherkin
When I execute command "echo Hello"
When I execute command "npm test" in directory "./frontend"
When I execute command "docker-compose up -d" with timeout 30
```

### Variable substitution
Use stored variables anywhere with ${variable_name} syntax:
```gherkin
# In URLs
When I send a "GET" request to "/users/${user_id}"

# In payloads
When I send a "POST" request to "/orders" with payload:
  """
  {
    "user_id": "${user_id}",
    "product_id": "${product_id}"
  }
  """

# In headers
Given I set header "Authorization" to "Bearer ${token}"`
```

## License
This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
