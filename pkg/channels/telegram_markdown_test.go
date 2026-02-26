package channels

import (
	"strings"
	"testing"
)

func TestMarkdownToTelegramHTML_MultilineFormatting(t *testing.T) {
	input := "# Title\n- one\n- two\n1. first\n2. second\n> quote"
	out := markdownToTelegramHTML(input)

	if strings.Contains(out, "# Title") {
		t.Fatalf("header not normalized: %q", out)
	}
	if strings.Count(out, "• ") < 4 {
		t.Fatalf("expected bullets for list items, got: %q", out)
	}
	if strings.Contains(out, "> quote") {
		t.Fatalf("blockquote marker not removed: %q", out)
	}
}

func TestMarkdownToTelegramHTML_LinkEscapesHrefQuotes(t *testing.T) {
	input := `[go](https://example.com?a=1&b="x")`
	out := markdownToTelegramHTML(input)

	if !strings.Contains(out, `<a href="https://example.com?a=1&amp;b=&quot;x&quot;">go</a>`) {
		t.Fatalf("unexpected link conversion: %q", out)
	}
}

func TestMarkdownToTelegramHTML_CodeIsPreserved(t *testing.T) {
	input := "before\n```bash\necho '<x>'\n```\n`a_b`"
	out := markdownToTelegramHTML(input)

	if !strings.Contains(out, "<pre><code>echo '&lt;x&gt;'\n</code></pre>") {
		t.Fatalf("code block not preserved: %q", out)
	}
	if !strings.Contains(out, "<code>a_b</code>") {
		t.Fatalf("inline code not preserved: %q", out)
	}
}
