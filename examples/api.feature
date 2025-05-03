Feature: API Testing
  Background:
    Given I execute command "task md" in directory "./"
    And I execute command "task mu" in directory "./"
    And I execute command "task sa" in directory "./"

  Scenario: Register and create a profile
    Given I generate fake data: "email={email}, password={password:true,true,true,true,false,20}"

    When I send a "POST" request to "/register" with payload:
      """
      {
        "email": "${email}",
        "password": "${password}",
        "confirm_password": "${password}"
      }
      """
    Then the response status should be 201
    And the response property "account_id" should not be empty
    And I store the response property "access_token" as "token"
    And I store the response property "account_id" as "account_id"
    And I set header "Authorization" to "Bearer ${token}"

    Given I generate fake data: "first_name={firstname}, last_name={lastname}, phone={phone}"

    When I send a "POST" request to "/profile" with payload:
      """
      {
        "first_name": "${first_name}",
        "last_name": "${last_name}",
        "phone": "${phone}"
      }
      """
    Then the response status should be 201
    And I store the response property "id" as "profile_id"

    When I send a "GET" request to "/profile"
    Then the response status should be 200
    And the response should contain JSON:
    """
    {
      "id": "${profile_id}",
      "first_name": "${first_name}",
      "last_name": "${last_name}",
      "phone": "${phone}"
    }
    """
