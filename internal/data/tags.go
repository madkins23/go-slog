package data

// -----------------------------------------------------------------------------

// TestTag is a unique name for a Benchmark or Verification test.
// The type is an alias for string so that types can't be confused.
type TestTag string

var aliasTestTag = make(map[string]TestTag)

// Alias defines an alias string for use by TestTagFor for this TestTag.
func (t TestTag) Alias(s string) {
	aliasTestTag[s] = t
}

// Tag returns any TestTag aliased to the string value of the current one,
// otherwise the current tag.
func (t TestTag) Tag() TestTag {
	if tag, found := aliasTestTag[string(t)]; found {
		return tag
	}
	return t
}

// -----------------------------------------------------------------------------

// HandlerTag is a unique name for a slog handler.
// The type is an alias for string so that types can't be confused.
type HandlerTag string

var aliasHandlerTag = make(map[string]HandlerTag)

// Alias defines an alias string for use by TestHandlerFor for this HandlerTag.
func (tag HandlerTag) Alias(s string) {
	aliasHandlerTag[s] = tag
}

// Tag returns any HandlerTag aliased to the string value of the current one,
// otherwise the current tag.
func (tag HandlerTag) Tag() HandlerTag {
	if tag, found := aliasHandlerTag[string(tag)]; found {
		return tag
	}
	return tag
}

// -----------------------------------------------------------------------------

func init() {
	// TODO: Is there a way to do this automagically instead of this awful hack?
	HandlerTag("darvaza/zerolog").Alias("darvaza_zerolog")
	HandlerTag("phsym/zeroslog").Alias("phsym_zerolog")
	HandlerTag("samber/slog-zap").Alias("samber_zap")
	HandlerTag("samber/slog-zerolog").Alias("samber_zerolog")
	HandlerTag("slog/JSONHandler").Alias("slog_JSONHandler")

}
