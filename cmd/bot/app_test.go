package main

import (
	"testing"

	"go.uber.org/fx"
)

func TestFxApp(t *testing.T) {
	if err := fx.ValidateApp(appOptions()...); err != nil {
		t.Error(err)
	}
}
