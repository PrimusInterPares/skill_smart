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

// Time complexity: O(K * N), где K - кол-во фильтров, N - размер фильтров. В данной реализации можно прированять к O(1)
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
