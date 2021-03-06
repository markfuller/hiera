package impl

import (
	"context"

	"github.com/lyraproj/hiera/lookup"
	"github.com/lyraproj/issue/issue"
	"github.com/lyraproj/pcore/pcore"
	"github.com/lyraproj/pcore/px"
)

var NoOptions = map[string]px.Value{}

func init() {
	lookup.TryWithParent = func(parent context.Context, tp lookup.LookupKey, options map[string]px.Value, consumer func(px.Context) error) error {
		return pcore.TryWithParent(parent, func(c px.Context) error {
			InitContext(c, tp, options)
			return consumer(c)
		})
	}

	lookup.DoWithParent = func(parent context.Context, tp lookup.LookupKey, options map[string]px.Value, consumer func(px.Context)) {
		pcore.DoWithParent(parent, func(c px.Context) {
			InitContext(c, tp, options)
			consumer(c)
		})
	}

	lookup.Lookup2 = func(
		ic lookup.Invocation,
		names []string,
		valueType px.Type,
		defaultValue px.Value,
		override px.OrderedMap,
		defaultValuesHash px.OrderedMap,
		options map[string]px.Value,
		block px.Lambda) px.Value {
		if override == nil {
			override = px.EmptyMap
		}
		if defaultValuesHash == nil {
			defaultValuesHash = px.EmptyMap
		}

		if options == nil {
			options = NoOptions
		}

		for _, name := range names {
			if ov, ok := override.Get4(name); ok {
				return ov
			}
			key := NewKey(name)
			if v, ok := ic.Check(key, func() (px.Value, bool) {
				return ic.(*invocation).lookupViaCache(key, options)
			}); ok {
				return v
			}
		}

		if defaultValuesHash.Len() > 0 {
			for _, name := range names {
				if dv, ok := defaultValuesHash.Get4(name); ok {
					return dv
				}
			}
		}

		if defaultValue == nil {
			// nil (as opposed to UNDEF) means that no default was provided.
			if len(names) == 1 {
				panic(px.Error(HieraNameNotFound, issue.H{`name`: names[0]}))
			}
			panic(px.Error(HieraNotAnyNameFound, issue.H{`name_list`: names}))
		}
		return defaultValue
	}
}
