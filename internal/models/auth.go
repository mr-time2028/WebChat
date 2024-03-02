package models

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/mr-time2028/WebChat/internal/helpers"
	"net/http"
	"strings"
	"time"
)

var (
	ErrNoAuthHeader      = errors.New(`no auth header`)
	ErrInvalidAuthHeader = errors.New("invalid auth header")
	ErrInvalidToken      = errors.New("invalid token")
)

type Auth struct {
	Issuer        string
	Audience      string
	Secret        string
	TokenExpiry   time.Duration
	RefreshExpiry time.Duration
}

type JwtUser struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

type TokenPairs struct {
	Token        string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type Claims struct {
	jwt.RegisteredClaims
}

func NewJWTAuth() *Auth {
	issuer := helpers.GetEnvOrDefaultString("ISSUER", "localhost")
	audience := helpers.GetEnvOrDefaultString("AUDIENCE", "localhost")
	secret := helpers.GetEnvOrDefaultString("SECRET", "")
	tokenExpiry := helpers.GetEnvOrDefaultInt("TOKEN_EXPIRY", 5)
	refreshExpiry := helpers.GetEnvOrDefaultInt("REFRESH_EXPIRY", 60)

	return &Auth{
		Issuer:        issuer,
		Audience:      audience,
		Secret:        secret,
		TokenExpiry:   time.Duration(tokenExpiry) * time.Minute,
		RefreshExpiry: time.Duration(refreshExpiry) * time.Minute,
	}
}

func (a *Auth) GenerateTokenPair(user *JwtUser) (TokenPairs, error) {
	// create a token
	token := jwt.New(jwt.SigningMethodHS256)

	// set the claims
	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = fmt.Sprintf("%s", user.Username)
	claims["sub"] = fmt.Sprint(user.ID)
	claims["aud"] = a.Audience
	claims["iss"] = a.Issuer
	claims["iat"] = time.Now().UTC().Unix()
	claims["typ"] = "JWT"

	// set the expiry for JWT
	claims["exp"] = time.Now().UTC().Add(a.TokenExpiry).Unix()

	// create a signed token
	signedAccessToken, err := token.SignedString([]byte(a.Secret))
	if err != nil {
		return TokenPairs{}, err
	}

	// create refresh token and set claims
	refreshToken := jwt.New(jwt.SigningMethodHS256)
	refreshTokenClaims := refreshToken.Claims.(jwt.MapClaims)
	refreshTokenClaims["sub"] = fmt.Sprint(user.ID)
	refreshTokenClaims["iat"] = time.Now().UTC().Unix()

	// set the expiry for refresh token
	refreshTokenClaims["exp"] = time.Now().UTC().Add(a.RefreshExpiry).Unix()

	// create signed refresh token
	signedRefreshToken, err := refreshToken.SignedString([]byte(a.Secret))
	if err != nil {
		return TokenPairs{}, err
	}

	// create token pairs and populate with signed tokens
	var tokenPairs = TokenPairs{
		Token:        signedAccessToken,
		RefreshToken: signedRefreshToken,
	}

	// return token pairs
	return tokenPairs, nil
}

func (a *Auth) GetAuthTokenFromHeader(w http.ResponseWriter, r *http.Request) (string, error) {
	w.Header().Add("Vary", "Authorization")

	// get auth header
	authToken := r.Header.Get("Authorization")

	// sanity check
	if authToken == "" {
		return "", ErrNoAuthHeader
	}

	return authToken, nil
}

func (a *Auth) VerifyAuthToken(authToken string) (*Claims, error) {
	// split the header on spaces
	headerParts := strings.Split(authToken, " ")
	if len(headerParts) != 2 {
		return nil, ErrInvalidAuthHeader
	}

	// check if we have the word "Bearer"
	if headerParts[0] != "Bearer" {
		return nil, ErrInvalidAuthHeader
	}
	token := headerParts[1]

	// declare an empty claims
	claims := &Claims{}

	// parse the token
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrTokenSignatureInvalid
		}
		return []byte(a.Secret), nil
	})

	if err != nil {
		if strings.HasPrefix(err.Error(), "token is expired by") {
			return nil, jwt.ErrTokenExpired
		}
		return nil, ErrInvalidToken
	}

	if claims.Issuer != a.Issuer {
		return nil, jwt.ErrTokenInvalidIssuer
	}

	// check expiration time
	if claims.ExpiresAt.Before(time.Now()) {
		return nil, jwt.ErrTokenExpired
	}

	return claims, nil
}
