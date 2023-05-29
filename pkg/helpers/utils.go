package helpers

import (
	"bytes"
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/kavenegar/kavenegar-go"
	"github.com/labstack/echo/v4"
	"github.com/machinebox/graphql"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/mail"
	//"net/smtp"
	env "raptor/config"
	"raptor/logger"
	"raptor/models"
	"regexp"
	"strings"
	"time"

	"gopkg.in/gomail.v2"
)

// RandomArray ...
func RandomArray(len int) []int {
	rand.Seed(time.Now().UnixNano())

	a := make([]int, len)
	for i := 0; i <= len-1; i++ {
		a[i] = rand.Intn(len)
	}
	return a
}

// ArrayToString ...
func ArrayToString(arr []int, dlm string) string {
	return strings.Trim(strings.Join(strings.Split(fmt.Sprint(arr), " "), dlm), "[]")
}

// RandomInt Returns an int >= min, < max
func RandomInt(min, max int) int {
	return min + rand.Intn(max-min)
}

// RandomString Generate a random string of A-Z chars with len = l
func RandomString(len int) string {
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		bytes[i] = byte(RandomInt(65, 90))
	}
	return string(bytes)
}

// RandomNumberString Generate a random string of A-Z chars with len = l
func RandomNumberString(len int) string {
	return ArrayToString(RandomArray(len), "")
}

// GenerateRandomString ...
func GenerateRandomString(n int) string {
	rand.Seed(time.Now().UnixNano())

	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}

	return string(b)
}

// GenerateCode ...
func GenerateCode() string {
	rand.Seed(time.Now().UnixNano())

	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	b := make([]byte, 20)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}

	return string(b)
}

// GenerateRandomInt ...
func GenerateRandomInt() string {
	return fmt.Sprintf("%v", rand.Intn(90000)+10000)
}

func SendMail(to string, sub string, msg string, templateName string) {

	username := env.GetString("courier.smtp.username")
	password := env.GetString("courier.smtp.password")
	host := env.GetString("courier.smtp.host")
	port := env.GetInt("courier.smtp.port")
	from := env.GetString("courier.smtp.from")

	message := gomail.NewMessage()
	message.SetHeader("From", from)
	message.SetHeader("To", to)
	message.SetHeader("Subject", sub)
	message.SetBody("text/plain", msg)

	// Configure the SMTP client to use Gmail
	dialer := gomail.NewDialer(host, port, username, password)

	// Send the email
	err := dialer.DialAndSend(message)
	if err != nil {
		panic(err)
	}
	fmt.Println("Email sent!")

	//
	//messageBody := []byte("test")
	//
	//username := env.GetString("courier.smtp.username")
	//password := env.GetString("courier.smtp.password")
	//host := env.GetString("courier.smtp.host")
	//auth := smtp.PlainAuth("", username, password, host)
	//address := env.GetString("courier.smtp.address")
	//from := env.GetString("courier.smtp.from")
	//destinationMail := []string{to}
	//err = smtp.SendMail(address, auth, from, destinationMail, messageBody)
	//if err != nil {
	//	log.Println(err)
	//}

}

func KavenegarSmsSender(code string, identity string) bool {

	api := kavenegar.New(env.GetString("courier.sms.kavenegar.api_key"))
	logger.Info("otp sms send with code: " + code + " to: " + identity)
	if res, err := api.Verify.Lookup(identity, env.GetString("courier.sms.kavenegar.pattern"), code, nil); err != nil {
		switch err := err.(type) {
		case *kavenegar.APIError:
			logger.Error(models.ErrKavenegarApi.Error() + err.Error())
		case *kavenegar.HTTPError:
			logger.Error(models.ErrKavenegarHttp.Error() + err.Error())
		default:
			logger.Error(models.ErrInternalServer.Error() + err.Error())
		}
		return false
	} else {
		logger.Debug("otp sms send with code: " + code + " to: " + identity + " by id: " + fmt.Sprintf("%v", res.MessageID) + " and status: " + fmt.Sprintf("%v", res.Status))
	}
	return true
}

func GenerateIdentityKeyToCacheOtpCode(str string) string {
	return env.GetString("name") + "_code_" + str
}

func GenerateIdentityKeyToCacheFlowID(str string) string {
	return env.GetString("name") + "_flow_" + str
}

func GenerateIdentityKeyToCacheState(str string) string {
	return env.GetString("name") + "_state_" + str
}

func GenerateIdentityKeyToCacheTimeForOtpCode(str string) string {
	return env.GetString("name") + "_time_" + str
}

