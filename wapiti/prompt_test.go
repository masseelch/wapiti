package wapiti

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRgx(t *testing.T) {
	require.Regexp(t, nodeNameRgx, "Pet")
	require.Regexp(t, nodeNameRgx, "PetOwner")
	require.NotRegexp(t, nodeNameRgx, "pet")
	require.NotRegexp(t, nodeNameRgx, "")
	require.NotRegexp(t, nodeNameRgx, "1pet")
	require.NotRegexp(t, nodeNameRgx, "_pet")

	require.Regexp(t, fieldNameRgx, "age")
	require.Regexp(t, fieldNameRgx, "fur_color")
	require.Regexp(t, fieldNameRgx, "fur-color")
	require.Regexp(t, fieldNameRgx, "fur_color_1")
	require.NotRegexp(t, fieldNameRgx, "_")
	require.NotRegexp(t, fieldNameRgx, "-")
	require.NotRegexp(t, fieldNameRgx, "1")
	require.NotRegexp(t, fieldNameRgx, "age?")
	require.NotRegexp(t, fieldNameRgx, "age_%")
}
