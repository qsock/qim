package password

import "testing"

func TestObscurePassword(t *testing.T) {
	pwd := ObscurePassword("123456")
	t.Log(pwd)
}

func TestValidatePassword(t *testing.T) {
	pwd := "d4dde13ff21039b3acdbfaf3b4950357e6adeca1c4eea4378944fc533a2b65f23b83eae9cb994d3fb679cdcf1394a15418686680378187100da1d0bbd62d8ee4$sha512$70b502c0-5478-4012-8708-8a50475e9a45"
	if !ValidatePassword("123456", pwd) {
		t.Fatal(pwd)
	}
}
