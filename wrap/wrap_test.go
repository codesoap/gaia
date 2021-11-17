package wrap

import (
	"reflect"
	"strings"
	"testing"
)

func TestWrapToWidth(t *testing.T) {
	type test struct {
		in          string
		prefix string
		lineWidth   int
		wantedOut   []string
	}
	tests := []test{
		{"foo", "", 8, []string{"foo"}},
		{"> foo", "> ", 4, []string{"> fo", "> o"}},
		//{"foo bar", "", 5, []string{"foo", "bar"}},
		//{"foo bar", "", 6, []string{"foo", "bar"}},
		//{"foo bar", "", 7, []string{"foo bar"}},
		//{"foo bar", "", 8, []string{"foo bar"}},
	}
	for i, testCase := range tests {
		gotOut := WrapToWidth(testCase.in, testCase.prefix, testCase.lineWidth)
		if !reflect.DeepEqual(testCase.wantedOut, gotOut) {
			t.Errorf("#%d Got: %s\nWant: %s", i, strings.Join(gotOut, "\n"), strings.Join(testCase.wantedOut, "\n"))
		}
	}
}

func BenchmarkWrapToWidth(b *testing.B) {
	in := `> Voluptas et non ab nihil mollitia pariatur. Aut voluptates autem
		mollitia accusamus ut nesciunt enim amet. Non officia nostrum quia eum
		vel. Vel ea ut exercitationem enim quas iusto. Et soluta ut omnis quos.
		Pariatur rerum veritatis consequatur qui.

		Rerum eaque est minima cumque est quidem in voluptatem. Quia vitae ea
		vero maxime. Facilis itaque quas dolorem et. Adipisci nihil autem magni
		autem sed dolorum beatae autem.

		Commodi est nisi ut nemo. Accusantium exercitationem quas ut nam ex
		dolorem. Provident est velit doloremque et aut qui eum.`
	b.SetBytes(int64(len(in)))
	for i := 0; i < b.N; i++ {
		WrapToWidth(in, "> ", 60)
	}
}
