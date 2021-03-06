/*
 * Arc - Copyleft of Simone 'evilsocket' Margaritelli.
 * evilsocket at protonmail dot com
 * https://www.evilsocket.net/
 *
 * See LICENSE.
 */
package tls

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"github.com/evilsocket/arc/arcd/config"
	"github.com/evilsocket/arc/arcd/log"
	"math/big"
	"os"
	"time"
)

func Generate(conf *config.Configuration) error {
	keyfile, err := os.Create(conf.Key)
	if err != nil {
		return err
	}
	defer keyfile.Close()

	certfile, err := os.Create(conf.Certificate)
	if err != nil {
		return err
	}
	defer certfile.Close()

	log.Debugf("Generating RSA key ...")

	priv, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return err
	}

	log.Debugf("Creating X509 certificate ...")

	notBefore := time.Now()
	notAfter := notBefore.Add(time.Duration(24*365) * time.Hour)
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return err
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			CommonName:         "Arc Server",
			Organization:       []string{"Arc Project"},
			OrganizationalUnit: []string{"RSA key generation module"},
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	cert_raw, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return err
	}

	log.Importantf("Saving key to %s ...", log.Bold(conf.Key))
	if err := pem.Encode(keyfile, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)}); err != nil {
		return err
	}

	log.Importantf("Saving certificate to %s ...", log.Bold(conf.Certificate))
	return pem.Encode(certfile, &pem.Block{Type: "CERTIFICATE", Bytes: cert_raw})
}
