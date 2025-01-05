package helpers

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/charmbracelet/log"

	"github.com/golang-jwt/jwt"
)

func EncryptPassword(data string) (string, error) {
	key := []byte(os.Getenv("ENCRYPTION_KEY"))[:32] // AES-256 needs 32-byte key
	c, err := aes.NewCipher(key)
	// if there are any errors, handle them
	if err != nil {
		return "", err
	}

	// gcm or Galois/Counter Mode, is a mode of operation
	// for symmetric key cryptographic block ciphers
	// - https://en.wikipedia.org/wiki/Galois/Counter_Mode
	gcm, err := cipher.NewGCM(c)
	// if any error generating new GCM
	// handle them
	if err != nil {
		return "", err
	}

	// creates a new byte array the size of the nonce
	// which must be passed to Seal
	nonce := make([]byte, gcm.NonceSize())
	// populates our nonce with a cryptographically secure
	// random sequence
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// here we encrypt our text using the Seal function
	// Seal encrypts and authenticates plaintext, authenticates the
	// additional data and appends the result to dst, returning the updated
	// slice. The nonce must be NonceSize() bytes long and unique for all
	// time, for a given key.
	eP := gcm.Seal(nonce, nonce, []byte(data), nil)

	return base64.StdEncoding.EncodeToString(eP), nil
}

func DecryptPassword(encrypted string) (string, error) {
	key := []byte(os.Getenv("ENCRYPTION_KEY"))[:32]
	eP, err := base64.StdEncoding.DecodeString(encrypted)

	if err != nil {
		return "", err

	}

	c, err := aes.NewCipher(key)
	if err != nil {
		return "", err

	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return "", err

	}

	nonceSize := gcm.NonceSize()
	if len(eP) < nonceSize {
		return "", err

	}

	nonce, eP := eP[:nonceSize], eP[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, eP, nil)
	if err != nil {
		log.Error(err)
	}

	return string(plaintext), nil
}

func GenerateJWT(ID uint, email string) (string, error) {
	var mySigningKey = []byte(os.Getenv("JWT_SECRET"))
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["email"] = email
	claims["ID"] = ID
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	tokenString, err := token.SignedString(mySigningKey)

	if err != nil {
		log.Fatal("Something Went Wrong: ", err.Error())
		return "", err
	}
	return tokenString, nil
}

func VerifyToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		// check expiration
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			return nil, fmt.Errorf("token Expired")
		}

		return claims, nil

	} else {
		return nil, fmt.Errorf("invalid Token")
	}
}
