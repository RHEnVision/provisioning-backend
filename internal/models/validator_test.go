package models_test

import (
	"context"
	"testing"

	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/RHEnVision/provisioning-backend/internal/testing/factories"
	"github.com/stretchr/testify/assert"
)

func TestPubkeyValid(t *testing.T) {
	pk := factories.NewPubkeyED25519()
	err := models.Validate(context.Background(), pk)
	assert.Nil(t, err)
}

func TestPubkeyMissingName(t *testing.T) {
	pk := factories.NewPubkeyED25519()
	pk.Name = ""
	err := models.Validate(context.Background(), pk)
	assert.EqualError(t, err, "Key: 'Pubkey.Name' Error:Field validation for 'Name' failed on the 'required' tag")
}

func TestPubkeyMissingBody(t *testing.T) {
	pk := factories.NewPubkeyED25519()
	pk.Body = ""
	err := models.Validate(context.Background(), pk)
	assert.EqualError(t, err, "Key: 'Pubkey.Body' Error:Field validation for 'Body' failed on the 'required' tag")
}

func TestPubkeyInvalidBody(t *testing.T) {
	pk := factories.NewPubkeyED25519()
	pk.Body = "xxx"
	err := models.Validate(context.Background(), pk)
	assert.EqualError(t, err, "Key: 'Pubkey.Body' Error:Field validation for 'Body' failed on the 'sshPubkey' tag")
}

func TestBodyWithUsername(t *testing.T) {
	t.Run("RSA", func(t *testing.T) {
		pk := factories.NewPubkeyRSA()
		pkBody, err := pk.BodyWithUsername(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, "gcp-user:ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC8w6DONv1qn3IdgxSpkYOClq7oe7davWFqKVHPbLoS6+dF"+
			"Inru7gdEO5byhTih6+PwRhHv/b1I+Mtt5MDZ8Sv7XFYpX/3P/u5zQiy1PkMSFSz0brRRUfEQxhXLW97FJa7l+bej2HJ"+
			"Dt7f9Gvcj+d/fNWC9Z58/GX11kWk4SIXaKotkN+kWn54xGGS7Zvtm86fP59Srt6wlklSsG8mZBF7jVUjyhAgm/V5gDF"+
			"b2/6jfiwSb2HyJ9/NbhLkWNdwrvpdGZqQlYhnwTfEZdpwizW/Mj3MxP5O31HN45aE0wog0UeWY4gvTl4Ogb6kescizA"+
			"M6pCff3RBslbFxLdOO7cR17", pkBody)
	})
	t.Run("DSS", func(t *testing.T) {
		pk := factories.NewPubkeyDSS()
		pkBody, err := pk.BodyWithUsername(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, "gcp-user:ssh-dss AAAAB3NzaC1kc3MAAACBAKqezP3rkK/NcWvMWqoP3qOggGG4QW1vhQJOfyH/l9CbdRxlrcTV9AD5"+
			"BYMcJNn3Ill0iu9d7gSQTZJu2cEWiE8yHJhWOerfPDB4R8BGQlMvbO+8rTplm1Eo3WxtYD0q45Urfh/Ej7HgliTsAYB"+
			"YrQZ0a09auzBlqR3XwH74MlPdAAAAFQClJSTbX6Hp9HzqXyw0P7HeXt0LrwAAAIAogU1yFPDn7xPPUEh16u3ceaZp5w"+
			"H2wDzPEjMHPv+GQd2/yiJB5TX5s9Z5HQax/r3NFhYKNzjyQf1alChS8M0ge9vtPx3oH3Q3NyJGo2wpyYzvDXzP9OHO6"+
			"Vh3PVVOcGL/TlYbFJUeJb8usjtpb4sLmUNuohwifXNAKzkFj/YpswAAAIAf96KvZqMC91JocjY0L09G2kH+v4Ax30VY"+
			"w3iFlA5LgYnbKEBxEvzM+xZ98uRT//Dmn76F5pFIk/QsHpDSHlx5TIuf1pIm6vzuWtUUQUYKTl+ljuft2FY+FfNW4MZ"+
			"ZKx52kr96AOGKKi+U/MAklf+obqf22XFGvNNu4KSjbxqxwg==", pkBody)
	})
	t.Run("ECDSA", func(t *testing.T) {
		pk := factories.NewPubkeyECDSA()
		pkBody, err := pk.BodyWithUsername(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, "gcp-user:ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBAaOrIRmMPX84l"+
			"YJ6y3mzH4gBLLCRdeAJX/lsImAn98u3wghha7pD+bp0O9d1iueMVcRpxfnOpxy3hBAoerDjOw=", pkBody)
	})
	t.Run("ED25519", func(t *testing.T) {
		pk := factories.NewPubkeyED25519()
		pkBody, err := pk.BodyWithUsername(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, "gcp-user:ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIEhnn80ZywmjeBFFOGm+cm+5HUwm62qTVnjKlOdYFLHN", pkBody)
	})

	t.Run("empty RSA", func(t *testing.T) {
		pk := factories.NewEmptyPubkeyRSA()
		pkBody, err := pk.BodyWithUsername(context.Background())
		assert.EqualError(t, err, "invalid public key format")
		assert.Equal(t, pkBody, "")
	})

	t.Run("RSA without a username", func(t *testing.T) {
		pk := factories.NewPubkeyRSAWithoutUsername()
		pkBody, err := pk.BodyWithUsername(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, pkBody, "gcp-user:ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC8w6DONv1qn3IdgxSpkYOClq7oe7davWFqKVHPbLoS6+dF"+
			"Inru7gdEO5byhTih6+PwRhHv/b1I+Mtt5MDZ8Sv7XFYpX/3P/u5zQiy1PkMSFSz0brRRUfEQxhXLW97FJa7l+bej2HJ"+
			"Dt7f9Gvcj+d/fNWC9Z58/GX11kWk4SIXaKotkN+kWn54xGGS7Zvtm86fP59Srt6wlklSsG8mZBF7jVUjyhAgm/V5gDF"+
			"b2/6jfiwSb2HyJ9/NbhLkWNdwrvpdGZqQlYhnwTfEZdpwizW/Mj3MxP5O31HN45aE0wog0UeWY4gvTl4Ogb6kescizA"+
			"M6pCff3RBslbFxLdOO7cR17")
	})

	t.Run("nil pk", func(t *testing.T) {
		pk := &models.Pubkey{}
		pkBody, err := pk.BodyWithUsername(context.Background())
		assert.EqualError(t, err, "invalid public key format")
		assert.Equal(t, pkBody, "")
	})
}
