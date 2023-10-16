package mattermost

// LoggerAPI isolates the logger methods from the plugin API
type LoggerAPI interface {
	LogDebug(msg string, keyValuePairs ...any)
	LogInfo(msg string, keyValuePairs ...any)
	LogError(msg string, keyValuePairs ...any)
	LogWarn(msg string, keyValuePairs ...any)
}