// GenerateTimeStampForFutureBaseOnMinute
// duration is int base on minutes number
func GenerateTimeStampForFutureBaseOnMinute(duration int) int64 {
	now := time.Now()
	return now.Add(time.Duration(duration) * time.Minute).Unix()
}

func GenerateCurrentTimeStamp() int64 {
	now := time.Now().Unix()
	return now
}

func CheckFormat(str string) (string, error) {

	email := models.Email
	sms := models.Sms
	//check email format
	ok := CheckEmailFormat(str)
	if ok {
		return string(email), nil
	}
	if CheckMobileFormat(str) {
		return string(sms), nil
	}
	logger.Error(models.ErrIdentityFormat.Error())
	return "", models.ErrIdentityFormat
}

func CheckEmailFormat(identity string) bool {
	_, err := mail.ParseAddress(identity)
	if err != nil {
		return false
	}
	return true
}

func CheckMobileFormat(identity string) bool {
	pattern := "^09\\d{9}$"
	matched, err := regexp.MatchString(pattern, identity)
	if err != nil || !matched {
		return false
	}
	return true
}

//func MobileFormatRegex(str string) bool {
//	re := regexp.MustCompile(`^(?:(?:\(?(?:00|\+)([1-4]\d\d|[1-9]\d?)\)?)?[\-\.\ \\\/]?)?((?:\(?\d{1,}\)?[\-\.\ \\\/]?){0,})(?:[\-\.\ \\\/]?(?:#|ext\.?|extension|x)[\-\.\ \\\/]?(\d+))?$`)
//	return re.MatchString(str)
//}

func SendHttpRequest(url string, method string, data []byte, header map[string]string) (interface{}, error) {
	var result map[string]interface{}
	request, err := http.NewRequest(method, url, bytes.NewBuffer(data))
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	if header != nil {
		for k, v := range header {
			request.Header.Set(k, v)

		}
	}
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		logger.Error(err)
		return nil, models.ErrHttpRequest
	}

	if response.StatusCode == 204 || response.StatusCode == 201 {
		return nil, nil
	}
	if err != nil {
		logger.Error(err)
		return nil, models.ErrHttpRequest
	}

	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)
	err = json.Unmarshal(body, &result)
	if err != nil {
		logger.Error(err)
		return nil, models.ErrGeneral
	}

	return result, nil
}

func SendGraphqlRequest(url string, query string, vars map[string]interface{}) (interface{}, error) {
	graphqlClient := graphql.NewClient(url)
	graphqlRequest := graphql.NewRequest(query)
	// set any variables
	for k, v := range vars {
		graphqlRequest.Var(k, v)
	}

	var graphqlResponse interface{}
	if err := graphqlClient.Run(context.Background(), graphqlRequest, &graphqlResponse); err != nil {
		logger.Error(err)
		return nil, models.ErrHttpRequest
	}

	return graphqlResponse, nil
}

func GenerateJwtToken(privateKet *rsa.PrivateKey, claims map[string]interface{}) (string, error) {
	var customPayload jwt.MapClaims
	customPayload = claims
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, customPayload)
	tokenString, err := token.SignedString(privateKet)
	if err != nil {
		return "", models.ErrInvalidToken
	}
	return tokenString, nil
}

func LoadPrivateKey(path string) *rsa.PrivateKey {
	keyBytes, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("Failed to read private key file: %v", err)
	}
	// Parse the PEM block
	block, _ := pem.Decode(keyBytes)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		log.Fatalf("Failed to decode private key")
	}
	// Parse the DER-encoded private key
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		panic("Failed to parse private key")
	}
	return privateKey
}

func LoadPublicKey(path string) *rsa.PublicKey {
	type CustomClaims struct {
		jwt.StandardClaims
	}

	// Read public key file
	keyBytes, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	// Parse public key
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(keyBytes)
	if err != nil {
		panic(err)
	}
	return publicKey
}

func GetJwtStringFromRequest(r *http.Request) (string, error) {
	headerJwt := r.Header.Get("Authorization")
	if headerJwt == "" {
		return "", models.ErrUnauthorized
	}
	headerJwtArr := strings.Split(headerJwt, "Bearer ")
	if len(headerJwtArr) != 2 {
		return "", models.ErrAuthorizedFormat
	}
	return headerJwtArr[1], nil
}

func ExtractClaimsFromToken(c echo.Context) (map[string]interface{}, error) {
	tokenString, err := GetJwtStringFromRequest(c.Request())
	if err != nil {
		return nil, models.ErrEmptyJwtToken
	}
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte{}, nil
	})

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, models.ErrJwtToken
	}
	return claims, nil
}

