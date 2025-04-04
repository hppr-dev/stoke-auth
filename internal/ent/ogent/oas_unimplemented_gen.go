// Code generated by ogen, DO NOT EDIT.

package ogent

import (
	"context"

	ht "github.com/ogen-go/ogen/http"
)

// UnimplementedHandler is no-op Handler which returns http.ErrNotImplemented.
type UnimplementedHandler struct{}

var _ Handler = UnimplementedHandler{}

// AvailableProviders implements available_providers operation.
//
// Get available providers.
//
// GET /available_providers
func (UnimplementedHandler) AvailableProviders(ctx context.Context) (r []AvailableProvidersOKItem, _ error) {
	return r, ht.ErrNotImplemented
}

// Capabilities implements capabilities operation.
//
// Get server capabilities.
//
// GET /capabilities
func (UnimplementedHandler) Capabilities(ctx context.Context) (r *CapabilitiesOK, _ error) {
	return r, ht.ErrNotImplemented
}

// CreateClaim implements createClaim operation.
//
// Creates a new Claim and persists it to storage.
//
// POST /admin/claims
func (UnimplementedHandler) CreateClaim(ctx context.Context, req *CreateClaimReq) (r CreateClaimRes, _ error) {
	return r, ht.ErrNotImplemented
}

// CreateClaimGroup implements createClaimGroup operation.
//
// Creates a new ClaimGroup and persists it to storage.
//
// POST /admin/claim-groups
func (UnimplementedHandler) CreateClaimGroup(ctx context.Context, req *CreateClaimGroupReq) (r CreateClaimGroupRes, _ error) {
	return r, ht.ErrNotImplemented
}

// CreateGroupLink implements createGroupLink operation.
//
// Creates a new GroupLink and persists it to storage.
//
// POST /admin/group-links
func (UnimplementedHandler) CreateGroupLink(ctx context.Context, req *CreateGroupLinkReq) (r CreateGroupLinkRes, _ error) {
	return r, ht.ErrNotImplemented
}

// CreateLocalUser implements createLocalUser operation.
//
// Create a new local user.
//
// POST /admin/localuser
func (UnimplementedHandler) CreateLocalUser(ctx context.Context, req *CreateLocalUserReq) (r CreateLocalUserRes, _ error) {
	return r, ht.ErrNotImplemented
}

// DeleteClaim implements deleteClaim operation.
//
// Deletes the Claim with the requested ID.
//
// DELETE /admin/claims/{id}
func (UnimplementedHandler) DeleteClaim(ctx context.Context, params DeleteClaimParams) (r DeleteClaimRes, _ error) {
	return r, ht.ErrNotImplemented
}

// DeleteClaimGroup implements deleteClaimGroup operation.
//
// Deletes the ClaimGroup with the requested ID.
//
// DELETE /admin/claim-groups/{id}
func (UnimplementedHandler) DeleteClaimGroup(ctx context.Context, params DeleteClaimGroupParams) (r DeleteClaimGroupRes, _ error) {
	return r, ht.ErrNotImplemented
}

// DeleteGroupLink implements deleteGroupLink operation.
//
// Deletes the GroupLink with the requested ID.
//
// DELETE /admin/group-links/{id}
func (UnimplementedHandler) DeleteGroupLink(ctx context.Context, params DeleteGroupLinkParams) (r DeleteGroupLinkRes, _ error) {
	return r, ht.ErrNotImplemented
}

// DeleteUser implements deleteUser operation.
//
// Deletes the User with the requested ID.
//
// DELETE /admin/users/{id}
func (UnimplementedHandler) DeleteUser(ctx context.Context, params DeleteUserParams) (r DeleteUserRes, _ error) {
	return r, ht.ErrNotImplemented
}

// ListClaim implements listClaim operation.
//
// List Claims.
//
// GET /admin/claims
func (UnimplementedHandler) ListClaim(ctx context.Context, params ListClaimParams) (r ListClaimRes, _ error) {
	return r, ht.ErrNotImplemented
}

