package cmd

func buildVersion(version, commit string) string {
	if commit != "" {
		if len(commit) > 8 {
			commit = commit[:8]
		}
		version += " (" + commit + ")"
	}
	return version
}
