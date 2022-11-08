package httplink

import (
	"fmt"
	"net/http"
	"strings"
)

const (
	Unreserved = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-._~"
	Delimiters = ":/?#[]@!$&'()*+,;="
	AllAllowed = Unreserved + Delimiters
	HexDigits  = "0123456789ABCDEFabcdef"
)

var CrossOriginValues = []string{"anonymous", "use-credentials"}

func Append(header http.Header, target string, rel string, opts ...Option) {
	args := &Options{}

	for _, opt := range opts {
		opt(args)
	}

	uriEncode := createStrEncoder(false, true)
	uriEncodeValue := createStrEncoder(true, true)

	if strings.Contains(rel, "//") {
		if strings.Contains(rel, " ") {
			var escaped []string
			for _, s := range strings.Split(rel, " ") {
				escaped = append(escaped, uriEncode(s))
			}
			rel = fmt.Sprintf(`"%s"`, strings.Join(escaped, " "))
		} else {
			rel = fmt.Sprintf(`"%s"`, uriEncode(rel))
		}
	}

	value := fmt.Sprintf("<%s>; rel=%s", uriEncode(target), rel)

	if args.Title != "" {
		value += fmt.Sprintf(`; title="%s"`, args.Title)
	}

	if len(args.TitleStar) > 0 {
		value += fmt.Sprintf(`; title*=UTF-8'%s'%s`, args.TitleStar[0], uriEncodeValue(args.TitleStar[1]))
	}

	if args.TypeHint != "" {
		value += fmt.Sprintf(`; type="%s"`, args.TypeHint)
	}

	if len(args.HREFLang) > 0 {
		var langs []string
		for _, s := range args.HREFLang {
			langs = append(langs, "hreflang="+s)
		}
		value += "; "
		value += strings.Join(langs, "; ")
	}

	if args.Anchor != "" {
		value += fmt.Sprintf(`; anchor="%s"`, uriEncode(args.Anchor))
	}

	if args.CrossOrigin != "" {
		crossOrigin := strings.ToLower(args.CrossOrigin)
		if !checkInSlice(crossOrigin, CrossOriginValues) {
			panic("crossorigin must be set to either 'anonymous' or 'use-credentials'")
		}
		if crossOrigin == "anonymous" {
			value += "; crossorigin"
		} else {
			value += `; crossorigin="use-credentials"`
		}
	}

	if len(args.LinkExtension) > 0 {
		var links []string
		for _, s := range args.LinkExtension {
			links = append(links, fmt.Sprintf("%s=%s", s[0], s[1]))
		}
		value += "; "
		value += strings.Join(links, "; ")
	}

	link := header.Get("Link")
	if link != "" {
		link += fmt.Sprintf(", %s", value)
		header.Set("Link", link)
	} else {
		header.Set("Link", value)
	}
}

func createCharLookup(allowedChars string) map[int]string {
	lookup := make(map[int]string)

	var encodedChar string
	for i := 0; i <= 255; i++ {
		if strings.Contains(allowedChars, fmt.Sprintf("%c", rune(i))) {
			encodedChar = fmt.Sprintf("%c", rune(i))
		} else {
			encodedChar = fmt.Sprintf("%%%0.2X", i)
		}

		lookup[i] = encodedChar
	}

	return lookup
}

func createStrEncoder(isValue bool, checkIsEscaped bool) func(uri string) string {
	var allowedChars string
	if isValue {
		allowedChars = Unreserved
	} else {
		allowedChars = AllAllowed
	}
	allowedCharsPlusPercent := allowedChars + "%"
	lookup := createCharLookup(allowedChars)

	return func(uri string) string {
		if trim := strings.TrimRight(uri, allowedChars); trim == "" {
			return uri
		}

		trim := strings.TrimRight(uri, allowedCharsPlusPercent)
		if checkIsEscaped && trim == "" {
			tokens := strings.Split(uri, "%")
			var broke bool
			for _, token := range tokens[1:] {
				hexOctet := token[:2]

				if len(hexOctet) != 2 {
					broke = true
					break
				}

				if !strings.Contains(string(hexOctet[0]), HexDigits) &&
					strings.Contains(string(hexOctet[1]), hexOctet) {
					broke = true
					break
				}
			}

			if broke {
				return uri
			}
		}

		var chars []string
		for _, i := range []byte(uri) {
			chars = append(chars, lookup[int(i)])
		}

		return strings.Join(chars, "")
	}
}

func checkInSlice(s string, slice []string) bool {
	for _, i := range slice {
		if i == s {
			return true
		}
	}

	return false
}
