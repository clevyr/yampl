package processor

import (
	"strings"
	"testing"

	"github.com/clevyr/yampl/internal/config"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

// Test_templateReader_differential renders an exhaustive corpus of template
// values through the real pipeline and verifies that yaml.v3 — standing in
// for the wider YAML parser ecosystem — reads every output back as the exact
// string that was templated. This guards the manual quoting rules in
// internal/visitor (needsDoubleQuote, isAmbiguousPlain) against gaps and
// against behavior changes in future goccy versions.
func Test_templateReader_differential(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping differential sweep in short mode")
	}

	check := func(t *testing.T, s string) {
		t.Helper()

		conf := config.New()
		conf.Vars = config.Vars{"newVal": s}

		got, err := templateReader(conf, "", strings.NewReader(`v: "" #yampl {{ .newVal }}`))
		require.NoError(t, err, "input %q", s)

		var decoded map[string]string
		require.NoError(t, yaml.Unmarshal([]byte(got), &decoded), "input %q rendered %q", s, got)
		require.Equal(t, s, decoded["v"], "input %q rendered %q", s, got)
	}

	alphabet := []rune{
		'a', '0', '1', '9', 'e', 'E', 'n', 'y', 'X', 'B', 'O',
		'.', '-', '+', ':', '?', '#', '@', '&', '*', '!', '%', '|', '>',
		'\'', '"', '[', ']', '{', '}', ',', ' ', '\t', '\n', '_', '~', '=', '`', '\\',
	}
	for _, a := range alphabet {
		check(t, string(a))
		for _, b := range alphabet {
			check(t, string(a)+string(b))
			for _, c := range alphabet {
				check(t, string(a)+string(b)+string(c))
			}
		}
	}

	numAlpha := []rune{'0', '9', 'e', 'E', '.', '-', '+', '_', 'x', 'o', 'b', 'X', 'A', 'f', 'n', 'i', ':', 'a'}
	var rec func(prefix string, depth int)
	rec = func(prefix string, depth int) {
		if prefix != "" {
			check(t, prefix)
		}
		if depth == 0 {
			return
		}
		for _, r := range numAlpha {
			rec(prefix+string(r), depth-1)
		}
	}
	rec("", 4)

	patterns := []string{
		"30E3",
		"1e10",
		"+30E3",
		"-30E3",
		"30E+3",
		"0x1F",
		"0o17",
		"0b101",
		"010",
		"1_000",
		"1__0",
		"60:30",
		"1:2:3",
		".inf",
		".Inf",
		".INF",
		"+.inf",
		"-.inf",
		".nan",
		".NaN",
		".NAN",
		"yes",
		"no",
		"on",
		"off",
		"Yes",
		"No",
		"On",
		"Off",
		"YES",
		"NO",
		"ON",
		"OFF",
		"true",
		"false",
		"True",
		"False",
		"null",
		"Null",
		"NULL",
		"~",
		"2001-12-14",
		"2001-12-14t21:59:43.10-05:00",
		"!!str",
		"!foo",
		"a: b",
		"- a",
		"? a",
		": a",
		"a #b",
		" a",
		"a ",
		"a b",
		" ",
		"\ufeffa",
		"a\u2028b",
		"\uffff",
		"\ufffe",
		"a\uffffb",
		"\ufffda",
		"a\u0085b",
		"a\ufdd0b",
		"a\U0001FFFFb",
		"a\U0010FFFFb",
		"a\u00a0b",
		"a\U0001F600b",
		"e5",
		"E5",
		"8e",
		"5e5e5",
		"+5",
		"-5",
		"0.",
		".5",
		"-.5",
		"+.5",
		"0.0.0",
		"127.0.0.1",
		"=",
		"0x_1",
		"0b_1",
		"-0b101",
		"0b-0",
		"0o+0",
		"-_09",
		"+_09",
		"0X83",
		"0B101",
		"0O17",
		"0XaF",
		"0X_1",
		"0X-1",
		"0Xg",
		"0B2",
		"0O8",
		"a\nb",
		"a\n\nb",
		"\n",
		"\n\n",
		"a\n",
		"a\n\n",
		"\na",
		" a\nb",
		"a\n b",
		"a\n\tb",
		"1\n2",
		"- a\n- b",
		"a:\nb",
		"#a\nb",
		"a\n#b",
		"|\na",
		">\na",
	}
	for _, p := range patterns {
		check(t, p)
	}
}
