Feature: Command Testing
  Scenario: Execute command store result
    Given I execute command "ls -l"
    And I store command output as "command_result"



