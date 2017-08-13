package common

import (
	"DazeProxy/config"
	"DazeProxy/disguise"
	"DazeProxy/encryption"
)

type ProxyUnit struct{
	Disguise disguise.DisguiseAction
	Encryption encryption.EncryptionAction
	Config config.ProxyUnitStruct
}
