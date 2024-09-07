Feature: Add Todo
  Scenario: Ensure todo is added
    Given: I open website "http://api:300"
    When I fill form in following information:
      | Field | Value |
      | Taks | New Todo Item |
    And I click sumbit button
    Then Verify result information "New Todo Item"
