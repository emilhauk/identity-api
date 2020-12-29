package store

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/emilhauk/identity-api/model"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"regexp"
	"strings"
)

type KeyMap map[string]model.RSAKeyPair

type RSAKeyStore struct {
	keyMap         KeyMap
	defaultKeyPair model.RSAKeyPair
	DefaultKeyId   string
}

func NewRSAKeyStore(path string, defaultKeyId string) (store RSAKeyStore) {
	keyMap := loadKeyMap(path)
	validateKeyMap(keyMap)

	defaultKeyPair, ok := keyMap[defaultKeyId]
	if !ok {
		logrus.Fatalln("Invalid default key id %s", defaultKeyId)
		return
	}
	return RSAKeyStore{
		keyMap,
		defaultKeyPair,
		defaultKeyId,
	}
}

func loadKeyMap(path string) KeyMap {
	dir, err := ioutil.ReadDir(path)
	if err != nil {
		logrus.Fatalln("Reading key store path", err)
	}
	reg, _ := regexp.Compile("\\.pub$")
	keyMap := KeyMap{}
	for _, file := range dir {
		keyId := strings.NewReplacer(".pub", "").Replace(file.Name())
		if _, ok := keyMap[keyId]; !ok {
			keyMap[keyId] = model.RSAKeyPair{}
		}
		currentKey := keyMap[keyId]

		data, err := ioutil.ReadFile(path + "/" + file.Name())
		if err != nil {
			logrus.Fatalln("Error reading key store file", err)
		}

		isPub := reg.Find([]byte(file.Name())) != nil
		if isPub {
			publicKey, err := jwt.ParseRSAPublicKeyFromPEM(data)
			if err != nil {
				logrus.Fatalln("Invalid public key for id (%s)", keyId)
			}
			currentKey.Public = publicKey
		} else {
			privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(data)
			if err != nil {
				logrus.Fatalln("Invalid private key for id (%s)", keyId)
			}
			currentKey.Private = privateKey
		}
		keyMap[keyId] = currentKey
	}
	return keyMap
}

func validateKeyMap(keyMap KeyMap) {
	for keyName, keyPair := range keyMap {
		if keyPair.Private == nil || keyPair.Public == nil {
			logrus.Fatalln("Unbalanced public/private key for (%s)", keyName)
		}
	}
}

func (c *RSAKeyStore) GetAllKeyPairs() (keyMap KeyMap) {
	return c.keyMap
}

func (c *RSAKeyStore) GetKeyPairById(id string) (keyPair model.RSAKeyPair, ok bool) {
	keyPair, ok = c.keyMap[id]
	return
}

func (c *RSAKeyStore) GetDefaultKeyPair() (keyPair model.RSAKeyPair) {
	return c.defaultKeyPair
}
