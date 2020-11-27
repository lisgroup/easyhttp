package easyhttp

// Options struct
type Options struct {
	BaseURI  string                 // request BaseURI
	Timeout  float32                // timeout (second)
	Headers  map[string]interface{} // request Headers
	BodyMaps map[string]interface{} // request body
	Query    interface{}            // request query
	JSON     interface{}            // request JSON
}
