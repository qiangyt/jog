package main

// LogParser ...
type LogParser interface {

	// Parse method parse and populate fields into the event. It should return amount of matched fields
	Parse(event LogEvent) int
}
