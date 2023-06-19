package factories

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"sync/atomic"
	"testing"

	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/ssh"
)

var nameSequence uint64 = 1

// SeqNameWithPrefix returns prefix with atomically increasing integer
func SeqNameWithPrefix(prefix string) string {
	return fmt.Sprintf("%s %d", prefix, atomic.AddUint64(&nameSequence, 1))
}

// GenerateRSAPubKey generates pubkey for use in tests
func GenerateRSAPubKey(t *testing.T) string {
	t.Helper()
	rsaKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("Failed to generate pubkey: %v", err)
	}
	pubkey, err := ssh.NewPublicKey(&rsaKey.PublicKey)
	if err != nil {
		t.Fatalf("Failed to generate pubkey: %v", err)
	}
	return string(ssh.MarshalAuthorizedKey(pubkey))
}

func PubkeyWithTrans(t *testing.T, ctx context.Context, pk *models.Pubkey) *models.Pubkey {
	terr := models.Transform(ctx, pk)
	require.NoError(t, terr)
	return pk
}

func NewPubkeyED25519() *models.Pubkey {
	return &models.Pubkey{
		AccountID: 1,
		Name:      SeqNameWithPrefix("lzap-ed25519-2021"),
		Body:      "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIEhnn80ZywmjeBFFOGm+cm+5HUwm62qTVnjKlOdYFLHN lzap-2021",
	}
}

func NewPubkeyRSA() *models.Pubkey {
	return &models.Pubkey{
		AccountID: 1,
		Name:      SeqNameWithPrefix("lzap-rsa-2013"),
		Body: "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC8w6DONv1qn3IdgxSpkYOClq7oe7davWFqKVHPbLoS6+dF" +
			"Inru7gdEO5byhTih6+PwRhHv/b1I+Mtt5MDZ8Sv7XFYpX/3P/u5zQiy1PkMSFSz0brRRUfEQxhXLW97FJa7l+bej2HJ" +
			"Dt7f9Gvcj+d/fNWC9Z58/GX11kWk4SIXaKotkN+kWn54xGGS7Zvtm86fP59Srt6wlklSsG8mZBF7jVUjyhAgm/V5gDF" +
			"b2/6jfiwSb2HyJ9/NbhLkWNdwrvpdGZqQlYhnwTfEZdpwizW/Mj3MxP5O31HN45aE0wog0UeWY4gvTl4Ogb6kescizA" +
			"M6pCff3RBslbFxLdOO7cR17 lzap-2013",
	}
}

func NewEmptyPubkeyRSA() *models.Pubkey {
	return &models.Pubkey{
		AccountID: 1,
		Name:      SeqNameWithPrefix("lzap-rsa-2013"),
		Body:      "",
	}
}

func NewPubkeyRSAWithoutUsername() *models.Pubkey {
	return &models.Pubkey{
		AccountID: 1,
		Name:      SeqNameWithPrefix("lzap-rsa-2013"),
		Body: "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC8w6DONv1qn3IdgxSpkYOClq7oe7davWFqKVHPbLoS6+dF" +
			"Inru7gdEO5byhTih6+PwRhHv/b1I+Mtt5MDZ8Sv7XFYpX/3P/u5zQiy1PkMSFSz0brRRUfEQxhXLW97FJa7l+bej2HJ" +
			"Dt7f9Gvcj+d/fNWC9Z58/GX11kWk4SIXaKotkN+kWn54xGGS7Zvtm86fP59Srt6wlklSsG8mZBF7jVUjyhAgm/V5gDF" +
			"b2/6jfiwSb2HyJ9/NbhLkWNdwrvpdGZqQlYhnwTfEZdpwizW/Mj3MxP5O31HN45aE0wog0UeWY4gvTl4Ogb6kescizA" +
			"M6pCff3RBslbFxLdOO7cR17",
	}
}

func NewPubkeyDSS() *models.Pubkey {
	return &models.Pubkey{
		AccountID: 1,
		Name:      SeqNameWithPrefix("lzap-dsa-2011"),
		Body: "ssh-dss AAAAB3NzaC1kc3MAAACBAKqezP3rkK/NcWvMWqoP3qOggGG4QW1vhQJOfyH/l9CbdRxlrcTV9AD5" +
			"BYMcJNn3Ill0iu9d7gSQTZJu2cEWiE8yHJhWOerfPDB4R8BGQlMvbO+8rTplm1Eo3WxtYD0q45Urfh/Ej7HgliTsAYB" +
			"YrQZ0a09auzBlqR3XwH74MlPdAAAAFQClJSTbX6Hp9HzqXyw0P7HeXt0LrwAAAIAogU1yFPDn7xPPUEh16u3ceaZp5w" +
			"H2wDzPEjMHPv+GQd2/yiJB5TX5s9Z5HQax/r3NFhYKNzjyQf1alChS8M0ge9vtPx3oH3Q3NyJGo2wpyYzvDXzP9OHO6" +
			"Vh3PVVOcGL/TlYbFJUeJb8usjtpb4sLmUNuohwifXNAKzkFj/YpswAAAIAf96KvZqMC91JocjY0L09G2kH+v4Ax30VY" +
			"w3iFlA5LgYnbKEBxEvzM+xZ98uRT//Dmn76F5pFIk/QsHpDSHlx5TIuf1pIm6vzuWtUUQUYKTl+ljuft2FY+FfNW4MZ" +
			"ZKx52kr96AOGKKi+U/MAklf+obqf22XFGvNNu4KSjbxqxwg== lzap-2011",
	}
}

func NewPubkeyECDSA() *models.Pubkey {
	return &models.Pubkey{
		AccountID: 1,
		Name:      SeqNameWithPrefix("avitova-nistp256-2021"),
		Body: "ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBAaOrIRmMPX84l" +
			"YJ6y3mzH4gBLLCRdeAJX/lsImAn98u3wghha7pD+bp0O9d1iueMVcRpxfnOpxy3hBAoerDjOw= avitova-2021",
	}
}
