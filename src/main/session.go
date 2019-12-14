package main
import(
	"encoding/base64"
	"encoding/json"
	"strings"
)

type Session struct{
	PhoneNumber string `json:phoneNumber`
	Id int64 `json:id`
	Iat int64 `json:iat`
	Exp int64 `json:exp`
}

func parseSession(sessionId string) *Session{
	var sessionBody = (strings.Split(sessionId,"."))[1]
	decoded, err := base64.RawURLEncoding.DecodeString(sessionBody)
	session := new(Session)
	err = json.Unmarshal(decoded, &session)
	if nil !=  err {
		sugar.Debugw("deocode user session failed")
		return &Session{}
	}else{
		return session
	}
}
