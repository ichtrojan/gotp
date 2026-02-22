package gotp

type Format int

const (
	Alpha Format = iota
	AlphaNumeric
	Numeric
)

// Deprecated: ALPHA is deprecated, use Alpha instead.
const ALPHA = Alpha

// Deprecated: ALPHA_NUMERIC is deprecated, use AlphaNumeric instead.
const ALPHA_NUMERIC = AlphaNumeric

// Deprecated: NUMERIC is deprecated, use Numeric instead.
const NUMERIC = Numeric
