package wrap

import (
	"reflect"
	"strings"
	"testing"
)

func TestWrap(t *testing.T) {
	type test struct {
		in          string
		lineWidth   int
		wantedOut   []string
	}
	tests := []test{
		{"foo", 8, []string{"foo"}},
		{"foo", 4, []string{"foo"}},
		{"foo", 2, []string{"fo", "o"}},
		{"foo bar", 5, []string{"foo", "bar"}},
		{"foo bar", 6, []string{"foo", "bar"}},
		{"foo bar", 7, []string{"foo bar"}},
		{"foo bar", 8, []string{"foo bar"}},
	}
	for i, testCase := range tests {
		gotOut := Wrap(testCase.in, testCase.lineWidth)
		if !reflect.DeepEqual(testCase.wantedOut, gotOut) {
			t.Errorf("#%d Got: %s\nWant: %s", i, strings.Join(gotOut, "\n"), strings.Join(testCase.wantedOut, "\n"))
		}
	}
}

func TestWrapWithPrefix(t *testing.T) {
	type test struct {
		in          string
		lineWidth   int
		wantedOut   []string
	}
	tests := []test{
		{"foo", 8, []string{"> foo"}},
		{"foo", 4, []string{"> fo", "> o"}},
		{"foo", 3, []string{"> f", "> o", "> o"}},
		{"foo bar", 5, []string{"> foo", "> bar"}},
		{"foo bar", 6, []string{"> foo", "> bar"}},
		{"foo bar", 8, []string{"> foo", "> bar"}},
		{"foo bar", 9, []string{"> foo bar"}},
	}
	for i, testCase := range tests {
		gotOut := WrapWithPrefix(testCase.in, "> ", testCase.lineWidth)
		if !reflect.DeepEqual(testCase.wantedOut, gotOut) {
			t.Errorf("#%d Got: %s\nWant: %s", i, strings.Join(gotOut, "\n"), strings.Join(testCase.wantedOut, "\n"))
		}
	}
}

func BenchmarkWrap(b *testing.B) {
	in := `Voluptas et non ab nihil mollitia pariatur. Aut voluptates autem mollitia accusamus ut nesciunt enim amet. Non officia nostrum quia eum vel. Vel ea ut exercitationem enim quas iusto. Et soluta ut omnis quos. Pariatur rerum veritatis consequatur qui. Rerum eaque est minima cumque est quidem in voluptatem. Quia vitae ea vero maxime. Facilis itaque quas dolorem et. Adipisci nihil autem magni autem sed dolorum beatae autem. Commodi est nisi ut nemo. Accusantium exercitationem quas ut nam ex dolorem. Provident est velit doloremque et aut qui eum.`
	b.SetBytes(int64(len(in)))
	for i := 0; i < b.N; i++ {
		Wrap(in, 60)
	}
}

func BenchmarkWrapWithPrefix(b *testing.B) {
	in := `Voluptas et non ab nihil mollitia pariatur. Aut voluptates autem mollitia accusamus ut nesciunt enim amet. Non officia nostrum quia eum vel. Vel ea ut exercitationem enim quas iusto. Et soluta ut omnis quos. Pariatur rerum veritatis consequatur qui. Rerum eaque est minima cumque est quidem in voluptatem. Quia vitae ea vero maxime. Facilis itaque quas dolorem et. Adipisci nihil autem magni autem sed dolorum beatae autem. Commodi est nisi ut nemo. Accusantium exercitationem quas ut nam ex dolorem. Provident est velit doloremque et aut qui eum.`
	b.SetBytes(int64(len(in)))
	for i := 0; i < b.N; i++ {
		WrapWithPrefix(in, "> ", 60)
	}
}
