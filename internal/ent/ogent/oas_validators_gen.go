// Code generated by ogen, DO NOT EDIT.

package ogent

import (
	"fmt"

	"github.com/go-faster/errors"

	"github.com/ogen-go/ogen/validate"
)

func (s ListClaimClaimGroupsOKApplicationJSON) Validate() error {
	alias := ([]ClaimClaimGroupsList)(s)
	if alias == nil {
		return errors.New("nil is invalid value")
	}
	return nil
}

func (s ListClaimGroupClaimsOKApplicationJSON) Validate() error {
	alias := ([]ClaimGroupClaimsList)(s)
	if alias == nil {
		return errors.New("nil is invalid value")
	}
	return nil
}

func (s ListClaimGroupGroupLinksOKApplicationJSON) Validate() error {
	alias := ([]ClaimGroupGroupLinksList)(s)
	if alias == nil {
		return errors.New("nil is invalid value")
	}
	return nil
}

func (s ListClaimGroupOKApplicationJSON) Validate() error {
	alias := ([]ClaimGroupList)(s)
	if alias == nil {
		return errors.New("nil is invalid value")
	}
	return nil
}

func (s ListClaimGroupUsersOKApplicationJSON) Validate() error {
	alias := ([]ClaimGroupUsersList)(s)
	if alias == nil {
		return errors.New("nil is invalid value")
	}
	return nil
}

func (s ListClaimOKApplicationJSON) Validate() error {
	alias := ([]ClaimList)(s)
	if alias == nil {
		return errors.New("nil is invalid value")
	}
	return nil
}

func (s ListGroupLinkOKApplicationJSON) Validate() error {
	alias := ([]GroupLinkList)(s)
	if alias == nil {
		return errors.New("nil is invalid value")
	}
	return nil
}

func (s ListPrivateKeyOKApplicationJSON) Validate() error {
	alias := ([]PrivateKeyList)(s)
	if alias == nil {
		return errors.New("nil is invalid value")
	}
	return nil
}

func (s ListUserClaimGroupsOKApplicationJSON) Validate() error {
	alias := ([]UserClaimGroupsList)(s)
	if alias == nil {
		return errors.New("nil is invalid value")
	}
	return nil
}

func (s ListUserOKApplicationJSON) Validate() error {
	alias := ([]UserList)(s)
	if alias == nil {
		return errors.New("nil is invalid value")
	}
	return nil
}

func (s *PkeysOK) Validate() error {
	if s == nil {
		return validate.ErrNilPointer
	}

	var failures []validate.FieldError
	if err := func() error {
		var failures []validate.FieldError
		for i, elem := range s.Keys {
			if err := func() error {
				if err := elem.Validate(); err != nil {
					return err
				}
				return nil
			}(); err != nil {
				failures = append(failures, validate.FieldError{
					Name:  fmt.Sprintf("[%d]", i),
					Error: err,
				})
			}
		}
		if len(failures) > 0 {
			return &validate.Error{Fields: failures}
		}
		return nil
	}(); err != nil {
		failures = append(failures, validate.FieldError{
			Name:  "keys",
			Error: err,
		})
	}
	if len(failures) > 0 {
		return &validate.Error{Fields: failures}
	}
	return nil
}

func (s *PkeysOKKeysItem) Validate() error {
	if s == nil {
		return validate.ErrNilPointer
	}

	var failures []validate.FieldError
	if err := func() error {
		if value, ok := s.Kty.Get(); ok {
			if err := func() error {
				if err := value.Validate(); err != nil {
					return err
				}
				return nil
			}(); err != nil {
				return err
			}
		}
		return nil
	}(); err != nil {
		failures = append(failures, validate.FieldError{
			Name:  "kty",
			Error: err,
		})
	}
	if err := func() error {
		if value, ok := s.Crv.Get(); ok {
			if err := func() error {
				if err := value.Validate(); err != nil {
					return err
				}
				return nil
			}(); err != nil {
				return err
			}
		}
		return nil
	}(); err != nil {
		failures = append(failures, validate.FieldError{
			Name:  "crv",
			Error: err,
		})
	}
	if len(failures) > 0 {
		return &validate.Error{Fields: failures}
	}
	return nil
}

func (s PkeysOKKeysItemCrv) Validate() error {
	switch s {
	case "P-256":
		return nil
	case "P-384":
		return nil
	case "P-521":
		return nil
	default:
		return errors.Errorf("invalid value: %v", s)
	}
}

func (s PkeysOKKeysItemKty) Validate() error {
	switch s {
	case "EC":
		return nil
	case "RSA":
		return nil
	case "OKP":
		return nil
	default:
		return errors.Errorf("invalid value: %v", s)
	}
}
