package wapiti

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNameRgx(t *testing.T) {
	require.Regexp(t, nameRgx, "Pet")
	require.Regexp(t, nameRgx, "PetOwner")
	require.NotRegexp(t, nameRgx, "pet")
	require.NotRegexp(t, nameRgx, "")
	require.NotRegexp(t, nameRgx, "1pet")
	require.NotRegexp(t, nameRgx, "_pet")
}
