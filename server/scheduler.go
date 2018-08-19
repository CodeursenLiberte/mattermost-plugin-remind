package main

import (
	"github.com/google/uuid"
	"time"
)

func (p *Plugin) Run() {

	if !p.running {
		p.running = true
		p.runner()
	}
}

func (p *Plugin) ScheduleReminder(request ReminderRequest) (string, error) {

	var when string
	var target string
	var message string
	var useTo bool
	useTo = false
	var useToString string
	if useTo {
		useToString = " to"
	} else {
		useToString = ""
	}

	guid, gErr := uuid.NewRandom()
	if gErr != nil {
		p.API.LogError("Failed to generate guid")
	}

	target, when, message, pErr := p.ParseRequest(request)
	if pErr != nil {
		return ExceptionText, nil
	}

	request.Reminder.TeamId = request.TeamId
	request.Reminder.Id = guid.String()
	request.Reminder.Username = request.Username
	request.Reminder.Target = target
	request.Reminder.Message = message
	request.Reminder.When = when
	request.Reminder.Occurrences = p.CreateOccurrences(request)

	//// TODO REMOVE THIS LATER
	//p.API.KVDelete(request.Username)
	//////////////

	p.UpsertReminder(request)

	response := ":thumbsup: I will remind " + target + useToString + " \"" + request.Reminder.Message + "\" " + when;
	return response, nil
}

func (p *Plugin) stop() {
	p.running = false
}

func (p *Plugin) runner() {

	go func() {
		<-time.NewTimer(time.Second).C
		p.TriggerReminders()
		if !p.running {
			return
		}
		p.runner()
	}()
}