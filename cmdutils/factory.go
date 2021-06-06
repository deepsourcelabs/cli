package cmdutils

type CLI struct {
	Config       ConfigMeta
	TokenExpired bool
}

type ConfigMeta struct {
	Token               string
	RefreshToken        string
	RefreshTokenExpiry  int64
	RefreshTokenSetTime int64
}
