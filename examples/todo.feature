Feature: Todo Testing
  Scenario: Get one todo
    When I send a "GET" request to "https://jsonplaceholder.typicode.com/todos/1"
    Then the response status should be 200
    And the response should match JSON:
      """
      {
        "userId": 1,
        "id": 1,
        "title": "delectus aut autem",
        "completed": false
      }
      """

    Then I start debugging
    And I store the response property "id" as "todo_id"
    Then I execute command "echo ${todo_id}"
    And I stop debugging
