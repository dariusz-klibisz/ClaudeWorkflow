package cmdid

import "testing"

func TestClassify(t *testing.T) {
	cases := []struct {
		name       string
		red, green string
		want       Pair
	}{
		{"same command", "go test ./...", "go test ./...", PairStrict},
		{"red selector, green suite", "pytest tests/test_x.py::test_new", "pytest", PairStrict},
		{"green selector superset", "pytest", "pytest tests/test_x.py", PairStrict},
		{"same runner different selector", "pytest tests/test_x.py", "pytest tests/test_y.py", PairWeak},
		{"go pkg vs all", "go test ./pkg", "go test ./...", PairWeak},
		{"cross runner", "gitleaks detect", "pytest", PairNone},
		{"cross runner cargo/go", "cargo test", "go test ./...", PairNone},
		{"env prefix ignored", "FOO=1 go test ./...", "go test ./...", PairStrict},
		{"leading cd stripped", "cd api && go test ./...", "go test ./...", PairStrict},
		{"generic interpreter same script", "python3 test_app.py", "python3 test_app.py -v", PairStrict},
		{"generic interpreter other script", "python3 test_app.py", "python3 other.py", PairWeak},
		{"empty red", "", "go test", PairNone},
		{"npm script variants", "npm run test:unit", "npm run test:unit", PairStrict},
		{"make targets diverge", "make test", "make test-unit", PairWeak},
	}
	for _, c := range cases {
		if got := Classify(c.red, c.green); got != c.want {
			t.Errorf("%s: Classify(%q, %q) = %v, want %v", c.name, c.red, c.green, got, c.want)
		}
	}
}

func TestLearnedHead(t *testing.T) {
	cases := map[string]string{
		"pytest tests/x.py":       "pytest",
		"go test ./...":           "go test",
		"python3 -m unittest x.y": "python3 -m unittest",
		"python3 test_app.py":     "", // bare generic interpreter
		"python3 -m":              "", // dangling module flag
		"FOO=1 cargo test":        "cargo test",
		"npx vitest run tests/x":  "npx vitest run",
	}
	for cmd, want := range cases {
		if got := LearnedHead(cmd); got != want {
			t.Errorf("LearnedHead(%q) = %q, want %q", cmd, got, want)
		}
	}
}
