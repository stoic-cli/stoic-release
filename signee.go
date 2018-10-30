package release

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

// Signee provides the available operations
// for interacting with a Signee
type Signee interface {
	// PublicKey fetches the public key
	// of a signee
	PublicKey() ([]byte, error)

	// Key returns the key associated with
	// the signee
	Key() string

	// User returns the user associated with
	// the signee
	User() string

	// Type returns the type of signee
	Type() SigneeType
}

// DefaultGithubAPIEndpoint can be overridden if a different access
// point is required
var DefaultGithubAPIEndpoint = "https://api.github.com"

// DefaultKeybaseEndpoint can be overridden if a different access
// point is required
var DefaultKeybaseEndpoint = "https://keybase.io/"

// SigneeType enumerates the available
// default sources for fetching a signee
type SigneeType string

const (
	// GithubSigneeType will interact with the
	// keys of a github user
	GithubSigneeType SigneeType = "github"

	// KeybaseSigneeType will interact with the
	// keys of a keybase user
	KeybaseSigneeType SigneeType = "keybase"
)

type signee struct {
	user       string
	key        string
	signeeType SigneeType
}

// NewSignee can fetch a public key of a signee for verification purposes
func NewSignee(user string, key string, signeeType SigneeType) Signee {
	return &signee{
		user:       user,
		key:        key,
		signeeType: signeeType,
	}
}

// PublicKey returns the public key of the signee
// or an error if it isn't able to fetch it
func (s *signee) PublicKey() ([]byte, error) {
	switch s.signeeType {
	case GithubSigneeType:
		return s.githubPublicKey()
	case KeybaseSigneeType:
		return s.keybasePublicKey()
	default:
		return nil, fmt.Errorf("unknown signee type: %s", s.signeeType)
	}
}

// GPGKey unmarshals a github response for gpg keys
// Proudly stolen from: https://github.com/google/go-github
// They currently don't support the `RawKey` field.
type GPGKey struct {
	ID                int64      `json:"id,omitempty"`
	PrimaryKeyID      int64      `json:"primary_key_id,omitempty"`
	KeyID             string     `json:"key_id,omitempty"`
	PublicKey         string     `json:"public_key,omitempty"`
	RawKey            string     `json:"raw_key,omitempty"`
	Emails            []GPGEmail `json:"emails,omitempty"`
	Subkeys           []GPGKey   `json:"subkeys,omitempty"`
	CanSign           bool       `json:"can_sign,omitempty"`
	CanEncryptComms   bool       `json:"can_encrypt_comms,omitempty"`
	CanEncryptStorage bool       `json:"can_encrypt_storage,omitempty"`
	CanCertify        bool       `json:"can_certify,omitempty"`
	CreatedAt         time.Time  `json:"created_at,omitempty"`
	ExpiresAt         time.Time  `json:"expires_at,omitempty"`
}

// GPGEmail represents the gpg email section
type GPGEmail struct {
	Email    string `json:"email,omitempty"`
	Verified bool   `json:"verified,omitempty"`
}

func (s *signee) githubPublicKey() ([]byte, error) {
	//FIXME: This is a horrible mess isn't it.
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/users/%s/gpg_keys", DefaultGithubAPIEndpoint, s.user), nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create gpg keys request")
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	res, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list users gpg keys")
	}

	var keys []GPGKey
	err = json.NewDecoder(res.Body).Decode(&keys)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode users gpg keys")
	}

	for _, key := range keys {
		if key.KeyID == s.key {
			verified := false
			for _, email := range key.Emails {
				if email.Verified {
					verified = true
					break
				}
			}
			if !verified {
				return nil, fmt.Errorf("gpg key: %s is not verified for user: %s", s.key, s.user)
			}

			return []byte(key.RawKey), nil
		}
	}

	return nil, fmt.Errorf("gpg key: %s not found for user: %s", s.key, s.user)
}

func (s *signee) keybasePublicKey() ([]byte, error) {
	return nil, fmt.Errorf("not implemented")
}

// Key returns the public key
func (s *signee) Key() string {
	return s.key
}

// User returns the user
func (s *signee) User() string {
	return s.user
}

// Type returns the signee type
func (s *signee) Type() SigneeType {
	return s.signeeType
}
