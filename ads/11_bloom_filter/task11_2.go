package main

import (
	"errors"
)

// Порядковый номер задания на курсе: 11
// Номер задачи из задания: 2
// Краткое название: алгоритм слияния нескольких фильтров Блюма (одинакового размера и с одинаковым набором хэш-функций)
// Вероятность ложных срабатываний при этом растёт, поскольку увеличивается кол-во 1 (0 не получить при побитовом ИЛИ).
// И вероятность растёт тем больше, чем больше уникальных элементов в сливаемых буферах (лучшее распределение 1).

var (
	ErrFilterSizesDiffer = errors.New("sizes of filters differ")
	ErrInvalidFilterSize = errors.New("invalid filter size")
	ErrIrrelevantLength  = errors.New("irrelevant bits length")
	ErrNilFilter         = errors.New("nil filter provided")
	ErrNoFiltersProvided = errors.New("no filters provided")
)

// Time complexity: O(K * N), где K - кол-во фильтров, N - размер фильтров. В данной реализации можно прираванять к O(1)
// Memory complexity: O(N), где N - размер фильтра.

func MergeBloomFilters(filters ...*BloomFilter) (*BloomFilter, error) {
	var result *BloomFilter

	if err := validateFilters(filters); err != nil {
		return result, err
	}

	result = NewBloomFilter(filters[0].filter_len)

	for _, filter := range filters {
		if len(filter.bits) == 0 {
			continue
		}

		orBits(result, filter)
	}

	return result, nil
}

func orBits(result, filter *BloomFilter) {
	for i, word := range filter.bits {
		result.bits[i] |= word
	}
}

func validateFilters(filters []*BloomFilter) error {
	if len(filters) == 0 {
		return ErrNoFiltersProvided
	}

	if filters[0] == nil {
		return ErrNilFilter
	}

	size := filters[0].filter_len
	for _, filter := range filters {
		if filter == nil {
			return ErrNilFilter
		}

		if filter.filter_len <= 0 {
			return ErrInvalidFilterSize
		}

		if size != filter.filter_len {
			return ErrFilterSizesDiffer
		}

		expectedChunks := (size + 31) >> 5
		if len(filter.bits) > 0 && len(filter.bits) != expectedChunks {
			return ErrIrrelevantLength
		}
	}

	return nil
}

// Номер задачи из задания: 3
// Краткое название: фильтр Блюма, предусматривающий удаление элементов

const MaxCounterVal = 15

type CountingBloomFilter struct {
	filterLen int
	counters  []uint32 // 8 counters per word
}

func NewCountingBloomFilter(flen int) *CountingBloomFilter {
	filterLen := flen
	if filterLen <= 0 {
		filterLen = 1
	}

	return &CountingBloomFilter{
		filterLen: filterLen,
		counters:  initCounters(countersLen(filterLen)),
	}
}

func countersLen(filterLen int) int {
	return (filterLen + 7) >> 3
}

func initCounters(size int) []uint32 {
	return make([]uint32, size)
}

func wordIdx(idx uint) uint {
	return idx >> 3
}

func shift(idx uint) uint {
	return (idx & 7) << 2
}

func (cbf *CountingBloomFilter) getCount(idx uint) uint8 {
	if len(cbf.counters) == 0 || idx >= uint(cbf.filterLen) {
		return 0
	}

	word := cbf.counters[wordIdx(idx)]
	val := (word >> shift(idx)) & 0xF

	return uint8(val)
}

func (cbf *CountingBloomFilter) setCount(idx uint, value uint8) {
	if idx >= uint(cbf.filterLen) {
		return
	}

	if len(cbf.counters) == 0 && cbf.filterLen > 0 {
		cbf.counters = initCounters(countersLen(cbf.filterLen))
	}

	word := cbf.counters[wordIdx(idx)]
	word &^= (0xF << shift(idx))
	word |= (uint32(value) << shift(idx))

	cbf.counters[wordIdx(idx)] = word
}

func (cbf *CountingBloomFilter) incCount(idx uint) {
	curr := cbf.getCount(idx)

	if curr < MaxCounterVal {
		cbf.setCount(idx, curr+1)
	}

}

func (cbf *CountingBloomFilter) decCount(idx uint) {
	curr := cbf.getCount(idx)

	if curr > 0 && curr < MaxCounterVal {
		cbf.setCount(idx, curr-1)
	}
}

func (cbf *CountingBloomFilter) commonHash(s string, magicMultiplier int) int {
	if cbf.filterLen <= 0 {
		return 0
	}

	sum := 0
	for _, char := range s {
		code := int(char)
		sum = (sum*magicMultiplier + code) % cbf.filterLen
	}

	return sum

}

func (cbf *CountingBloomFilter) Hash1(s string) int {
	return cbf.commonHash(s, hash1MagicMultiplier)
}

