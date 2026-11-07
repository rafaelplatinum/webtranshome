package handler

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestPublicArticleFilter(t *testing.T) {
	cases := []struct {
		input string
		valid bool
		count int
	}{
		{"", true, 0},
		{"ARTIKEL", true, 1},
		{"PROMO,EVENT", true, 2},
		{" PROMO, PROMO , EVENT ", true, 2},
		{"HALAMAN", true, 1},
		{"SECRET", false, 0},
		{"PROMO,,EVENT", false, 0},
	}
	for _, tc := range cases {
		items, valid := publicArticleFilter(tc.input)
		if valid != tc.valid || len(items) != tc.count {
			t.Fatalf("input %q: valid=%v items=%v", tc.input, valid, items)
		}
	}
}

func TestPublicArticleExcerpt(t *testing.T) {
	cases := []struct{ input, want string }{
		{"## Judul\n\nKalimat pertama berisi info.\n\nKedua.", "Kalimat pertama berisi info."},
		{"- Item\n\nParagraf teks biasa.", "Paragraf teks biasa."},
		{"", ""},
	}
	for _, tc := range cases {
		got := publicArticleExcerpt(tc.input)
		if got != tc.want {
			t.Fatalf("got %q want %q", got, tc.want)
		}
	}
	got := publicArticleExcerpt(strings.Repeat("kata ", 60))
	if len([]rune(got)) > 160 || !strings.HasSuffix(got, "…") {
		t.Fatalf("bad truncation: %q", got)
	}
}

func TestPublicArticlePagination(t *testing.T) {
	if p, ok := publicArticlePage("", 9); !ok || p != 9 {
		t.Fatal("fallback pagination")
	}
	if _, ok := publicArticlePage("0", 9); ok {
		t.Fatal("should reject 0")
	}
	if _, ok := publicArticlePage("xyz", 9); ok {
		t.Fatal("should reject non-number")
	}
}

func TestPublicSettingsWhatsAppAndNulls(t *testing.T) {
	if v := publicSettingWA("0877 0877-7989"); v == nil || *v != "6287708777989" {
		t.Fatalf("number=%v", v)
	}
	if v := publicSettingWA(""); v != nil {
		t.Fatalf("empty number=%v", v)
	}
	if v := publicSettingString(map[string]string{}, "address"); v != nil {
		t.Fatalf("missing key=%v", v)
	}
}

func TestPublicContentJSONShape(t *testing.T) {
	body := publicArticleDetailResponse{Data: publicArticleDetail{
		publicArticleItem: publicArticleItem{ID: 1, Title: "Contoh", Type: "PROMO", Slug: "contoh"},
		Content:           "Isi artikel", Status: "PUBLISHED",
	}}
	raw, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}
	for _, key := range []string{`"data"`, `"title"`, `"content"`, `"status"`} {
		if !strings.Contains(string(raw), key) {
			t.Fatalf("missing %s: %s", key, raw)
		}
	}
}