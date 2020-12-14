package logic

import (
	"context"
	"encoding/json"
	"testing"
)

func TestGetLvlocation(t *testing.T) {
	res, err := GetLvlocation(context.TODO(), 0)
	if err != nil {
		t.Fatal(err)
	}
	b, _ := json.Marshal(res)
	t.Log(string(b))
}
