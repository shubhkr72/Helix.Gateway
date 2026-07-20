package jwt

import (
	"crypto/rsa"
	"os"
	"time"

	jwtlib "github.com/golang-jwt/jwt/v5"
)

type Manager struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey

	issuer   string
	audience string
	expiry   time.Duration
}

type Claims struct {
	Email string   `json:"email"`
	Roles []string `json:"roles"`
	jwtlib.RegisteredClaims
}

// NewManager loads the RSA keys.
func NewAuthManager(
	privateKeyPath string,
	publicKeyPath string,
	issuer string,
	audience string,
	expiry time.Duration,
) (*Manager, error) {

	privateBytes, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return nil, err
	}

	publicBytes, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return nil, err
	}

	privateKey, err := jwtlib.ParseRSAPrivateKeyFromPEM(privateBytes)
	if err != nil {
		return nil, err
	}

	publicKey, err := jwtlib.ParseRSAPublicKeyFromPEM(publicBytes)
	if err != nil {
		return nil, err
	}

	return &Manager{
		privateKey: privateKey,
		publicKey:  publicKey,
		issuer:     issuer,
		audience:   audience,
		expiry:     expiry,
	}, nil
}

//gatwey manager
func NewGatewayManager(
	publicKeyPath string,
	issuer string,
	audience string,
) (*Manager, error) {

	publicBytes, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return nil, err
	}

	publicKey, err := jwtlib.ParseRSAPublicKeyFromPEM(publicBytes)
	if err != nil {
		return nil, err
	}

	return &Manager{
		publicKey: publicKey,
		issuer:    issuer,
		audience:  audience,
	}, nil
}

// GenerateToken creates a signed JWT.
func (m *Manager) GenerateToken(
	userID string,
	email string,
	roles []string,
) (string, error) {

	if m.privateKey == nil {
		return "", jwtlib.ErrInvalidKey
	}

	now := time.Now()

	claims := Claims{
		Email: email,
		Roles: roles,
		RegisteredClaims: jwtlib.RegisteredClaims{
			Subject:   userID,
			Issuer:    m.issuer,
			Audience:  []string{m.audience},
			IssuedAt:  jwtlib.NewNumericDate(now),
			NotBefore: jwtlib.NewNumericDate(now),
			ExpiresAt: jwtlib.NewNumericDate(now.Add(m.expiry)),
		},
	}

	token := jwtlib.NewWithClaims(
		jwtlib.SigningMethodRS256,
		claims,
	)

	return token.SignedString(m.privateKey)
}

// VerifyToken validates a JWT.
func (m *Manager) VerifyToken(tokenString string) (*Claims, error) {

	if m.publicKey == nil {
		return nil, jwtlib.ErrInvalidKey
	}
	token, err := jwtlib.ParseWithClaims(
		tokenString,
		&Claims{},
		func(token *jwtlib.Token) (interface{}, error) {
			if token.Method != jwtlib.SigningMethodRS256 {
				return nil, jwtlib.ErrTokenSignatureInvalid
			}

			return m.publicKey, nil
		},
		jwtlib.WithIssuer(m.issuer),
		jwtlib.WithAudience(m.audience),
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, jwtlib.ErrTokenInvalidClaims
	}

	return claims, nil
}
