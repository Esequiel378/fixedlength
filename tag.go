package fixedlength

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var (
	ErrTagInefectualRange    = errors.New("fixedlength: inefectual range")
	ErrTagEmpty              = errors.New("fixedlength: tag is empty")
	ErrTagInvalidRangeValues = errors.New("fixedlength: invalid range values")
	ErrTagInvalidUpperBound  = errors.New("fixedlength: invalid upper bound")
)

// parseTag splits a struct field's json tag into its name and
// comma-separated options.
func parseTag(tag string, upperBound int) (int, int, error) {
	if tag == "" {
		return 0, 0, ErrTagEmpty
	}

	if upperBound == 0 {
		return 0, 0, fmt.Errorf("%w: %d", ErrTagInvalidUpperBound, upperBound)
	}

	parts := strings.Split(tag, ",")
	x, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, errors.Join(ErrTagInvalidRangeValues, err)
	}

	y, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, errors.Join(ErrTagInvalidRangeValues, err)
	}

	start := max(x, 0)
	end := min(y, upperBound)

	// -1 is used to indicate that the end of the string should be used
	if end == -1 {
		end = upperBound
	}

	if start == end {
		return 0, 0, fmt.Errorf("%w: %s", ErrTagInefectualRange, tag)
	}

	if start > end {
		return 0, 0, fmt.Errorf("%w: x > y from %s", ErrTagInefectualRange, tag)
	}

	return start, end, nil
}
