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

type Users struct {
	// Enable authenticating with an LDAP server
	EnableLDAP            bool `json:"enable_ldap"`
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
	LDAPCACert        string `json:"ldap_ca_cert"`
	// Skip verifying the certificate
	SkipCertificateVerify bool   `json:"skip_certificate_verify"`
}

func (u Users) withContext(ctx context.Context) context.Context {
	logger := zerolog.Ctx(ctx)
	var provider usr.Provider

	if u.EnableLDAP {
		groupFilterTemplate := template.New("group-filter")
		groupFilterTemplate, err := groupFilterTemplate.Parse(u.GroupFilter)
		if err != nil {
			logger.Fatal().
				Err(err).
				Str("groupFilterTemplate", u.GroupFilter).
				Msg("Could not parse group filter template")
		}

		userFilterTemplate := template.New("user-filter")
		userFilterTemplate, err = userFilterTemplate.Parse(u.UserFilter)
		if err != nil {
			logger.Fatal().
				Err(err).
				Str("userFilterTemplate", u.UserFilter).
				Msg("Could not parse user filter template")
		}

		var dialOpts []ldap.DialOpt

		if strings.HasPrefix(u.ServerURL, "ldaps://") {
			certPool, err := x509.SystemCertPool()
			if err != nil {
				logger.Error().
					Err(err).
					Str("serverUrl", u.ServerURL).
					Str("ldapCert", u.LDAPCACert).
					Msg("Could not load system cert pool. Using new empty pool.")
				certPool = x509.NewCertPool()
			}
			if u.LDAPCACert != "" {
				publicCerts, err := readPublicCertFile(u.LDAPCACert)
				if err != nil {
					logger.Fatal().
						Err(err).
						Str("ldapCert", u.LDAPCACert).
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
						InsecureSkipVerify: u.SkipCertificateVerify,
					},
				),
			)
		}

		provider = usr.LDAPUserProvider{
			ServerURL: u.ServerURL,
			BindUserDN: u.BindUserDN,
			BindUserPassword: u.BindUserPassword,

			GroupSearchRoot: u.GroupSearchRoot,
			GroupFilter: groupFilterTemplate,
			GroupAttribute: u.GroupNameField,

			UserSearchRoot: u.UserSearchRoot,
			UserFilter: userFilterTemplate,

			FirstNameField: u.FirstNameField,
			LastNameField: u.LastNameField,
			EmailField: u.EmailField,

			SearchTimeout: u.SearchTimeout,
			DialOpts: dialOpts,
		}
	} else {
		provider = usr.LocalProvider{}
	}

	if err := provider.CheckCreateForSuperUser(ctx); err != nil {
		logger.Error().Err(err).Msg("Could not check/create superuser.")
	}

	return context.WithValue(ctx, "user-provider", provider)
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
