package tls

import (
	"crypto/cipher"
	"errors"
)

type MockBoring struct{
	Enabled bool
}


const Enabled bool = false

func (_ *MockBoring) NewGCMTLS(_ cipher.Block) (cipher.AEAD, error) {
	return nil, errors.New("boring not implemented")
}

func (_ *MockBoring) NewGCMTLS13(_ cipher.Block) (cipher.AEAD, error) {
	return nil, errors.New("boring not implemented")
}

func (_ *MockBoring) Unreachable() {
	// do nothing
}

boring := MockBoring{Enabled: false}