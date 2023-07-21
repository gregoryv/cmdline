// GENERATED CODE, DO NOT EDIT

package cmdline


// StringVar same as String but uses dst as default and destination.
func (opt *Option) StringVar(dst *string) string {
	*dst = opt.String(*dst)
	return *dst
}

// IntVar same as Int but uses dst as default and destination.
func (opt *Option) IntVar(dst *int) int {
	*dst = opt.Int(*dst)
	return *dst
}

// Uint8Var same as Uint8 but uses dst as default and destination.
func (opt *Option) Uint8Var(dst *uint8) uint8 {
	*dst = opt.Uint8(*dst)
	return *dst
}

// Uint16Var same as Uint16 but uses dst as default and destination.
func (opt *Option) Uint16Var(dst *uint16) uint16 {
	*dst = opt.Uint16(*dst)
	return *dst
}

// Float64Var same as Float64 but uses dst as default and destination.
func (opt *Option) Float64Var(dst *float64) float64 {
	*dst = opt.Float64(*dst)
	return *dst
}

