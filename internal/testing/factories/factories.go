package factories

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"sync/atomic"
	"testing"

	"golang.org/x/crypto/ssh"
)

var nameSequence uint64 = 1

func GetSequenceName(prefix string) string {
	return fmt.Sprintf("%s %d", prefix, atomic.AddUint64(&nameSequence, 1))
}

// GenerateRSAPubKey generates pubkey for use in tests
func GenerateRSAPubKey(t *testing.T) string {
	t.Helper()
	rsaKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("Failed to generate pubky: %v", err)
	}
	pubkey, err := ssh.NewPublicKey(&rsaKey.PublicKey)
	if err != nil {
		t.Fatalf("Failed to generate pubky: %v", err)
	}
	return string(ssh.MarshalAuthorizedKey(pubkey))
}
