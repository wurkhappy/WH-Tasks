#
Feature: Create tasks
	In order for users to create an agreement there need to be terms for task.
	Creating tasks allows freelancers to communicate with their about what the expectations
	are around when they will get paid, how much and for what services.

	Scenario: Create and Save
		Given a new set of tasks
		When I save them
		Then they each have an id
