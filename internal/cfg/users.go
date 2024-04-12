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
	// Skip verifying the certificate TODO implement
	SkipCertificateVerify bool   `json:"skip_certificate_verify"`
	// LDAP root search query
	GroupSearchRoot       string `json:"group_search_root"`
	// LDAP group filter template string. Use {{ .Username }} or {{ .Email }} to fill in username or email.
	GroupFilter           string `json:"group_filter_template"`
	// LDAP attribute to use as the group identifier. Used in group links to link local groups to ldap groups
	GroupAttribute        string `json:"group_attribute"`
	// Automatically add groups when users log in
	AutoAddGroups         bool   `json:"audo_add_groups"`
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

		multiProvider.Add(usr.LDAP, usr.LDAPUserProvider{
			ServerURL: u.ServerURL,
			SkipCertificateVerify: u.SkipCertificateVerify,
			GroupSearchRoot: u.GroupSearchRoot,
			GroupFilter: groupFilterTemplate,
			GroupAttribute: u.GroupAttribute,
			AutoAddGroups: u.AutoAddGroups,
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
