package common

import (
	"github.com/crabkun/DazeProxy/obscure"
	"github.com/crabkun/DazeProxy/encryption"
)

type S_config struct{
	Debug bool
	Proxy []S_proxy
}
type S_proxy struct{
	Port string

	//加密方式与参数
	Encryption string
	EncryptionParam string

	//混淆方式与参数
	Obscure string
	ObscureParam string

	//加密与混淆的接口
	Ob obscure.Action `json:"-"`
	E encryption.Action `json:"-"`
	EReserved interface{} `json:"-"`
	ObReserved interface{} `json:"-"`
}