func ExtractClaimsFromTokenString(tokenString string) (map[string]interface{}, error) {
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte{}, nil
	})

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		logger.Error(models.ErrJwtToken)
		return nil, models.ErrJwtToken
	}
	return claims, nil
}

func GenerateAccessTokenClaims(res interface{}) (map[string]interface{}, error) {

	var (
		identity map[string]interface{}
		metadata map[string]interface{}
	)

	// fetch data from response
	data, _ := res.(map[string]interface{})
	//check error from response
	if data["error"] != nil {
		return nil, models.ErrSubmitLoginFlow
	}

	//fetch needed attributes from response
	sessionToken, ok := data["session_token"].(string)
	if !ok {
		return nil, models.ErrSubmitLoginFlow
	}
	session, _ := data["session"].(map[string]interface{})
	if session == nil {
		identity = data
	}
	if session != nil {
		identity, _ = session["identity"].(map[string]interface{})
	}

	sub, _ := identity["id"]
	identityTraits, _ := identity["traits"].(map[string]interface{})

	avatars := env.GetStringSlice("identity.public_avatars")
	avatarsLen := len(avatars)

	var (
		displayName = env.GetString("identity.default_display_name")
		avatar      = avatars[rand.Intn(avatarsLen)]
		uuid        = sub
	)

	if identityTraits["avatar"] != nil {
		avatar = fmt.Sprintf("%v", identityTraits["avatar"])
	}
	if identity["metadata_public"] != nil {
		metadata = identity["metadata_public"].(map[string]interface{})
		if metadata["user_id"] != nil {
			uuid = metadata["user_id"]
		}

	}

	//fill data and put uuid in iy
	dataInClaims := make(map[string]interface{})
	user := make(map[string]interface{})
	user["id"] = uuid
	dataInClaims["user"] = user

	if identityTraits["displayName"] != nil {
		displayName = fmt.Sprintf("%v", identityTraits["displayName"])
	}

	//create claims to generate token

	aud := env.GetStringSlice("security.jwt.aud")
	//set claims for access_token
	claims := make(map[string]interface{})
	claims["iss"] = env.GetString("security.jwt.iss")
	claims["aud"] = aud
	claims["sub"] = sub
	claims["uuid"] = uuid
	claims["data"] = dataInClaims
	claims["sid"] = sessionToken
	claims["display_name"] = displayName
	claims["avatar"] = avatar
	//claims["metadata"] = metadata
	claims["iat"] = GenerateCurrentTimeStamp()
	claims["exp"] = GenerateTimeStampForFutureBaseOnMinute(env.GetInt("security.jwt.access_token.expire_time"))

	return claims, nil
}

// GenerateRandInt ...
func GenerateRandInt() string {
	return fmt.Sprintf("%v", rand.Intn(90000)+10000)
}

func GenerateRefreshTokenClaims(res interface{}) (map[string]interface{}, error) {

	// fetch data from response
	data, _ := res.(map[string]interface{})
	//check error from response
	if data["error"] != nil {
		return nil, models.ErrExpireCodeTime
	}
	aud := env.GetStringSlice("security.jwt.aud")
	//fetch needed attributes from response
	session, _ := data["session"].(map[string]interface{})
	identity, _ := session["identity"].(map[string]interface{})
	sub, _ := identity["id"]
	//set claims for access_token
	claims := make(map[string]interface{})
	claims["iss"] = env.GetString("security.jwt.iss")
	claims["aud"] = aud
	claims["sub"] = sub
	//claims["traits"] = identityTraits
	claims["iat"] = GenerateCurrentTimeStamp()
	claims["exp"] = GenerateTimeStampForFutureBaseOnMinute(env.GetInt("security.jwt.refresh_token.expire_time"))

	return claims, nil
}

func ExportKratosFlowsErr(response interface{}) error {
	res, _ := response.(map[string]interface{})
	if res["ui"] != nil {
		ui := res["ui"].(map[string]interface{})
		if ui["messages"] != nil {
			message := ui["messages"].([]interface{})
			messageInfo := message[0].(map[string]interface{})
			if messageInfo["type"] != nil {
				if messageInfo["type"] == "error" {
					logger.Error(fmt.Sprintf("%v", messageInfo["text"]))
					return errors.New(fmt.Sprintf("%v", messageInfo["text"]))
				}
			}
		}
	}
	return nil
}

