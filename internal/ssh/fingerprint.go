package ssh

import (
	"crypto/md5"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"strings"

	"golang.org/x/crypto/ssh"
)

// OpenSSHFingerprints is the de-facto standard OpenSSH fingerprints for SSH public keys:
// SHA256 (used for ED type keys) and MD5 (used for RSA keys). Fingerprints are returned as
// string encoded into base64 or hex respectively. Additionally, type and comment are also
// returned. Type as one of the: "ssh-ed25519" or "ssh-rsa".
type OpenSSHFingerprints struct {
	Type    string
	SHA256  string
	MD5     string
	Comment string
}

// AWSFingerprint is a special Amazon EC2 MD5 fingerprint calculated from the PEM-encoded
// public key. It is used in AWS for ssh keys.
type AWSFingerprint string

//#nosec
func md5Separator(data []byte, separator byte) string {
	hashSlice := md5.Sum(data)
	var sb strings.Builder
	buf := make([]byte, 2)
	sb.Grow(47)
	for i := 1; i <= len(hashSlice); i++ {
		hex.Encode(buf, hashSlice[i-1:i])
		sb.Write(buf)
		if i < 16 {
			sb.WriteByte(separator)
		}
	}
	return sb.String()
}

// GenerateOpenSSHFingerprints parses a public key and returns OpenSSH fingerprints.
func GenerateOpenSSHFingerprints(pubkeyBody []byte) (OpenSSHFingerprints, error) {
	fps := OpenSSHFingerprints{}

	pkey, cmt, _, _, err := ssh.ParseAuthorizedKey(pubkeyBody)
	if err != nil {
		return OpenSSHFingerprints{}, fmt.Errorf("unable to parse public key %s: %w", pubkeyBody, err)
	}

	fps.Comment = cmt
	fps.Type = pkey.Type()
	fps.SHA256 = strings.TrimLeft(ssh.FingerprintSHA256(pkey), "SHA256:") + "="
	fps.MD5 = strings.TrimLeft(ssh.FingerprintLegacyMD5(pkey), "MD5:")

	return fps, nil
}

// GenerateAWSFingerprint parses a public key and returns AWS PEM fingerprint used for RSA keys.
// MD5 fingerprint stored as hexadecimal with colons without any prefix from key in PEM format.
// This format is specific to AWS. To generate such fingerprint:
//
// ssh-keygen -e -f $HOME/.ssh/key.pub -m pkcs8 | openssl pkey -pubin -outform der | openssl md5 -c
//
// Example: "c4:ba:72:45:16:a9:2c:39:c3:99:8d:e7:16:01:9c:77"
func GenerateAWSFingerprint(pubkeyBody []byte) (AWSFingerprint, error) {
	pkey, _, _, _, err := ssh.ParseAuthorizedKey(pubkeyBody)
	if err != nil {
		return "", fmt.Errorf("unable to parse public key %s: %w", pubkeyBody, err)
	}

	parsedCryptoKey := pkey.(ssh.CryptoPublicKey)
	pub := parsedCryptoKey.CryptoPublicKey()
	der, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		return "", fmt.Errorf("failed to calculate legacy fingerprint: %w", err)
	}

	return AWSFingerprint(md5Separator(der, ':')), nil
}
