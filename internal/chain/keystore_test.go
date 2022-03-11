package chain

import (
	"crypto/ecdsa"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
)

func TestDecryptKeyfile(t *testing.T) {
	privateKey, _ := crypto.HexToECDSA("976f9f7772781ff6d1c93941129d417c49a209c674056a3cf5e27e225ee55fa8")
	type args struct {
		keyfile  string
		password string
	}
	tests := []struct {
		name    string
		args    args
		want    *ecdsa.PrivateKey
		wantErr bool
	}{
		{
			name: "empty",
			args: args{
				keyfile:  "testdata/keystore/empty",
				password: "foobar",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "success",
			args: args{
				keyfile:  "testdata/keystore/UTC--2016-03-22T12-57-55.920751759Z--7ef5a6135f1fd6a02593eedc869c6d41d934aef8",
				password: "foobar",
			},
			want:    privateKey,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DecryptKeyfile(tt.args.keyfile, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecryptPrivateKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DecryptPrivateKey() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResolveKeyfilePath(t *testing.T) {
	tests := []struct {
		name    string
		keydir  string
		want    string
		wantErr bool
	}{
		{
			name:    "directory",
			keydir:  "testdata/keystore",
			want:    "UTC--2016-03-22T12-57-55.920751759Z--7ef5a6135f1fd6a02593eedc869c6d41d934aef8",
			wantErr: false,
		},
		{
			name:    "notfound",
			keydir:  "testdata/keystore/null",
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ResolveKeyfilePath(tt.keydir)
			if tt.name != "notfound" {
				got = filepath.Base(got)
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("ResolveKeyfilePath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ResolveKeyfilePath() got = %v, want %v", got, tt.want)
			}
		})
	}
}
