package config

import "testing"

func TestNodeResourceConfigParsing(t *testing.T) {
	yaml := `
Backends:
  Basic:
    Node:
      Resources:
        Cpus: 42
        RamGb: 2.5
        DiskGb: 50.0
  `
	conf := Config{}
	Parse([]byte(yaml), &conf)

	if conf.Backends.Basic.Node.Resources.Cpus != 42 {
		t.Fatal("unexpected cpus")
	}
	if conf.Backends.Basic.Node.Resources.RamGb != 2.5 {
		t.Fatal("unexpected ram")
	}
	if conf.Backends.Basic.Node.Resources.DiskGb != 50.0 {
		t.Fatal("unexpected disk")
	}
}
