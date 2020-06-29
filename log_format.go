package main

// LogFormat ...
type LogFormat interface {
	Parse(raw string) LogEvent
}
