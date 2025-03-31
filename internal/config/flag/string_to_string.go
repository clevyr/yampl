package flag

import (
	"bytes"
	"encoding/csv"
	"errors"
	"io"
	"maps"
	"strings"
)

var ErrStringToStringFormat = errors.New("must be formatted as key=value")

type StringToString struct {
	value   map[string]string
	changed bool
}

// Set Format: a=1,b=2.
func (s *StringToString) Set(val string) error {
	val = strings.TrimSpace(val)
	count := strings.Count(val, "=")
	records := make([]string, 0, count)
	switch count {
	case 0:
		return ErrStringToStringFormat
	case 1:
		records = append(records, val)
	default:
		r := csv.NewReader(strings.NewReader(val))
		r.TrimLeadingSpace = true
		for {
			line, err := r.Read()
			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				}
				return err
			}

			r.FieldsPerRecord = 0 // Prevent wrong number of fields error

			for _, v := range line {
				switch {
				case strings.Contains(v, "="):
					records = append(records, v)
				case len(records) != 0:
					records[len(records)-1] += "\n" + v
				default:
					return ErrStringToStringFormat
				}
			}
		}
	}

	result := make(map[string]string, len(records))
	for _, pair := range records {
		kv := strings.SplitN(pair, "=", 2)
		if len(kv) != 2 {
			return ErrStringToStringFormat
		}
		result[kv[0]] = kv[1]
	}

	if s.changed {
		for k, v := range result {
			s.value[k] = v
		}
	} else {
		s.changed = true
		s.value = result
	}

	return nil
}

func (s *StringToString) Type() string {
	return "stringToString"
}

func (s *StringToString) String() string {
	records := make([]string, 0, len(s.value))
	for k, v := range s.value {
		records = append(records, k+"="+v)
	}

	var buf bytes.Buffer
	w := csv.NewWriter(&buf)
	if err := w.Write(records); err != nil {
		panic(err)
	}
	w.Flush()
	return "[" + strings.TrimSpace(buf.String()) + "]"
}

func (s *StringToString) Value() map[string]string {
	return maps.Clone(s.value)
}
