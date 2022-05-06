package cmd

import "testing"

func Test_buildVersion(t *testing.T) {
	t.Run("no commit", func(t *testing.T) {
		Version = "0.0.0-next"
		Commit = ""
		v := buildVersion()
		want := Version
		if v != want {
			t.Errorf("buildVersion() got = %v, want %v", v, want)
		}
	})

	t.Run("commit", func(t *testing.T) {
		Version = "0.0.0-next"
		Commit = "123"
		v := buildVersion()
		want := Version + " (" + Commit + ")"
		if v != want {
			t.Errorf("buildVersion() got = %v, want %v", v, want)
		}
	})
}
