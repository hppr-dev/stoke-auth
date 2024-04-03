// Code generated by ogen, DO NOT EDIT.

package ogent

import (
	"context"
)

// Handler handles operations described by OpenAPI v3 specification.
type Handler interface {
	// CreateClaim implements createClaim operation.
	//
	// Creates a new Claim and persists it to storage.
	//
	// POST /claims
	CreateClaim(ctx context.Context, req *CreateClaimReq) (CreateClaimRes, error)
	// CreateClaimGroup implements createClaimGroup operation.
	//
	// Creates a new ClaimGroup and persists it to storage.
	//
	// POST /claim-groups
	CreateClaimGroup(ctx context.Context, req *CreateClaimGroupReq) (CreateClaimGroupRes, error)
	// CreateGroupLink implements createGroupLink operation.
	//
	// Creates a new GroupLink and persists it to storage.
	//
	// POST /group-links
	CreateGroupLink(ctx context.Context, req *CreateGroupLinkReq) (CreateGroupLinkRes, error)
	// DeleteClaim implements deleteClaim operation.
	//
	// Deletes the Claim with the requested ID.
	//
	// DELETE /claims/{id}
	DeleteClaim(ctx context.Context, params DeleteClaimParams) (DeleteClaimRes, error)
	// DeleteClaimGroup implements deleteClaimGroup operation.
	//
	// Deletes the ClaimGroup with the requested ID.
	//
	// DELETE /claim-groups/{id}
	DeleteClaimGroup(ctx context.Context, params DeleteClaimGroupParams) (DeleteClaimGroupRes, error)
	// DeleteGroupLink implements deleteGroupLink operation.
	//
	// Deletes the GroupLink with the requested ID.
	//
	// DELETE /group-links/{id}
	DeleteGroupLink(ctx context.Context, params DeleteGroupLinkParams) (DeleteGroupLinkRes, error)
	// ListClaim implements listClaim operation.
	//
	// List Claims.
	//
	// GET /claims
	ListClaim(ctx context.Context, params ListClaimParams) (ListClaimRes, error)
	// ListClaimClaimGroups implements listClaimClaimGroups operation.
	//
	// List attached ClaimGroups.
	//
	// GET /claims/{id}/claim-groups
	ListClaimClaimGroups(ctx context.Context, params ListClaimClaimGroupsParams) (ListClaimClaimGroupsRes, error)
	// ListClaimGroup implements listClaimGroup operation.
	//
	// List ClaimGroups.
	//
	// GET /claim-groups
	ListClaimGroup(ctx context.Context, params ListClaimGroupParams) (ListClaimGroupRes, error)
	// ListClaimGroupClaims implements listClaimGroupClaims operation.
	//
	// List attached Claims.
	//
	// GET /claim-groups/{id}/claims
	ListClaimGroupClaims(ctx context.Context, params ListClaimGroupClaimsParams) (ListClaimGroupClaimsRes, error)
	// ListClaimGroupGroupLinks implements listClaimGroupGroupLinks operation.
	//
	// List attached GroupLinks.
	//
	// GET /claim-groups/{id}/group-links
	ListClaimGroupGroupLinks(ctx context.Context, params ListClaimGroupGroupLinksParams) (ListClaimGroupGroupLinksRes, error)
	// ListClaimGroupUsers implements listClaimGroupUsers operation.
	//
	// List attached Users.
	//
	// GET /claim-groups/{id}/users
	ListClaimGroupUsers(ctx context.Context, params ListClaimGroupUsersParams) (ListClaimGroupUsersRes, error)
	// ListGroupLink implements listGroupLink operation.
	//
	// List GroupLinks.
	//
	// GET /group-links
	ListGroupLink(ctx context.Context, params ListGroupLinkParams) (ListGroupLinkRes, error)
	// ListPrivateKey implements listPrivateKey operation.
	//
	// List PrivateKeys.
	//
	// GET /private-keys
	ListPrivateKey(ctx context.Context, params ListPrivateKeyParams) (ListPrivateKeyRes, error)
	// ListUser implements listUser operation.
	//
	// List Users.
	//
	// GET /users
	ListUser(ctx context.Context, params ListUserParams) (ListUserRes, error)
	// ListUserClaimGroups implements listUserClaimGroups operation.
	//
	// List attached ClaimGroups.
	//
	// GET /users/{id}/claim-groups
	ListUserClaimGroups(ctx context.Context, params ListUserClaimGroupsParams) (ListUserClaimGroupsRes, error)
	// ReadClaim implements readClaim operation.
	//
	// Finds the Claim with the requested ID and returns it.
	//
	// GET /claims/{id}
	ReadClaim(ctx context.Context, params ReadClaimParams) (ReadClaimRes, error)
	// ReadClaimGroup implements readClaimGroup operation.
	//
	// Finds the ClaimGroup with the requested ID and returns it.
	//
	// GET /claim-groups/{id}
	ReadClaimGroup(ctx context.Context, params ReadClaimGroupParams) (ReadClaimGroupRes, error)
	// ReadGroupLink implements readGroupLink operation.
	//
	// Finds the GroupLink with the requested ID and returns it.
	//
	// GET /group-links/{id}
	ReadGroupLink(ctx context.Context, params ReadGroupLinkParams) (ReadGroupLinkRes, error)
	// ReadGroupLinkClaimGroups implements readGroupLinkClaimGroups operation.
	//
	// Find the attached ClaimGroup of the GroupLink with the given ID.
	//
	// GET /group-links/{id}/claim-groups
	ReadGroupLinkClaimGroups(ctx context.Context, params ReadGroupLinkClaimGroupsParams) (ReadGroupLinkClaimGroupsRes, error)
	// ReadPrivateKey implements readPrivateKey operation.
	//
	// Finds the PrivateKey with the requested ID and returns it.
	//
	// GET /private-keys/{id}
	ReadPrivateKey(ctx context.Context, params ReadPrivateKeyParams) (ReadPrivateKeyRes, error)
	// ReadUser implements readUser operation.
	//
	// Finds the User with the requested ID and returns it.
	//
	// GET /users/{id}
	ReadUser(ctx context.Context, params ReadUserParams) (ReadUserRes, error)
	// UpdateClaim implements updateClaim operation.
	//
	// Updates a Claim and persists changes to storage.
	//
	// PATCH /claims/{id}
	UpdateClaim(ctx context.Context, req *UpdateClaimReq, params UpdateClaimParams) (UpdateClaimRes, error)
	// UpdateClaimGroup implements updateClaimGroup operation.
	//
	// Updates a ClaimGroup and persists changes to storage.
	//
	// PATCH /claim-groups/{id}
	UpdateClaimGroup(ctx context.Context, req *UpdateClaimGroupReq, params UpdateClaimGroupParams) (UpdateClaimGroupRes, error)
	// UpdateGroupLink implements updateGroupLink operation.
	//
	// Updates a GroupLink and persists changes to storage.
	//
	// PATCH /group-links/{id}
	UpdateGroupLink(ctx context.Context, req *UpdateGroupLinkReq, params UpdateGroupLinkParams) (UpdateGroupLinkRes, error)
	// UpdateUser implements updateUser operation.
	//
	// Updates a User and persists changes to storage.
	//
	// PATCH /users/{id}
	UpdateUser(ctx context.Context, req *UpdateUserReq, params UpdateUserParams) (UpdateUserRes, error)
}

// Server implements http server based on OpenAPI v3 specification and
// calls Handler to handle requests.
type Server struct {
	h Handler
	baseServer
}

// NewServer creates new Server.
func NewServer(h Handler, opts ...ServerOption) (*Server, error) {
	s, err := newServerConfig(opts...).baseServer()
	if err != nil {
		return nil, err
	}
	return &Server{
		h:          h,
		baseServer: s,
	}, nil
}
