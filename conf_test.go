package go_conf_test

import (
	"github.com/sohaha/go_conf"
	"github.com/sohaha/zlsgo"
	"os"
	"testing"
)

func TestDef(t *testing.T) {
	tt := zlsgo.NewTest(t)

	c := go_conf.New("zls")
	c.SetDefault("def", 1)
	err := c.Read()
	tt.Equal(true, err == nil)

	c2 := go_conf.New("zls.yaml")
	err = c2.Read()
	tt.Equal(true, err == nil)
	tt.Equal(c.Core.GetInt("def"), c2.Core.GetInt("def"))

	os.Remove("zls.toml")
}
