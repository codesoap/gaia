package wrap

import (
	"reflect"
	"strings"
	"testing"
)

func TestFurtherSplitLongWords(t *testing.T) {
	type test struct {
		in          []string
		prefixWidth int
		lineWidth   int
		wantedOut   []string
	}
	tests := []test{
		{[]string{"foo"}, 1, 8, []string{"foo"}},
		{[]string{"foo"}, 1, 2, []string{"fo", "o"}},
		{[]string{"foo"}, 0, 2, []string{"fo", "o"}},
		{[]string{"foo", "bar"}, 3, 4, []string{"foo", "b", "a", "r"}},
		{[]string{"foo", "bar"}, 2, 4, []string{"foo", "ba", "r"}},
		{[]string{"foo", "bar"}, 3, 8, []string{"foo", "bar"}},
		{[]string{"foo", "bar"}, 1, 2, []string{"fo", "o", "b", "a", "r"}},
		{[]string{"fooo", "ba"}, 1, 2, []string{"fo", "o", "o", "b", "a"}},
	}
	for i, testCase := range tests {
		gotOut := furtherSplitLongWords(testCase.in, testCase.prefixWidth, testCase.lineWidth)
		if !reflect.DeepEqual(testCase.wantedOut, gotOut) {
			t.Errorf("#%d Got: %s\nWant: %s", i, strings.Join(gotOut, " "), strings.Join(testCase.wantedOut, " "))
		}
	}
}

func BenchmarkFurtherSplitLongWords(b *testing.B) {
	in := strings.Fields(`
		Voluptas et non ab nihil mollitia pariatur. Aut voluptates autem
		mollitia accusamus ut nesciunt enim amet. Non officia nostrum quia eum
		vel. Vel ea ut exercitationem enim quas iusto. Et soluta ut omnis quos.
		Pariatur rerum veritatis consequatur qui.

		Rerum eaque est minima cumque est quidem in voluptatem. Quia vitae ea
		vero maxime. Facilis itaque quas dolorem et. Adipisci nihil autem magni
		autem sed dolorum beatae autem.

		Commodi est nisi ut nemo. Accusantium exercitationem quas ut nam ex
		dolorem. Provident est velit doloremque et aut qui eum.`)
	var bytes int64 = 0
	for _, word := range in {
		bytes += int64(len(word))
	}
	b.SetBytes(bytes)
	for i := 0; i < b.N; i++ {
		furtherSplitLongWords(in, 4, 60)
	}
}
