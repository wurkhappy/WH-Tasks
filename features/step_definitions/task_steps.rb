#encoding: utf-8
require "rubygems"
require "majordomo"
require "json"
require "securerandom"

Before do
	@client = Majordomo::Client.new("tcp://localhost:5555", false)
end

After do
	@client.close
end

Given /^a new set of tasks$/ do
	@versionID = SecureRandom.uuid
	@tasks = [{:title=>"1", :versionID => @versionID}, {:title=>"2", :versionID => @versionID}]
end

Given /^an existing set of tasks$/ do
	@versionID = SecureRandom.uuid
	tasks = [{:title=>"1", :versionID => @versionID}, {:title=>"2", :versionID => @versionID}]
	jsonString = [JSON.generate(tasks)].pack('m0')
	reply = @client.send_and_receive("Tasks", '{"Method":"POST", "Path":"/agreements/v/'+@versionID+'/tasks", "Body":"'+jsonString+'"}')
	@savedTasks = JSON.parse(JSON.parse(reply[0])["body"].unpack('m0')[0])
end

When /^I save them$/ do
	jsonString = [JSON.generate(@tasks)].pack('m0')
	reply = @client.send_and_receive("Tasks", '{"Method":"POST", "Path":"/agreements/v/'+@versionID+'/tasks", "Body":"'+jsonString+'"}')
	@savedTasks = JSON.parse(JSON.parse(reply[0])["body"].unpack('m0')[0])
end

When /^I fetch them by the agreements version ID$/ do
	reply = @client.send_and_receive("Tasks", '{"Method":"GET", "Path":"/agreements/v/'+ @versionID+'/tasks"}')
	@savedTasks = JSON.parse(JSON.parse(reply[0])["body"].unpack('m0')[0])
end

When /^I update the task items$/ do
	@updatedTask = @savedTasks[0]
	@updatedTask["subTasks"] = [{:title => "test1"}]
	jsonString = [JSON.generate(@updatedTask)].pack('m0')
	reply = @client.send_and_receive("Tasks", '{"Method":"PUT", "Path":"/tasks/'+@updatedTask["id"]+'", "Body":"'+jsonString+'"}')
	@savedTask = JSON.parse(JSON.parse(reply[0])["body"].unpack('m0')[0])
end

When /^I take one of them$/ do
	@savedTask = @savedTasks[0]
end

When /^I create a new action for that task$/ do
	jsonString = [JSON.generate({:name => "completed"})].pack('m0')
	reply = @client.send_and_receive("Tasks", '{"Method":"POST", "Path":"/tasks/'+@savedTask["id"]+'/action", "Body":"'+jsonString+'"}')
	@action = JSON.parse(JSON.parse(reply[0])["body"].unpack('m0')[0])
end

Then /^they each have an id$/ do
	@savedTasks.should be_a_kind_of(Array)

	@savedTasks.each {|x| 
		(!x["id"].nil? or (!x["id"].nil? and !x["id"].empty?)).should be_true
	}
end

Then /^at least one is returned$/ do
	@savedTasks.should be_a_kind_of(Array)
	@savedTasks.should have_at_least(1).items
end

Then /^a task is returned with the same sub items$/ do
	@savedTask["subTasks"].should be_a_kind_of(Array)
	@savedTask["subTasks"].length.should eql(@updatedTask["subTasks"].length)
end

Then /^an action is returned$/ do
	@action.should be_a_kind_of(Hash)
	@action["name"].should_not be_nil
	@action["name"].should_not be_empty
end