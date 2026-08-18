package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"strings"
)

type Role string

const (
	PlatformAdmin Role = "PLATFORM_ADMIN"
	TenantAdmin   Role = "TENANT_ADMIN"
	Operator      Role = "OPERATOR"
	Viewer        Role = "VIEWER"
)

type DevicePrincipal struct{ TenantID, DeviceID, CredentialID, DeviceType, ProjectID string }
type OperatorPrincipal struct {
	APIKeyID, TenantID string
	Role               Role
}

var ErrInvalidCredential = errors.New("invalid credential")

func NewID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	h := hex.EncodeToString(b)
	return h[0:8] + "-" + h[8:12] + "-" + h[12:16] + "-" + h[16:20] + "-" + h[20:32]
}

func GenerateToken(kind string) (raw, prefix string, hash []byte, err error) {
	public := make([]byte, 8)
	secret := make([]byte, 32)
	if _, err = rand.Read(public); err != nil {
		return
	}
	if _, err = rand.Read(secret); err != nil {
		return
	}
	prefix = hex.EncodeToString(public)
	raw = "pol_" + kind + "_" + prefix + "." + hex.EncodeToString(secret)
	sum := sha256.Sum256([]byte(raw))
	hash = sum[:]
	return
}
func TokenPrefix(raw string) (string, error) {
	parts := strings.Split(raw, ".")
	if len(parts) != 2 {
		return "", ErrInvalidCredential
	}
	head := parts[0]
	i := strings.LastIndex(head, "_")
	if i < 0 || i == len(head)-1 {
		return "", ErrInvalidCredential
	}
	return head[i+1:], nil
}
func Verify(raw string, expected []byte) bool {
	sum := sha256.Sum256([]byte(raw))
	return len(expected) == len(sum) && subtle.ConstantTimeCompare(sum[:], expected) == 1
}
func Hash(raw string) []byte { sum := sha256.Sum256([]byte(raw)); return sum[:] }

func Can(role Role, permission string) bool {
	if role == PlatformAdmin {
		return true
	}
	switch permission {
	case "read":
		return role == TenantAdmin || role == Operator || role == Viewer
	case "mutate":
		return role == TenantAdmin
	case "audit":
		return role == TenantAdmin
	default:
		return false
	}
}
