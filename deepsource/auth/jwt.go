package auth

type JWT struct {
	Payload struct {
		Email   string `json:"email"`
		Exp     string `json:"exp"`
		Origiat int64  `json:"origIat"`
	} `json:"payload"`
	Token            string `json:"token"`
	Refreshtoken     string `json:"refreshToken"`
	TokenExpiresIn   int64  `json:"tokenExpiresIn"`
	RefreshExpiresIn int64  `json:"refreshExpiresIn"`
}
