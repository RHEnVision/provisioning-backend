package stubs

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"fmt"

	"github.com/RHEnVision/provisioning-backend/internal/models"
	"golang.org/x/crypto/ssh"
)

// GeneratePubkey generates a pubkey and stores it in the context stub
func GeneratePubkey(ctx context.Context, pubkey models.Pubkey) error {
	pubkeyDao, err := getPubkeyDaoStub(ctx)
	if err != nil {
		return err
	}
	if pubkey.AccountID == 0 {
		pubkey.AccountID = 1
	}
	if pubkey.Name == "" {
		pubkey.Name = fmt.Sprintf("pubkey %d", pubkeyDao.lastId+1)
	}
	if pubkey.Body == "" {
		if pubkey.Body, err = generateRandomRSAPubKey(ctx); err != nil {
			return err
		}
	}
	if err = pubkeyDao.Create(ctx, &pubkey); err != nil {
		return err
	}
	return nil
}

func generateRandomRSAPubKey(_ context.Context) (string, error) {
	rsaKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return "", RSAGenerationError
	}
	pubkey, err := ssh.NewPublicKey(&rsaKey.PublicKey)
	if err != nil {
		return "", RSAGenerationError
	}
	return string(ssh.MarshalAuthorizedKey(pubkey)), nil
}
