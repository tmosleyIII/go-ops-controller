package cmd

// Error represents information regarding errors and status
type Error struct {
	Error  string `json:"error"`
	Status int    `json:"status"`
}

var codes = map[int]string{
	400: "Bad Request",
	404: "Not Found",
	405: "Invalid Input",
}
