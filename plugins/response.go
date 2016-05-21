package plugins

type Response struct {
	Data         map[string]interface{}
	Headers      map[string][]string
	ResponseCode int
}