func (cbf *CountingBloomFilter) Hash2(s string) int {
	return cbf.commonHash(s, hash2MagicMultiplier)
}

func (cbf *CountingBloomFilter) Add(s string) {
	idx1, idx2 := cbf.Hash1(s), cbf.Hash2(s)

	cbf.incCount(uint(idx1))
	cbf.incCount(uint(idx2))
}

func (cbf *CountingBloomFilter) IsValue(s string) bool {
	if len(cbf.counters) == 0 {
		return false
	}

	idx1, idx2 := cbf.Hash1(s), cbf.Hash2(s)

	return cbf.getCount(uint(idx1)) > 0 && cbf.getCount(uint(idx2)) > 0
}

func (cbf *CountingBloomFilter) Remove(s string) {
	if !cbf.IsValue(s) {
		return
	}

	idx1, idx2 := cbf.Hash1(s), cbf.Hash2(s)

	cbf.decCount(uint(idx1))
	cbf.decCount(uint(idx2))
}

// Номер задачи из задания: 4
// Краткое название: алгоритм, который анализирует конфигурацию фильтра Блюма и пытается,
// насколько возможно, восстановить исходное множество
//
// Размышления.
// 1. Мы отображаем бесконечное множество строк в конечное множество состояний фильтра (2^32).
// 2. В битовом массиве не хранится информация о том, сколько было добавлено строк и в каком порядке.
// Получатся, без доп. условий, однозначно восстановить исходный словарь невозможно.
// Но, как я это понимаю, можно найти подмножество входных множеств,
// которые при добавлении в фильтр дадут одинаковую битовую маску (или близкую в ней).

var (
	ErrIncompleteUniverse = errors.New("universe is incomplete")
)

type candidate struct {
	word string
	h1   uint
	h2   uint
}

func ReconstructSet(target *BloomFilter, universe []string) ([]string, error) {
	var result []string

	if target == nil {
		return result, ErrNilFilter
	}

	if target.filter_len <= 0 {
		return result, ErrInvalidFilterSize
	}

	seen := make(map[string]struct{}, len(universe)) // страховка от дубликатов в universe

	// Сформировать список кандидатов (в т.ч. ложноположительные)
	var candidates []candidate

	for _, word := range universe {
		if _, exists := seen[word]; exists {
			continue
		}

		seen[word] = struct{}{}

		if !target.IsValue(word) {
			continue
		}

		h1 := uint(target.Hash1(word))
		h2 := uint(target.Hash2(word))
		candidates = append(candidates, candidate{word: word, h1: h1, h2: h2})
	}

	// Найти обязательные элементы
	bitOwners := make([][]int, target.filter_len)

	for candIdx, cand := range candidates {
		bitOwners[cand.h1] = append(bitOwners[cand.h1], candIdx)

		if cand.h1 != cand.h2 {
			bitOwners[cand.h2] = append(bitOwners[cand.h2], candIdx)
		}
	}

	// Найти минимальное покрытие множества
	covered := make([]bool, target.filter_len) // какие биты target уже закрыты
	selected := make([]bool, len(candidates))  // выбранные как результат

	for b := range target.filter_len {
		// этот бит "поджигает" больше одного кандидата (или ни один)
		if !target.testBit(uint(b)) || len(bitOwners[b]) != 1 {
			continue
		}

		candIdx := bitOwners[b][0]
		if selected[candIdx] {
			continue
		}

		selected[candIdx] = true

		cand := candidates[candIdx]
		result = append(result, cand.word)
		covered[cand.h1] = true
		covered[cand.h2] = true
	}

	// Найти кандидатов на непокрытые биты
	for {
		bestIdx := findBestCandidate(candidates, selected, covered)

		if bestIdx == -1 {
			break
		}

		selected[bestIdx] = true

		cand := candidates[bestIdx]
		result = append(result, cand.word)
		covered[cand.h1] = true
		covered[cand.h2] = true
	}

	// Проверка полноты покрытия
	if hasUncoveredBits(target, covered) {
		return result, ErrIncompleteUniverse
	}

	return result, nil
}

func findBestCandidate(candidates []candidate, selected, covered []bool) int {
	bestIdx := -1
	bestGain := 0

	for i, cand := range candidates {
		if selected[i] {
			continue
		}

		gain := countUncoveredBits(cand, covered)
		if gain > bestGain {
			bestGain = gain
			bestIdx = i
		}
	}

	return bestIdx
}

func countUncoveredBits(cand candidate, covered []bool) int {
	gain := 0

	if !covered[cand.h1] {
		gain++
	}

	if cand.h1 != cand.h2 && !covered[cand.h2] {
		gain++
	}

	return gain
}

func hasUncoveredBits(target *BloomFilter, covered []bool) bool {
	for b := range target.filter_len {
		if target.testBit(uint(b)) && !covered[b] {
			return true
		}
	}

	return false
}

// Рефлексия
//
