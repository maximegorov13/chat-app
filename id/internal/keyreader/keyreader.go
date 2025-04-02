package keyreader

import (
	"fmt"
	"os"
	"path/filepath"
)

type KeyReader struct {
	BasePath string
}

func NewKeyReader(basePath string) *KeyReader {
	return &KeyReader{
		BasePath: basePath,
	}
}

func (r *KeyReader) ReadPrivateKey(keyID string) ([]byte, error) {
	return os.ReadFile(filepath.Join(r.BasePath, fmt.Sprintf("private_%s.pem", keyID)))
}

func (r *KeyReader) ReadPublicKey(keyID string) ([]byte, error) {
	return os.ReadFile(filepath.Join(r.BasePath, fmt.Sprintf("public_%s.pem", keyID)))
}
