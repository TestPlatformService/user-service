package token

import (
	"time"
	"user/config"
	pb "user/genproto/user"

	"github.com/dgrijalva/jwt-go"
)

func GeneratedRefreshJWTToken(req *pb.LoginResponse) error {
	conf := config.Load()
	token := *jwt.New(jwt.SigningMethodHS256)
	//payload
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = req.Id
	claims["role"] = req.Role
	claims["iat"] = time.Now().Unix()
	claims["exp"] = time.Now().Add(24 * time.Hour).Unix()

	newToken, err := token.SignedString([]byte(conf.Token.REFRESH_KEY))
	if err != nil {
		return err
	}

	req.Refresh = newToken
	return nil
}

func ValidateRefreshToken(tokenStr string) (bool, error) {
	_, err := ExtractRefreshClaim(tokenStr)
	if err != nil {
		return false, err
	}
	return true, nil
}

func ExtractRefreshClaim(tokenStr string) (*jwt.MapClaims, error) {
	conf := config.Load()
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return []byte(conf.Token.REFRESH_KEY), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !(ok && token.Valid) {
		return nil, err
	}

	return &claims, nil
}

func GetUserIdFromRefreshToken(token string) (string, string, error) {
	conf := config.Load()
	Refresh, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) { return []byte(conf.Token.ACCES_KEY), nil })
	if err != nil || !Refresh.Valid {
		return "", "", err
	}
	claims, ok := Refresh.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", err
	}
	id := claims["user_id"].(string)
	role := claims["role"].(string)

	return id, role, nil
}
