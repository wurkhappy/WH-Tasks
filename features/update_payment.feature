#
Feature: Update task
	In order for users to add sub items to the task it must be updated.
	Sub items help communicate to the client what they are paying for.
	They also help keep track of what work has been completed.

	Scenario: Update
		Given an existing set of tasks
		When I update the task items
		Then a task is returned with the same sub items
