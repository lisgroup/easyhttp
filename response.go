package easyhttp

// Response struct
type Response struct {
	Content  string  // Http response content
	HttpCode int     // http code
	Status   string  // http status
	Cost     float32 // http cost Microsecond
}
