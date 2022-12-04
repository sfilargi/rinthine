package coreops

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/rinthine/pkg/model"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

func SecureRandUint32() uint32 {
	buffer := make([]byte, 4)
	n, err := rand.Read(buffer)
	if n != 4 || err != nil {
		panic(err)
	}
	return binary.LittleEndian.Uint32(buffer)
}

func GenId() []byte {
	var start = uint64(time.Date(2020, 02, 02, 02, 02, 02, 0, time.UTC).UnixNano() / 1000)
	var now = uint64(time.Now().UnixNano() / 1000)
	id := make([]byte, 12)
	binary.BigEndian.PutUint64(id[0:], now-start)
	binary.BigEndian.PutUint32(id[8:], SecureRandUint32())
	return id
}

func RandomString(count int) string {
	buf := make([]byte, count)
	_, err := rand.Read(buf)
	if err != nil {
		panic(err)
	}
	return base64.RawURLEncoding.EncodeToString(buf)
}

func VerifyPassword(user *model.User, password string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return false
	}
	return true
}

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

func CreateToken(userId []byte, appId []byte) (string, error) {
	token := model.Token{
		Token:  RandomString(36),
		UserId: userId,
		AppId:  appId,
	}

	err := model.TokenPut(&token)
	if err != nil {
		zap.S().Error(err.Error())
		return "", err
	}

	return token.Token, nil
}

func GetUserFromHandle(handle string) (*model.User, error) {
	user, err := model.UserGetFromHandle(handle)
	if err != nil {
		zap.S().Error(err.Error())
		return nil, err
	}
	return user, nil
}

func GetUser(userId []byte) (*model.User, error) {
	user, err := model.UserGet(userId)
	if err != nil {
		zap.S().Error(err.Error())
		return nil, err
	}
	return user, nil
}

func CreateUser(user *model.User) error {
	return model.UserPut(user)
}

func BearerToken(authstring string) (string, error) {
	ss := strings.SplitN(authstring, " ", 2)
	if len(ss) != 2 {
		return "", fmt.Errorf("Couldn't parse Authorization header %s, %v, %d",
			authstring, ss, len(ss))
	}
	if !strings.EqualFold(ss[0], "Bearer") {
		return "", fmt.Errorf("Not a Bearer Authorization")
	}
	return ss[1], nil
}

func BasicCredentials(authstring string) (string, string, error) {
	ss := strings.SplitN(authstring, " ", 2)
	if len(ss) != 2 {
		return "", "", fmt.Errorf("Couldn't parse Authorization header %s, %v, %d",
			authstring, ss, len(ss))
	}
	if !strings.EqualFold(ss[0], "Basic") {
		return "", "", fmt.Errorf("Not a Basic Authorization")
	}

	pair, err := base64.StdEncoding.DecodeString(ss[1])
	if err != nil {
		return "", "", fmt.Errorf("Failed to decode Basic Authentication username/password")
	}

	up := strings.SplitN(string(pair), ":", 2)
	if len(up) != 2 {
		zap.S().Warn("Couldn't find : in username/password string",
			"authstring", authstring,
			"decoded", string(pair))
		return "", "", fmt.Errorf("Couldn't find : in username/password string")
	}

	return up[0], up[1], nil
}

func GetToken(token string) (*model.Token, error) {
	return model.TokenGet(token)
}

func BearerAuthenticate(token string) (*model.User, error) {
	userToken, err := GetToken(token)
	if userToken == nil || err != nil {
		return nil, err
	}

	user, err := GetUser(userToken.UserId)
	return user, nil
}

func GetApp(appId []byte) (*model.App, error) {
	return model.AppGet(appId)
}

func GetAppFromBase64(encoded string) (*model.App, error) {
	decoded, err := base64.UrlEncoding.DecodeString(encoded)
	if err != nil {
		return nil, err
	}
	return GetApp(decoded)
}

func CreateApp(app *model.App) error {
	return model.AppPut(app)
}

func CreateAppCode(app *model.App, user *model.User) (string, error) {

	code := RandomString(36)

	err := model.OauthCodePut(&model.OauthCode{
		Code:   code,
		AppId:  app.AppId,
		UserId: user.UserId,
	})
	if err != nil {
		return "", err
	}

	return code, nil
}

func GetOauthCode(app *model.App, code string) (*model.OauthCode, error) {
	oauth, err := model.OauthCodeGet(code)
	if err != nil {
		return nil, err
	}

	if bytes.Compare(oauth.AppId, app.AppId) != 0 {
		return nil, fmt.Errorf("Not found")
	}

	return oauth, nil
}

func CreateOauthUsedCode(code string) error {
	err := model.UsedOauthCodePut(&model.UsedOauthCode{
		Code: code,
	})

	return err
}

func CreateAppToken(app *model.App, code string) (string, error) {

	oauth, err := GetOauthCode(app, code)
	if err != nil {
		return "", err
	}

	// Double-check
	if bytes.Compare(oauth.AppId, app.AppId) != 0 {
		return "", fmt.Errorf("unlucky")
	}

	// Guard for race conditions. In worst case, if something fails
	// after this step, app will have to do the dance again
	err = CreateOauthUsedCode(code)
	if err != nil {
		return "", err
	}

	token, err := CreateToken(oauth.UserId, oauth.AppId)
	if err != nil {
		return "", err
	}

	return token, nil
}

func Authorize(appstring, handle string, password string) (*model.App, string, error) {
	app, err := GetAppFromBase64(appstring)
	if app == nil || err != nil {
		return nil, "", err
	}

	user, err := GetUserFromHandle(handle)
	if user == nil || err != nil {
		return nil, "", err
	}

	if !VerifyPassword(user, password) {
		log.Printf("Password error for %s", user)
		return nil, "", fmt.Errorf("Invalid password")
	}

	code, err := CreateAppCode(app, user)
	if code == "" || err != nil {
		return nil, "", fmt.Errorf("Internal Error")
	}

	return app, code, nil
}
