package stream

func Sum[T Number | ~string](cumulative, next T) T { return cumulative + next }
func Difference[T Number](cumulative, next T) T    { return cumulative - next }
func Product[T Number](cumulative, next T) T       { return cumulative * next }
func Quotient[T Number](cumulative, next T) T      { return cumulative / next }
func Remainder[T Integer](cumulative, next T) T    { return cumulative % next }

func BitOr[T Integer](cumulative, next T) T      { return cumulative | next }
func BitAnd[T Integer](cumulative, next T) T     { return cumulative & next }
func BitXor[T Integer](cumulative, next T) T     { return cumulative ^ next }
func BitClear[T Integer](cumulative, next T) T   { return cumulative &^ next }
func LeftShift[T Integer](cumulative, next T) T  { return cumulative << next }
func RightShift[T Integer](cumulative, next T) T { return cumulative >> next }

func And[T ~bool](cumulative, next T) T { return cumulative && next }
func Or[T ~bool](cumulative, next T) T  { return cumulative || next }

func Append[T any, S ~[]T](total, next S) S { return append(total, next...) }
