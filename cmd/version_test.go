package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_buildVersion(t *testing.T) {
	t.Run("no commit", func(t *testing.T) {
		version := "0.0.0-next"
		commit := ""
		v := buildVersion(version, commit)
		want := version
		assert.Equal(t, want, v)
	})

	t.Run("commit", func(t *testing.T) {
		version := "0.0.0-next"
		commit := "123"
		v := buildVersion(version, commit)
		want := version + " (" + commit + ")"
		assert.Equal(t, want, v)
	})
}
