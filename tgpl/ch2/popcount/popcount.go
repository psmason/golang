package popcount

var pc [256]byte

func init() {
	for i := range pc {
		pc[i] = pc[i/2] + byte(i&1)
	}
}

func PopCount(x uint64) int {
	return int(pc[byte(x>>(0*8))] +
		pc[byte(x>>(1*8))] +
		pc[byte(x>>(2*8))] +
		pc[byte(x>>(3*8))] +
		pc[byte(x>>(4*8))] +
		pc[byte(x>>(5*8))] +
		pc[byte(x>>(6*8))] +
		pc[byte(x>>(7*8))])
}

func PopCountLooped(x uint64) int {
	var result int = 0
	var i uint64 = 0
	for ; i < 8; i++ {
		result = result + int(pc[byte(x>>(i*8))])
	}
	return result
}

func PopCountSlow(x uint64) int {
	var result int = int(byte(x & 1))
	for i := 0; i < 64; i++ {
		x = x >> 1
		result = result + int(byte(x&1))
	}
	return result
}

func PopCountBitTrick(x uint64) int {
	var result int = 0
	for ; x != 0; result++ {
		x &= (x - 1)
	}
	return result
}
