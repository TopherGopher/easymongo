package easymongo_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIndex(t *testing.T) {
	setup(t)
	t.Cleanup(func() { teardown(t) })
	t.Run("Ensure", func(t *testing.T) {
		is := assert.New(t)
		_ = is
		t.Skipf("TODO: Add Index.Ensure() test")
	})
	t.Run("List", func(t *testing.T) {
		is := assert.New(t)
		_ = is
		t.Skipf("TODO: Add Index.List() test")
	})
	t.Run("Drop", func(t *testing.T) {
		is := assert.New(t)
		_ = is
		t.Skipf("TODO: Add Index.Drop() test")
	})
}
