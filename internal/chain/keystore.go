package chain

import (
	"crypto/ecdsa"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/keystore"
)

func DecryptKeyfile(keyfile, password string) (*ecdsa.PrivateKey, error) {
	jsonBytes, err := os.ReadFile(keyfile)
	if err != nil {
		return nil, err
	}

	key, err := keystore.DecryptKey(jsonBytes, password)
	if err != nil {
		return nil, err
	}

	return key.PrivateKey, nil
}

func ResolveKeyfilePath(keydir string) (string, error) {
	keydir, _ = filepath.Abs(keydir)
	fileInfo, err := os.Stat(keydir)
	if err != nil {
		return "", err
	}
	if !fileInfo.IsDir() {
		return keydir, nil
	}

	files, _ := os.ReadDir(keydir)
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if strings.HasPrefix(file.Name(), "UTC--") {
			return filepath.Join(keydir, file.Name()), nil
		}
	}

	return "", fmt.Errorf("keyfile is not in %s", keydir)
}
