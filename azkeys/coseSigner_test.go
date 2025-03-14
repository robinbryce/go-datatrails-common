package azkeys

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"math/big"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/2016-10-01/keyvault"
	"github.com/stretchr/testify/assert"
)

// Test_publicKey tests:
//
// 1. a known public key jwk can be converted to a crypto.PublicKey
func Test_publicKey(t *testing.T) {

	kid := "merkle-log-signing/044da0af7a574acab984bee54a9946bc"
	keyops := []string{"verify"}
	x := "bBPhBjEt9H-AA6oD_D3Tkn3q5DEz2enBX4dTkh0_Rr1KRYGKU_i84ELM-jgJAKuH"
	y := "NRr8Nj8OoUf_voYkuTJDv_FFx6xZxyLmurBdpiimXreuBlTgdzE6AlZNBLp6Empg"

	// get the known values as big ints
	xBig := new(big.Int)
	xBig.SetBytes([]byte{108, 19, 225, 6, 49, 45, 244, 127, 128, 3, 170, 3, 252, 61, 211, 146, 125, 234, 228, 49, 51, 217, 233, 193, 95, 135, 83, 146, 29, 63, 70, 189, 74, 69, 129, 138, 83, 248, 188, 224, 66, 204, 250, 56, 9, 0, 171, 135})

	yBig := new(big.Int)
	yBig.SetBytes([]byte{53, 26, 252, 54, 63, 14, 161, 71, 255, 190, 134, 36, 185, 50, 67, 191, 241, 69, 199, 172, 89, 199, 34, 230, 186, 176, 93, 166, 40, 166, 94, 183, 174, 6, 84, 224, 119, 49, 58, 2, 86, 77, 4, 186, 122, 18, 106, 96})

	type args struct {
		key keyvault.KeyBundle
	}
	tests := []struct {
		name     string
		args     args
		expected *ecdsa.PublicKey
		err      error
	}{
		{
			name: "positive",
			args: args{
				key: keyvault.KeyBundle{
					Key: &keyvault.JSONWebKey{
						Kid:    &kid,
						KeyOps: &keyops,
						Kty:    keyvault.EC,
						Crv:    keyvault.P384,
						X:      &x,
						Y:      &y,
					},
				},
			},
			expected: &ecdsa.PublicKey{
				Curve: elliptic.P384(),
				X:     xBig,
				Y:     yBig,
			},
			err: nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := publicKey(test.args.key)

			assert.Equal(t, test.err, err)
			assert.Equal(t, test.expected, actual)
		})
	}
}
