package models

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
)

func RsaDecrypt(base64Data string) []byte {
	rsaData, _ := base64.StdEncoding.DecodeString(base64Data)
	block, _ := pem.Decode(privateKey)
	priv, _ := x509.ParsePKCS1PrivateKey(block.Bytes)
	result ,_ := rsa.DecryptPKCS1v15(rand.Reader, priv, rsaData)
	return result
}

var privateKey = []byte(`
-----BEGIN RSA PRIVATE KEY-----
MIICXgIBAAKBgQCyH/slPc3xg1AkBce0hhNjKxz5HoQUSevtALexFOEBj81IIo2v
l8DZJy8iRBgFUD673spr1BngolTBiOPCJkwVcAbHlHFUJ5DL1Rduiyc9fwMyl7Eg
6lEWXev7Du6EPGrT4tEz82AptJ6qoWuenanDJDYg3jEEAk6734VQRsCYawIDAQAB
AoGBAJFMOx4Gyz4tgirQODYOhDQJkAm6Fb1DC1r5kd22DVCrz6T+4pqQbDP2naES
8JEtAu9W7cGFc1JkuERieH7/pGEfqUAPgJs6fEcY0co+FE6C17QQZc6iB9nB3WXV
bQDae++2zKB6pAIR/qxDdGuE9CsG/cPRtzXxVUpclxHxD5EpAkEA3h687NkJ0NbY
KYmbLj3Pcfcu692lXgyIQFZiGqO8WIqHxTDATAr2ln7U8bP6awUWN5QCtAF21wyk
46mw7Qb5RQJBAM1LVkUKNxogpdpfs07SIArEN7IBEPYhL8snyyykR9UuO5Tg6BSW
/37Ln7eJ8VtL0ikgFf9q5k7u2B0DTdOl7e8CQDmQDbXzqS+N/gcFukmJizEltes6
TZjJ9qV1vYbZ1/26KOVZdPw/+xeVVuoskkEZ2GAe43RyzLF+fVzipQ9IN2ECQQCx
X9Y58ImLWYnzE5ypDYQByWcVtTYicqoIrWkuOQKXfkqcZ3Yd1BkMRILK4bRXXTtH
rSFUfdhfep3e82va4hKhAkEAg86PYaEk1+aeeMayPckifsPzkvQt5WOZ2VXFRQQI
bvMxhVc7NNIn2kgUI5jCL8/smB4PINowlnfQhT3Bu/9qNg==
-----END RSA PRIVATE KEY-----
`)
