package dto

import "TTCS/src/common/crypto"

type TokenDto struct {
	AccessToken  *crypto.Token
	RefreshToken *crypto.Token
}
