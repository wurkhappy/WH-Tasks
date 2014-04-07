#
Feature: Update last action
	In order for to know where an agreement stands in terms of progress
	we need to keep track of the most current actions taken. In order to keep track we must be able to update the last action of a task.

	Scenario: Update
		Given an existing set of tasks
		When I take one of them
		When I create a new action for that task
		Then an action is returned
