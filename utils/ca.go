package utils

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"time"
)

func LocalIP() (net.IP, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP, nil
}

func InitCA() error {
	ip, err := LocalIP()
	if err != nil {
		return err
	}
	directory := ".ca"
	_, err = os.Stat(directory)
	if os.IsNotExist(err) {
		os.Mkdir(directory, 0755)

		duration, _ := time.ParseDuration(fmt.Sprintf("%vs", 3*365*24*60*60))
		easyCert := EasyCert{
			org:      "IPAPK Generated CA " + ip.String(),
			duration: duration,
			rsaBits:  2048,
			hosts:    []string{ip.String()},
		}

		ca, caKey, err := easyCert.generateCA(filepath.Join(directory, "myCA.cer"))
		if err != nil {
			return err
		}
		certFile := filepath.Join(directory, "mycert1.cer")
		keyFile := filepath.Join(directory, "mycert1.key")
		if err := easyCert.generateCert(certFile, keyFile, ca, caKey); err != nil {
			return err
		}
	}
	return nil
}

type EasyCert struct {
	org      string
	duration time.Duration
	rsaBits  int
	hosts    []string
	ec       string
}

func (c EasyCert) newCertificate() *x509.Certificate {
	notBefore := time.Now()
	notAfter := notBefore.Add(c.duration)

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		log.Fatalf("failed to generate serial number: %s", err)
	}

	return &x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{c.org},
		},
		NotBefore: notBefore,
		NotAfter:  notAfter,

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		BasicConstraintsValid: true,
	}
}

// newCertificate creates a new template
func (c EasyCert) newPrivateKey() (crypto.PrivateKey, error) {
	if c.ec != "" {
		var curve elliptic.Curve
		switch c.ec {
		case "224":
			curve = elliptic.P224()
		case "384":
			curve = elliptic.P384()
		case "521":
			curve = elliptic.P521()
		default:
			return nil, fmt.Errorf("Unknown elliptic curve: %q", c.ec)
		}
		return ecdsa.GenerateKey(curve, rand.Reader)
	}
	return rsa.GenerateKey(rand.Reader, c.rsaBits)
}

// newPrivateKey creates a new private key depending
// on the input flags
func (c EasyCert) generateCA(caFile string) (*x509.Certificate, crypto.PrivateKey, error) {
	template := c.newCertificate()
	template.IsCA = true
	template.KeyUsage |= x509.KeyUsageCertSign
	template.Subject.CommonName = c.org

	priv, err := c.newPrivateKey()
	if err != nil {
		return nil, nil, err
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, template, template, priv.(crypto.Signer).Public(), priv)
	if err != nil {
		return nil, nil, err
	}

	ca, err := x509.ParseCertificate(derBytes)
	if err != nil {
		return nil, nil, err
	}

	certOut, err := os.Create(caFile)
	if err != nil {
		return nil, nil, err
	}
	defer certOut.Close()
	if err := pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes}); err != nil {
		return nil, nil, err
	}

	return ca, priv, nil
}

// generateCA creates a new CA certificate, saves the certificate
// and returns the x509 certificate and crypto private key. This
// private key should never be saved to disk, but rather used to
// immediately generate further certificates.
func (c EasyCert) generateCert(certFile, keyFile string, ca *x509.Certificate, caKey crypto.PrivateKey) error {
	template := c.newCertificate()
	for _, h := range c.hosts {
		if ip := net.ParseIP(h); ip != nil {
			template.IPAddresses = append(template.IPAddresses, ip)
		} else {
			template.DNSNames = append(template.DNSNames, h)
			if template.Subject.CommonName == "" {
				template.Subject.CommonName = h
			}
		}
	}

	priv, err := c.newPrivateKey()
	if err != nil {
		return err
	}

	return c.generateFromTemplate(certFile, keyFile, template, ca, priv, caKey)
}

// generateCert generates a new certificate for the given hosts using the
// provided certificate authority. The cert and key files are stored in the
// the provided files.
func (c EasyCert) generateClient(certFile, keyFile string, ca *x509.Certificate, caKey crypto.PrivateKey) error {
	template := c.newCertificate()
	template.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth}

	priv, err := c.newPrivateKey()
	if err != nil {
		return err
	}

	return c.generateFromTemplate(certFile, keyFile, template, ca, priv, caKey)
}

// generateFromTemplate generates a certificate from the given template and signed by
// the given parent, storing the results in a certificate and key file.
func (c EasyCert) generateFromTemplate(certFile, keyFile string, template, parent *x509.Certificate, key crypto.PrivateKey, parentKey crypto.PrivateKey) error {
	derBytes, err := x509.CreateCertificate(rand.Reader, template, parent, key.(crypto.Signer).Public(), parentKey)
	if err != nil {
		return err
	}

	certOut, err := os.Create(certFile)
	if err != nil {
		return err
	}
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	certOut.Close()

	return c.savePrivateKey(key, keyFile)
}

// savePrivateKey saves the private key to a PEM file
func (c EasyCert) savePrivateKey(key crypto.PrivateKey, file string) error {
	keyOut, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer keyOut.Close()

	switch v := key.(type) {
	case *rsa.PrivateKey:
		keyBytes := x509.MarshalPKCS1PrivateKey(v)
		pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: keyBytes})
	case *ecdsa.PrivateKey:
		keyBytes, err := x509.MarshalECPrivateKey(v)
		if err != nil {
			return err
		}
		pem.Encode(keyOut, &pem.Block{Type: "EC PRIVATE KEY", Bytes: keyBytes})
	default:
		return fmt.Errorf("Unsupport private key type: %#v", key)
	}

	return nil
}
