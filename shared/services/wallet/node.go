package wallet

import (
	"crypto/ecdsa"
	"errors"
	"fmt"

	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
)

// Config
const NodeKeyPath = "m/44'/60'/0'/0/%d"


// Get the node account
func (w *Wallet) GetNodeAccount() (accounts.Account, error) {

    // Check wallet is initialized
    if !w.IsInitialized() {
        return accounts.Account{}, errors.New("Wallet is not initialized")
    }

    // Get private key
    privateKey, path, err := w.getNodePrivateKey()
    if err != nil {
        return accounts.Account{}, err
    }

    // Get public key
    publicKey := privateKey.Public()
    publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
    if !ok {
        return accounts.Account{}, errors.New("Could not get node public key")
    }

    // Create & return account
    return accounts.Account{
        Address: crypto.PubkeyToAddress(*publicKeyECDSA),
        URL: accounts.URL{
            Scheme: "",
            Path: path,
        },
    }, nil

}


// Get a transactor for the node account
func (w *Wallet) GetNodeAccountTransactor() (*bind.TransactOpts, error) {

    // Check wallet is initialized
    if !w.IsInitialized() {
        return nil, errors.New("Wallet is not initialized")
    }

    // Get private key
    privateKey, _, err := w.getNodePrivateKey()
    if err != nil {
        return nil, err
    }

    // Create & return transactor
    transactor, err := bind.NewKeyedTransactorWithChainID(privateKey, w.chainID)
    transactor.GasPrice = w.gasPrice
    transactor.GasLimit = w.gasLimit
    return transactor, err

}


// Get the node account private key bytes
func (w *Wallet) GetNodePrivateKeyBytes() ([]byte, error) {

    // Check wallet is initialized
    if !w.IsInitialized() {
        return nil, errors.New("Wallet is not initialized")
    }

    // Get private key
    privateKey, _, err := w.getNodePrivateKey()
    if err != nil {
        return nil, err
    }

    // Return private key bytes
    return crypto.FromECDSA(privateKey), nil

}


// Get the node private key
func (w *Wallet) getNodePrivateKey() (*ecdsa.PrivateKey, string, error) {

    // Check for cached node key
    if w.nodeKey != nil {
        return w.nodeKey, w.nodeKeyPath, nil
    }

    // Get derived key
    derivedKey, path, err := w.getNodeDerivedKey(0)
    if err != nil {
        return nil, "", err
    }

    // Get private key
    privateKey, err := derivedKey.ECPrivKey()
    if err != nil {
        return nil, "", fmt.Errorf("Could not get node private key: %w", err)
    }
    privateKeyECDSA := privateKey.ToECDSA()

    // Cache node key
    w.nodeKey = privateKeyECDSA
    w.nodeKeyPath = path

    // Return
    return privateKeyECDSA, path, nil

}


// Get the derived key & derivation path for the node account at the index
func (w *Wallet) getNodeDerivedKey(index uint) (*hdkeychain.ExtendedKey, string, error) {

    // Get derivation path
    derivationPath := fmt.Sprintf(NodeKeyPath, index)

    // Parse derivation path
    path, err := accounts.ParseDerivationPath(derivationPath)
    if err != nil {
        return nil, "", fmt.Errorf("Invalid node key derivation path '%s': %w", derivationPath, err)
    }

    // Follow derivation path
    key := w.mk
    for i, n := range path {
        key, err = key.Child(n)
        if err == hdkeychain.ErrInvalidChild {
            return w.getNodeDerivedKey(index + 1)
        } else if err != nil {
            return nil, "", fmt.Errorf("Invalid child key at depth %d: %w", i, err)
        }
    }

    // Return
    return key, derivationPath, nil

}

