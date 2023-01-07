package rules

import (
	"testing"
)

const (
	yes = "\u2713"
	no = "\u2717"
)

func TestReadYamlConfig(t *testing.T) {
	rules, err := ReadYamlConfig()
	if err != nil {
		t.Fatal("读取 rules.yaml 文件失败", err, no)
	}
	t.Log("读取 rules.yaml 文件成功", rules, yes)
}
