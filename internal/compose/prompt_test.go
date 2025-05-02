package compose

import (
	"strings"
	"testing"
)

func TestMakePromp(t *testing.T) {
	prompt, err := makePrompt("English", []string{})

	if len(prompt) == 0 {
		t.Errorf("Expected prompt to be non-empty, got: %s", prompt)
	}

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	const expectedLangCount = 5
	if langCount := strings.Count(prompt, "English"); langCount != expectedLangCount {
		t.Errorf("Expected prompt to contain 'English' %v times, got %v times instead", expectedLangCount, langCount)
	}
}

func TestMakeTagsString(t *testing.T) {
	cases := []struct {
		name string
		tags []string
		want string
	}{
		{
			name: "natural string when no tags provided",
			tags: []string{},
			want: "No tags provided",
		},
		{
			name: "single tag",
			tags: []string{"Funny"},
			want: "Funny",
		},
		{
			name: "multiple tags comma separated",
			tags: []string{"Funny", "Whimsical"},
			want: "Funny, Whimsical",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			got := makeTagsString(c.tags)

			if got != c.want {
				t.Errorf("Expected %q, got %q", c.want, got)
			}
		})
	}
}
