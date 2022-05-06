package cmd

var (
	Version = "0.0.0-next"
	Commit  = ""
)

func buildVersion() string {
	result := Version
	if Commit != "" {
		result += " (" + Commit + ")"
	}
	return result
}
