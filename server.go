package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"math/big"
	"net"
	"time"
)

func generateTLSConfig() (*tls.Config, error) {
	// RSAキーペアを作成
	key, e := rsa.GenerateKey(rand.Reader, 4096)
	name := "test"
	// 証明書を作成
	der, e := createCACertificate(name, key)
	if e != nil {
		return nil, e
	}
	return &tls.Config{
		ServerName: name,
		Certificates: []tls.Certificate{{
			Certificate: [][]byte{der},
			PrivateKey:  key,
		}},
		// HTTP/2をサポートするために必要
		NextProtos: []string{"h2"},
	}, e
}

func createCACertificate(name string, k *rsa.PrivateKey) ([]byte, error) {
	// X.509 識別名の構造体を作成
	subject := pkix.Name{CommonName: name}
	// x.509証明書の構造体を作成
	cert := &x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               subject,
		NotBefore:             time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
		NotAfter:              time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
		BasicConstraintsValid: true,
		IsCA:                  true,
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
	}
	// DERエンコーディングされた証明書を作成
	b, e := x509.CreateCertificate(rand.Reader, cert, cert, &k.PublicKey, k)
	if e != nil {
		return nil, fmt.Errorf("failed to create certificate. %v", e)
	}
	return b, nil
}

func newListener(addr string, tlsCfg *tls.Config) (net.Listener, error) {
	l, e := net.Listen("tcp", addr)
	if e != nil {
		return nil, e
	}
	// tlsconfigを使用した新たなネットワークリスナーを作成
	return tls.NewListener(l, tlsCfg), nil
}
