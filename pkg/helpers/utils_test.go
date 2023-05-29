package helpers

import (
	cRand "crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"github.com/lestrrat-go/jwx/jwk"
	env "raptor/config"
	"testing"
)

func TestRandomNumberString(t *testing.T) {
	str := RandomNumberString(10)
	fmt.Println(str)
}

func TestGenerateJwks(t *testing.T) {
	var err error
	PrivateKey, err = rsa.GenerateKey(cRand.Reader, 2048)
	if err != nil {
		panic(err)
	}
	//extract public key
	publicKey := &PrivateKey.PublicKey

	// Create a new JWK for the public key
	jwkKey, err := jwk.New(publicKey)
	if err != nil {
		panic(err)
	}
	jwkKey.Set("kid", env.GetString("oauth2.secret.jwks.kid"))
	jwkKey.Set("alg", "RS256")
	jwkKey.Set("use", "sig")
	jwkSet := jwk.NewSet()

	// Add the JWK to the set
	jwkSet.Add(jwkKey)

	// Convert the JWK set to JSON
	jwkSetJSON, err := json.MarshalIndent(jwkSet, "", "")
	if err != nil {
		panic(err)
	}
	println(jwkSetJSON)
	//fmt.Sprint(jwkSetJSON)

}

func TestKavenegarSmsSender(t *testing.T) {
	env.Init("./config.json")
	KavenegarSmsSender("1234", "09306893359")
}

func TestGenerateTimeStampForFutureBaseOnMinute(t *testing.T) {
	s := GenerateTimeStampForFutureBaseOnMinute(2)
	fmt.Println(s)
}

func TestSendMail(t *testing.T) {
	env.Init("./config_sample.json")
	SendMail("javiddastgoshadeh@gmail.com", "send message from arz", "Test", "test")
}
