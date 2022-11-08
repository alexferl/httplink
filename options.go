package httplink

type Options struct {
	Title         string
	TitleStar     []string
	Anchor        string
	HREFLang      []string
	TypeHint      string
	CrossOrigin   string
	LinkExtension [][]string
}

type Option func(*Options)

func Title(s string) Option {
	return func(args *Options) {
		args.Title = s
	}
}

func TitleStar(s []string) Option {
	return func(args *Options) {
		args.TitleStar = s
	}
}

func Anchor(s string) Option {
	return func(args *Options) {
		args.Anchor = s
	}
}

func HREFLang(s []string) Option {
	return func(args *Options) {
		args.HREFLang = s
	}
}

func TypeHint(s string) Option {
	return func(args *Options) {
		args.TypeHint = s
	}
}

func CrossOrigin(s string) Option {
	return func(args *Options) {
		args.CrossOrigin = s
	}
}

func LinkExtension(s [][]string) Option {
	return func(args *Options) {
		args.LinkExtension = s
	}
}
