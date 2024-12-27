package main

var comparisonResults struct {
	isGreater uint8
	isEqual   uint8
} = struct {
	isGreater uint8
	isEqual   uint8
}{
	0b00000001,
	0b00000010,
}

func makeComparisonResult(val1, val2 uint64) uint8 {
	res := uint8(0)

	if val1 > val2 {
		res |= comparisonResults.isGreater
	}
	if val1 == val2 {
		res |= comparisonResults.isEqual
	}

	return res
}