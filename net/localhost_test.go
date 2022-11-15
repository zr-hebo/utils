package net

import "testing"

func TestLocalIPAddr(t *testing.T) {
	ipAddr, err := LocalIPAddr()
	if err != nil {
		t.Errorf("got local IP add failed <-- %s", err.Error())
		return
	}

	t.Logf("got local IP: %s", ipAddr)
}
