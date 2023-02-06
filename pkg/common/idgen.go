package common

import gonanoid "github.com/matoous/go-nanoid/v2"

func GenerateId() string {
	return gonanoid.MustGenerate("qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM123456789", 10)
}
