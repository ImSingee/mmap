package mmap

type Grower func(current int, atLeast int) (next int)

const oneMB = 1024 * 1024
const oneGB = 1024 * oneMB
const twoGB = 2 * oneGB

func DefaultGrowPolicy(current int, atLeast int) (next int) {
	var fac int
	if current < twoGB {
		fac = current * 2
	} else {
		fac = current + oneGB
	}

	if fac < atLeast {
		fac = atLeast
	}

	return alignOneMB(fac)
}

func alignOneMB(n int) int {
	return align(n, oneMB)
}

func align(n, m int) int {
	return ((n) + ((m) - 1)) & ^((m) - 1)
}

func (m *Mmap) ChangeGrowPolicy(newGrowPolicy Grower) {
	m.grow = newGrowPolicy
}
