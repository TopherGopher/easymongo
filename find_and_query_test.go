package easymongo_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindAnd(t *testing.T) {
	setup(t)
	coll := createBatmanArchive(t)
	t.Run("FindAnd().Delete()", func(t *testing.T) {
		is := assert.New(t)
		_ = is
		_ = coll
		t.Skipf("TODO: FindAnd().Delete()")
	})
	t.Run("FindAnd().Replace()", func(t *testing.T) {
		is := assert.New(t)
		_ = is
		_ = coll
		t.Skipf("TODO: FindAnd().Replace()")
	})
	t.Run("FindAnd().Delete()", func(t *testing.T) {
		is := assert.New(t)
		_ = is
		_ = coll
		t.Skipf("TODO: FindAnd().Delete()")
	})
}
