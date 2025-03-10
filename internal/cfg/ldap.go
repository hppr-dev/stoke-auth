package cfg

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"io"
	"os"
	"stoke/internal/usr"
	"strings"
	"text/template"

	"github.com/go-ldap/ldap/v3"
	"github.com/rs/zerolog"
)

type LDAPProviderConfig struct {
	// Name of this ldap provider
	Name                  string
	// URL (starting with ldap:// or ldaps://) of the ldap server
	ServerURL             string `json:"server_url"`
	// Readonly bind user distinguished name used to look up users in ldap
	BindUserDN            string `json:"bind_user_dn"`
	// Read-only bind user password
	BindUserPassword      string `json:"bind_user_password"`


	// LDAP root group search query
	GroupSearchRoot       string `json:"group_search_root"`
	// LDAP group filter template string. Use {{ .Username }} or {{ .Email }} to fill in username or email.
	GroupFilter           string `json:"group_filter_template"`
	// LDAP group name attribute used to match groups
	GroupNameField        string `json:"ldap_group_name_field"`

	// LDAP root user search query
	UserSearchRoot        string `json:"user_search_root"`
	// LDAP user filter template string. Use {{ .Username }} to fill in the the username
	UserFilter            string `json:"user_filter_template"`

	// LDAP field to pull from ldap as the first name
	FirstNameField        string `json:"ldap_first_name_field"`
	// LDAP field to pull from ldap as the last name
	LastNameField         string `json:"ldap_last_name_field"`
	// LDAP field to pull from ldap as the email
	EmailField            string `json:"ldap_email_field"`

	// Timeout for LDAP searches
	SearchTimeout         int    `json:"search_timeout"`

	// LDAP public ca certificate file
	LDAPCACert            string `json:"ldap_ca_cert"`
	// Skip verifying the certificate
	SkipCertificateVerify bool   `json:"skip_certificate_verify"`
}

func (l LDAPProviderConfig) TypeSpec() string {
	return "LDAP:" + l.Name
}

func (l LDAPProviderConfig) CreateProvider(ctx context.Context) foreignProvider {
	logger := zerolog.Ctx(ctx).With().
		Str("component", "cfg.LdapProviderConfig.CreateProvider").
		Logger()

	groupFilterTemplate := template.New("group-filter")
	groupFilterTemplate, err := groupFilterTemplate.Parse(l.GroupFilter)
	if err != nil {
		logger.Fatal().
			Err(err).
			Str("groupFilterTemplate", l.GroupFilter).
			Msg("Could not parse group filter template")
	}

	userFilterTemplate := template.New("user-filter")
	userFilterTemplate, err = userFilterTemplate.Parse(l.UserFilter)
	if err != nil {
		logger.Fatal().
			Err(err).
			Str("userFilterTemplate", l.UserFilter).
			Msg("Could not parse user filter template")
	}

	var dialOpts []ldap.DialOpt

	if strings.HasPrefix(l.ServerURL, "ldaps://") {
		certPool, err := x509.SystemCertPool()
		if err != nil {
			logger.Error().
				Err(err).
				Str("serverUrl", l.ServerURL).
				Str("ldapCert", l.LDAPCACert).
				Msg("Could not load system cert pool. Using new empty pool.")
			certPool = x509.NewCertPool()
		}
		if l.LDAPCACert != "" {
			publicCerts, err := readPublicCertFile(l.LDAPCACert)
			if err != nil {
				logger.Fatal().
					Err(err).
					Str("ldapCert", l.LDAPCACert).
					Msg("Could not read ldap cert file")
			}
			for _, cert := range publicCerts {
				certPool.AddCert(cert)
			}
		}
		dialOpts = append(dialOpts,
			ldap.DialWithTLSConfig(
				&tls.Config{
					ClientCAs: certPool,
					InsecureSkipVerify: l.SkipCertificateVerify,
				},
			),
		)
	}

	return usr.NewLDAPUserProvider(
		l.Name,
		l.ServerURL,
		l.BindUserDN,
		l.BindUserPassword,
		l.GroupSearchRoot,
		l.GroupNameField,
		l.UserSearchRoot,
		l.FirstNameField,
		l.LastNameField,
		l.EmailField,
		l.SearchTimeout,
		groupFilterTemplate,
		userFilterTemplate,
		dialOpts...,
	)

}

func readPublicCertFile(name string) ([]*x509.Certificate, error) {
	certFile, err := os.Open(name)
	if err != nil {
		return nil, err
	}

	certBytes, err := io.ReadAll(certFile)
	if err != nil {
		return nil, err
	}

	return x509.ParseCertificates(certBytes)
}
