package cfg

type Server struct {
	// Address to serve the application on
	Address string `json:"address"`
	// Port to serve the application on
	Port    int    `json:"port"`
	// Request Timeout in milliseconds
	Timeout int    `json:"timeout"`

	// Private TLS Key
	TLSPrivateKey string `json:"tls_private_key"`
	// Public TLS Cert
	TLSPublicCert string `json:"tls_public_cert"`

	// Allowed Hosts
	AllowedHosts []string `json:"allowed_hosts"`
}
