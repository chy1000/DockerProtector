package iptables

import "testing"

const (
	yes = "\u2713"
	no = "\u2717"
)

func TestIptablesCreate(t *testing.T) {
	iptables, err := Create()
	if err != nil {
		t.Fatal("创建 Iptables 对象失败", err, no)
	}
	t.Log("创建 Iptables 对象成功", iptables, yes)
}