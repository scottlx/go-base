package design_pattern

import (
	"testing"
)

func TestNewFish(t *testing.T) {
	t.Run("testTemplate", func(t *testing.T) {
		fish := NewFish()
		fish.MakeDish()
		chicken := NewChickBreast()
		chicken.MakeDish()
	})
}
