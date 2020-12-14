package common

func GetQiniuConf(env string) (string, string) {
	if env == "test" {
		return "", ""
	} else {
		return "", ""
	}
}

func GetQiniuBucketByPath(path string) string {
	return ""
}

func GetQiniuCallbackUrl() string {
	return ""
}

func GetQiniuUrlByPath(path string) string {
	return ""
}
