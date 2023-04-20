package cmd

import "testing"

func Test_buildVersion(t *testing.T) {
	t.Run("no commit", func(t *testing.T) {
		version := "0.0.0-next"
		commit := ""
		v := buildVersion(version, commit)
		want := version
		if v != want {
			t.Errorf("buildVersion() got = %v, want %v", v, want)
		}
	})

	t.Run("commit", func(t *testing.T) {
		version := "0.0.0-next"
		commit := "123"
		v := buildVersion(version, commit)
		want := version + " (" + commit + ")"
		if v != want {
			t.Errorf("buildVersion() got = %v, want %v", v, want)
		}
	})
}
