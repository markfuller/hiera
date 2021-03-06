package impl_test

import (
	"context"
	"fmt"
	"github.com/lyraproj/hiera/impl"
	"github.com/lyraproj/pcore/pcore"
	"github.com/lyraproj/pcore/px"
	"strings"
)

func ExampleNewKey_simple() {
	key := impl.NewKey(`simple`)
	fmt.Printf(`%s, %d`, key, len(key.Parts()))
	// Output: simple, 1
}

func ExampleNewKey_dotted() {
	key := impl.NewKey(`a.b.c`)
	fmt.Printf(`%s, %d`, key, len(key.Parts()))
	// Output: a.b.c, 3
}

func ExampleNewKey_dotted_int() {
	key := impl.NewKey(`a.3`)
	fmt.Printf(`%T`, key.Parts()[1])
	// Output: int
}

func ExampleNewKey_quoted() {
	key := impl.NewKey(`'a.b.c'`)
	fmt.Printf(`%s, %d`, key, len(key.Parts()))
	// Output: 'a.b.c', 1
}

func ExampleNewKey_doubleQuoted() {
	key := impl.NewKey(`"a.b.c"`)
	fmt.Printf(`%s, %d`, key, len(key.Parts()))
	// Output: "a.b.c", 1
}

func ExampleNewKey_quotedDot() {
	key := impl.NewKey(`a.'b.c'`)
	fmt.Printf(`%s, %d, %s`, key, len(key.Parts()), key.Parts()[1])
	// Output: a.'b.c', 2, b.c
}

func ExampleNewKey_quotedDotX() {
	key := impl.NewKey(`a.'b.c'.d`)
	fmt.Printf(`%s, %d, %s`, key, len(key.Parts()), key.Parts()[1])
	// Output: a.'b.c'.d, 3, b.c
}

func ExampleNewKey_quotedQuote() {
	key := impl.NewKey(`a.b.'c"d"e'`)
	fmt.Printf(`%s, %d, %s`, key, len(key.Parts()), key.Parts()[2])
	// Output: a.b.'c"d"e', 3, c"d"e
}

func ExampleNewKey_doubleQuotedQuote() {
	key := impl.NewKey(`a.b."c'd'e"`)
	fmt.Printf(`%s, %d, %s`, key, len(key.Parts()), key.Parts()[2])
	// Output: a.b."c'd'e", 3, c'd'e
}

func printErr(e error) {
	s := e.Error()
	if ix := strings.Index(s, ` (file: `); ix > 0 {
		s = s[0:ix]
	}
	fmt.Println(s)
}

func ExampleNewKey_unterminatedQuoted() {
	printErr(pcore.TryWithParent(context.Background(), func(c px.Context) error {
		impl.NewKey(`a.b."c`)
		return nil
	}))
	// Output: Unterminated quote in key 'a.b."c'
}

func ExampleNewKey_empty() {
	printErr(pcore.TryWithParent(context.Background(), func(c px.Context) error {
		impl.NewKey(``)
		return nil
	}))
	// Output: lookup() key '' contains an empty segment
}

func ExampleNewKey_emptySegment() {
	printErr(pcore.TryWithParent(context.Background(), func(c px.Context) error {
		impl.NewKey(`a..b`)
		return nil
	}))
	// Output: lookup() key 'a..b' contains an empty segment
}

func ExampleNewKey_emptySegmentStart() {
	printErr(pcore.TryWithParent(context.Background(), func(c px.Context) error {
		impl.NewKey(`.b`)
		return nil
	}))
	// Output: lookup() key '.b' contains an empty segment
}

func ExampleNewKey_emptySegmentEnd() {
	printErr(pcore.TryWithParent(context.Background(), func(c px.Context) error {
		impl.NewKey(`a.`)
		return nil
	}))
	// Output: lookup() key 'a.' contains an empty segment
}

func ExampleNewKey_firstSegmentIndex() {
	printErr(pcore.TryWithParent(context.Background(), func(c px.Context) error {
		impl.NewKey(`1.a`)
		return nil
	}))
	// Output: lookup() key '1.a' first segment cannot be an index
}