func GetTraitsFromKratosResponse(identityInfo interface{}) (*models.Traits, error) {

	var (
		name map[string]interface{}
		//verification map[string]interface{}
	)

	info := identityInfo.(map[string]interface{})
	identityTraits, ok := info["traits"].(map[string]interface{})
	if !ok {
		return nil, models.ErrTraitsFormat
	}
	if identityTraits["name"] != nil {
		name, ok = identityTraits["name"].(map[string]interface{})
		if !ok {
			return nil, models.ErrTraitsFormat
		}
	}
	//if identityTraits["verification"] != nil {
	//	verification, ok = identityTraits["verification"].(map[string]interface{})
	//	if !ok {
	//		return nil, models.ErrTraitsFormat
	//	}
	//}

	traits := &models.Traits{
		Email:               identityTraits["email"],
		PhoneNumber:         identityTraits["phone_number"],
		Username:            identityTraits["username"],
		EmailVerified:       identityTraits["email_verified"],
		PhoneNumberVerified: identityTraits["phone_number_verified"],
		Name: &models.Name{
			First: name["first"],
			Last:  name["last"],
		},
		DisplayName: identityTraits["display_name"],
		Nickname:    identityTraits["nickname"],
		Avatar:      identityTraits["avatar"],
		//Verification: &models.Verification{
		//	Birthdate:    verification["birthdate"],
		//	NationalCode: verification["national_code"],
		//	DisplayName:  verification["display_name"],
		//	State:        verification["state"],
		//},
	}
	return traits, nil
}

func InjectDataToIdentity(data interface{}, inject map[string]interface{}) (interface{}, error) {

	sessionInfo, ok := data.(map[string]interface{})
	if !ok {
		logger.Error(models.ErrCreatingRegisterFlow.Error())
		return nil, models.ErrIdentityFormat
	}
	metaData := make(map[string]interface{})
	metaData["status"] = "active"
	var identityInfo map[string]interface{}

	if sessionInfo["session"] != nil {
		dataInfo := sessionInfo["session"].(map[string]interface{})
		identityInfo = dataInfo["identity"].(map[string]interface{})
	} else {
		identityInfo = sessionInfo
	}

	if identityInfo["metadata_public"] == nil {
		identityInfo["metadata_public"] = metaData
	} else {
		pMetadata := identityInfo["metadata_public"].(map[string]interface{})
		pMetadata["status"] = "active"
	}
	//extract traits
	userTraits, ok := identityInfo["traits"].(map[string]interface{})
	if !ok {
		logger.Error(models.ErrCreatingRegisterFlow.Error())
		return nil, models.ErrIdentityFormat
	}
	identityInfo["credentials"] = nil
	for k, v := range inject {
		userTraits[k] = v
	}
	return identityInfo, nil

}
func InjectVerificationDataToMeta(data interface{}, inject map[string]interface{}) (interface{}, error) {
	sessionInfo, ok := data.(map[string]interface{})
	if !ok {
		return nil, models.ErrCreatingRegisterFlow
	}
	metaData := make(map[string]interface{})
	metaData["status"] = "active"
	var identityInfo map[string]interface{}

	if sessionInfo["session"] != nil {
		dataInfo := sessionInfo["session"].(map[string]interface{})
		identityInfo = dataInfo["identity"].(map[string]interface{})
	} else {
		identityInfo = sessionInfo
	}
	publicMeta := identityInfo["metadata_public"].(map[string]interface{})
	publicMeta["verification"] = inject
	return identityInfo, nil

}

func ExportIdentityFromKratosResponse(req interface{}) (interface{}, error) {
	data, ok := req.(map[string]interface{})

	if !ok {
		logger.Error(models.ErrIdentityFormat.Error())
		return nil, models.ErrIdentityFormat
	}
	identity, ok := data["identity"].(map[string]interface{})
	if ok {
		return identity["traits"], nil
	}
	traits, ok := data["traits"].(map[string]interface{})
	if ok {
		return traits, nil
	}

	logger.Error(models.ErrIdentityFormat.Error())
	return nil, models.ErrIdentityFormat
}

func UploadImage(jwt string, imageUrl string) error {
	var bearer = "Bearer " + jwt
	header := make(map[string]string)
	header["Authorization"] = bearer
	body := make(map[string]interface{})
	body["image"] = imageUrl
	body["size"] = env.GetString("apis.uploader.config.size")
	payload, err := json.Marshal(body)
	if err != nil {
		logger.Error(err.Error())
	}
	res, err := SendHttpRequest(env.GetString("apis.uploader.url"), env.GetString("apis.uploader.method"), payload, header)
	result := res.(map[string]interface{})
	if result["status"] == "fail" {
		return models.ErrUploadImage
	}
	return err
}