// ListClaimClaimGroups implements listClaimClaimGroups operation.
//
// List attached ClaimGroups.
//
// GET /admin/claims/{id}/claim-groups
func (UnimplementedHandler) ListClaimClaimGroups(ctx context.Context, params ListClaimClaimGroupsParams) (r ListClaimClaimGroupsRes, _ error) {
	return r, ht.ErrNotImplemented
}

// ListClaimGroup implements listClaimGroup operation.
//
// List ClaimGroups.
//
// GET /admin/claim-groups
func (UnimplementedHandler) ListClaimGroup(ctx context.Context, params ListClaimGroupParams) (r ListClaimGroupRes, _ error) {
	return r, ht.ErrNotImplemented
}

// ListClaimGroupClaims implements listClaimGroupClaims operation.
//
// List attached Claims.
//
// GET /admin/claim-groups/{id}/claims
func (UnimplementedHandler) ListClaimGroupClaims(ctx context.Context, params ListClaimGroupClaimsParams) (r ListClaimGroupClaimsRes, _ error) {
	return r, ht.ErrNotImplemented
}

// ListClaimGroupGroupLinks implements listClaimGroupGroupLinks operation.
//
// List attached GroupLinks.
//
// GET /admin/claim-groups/{id}/group-links
func (UnimplementedHandler) ListClaimGroupGroupLinks(ctx context.Context, params ListClaimGroupGroupLinksParams) (r ListClaimGroupGroupLinksRes, _ error) {
	return r, ht.ErrNotImplemented
}

// ListClaimGroupUsers implements listClaimGroupUsers operation.
//
// List attached Users.
//
// GET /admin/claim-groups/{id}/users
func (UnimplementedHandler) ListClaimGroupUsers(ctx context.Context, params ListClaimGroupUsersParams) (r ListClaimGroupUsersRes, _ error) {
	return r, ht.ErrNotImplemented
}

// ListGroupLink implements listGroupLink operation.
//
// List GroupLinks.
//
// GET /admin/group-links
func (UnimplementedHandler) ListGroupLink(ctx context.Context, params ListGroupLinkParams) (r ListGroupLinkRes, _ error) {
	return r, ht.ErrNotImplemented
}

// ListPrivateKey implements listPrivateKey operation.
//
// List PrivateKeys.
//
// GET /admin/private-keys
func (UnimplementedHandler) ListPrivateKey(ctx context.Context, params ListPrivateKeyParams) (r ListPrivateKeyRes, _ error) {
	return r, ht.ErrNotImplemented
}

// ListUser implements listUser operation.
//
// List Users.
//
// GET /admin/users
func (UnimplementedHandler) ListUser(ctx context.Context, params ListUserParams) (r ListUserRes, _ error) {
	return r, ht.ErrNotImplemented
}

// ListUserClaimGroups implements listUserClaimGroups operation.
//
// List attached ClaimGroups.
//
// GET /admin/users/{id}/claim-groups
func (UnimplementedHandler) ListUserClaimGroups(ctx context.Context, params ListUserClaimGroupsParams) (r ListUserClaimGroupsRes, _ error) {
	return r, ht.ErrNotImplemented
}

// Login implements login operation.
//
// Request a token.
//
// POST /login
func (UnimplementedHandler) Login(ctx context.Context, req *LoginReq) (r LoginRes, _ error) {
	return r, ht.ErrNotImplemented
}

// Pkeys implements pkeys operation.
//
// Get current valid public keys.
//
// GET /pkeys
func (UnimplementedHandler) Pkeys(ctx context.Context) (r *PkeysOK, _ error) {
	return r, ht.ErrNotImplemented
}

// ReadClaim implements readClaim operation.
//
// Finds the Claim with the requested ID and returns it.
//
// GET /admin/claims/{id}
func (UnimplementedHandler) ReadClaim(ctx context.Context, params ReadClaimParams) (r ReadClaimRes, _ error) {
	return r, ht.ErrNotImplemented
}

// ReadClaimGroup implements readClaimGroup operation.
//
// Finds the ClaimGroup with the requested ID and returns it.
//
// GET /admin/claim-groups/{id}
func (UnimplementedHandler) ReadClaimGroup(ctx context.Context, params ReadClaimGroupParams) (r ReadClaimGroupRes, _ error) {
	return r, ht.ErrNotImplemented
}

