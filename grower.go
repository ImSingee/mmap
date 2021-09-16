package mmap

type Grower func(current int, atLeast int) (next int)

const OneMB = 1024 * 1024
const OneGB = 1024 * OneMB
const TwoGB = 2 * OneGB

func DefaultGrowPolicy(current int, atLeast int) (next int) {
	var fac int
	if current < TwoGB {
		fac = current * 2
	} else {
		fac = current + OneGB
	}

	if fac < atLeast {
		fac = atLeast
	}

	return alignOneMB(fac)
}

func alignOneMB(n int) int {
	return align(n, OneMB)
}

func align(n, m int) int {
	return ((n) + ((m) - 1)) & ^((m) - 1)
}

func (m *Mmap) ChangeGrowPolicy(newGrowPolicy Grower) {
	m.grow = newGrowPolicy
}
