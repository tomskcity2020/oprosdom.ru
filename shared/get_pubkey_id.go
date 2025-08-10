package shared

import (
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/hex"
	"log"
	"strings"
)

func GetPubkeyId(pubKey *rsa.PublicKey) string {

	pubDER, err := x509.MarshalPKIXPublicKey(pubKey)
	if err != nil {
		log.Fatalf("Failed to marshal public key: %v", err)
	}

	hash := sha1.Sum(pubDER)

	fullFp := hex.EncodeToString(hash[:])
	return strings.ToLower(fullFp[len(fullFp)-8:]) // делаем как в gnupg
}
