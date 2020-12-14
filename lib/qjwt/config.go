package qjwt

type Config struct {
	// 加密key
	Signkey string `toml:"signkey"`
	AesKey  string `toml:"aeskey"`
	AesIv   string `toml:"aesiv"`
}

func (c *Config) Check() bool {
	if len(c.AesKey) != 32 ||
		len(c.AesIv) != 16 ||
		len(c.Signkey) == 0 {
		return false
	}
	return true
}
