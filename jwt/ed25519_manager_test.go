package jwt

import (
	"errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	priPem = `-----BEGIN PRIVATE KEY-----
MC4CAQAwBQYDK2VwBCIEIHXEAUN6Lp8Hdq8P0Mcv9mjIG1sgPWBf1Mh+OKP5HXvC
-----END PRIVATE KEY-----`
	pubPem = `-----BEGIN PUBLIC KEY-----
MCowBQYDK2VwAyEAxCSxEyY/+A7T7EtXF7AHw4Zfklh/QdjG8fxfRFYZgY8=
-----END PUBLIC KEY-----`
)

type user struct {
	Id uint64
}

func TestEd25519Manager(t *testing.T) {
	tcs := []struct {
		name           string
		user           user
		manager        Manager[user]
		wantEncryptErr error
		wantDecryptErr error
	}{
		{
			name: "basic",
			user: user{Id: 1},
			manager: func() *Ed25519Manager[user] {
				manager, err := NewEd25519Manager[user](priPem, pubPem)
				require.NoError(t, err)
				return manager
			}(),
			wantEncryptErr: nil,
			wantDecryptErr: nil,
		}, {
			name: "expired",
			user: user{Id: 1},
			manager: func() *Ed25519Manager[user] {
				cfg := NewClaimsConfig(time.Millisecond)

				manager, err := NewEd25519Manager[user](priPem, pubPem, WithClaimsConfig[user](cfg))
				require.NoError(t, err)
				return manager
			}(),
			wantEncryptErr: nil,
			wantDecryptErr: jwt.ErrTokenExpired,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			token, err := tc.manager.Encrypt(tc.user)
			assert.Equal(t, tc.wantEncryptErr, err)
			if err != nil {
				return
			}

			time.Sleep(time.Millisecond * 2)
			decrypted, err := tc.manager.Decrypt(token)
			assert.Truef(t, errors.Is(err, tc.wantDecryptErr), "want: %v, got: %v", tc.wantDecryptErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.user, decrypted.Data)
		})
	}
}

func TestEd25519Manager_InvalidToken(t *testing.T) {
	manager, err := NewEd25519Manager[user](priPem, pubPem)
	require.NoError(t, err)

	_, err = manager.Encrypt(user{Id: 1})
	require.NoError(t, err)

	_, err = manager.Decrypt("invalid token")
	wantErr := jwt.ErrTokenMalformed
	assert.Truef(t, errors.Is(err, wantErr), "want: %v, got: %v", wantErr, err)
}
