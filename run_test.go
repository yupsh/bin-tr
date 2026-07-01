package main

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/spf13/afero"
)

func TestRun(t *testing.T) {
	cases := []struct {
		name       string
		version    string
		stdin      string
		wantOut    string
		wantErrSub string
		args       []string
		wantCode   int
	}{
		{
			name:    "translate set1 to set2",
			args:    []string{"tr", "abc", "xyz"},
			stdin:   "abc\n",
			wantOut: "xyz\n",
		},
		{
			name:    "translate range",
			args:    []string{"tr", "a-z", "A-Z"},
			stdin:   "hello\n",
			wantOut: "HELLO\n",
		},
		{
			name:    "delete set1 only",
			args:    []string{"tr", "-d", "aeiou"},
			stdin:   "hello world\n",
			wantOut: "hll wrld\n",
		},
		{
			name:    "squeeze repeats with single operand",
			args:    []string{"tr", "-s", " "},
			stdin:   "hello   world\n",
			wantOut: "hello world\n",
		},
		{
			name:    "squeeze repeats with translate operands",
			args:    []string{"tr", "-s", "ab", "xy"},
			stdin:   "aabb\n",
			wantOut: "xy\n",
		},
		{
			name:    "complement translate",
			args:    []string{"tr", "-c", "a-z", "_"},
			stdin:   "Hello World 123\n",
			wantOut: "_ello__orld____\n",
		},
		{
			name:    "version flag reports injected version",
			version: "1.2.3",
			args:    []string{"tr", "--version"},
			wantOut: "tr version 1.2.3\n",
		},
		{
			name:       "missing all operands errors",
			args:       []string{"tr"},
			wantCode:   1,
			wantErrSub: "tr: missing operand",
		},
		{
			name:       "missing set2 without delete errors",
			args:       []string{"tr", "abc"},
			wantCode:   1,
			wantErrSub: "tr: missing operand",
		},
		{
			name:       "unknown flag errors",
			args:       []string{"tr", "--nope", "abc", "xyz"},
			wantCode:   1,
			wantErrSub: "tr:",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			var out, errOut bytes.Buffer
			code := run(tc.version, tc.args, strings.NewReader(tc.stdin), &out, &errOut, fs)

			if code != tc.wantCode {
				t.Fatalf("exit code = %d, want %d (stderr=%q)", code, tc.wantCode, errOut.String())
			}
			if tc.wantErrSub == "" && out.String() != tc.wantOut {
				t.Fatalf("stdout = %q, want %q", out.String(), tc.wantOut)
			}
			if tc.wantErrSub != "" && !strings.Contains(errOut.String(), tc.wantErrSub) {
				t.Fatalf("stderr = %q, want substring %q", errOut.String(), tc.wantErrSub)
			}
		})
	}
}

func Test_main(t *testing.T) {
	origExit, origRun := osExit, runCLI
	t.Cleanup(func() { osExit, runCLI = origExit, origRun })

	gotCode := -1
	osExit = func(code int) { gotCode = code }
	runCLI = func(string, []string, io.Reader, io.Writer, io.Writer, afero.Fs) int { return 7 }

	main()

	if gotCode != 7 {
		t.Fatalf("main propagated exit code %d, want 7", gotCode)
	}
}
