package cfg

type Server struct {
	// Address to serve the application on
	Address string `json:"address"`
	// Port to serve the application on
	Port    int    `json:"port"`
}