// ReadGroupLink implements readGroupLink operation.
//
// Finds the GroupLink with the requested ID and returns it.
//
// GET /admin/group-links/{id}
func (UnimplementedHandler) ReadGroupLink(ctx context.Context, params ReadGroupLinkParams) (r ReadGroupLinkRes, _ error) {
	return r, ht.ErrNotImplemented
}

// ReadGroupLinkClaimGroup implements readGroupLinkClaimGroup operation.
//
// Find the attached ClaimGroup of the GroupLink with the given ID.
//
// GET /admin/group-links/{id}/claim-group
func (UnimplementedHandler) ReadGroupLinkClaimGroup(ctx context.Context, params ReadGroupLinkClaimGroupParams) (r ReadGroupLinkClaimGroupRes, _ error) {
	return r, ht.ErrNotImplemented
}

// ReadPrivateKey implements readPrivateKey operation.
//
// Finds the PrivateKey with the requested ID and returns it.
//
// GET /admin/private-keys/{id}
func (UnimplementedHandler) ReadPrivateKey(ctx context.Context, params ReadPrivateKeyParams) (r ReadPrivateKeyRes, _ error) {
	return r, ht.ErrNotImplemented
}

// ReadUser implements readUser operation.
//
// Finds the User with the requested ID and returns it.
//
// GET /admin/users/{id}
func (UnimplementedHandler) ReadUser(ctx context.Context, params ReadUserParams) (r ReadUserRes, _ error) {
	return r, ht.ErrNotImplemented
}

// Refresh implements refresh operation.
//
// Request a refreshed token.
//
// POST /refresh
func (UnimplementedHandler) Refresh(ctx context.Context, req *RefreshReq) (r RefreshRes, _ error) {
	return r, ht.ErrNotImplemented
}

// Totals implements totals operation.
//
// Get entity count totals.
//
// GET /admin/totals
func (UnimplementedHandler) Totals(ctx context.Context) (r *TotalsOK, _ error) {
	return r, ht.ErrNotImplemented
}

// UpdateClaim implements updateClaim operation.
//
// Updates a Claim and persists changes to storage.
//
// PATCH /admin/claims/{id}
func (UnimplementedHandler) UpdateClaim(ctx context.Context, req *UpdateClaimReq, params UpdateClaimParams) (r UpdateClaimRes, _ error) {
	return r, ht.ErrNotImplemented
}

// UpdateClaimGroup implements updateClaimGroup operation.
//
// Updates a ClaimGroup and persists changes to storage.
//
// PATCH /admin/claim-groups/{id}
func (UnimplementedHandler) UpdateClaimGroup(ctx context.Context, req *UpdateClaimGroupReq, params UpdateClaimGroupParams) (r UpdateClaimGroupRes, _ error) {
	return r, ht.ErrNotImplemented
}

// UpdateGroupLink implements updateGroupLink operation.
//
// Updates a GroupLink and persists changes to storage.
//
// PATCH /admin/group-links/{id}
func (UnimplementedHandler) UpdateGroupLink(ctx context.Context, req *UpdateGroupLinkReq, params UpdateGroupLinkParams) (r UpdateGroupLinkRes, _ error) {
	return r, ht.ErrNotImplemented
}

// UpdateLocalUserPassword implements updateLocalUserPassword operation.
//
// Update local user's password.
//
// PATCH /admin/localuser
func (UnimplementedHandler) UpdateLocalUserPassword(ctx context.Context, req *UpdateLocalUserPasswordReq) (r UpdateLocalUserPasswordRes, _ error) {
	return r, ht.ErrNotImplemented
}

// UpdateUser implements updateUser operation.
//
// Updates a User and persists changes to storage.
//
// PATCH /admin/users/{id}
func (UnimplementedHandler) UpdateUser(ctx context.Context, req *UpdateUserReq, params UpdateUserParams) (r UpdateUserRes, _ error) {
	return r, ht.ErrNotImplemented
}
