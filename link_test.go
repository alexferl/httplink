package httplink

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(""))
}

func checkLinkHeader(t *testing.T, resp *httptest.ResponseRecorder, expected string) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	handler(resp, req)
	assert.Equal(t, expected, resp.Result().Header.Get("Link"))
}

func Test_AppendLink_Single(t *testing.T) {
	expected := "</things/2842>; rel=next"
	resp := httptest.NewRecorder()

	Append(resp.Header(), "/things/2842", "next")

	checkLinkHeader(t, resp, expected)
}

func Test_AppendLink_Multiple(t *testing.T) {
	expected := `</things/2842>; rel=next, ` +
		`<http://%C3%A7runchy/bacon>; rel=contents, ` +
		`<ab%C3%A7>; rel="http://example.com/ext-type", ` +
		`<ab%C3%A7>; rel="http://example.com/%C3%A7runchy", ` +
		`<ab%C3%A7>; rel="https://example.com/too-%C3%A7runchy", ` +
		`</alt-thing>; rel="alternate http://example.com/%C3%A7runchy"`
	resp := httptest.NewRecorder()

	uri := "abç"

	Append(resp.Header(), "/things/2842", "next")
	Append(resp.Header(), "http://çrunchy/bacon", "contents")
	Append(resp.Header(), uri, "http://example.com/ext-type")
	Append(resp.Header(), uri, "http://example.com/çrunchy")
	Append(resp.Header(), uri, "https://example.com/too-çrunchy")
	Append(resp.Header(), "/alt-thing", "alternate http://example.com/çrunchy")

	checkLinkHeader(t, resp, expected)
}

func Test_AppendLink_With_Title(t *testing.T) {
	expected := `</related/thing>; rel=item; title="A related thing"`
	resp := httptest.NewRecorder()

	Append(resp.Header(), "/related/thing", "item", Title("A related thing"))

	checkLinkHeader(t, resp, expected)
}

func Test_AppendLink_With_Title_Star(t *testing.T) {
	expected := `</related/thing>; rel=item; ` +
		`title*=UTF-8''A%20related%20thing, ` +
		`</%C3%A7runchy/thing>; rel=item; ` +
		`title*=UTF-8'en'A%20%C3%A7runchy%20thing`
	resp := httptest.NewRecorder()

	Append(resp.Header(), "/related/thing", "item", TitleStar([]string{"", "A related thing"}))

	Append(resp.Header(), "/çrunchy/thing", "item", TitleStar([]string{"en", "A çrunchy thing"}))

	checkLinkHeader(t, resp, expected)
}

func Test_AppendLink_With_Anchor(t *testing.T) {
	expected := `</related/thing>; rel=item; anchor="/some%20thing/or-other"`
	resp := httptest.NewRecorder()

	Append(resp.Header(), "/related/thing", "item", Anchor("/some thing/or-other"))

	checkLinkHeader(t, resp, expected)
}

func Test_AppendLink_With_HREFLang(t *testing.T) {
	expected := `</related/thing>; rel=about; hreflang=en`
	resp := httptest.NewRecorder()

	Append(resp.Header(), "/related/thing", "about", HREFLang([]string{"en"}))

	checkLinkHeader(t, resp, expected)
}

func Test_AppendLink_With_HREFLang_Multi(t *testing.T) {
	expected := `</related/thing>; rel=about; hreflang=en-GB; hreflang=de`
	resp := httptest.NewRecorder()

	Append(resp.Header(), "/related/thing", "about", HREFLang([]string{"en-GB", "de"}))

	checkLinkHeader(t, resp, expected)
}

func Test_AppendLink_With_TypeHint(t *testing.T) {
	expected := `</related/thing>; rel=alternate; type="video/mp4; codecs=avc1.640028"`
	resp := httptest.NewRecorder()

	Append(resp.Header(), "/related/thing", "alternate", TypeHint("video/mp4; codecs=avc1.640028"))

	checkLinkHeader(t, resp, expected)
}

func Test_AppendLink_Complex(t *testing.T) {
	expected := `</related/thing>; rel=alternate; ` +
		`title="A related thing"; ` +
		`title*=UTF-8'en'A%20%C3%A7runchy%20thing; ` +
		`type="application/json"; ` +
		`hreflang=en-GB; hreflang=de`
	resp := httptest.NewRecorder()

	Append(
		resp.Header(),
		"/related/thing",
		"alternate",
		Title("A related thing"),
		HREFLang([]string{"en-GB", "de"}),
		TypeHint("application/json"),
		TitleStar([]string{"en", "A çrunchy thing"}),
	)

	checkLinkHeader(t, resp, expected)
}

func Test_AppendLink_CrossOrigin(t *testing.T) {
	testCases := []struct {
		crossOrigin string
		expected    string
	}{
		{"", `</related/thing>; rel=alternate`},
		{"anonymous", `</related/thing>; rel=alternate; crossorigin`},
		{"AnOnYmOUs", `</related/thing>; rel=alternate; crossorigin`},
		{"Use-Credentials", `</related/thing>; rel=alternate; crossorigin="use-credentials"`},
		{"use-credentials", `</related/thing>; rel=alternate; crossorigin="use-credentials"`},
	}

	for _, tc := range testCases {
		t.Run(tc.crossOrigin, func(t *testing.T) {
			resp := httptest.NewRecorder()
			Append(resp.Header(), "/related/thing", "alternate", CrossOrigin(tc.crossOrigin))

			checkLinkHeader(t, resp, tc.expected)
		})
	}
}

func Test_AppendLink_Invalid_CrossOrigin_Value(t *testing.T) {
	testCases := []struct {
		crossOrigin string
	}{
		{"*"},
		{"Allow-all"},
		{"Lax"},
		{"MUST-REVALIDATE"},
		{"Strict"},
		{"deny"},
	}

	for _, tc := range testCases {
		t.Run(tc.crossOrigin, func(t *testing.T) {
			resp := httptest.NewRecorder()

			assert.Panics(t, func() {
				Append(resp.Header(), "/related/thing", "alternate", CrossOrigin(tc.crossOrigin))
			})
		})
	}
}

func Test_AppendLink_With_Link_Extension(t *testing.T) {
	expected := `</related/thing>; rel=item; sizes=72x72`
	resp := httptest.NewRecorder()

	Append(resp.Header(), "/related/thing", "item", LinkExtension([][]string{{"sizes", "72x72"}}))

	checkLinkHeader(t, resp, expected)
}
