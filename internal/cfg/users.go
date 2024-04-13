package cfg

import (
	"context"
	"stoke/internal/usr"
	"text/template"

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

	// Skip verifying the certificate TODO implement
	SkipCertificateVerify bool   `json:"skip_certificate_verify"`
}

func (u Users) withContext(ctx context.Context) context.Context {
	logger := zerolog.Ctx(ctx)
	multiProvider := &usr.MultiProvider{}

	// Will always have local user database
	multiProvider.Add(usr.LOCAL, usr.LocalProvider{})

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

		multiProvider.Add(usr.LDAP, usr.LDAPUserProvider{
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
			SkipCertificateVerify: u.SkipCertificateVerify,
		})
	}

	err := multiProvider.Init(ctx)
	if err != nil {
		logger.Fatal().
			Err(err).
			Msg("Could not initialize user providers")
	}

	return context.WithValue(ctx, "user-provider", multiProvider)
}
