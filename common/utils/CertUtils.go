package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/hex"
	"encoding/pem"
	"github.com/kris-nova/logger"
	"io/ioutil"
	"math/big"
	"os"
	"time"
)

func GenerateCert(org string, orgUnit string, country string, locality string, cname string, alias []string, targetDir string) (string, string, error) {
	ca := &x509.Certificate{
		SerialNumber: big.NewInt(1653),
		Subject: pkix.Name{
			Organization:       []string{org},
			OrganizationalUnit: []string{orgUnit},
			Country:            []string{country},
			Locality:           []string{locality},
			CommonName:         cname,
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		IsCA:                  false,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
		DNSNames:              alias,
	}

	priv, _ := rsa.GenerateKey(rand.Reader, 2048)
	pub := &priv.PublicKey
	ca_b, err := x509.CreateCertificate(rand.Reader, ca, ca, pub, priv)
	if err != nil {
		logger.Critical("create ca failed", err)
		return "", "", nil
	}

	// Public key
	certOut, err := os.Create(targetDir + "ca.crt")
	_ = pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: ca_b})
	_ = certOut.Close()

	// Private key
	keyOut, err := os.OpenFile(targetDir+"ca.key", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	_ = pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})
	_ = keyOut.Close()

	crt, _ := ioutil.ReadFile(targetDir + "/ca.crt")
	key, _ := ioutil.ReadFile(targetDir + "/ca.key")

	return string(crt), string(key), nil
}

func GenerateRandomKey() string {
	key := make([]byte, 16)
	_, err := rand.Read(key)
	if err != nil {
		// todo - handle error
		logger.Critical("got error " ,err)
	}

	return hex.EncodeToString(key)
}
