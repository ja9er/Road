package view

import (
	"Road/moudle/sqlmoudle"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"time"
)

//rsa解密
func rsadecry(encryptString string) (result string, err error) {
	var privatePEMData = []byte(`-----BEGIN PRIVATE KEY-----
MIIBVAIBADANBgkqhkiG9w0BAQEFAASCAT4wggE6AgEAAkEAtVa12acTsp/Ozmrx
fmkUlmC1POB72aDfiKEgN9tIiOl22B/sD11S9Hjg3VmpFj3TT8dO75f4XLW4r+J5
f6cMxQIDAQABAkEAqc8K90gnf+t6Q32NquxHpRHmZZ1pHMAy0sTfYK7tW5Z9yadJ
Y85u1hdPujSrgFGTUQxblmY5NfW9bKCSNeCzAQIhAPG/nDpE1aY6/bH/I+Itrzah
tLJT0UxISMwF9AcsjmThAiEAwAdqXH1IzOaS2Bid2ukR4TjdktxkYwd07YWF9zlh
QGUCIHnahVrxm2eA0KPJ4UJ+mJTHCZfhm9wBi4AbeBeto9DBAiBbLao1DE/a6shi
zx106iHRPP0IVJld5BaDCVlYz+f7eQIgIuwmIQmIx4vO+R5tEeKI1MvTtqNxEPo2
80lV7QWE/NY=
-----END PRIVATE KEY-----`)
	block2, _ := pem.Decode(privatePEMData)
	private, _ := x509.ParsePKCS8PrivateKey(block2.Bytes)
	privatekey := private.(*rsa.PrivateKey)
	decodeString, _ := base64.StdEncoding.DecodeString(encryptString)
	decryptstring, err := rsa.DecryptPKCS1v15(rand.Reader, privatekey, decodeString)
	if err != nil {
		return "", err
	}
	return string(decryptstring), nil
}

func rsaencryt(encrystring string) (result string) {
	var pubPEMData = []byte(`-----BEGIN PUBLIC KEY-----
MFwwDQYJKoZIhvcNAQEBBQADSwAwSAJBALVWtdmnE7Kfzs5q8X5pFJZgtTzge9mg
34ihIDfbSIjpdtgf7A9dUvR44N1ZqRY900/HTu+X+Fy1uK/ieX+nDMUCAwEAAQ==
-----END PUBLIC KEY-----`)
	block, _ := pem.Decode(pubPEMData)
	if block == nil || block.Type != "PUBLIC KEY" {
		log.Fatal("failed to decode PEM block containing public key")
		os.Exit(1)
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		log.Fatal(err)
	}
	publicKey := pub.(*rsa.PublicKey)
	encryptPKCS1v15, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, []byte(encrystring))
	encrystring = base64.StdEncoding.EncodeToString(encryptPKCS1v15)
	return encrystring
}

func CheckAuth(c *gin.Context) {
	hash, _ := c.Cookie("HMACCOUNT")
	res, err := rsadecry(hash)
	if err != nil {
		c.Redirect(http.StatusFound, "/login")
		c.Abort()
	}
	log.Println("[+] ", res)
	c.Next()
}
func commentHandler(c *gin.Context) {
	username := c.PostForm("username")
	passwd := c.PostForm("password")
	userres := sqlmoudle.Queryuser(username)
	if userres.Passwd == passwd {
		formatTimeStr := time.Now().Format("2006-01-02 15-04-05")
		fmt.Println(formatTimeStr)
		c.SetCookie("HMACCOUNT", rsaencryt(username+":"+formatTimeStr), 0, "", "", false, true)
		c.SetCookie("name", username, 0, "", "", false, true)
		c.Redirect(http.StatusFound, "/index")
	} else {
		c.HTML(http.StatusOK, "test.html", nil)
	}

}

func forwardHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "test.html", nil)
}

func Loadlogin(e *gin.Engine) {
	e.POST("/login", commentHandler)
	e.GET("/login", forwardHandler)

}
