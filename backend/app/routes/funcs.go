package routes

import (
	"github.com/out-of-mind/catalog/structures"
	vars "github.com/out-of-mind/catalog/variables"

	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"html/template"
	"math/big"
	"net/http"
	"net/url"
	"strings"
	"time"
	"unicode"
)

func showHTML(w http.ResponseWriter, htmlFile string, data interface{}) {
	tmpl, err := template.ParseFiles(vars.TemplateDir + htmlFile)
	if err != nil {
		vars.Log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 Internal Server Error"))
		return
	}

	tmpl.Execute(w, data)
}

func verifyPassword(s string) bool {
	var (
		hasMinLen = false
		hasLower  = false
		hasNumber = false
	)
	if len(s) >= 8 {
		hasMinLen = true
	}
	for _, char := range s {
		switch {
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		}
	}
	return hasMinLen && hasLower && hasNumber
}

func unescapeUrl(path string) (string, error) {
	unescapedPath, err := url.PathUnescape(path)

	return unescapedPath, err
}
func validateAndParseJWT(jwt string) (structures.JWT, error) {
	jwtPieces := strings.Split(jwt, ".")
	header := jwtPieces[0]
	payload := jwtPieces[1]
	signature := jwtPieces[2]

	signatureBytes, err := base64.RawURLEncoding.DecodeString(signature)
	if err != nil {
		return structures.JWT{}, err
	}
	message := header + "." + payload

	if validMAC([]byte(message), signatureBytes, []byte(vars.Secret)) {
		headerBytes, err := base64.RawURLEncoding.DecodeString(header)
		if err != nil {
			return structures.JWT{}, err
		}

		payloadBytes, err := base64.RawURLEncoding.DecodeString(payload)
		if err != nil {
			return structures.JWT{}, err
		}

		var (
			jwt          structures.JWT
			jwtHeader    structures.JWTHeader
			jwtPayload   structures.JWTPayload
			jwtSignature structures.JWTSignature
		)

		err = json.Unmarshal(headerBytes, &jwtHeader)
		if err != nil {
			return structures.JWT{}, err
		}

		err = json.Unmarshal(payloadBytes, &jwtPayload)
		if err != nil {
			return structures.JWT{}, err
		}

		if jwtPayload.Exp.Sub(time.Now()).Minutes() <= 0 {
			return structures.JWT{}, errors.New("jwt: jwt token expired")
		}

		jwtSignature.Hash = string(signatureBytes)
		jwt.Header = jwtHeader
		jwt.Payload = jwtPayload

		jwt.Signature = jwtSignature

		return jwt, nil
	} else {
		return structures.JWT{}, errors.New("jwt: signatures isn't matched")
	}
}
func validMAC(message, messageMAC, key []byte) bool {
	mac := hmac.New(sha256.New, key)
	mac.Write(message)
	expectedMAC := mac.Sum(nil)

	return hmac.Equal(messageMAC, expectedMAC)
}
func newJWT(jwt structures.JWT) (string, error) {
	header, err := json.Marshal(jwt.Header)
	if err != nil {
		return "", err
	}
	payload, err := json.Marshal(jwt.Payload)
	if err != nil {
		return "", err
	}

	headerStr := base64.RawURLEncoding.EncodeToString(header)
	payloadStr := base64.RawURLEncoding.EncodeToString(payload)

	message := headerStr + "." + payloadStr

	mac := hmac.New(sha256.New, []byte(vars.Secret))
	mac.Write([]byte(message))
	signature := mac.Sum(nil)

	signatureStr := base64.RawURLEncoding.EncodeToString(signature)

	JWT := message + "." + signatureStr

	return JWT, nil
}

func randomLink() (string, error) {
	chars := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	link := make([]string, 32)
	charsArr := strings.Split(chars, "")

	for i := 0; i < 32; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(32)))
		if err != nil {
			return "", err
		}
		link[i] = charsArr[n.Int64()]
	}

	return strings.Join(link[:], ""), nil
}
