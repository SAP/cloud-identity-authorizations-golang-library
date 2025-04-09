package expression

import (
	"regexp"
	"strings"
)

type Like struct {
	Arg     Expression
	Pattern String
	Escape  String
	regex   *regexp.Regexp
}

type NotLike struct {
	Like
}

const (
	placeholder1 = "\x1c"
	placeholder2 = "\x1e"
	placeholder3 = "\x1f"
)

func NewNotLike(arg Expression, pattern String, escape String) NotLike {
	return NotLike{
		Like: NewLike(arg, pattern, escape),
	}
}

func NewLike(arg Expression, pattern String, escape String) Like {
	return Like{
		Arg:     arg,
		Pattern: pattern,
		Escape:  escape,
		regex:   createLikeRegex(pattern, escape),
	}
}

func (l NotLike) Evaluate(input Input) Expression {
	arg := l.Arg.Evaluate(input)
	if arg == UNSET {
		return UNSET
	}
	if arg == IGNORE {
		return IGNORE
	}
	str, ok := arg.(String)
	if !ok {
		return NotLike{
			Like: Like{
				Arg:     arg,
				Pattern: l.Pattern,
				Escape:  l.Escape,
				regex:   l.regex,
			},
		}
	}
	return Bool(!l.regex.MatchString(string(str)))
}

func (l Like) Evaluate(input Input) Expression {
	arg := l.Arg.Evaluate(input)
	if arg == UNSET {
		return UNSET
	}
	if arg == IGNORE {
		return IGNORE
	}
	str, ok := arg.(String)
	if !ok {
		return Like{
			Arg:     arg,
			Pattern: l.Pattern,
			Escape:  l.Escape,
			regex:   l.regex,
		}
	}
	return Bool(l.regex.MatchString(string(str)))
}

func createLikeRegex(pattern, escape String) *regexp.Regexp {
	p := string(pattern)
	e := string(escape)
	if e != "" {
		p = strings.ReplaceAll(p, e+e, placeholder1)
		p = strings.ReplaceAll(p, e+"_", placeholder2)
		p = strings.ReplaceAll(p, e+"%", placeholder3)
	}
	// no we need to escape the regex characters
	p = regexp.QuoteMeta(p)
	p = strings.ReplaceAll(p, "%", ".*")
	p = strings.ReplaceAll(p, "_", ".")
	if escape != "" {
		p = strings.ReplaceAll(p, placeholder1, e)
		p = strings.ReplaceAll(p, placeholder2, "_")
		p = strings.ReplaceAll(p, placeholder3, "%")
	}
	return regexp.MustCompile(p)
}
