package types

import (
	"testing"
	"time"

	_pgtype "github.com/jackc/pgtype"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func TestTypes(t *testing.T) {
	t.Run("FormatDuration", func(t *testing.T) {
		inputs := []time.Duration{
			time.Duration(0),
			// time.Nanosecond * 10,
			time.Microsecond * 10,
			time.Millisecond * 10,
			time.Second * 10,
			time.Minute * 10,
			time.Hour * 10,
		}

		expectedOutputs := []pgtype.Interval{
			{Microseconds: 0, Days: 0, Months: 0, Valid: true},
			// {Microseconds: 0, Days: 0, Months: 0, Valid: true},
			{Microseconds: 10, Days: 0, Months: 0, Valid: true},
			{Microseconds: 10 * 1000, Days: 0, Months: 0, Valid: true},
			{Microseconds: 10 * 1000 * 1000, Days: 0, Months: 0, Valid: true},
			{Microseconds: 10 * 1000 * 1000 * 60, Days: 0, Months: 0, Valid: true},
			{Microseconds: 10 * 1000 * 1000 * 60 * 60, Days: 0, Months: 0, Valid: true},
		}

		for i, input := range inputs {
			expectedOutput := expectedOutputs[i]

			actualOutput, err := FormatDuration(input)
			require.NoError(t, err)
			require.Equal(t, expectedOutput, actualOutput)
		}
	})

	t.Run("ParseDuration1", func(t *testing.T) {
		inputs := []pgtype.Interval{
			{Microseconds: 0, Days: 0, Months: 0, Valid: true},
			// {Microseconds: 0, Days: 0, Months: 0, Valid: true},
			{Microseconds: 10, Days: 0, Months: 0, Valid: true},
			{Microseconds: 10 * 1000, Days: 0, Months: 0, Valid: true},
			{Microseconds: 10 * 1000 * 1000, Days: 0, Months: 0, Valid: true},
			{Microseconds: 10 * 1000 * 1000 * 60, Days: 0, Months: 0, Valid: true},
			{Microseconds: 10 * 1000 * 1000 * 60 * 60, Days: 0, Months: 0, Valid: true},
		}

		expectedOutputs := []time.Duration{
			time.Duration(0),
			// time.Nanosecond * 10,
			time.Microsecond * 10,
			time.Millisecond * 10,
			time.Second * 10,
			time.Minute * 10,
			time.Hour * 10,
		}

		for i, input := range inputs {
			expectedOutput := expectedOutputs[i]

			actualOutput, err := ParseDuration(input)
			require.NoError(t, err)
			require.Equal(t, expectedOutput, actualOutput)
		}
	})

	t.Run("ParseDuration2", func(t *testing.T) {
		inputs := []_pgtype.Interval{
			{Microseconds: 0, Days: 0, Months: 0},
			// {Microseconds: 0, Days: 0, Months: 0},
			{Microseconds: 10, Days: 0, Months: 0},
			{Microseconds: 10 * 1000, Days: 0, Months: 0},
			{Microseconds: 10 * 1000 * 1000, Days: 0, Months: 0},
			{Microseconds: 10 * 1000 * 1000 * 60, Days: 0, Months: 0},
			{Microseconds: 10 * 1000 * 1000 * 60 * 60, Days: 0, Months: 0},
		}

		expectedOutputs := []time.Duration{
			time.Duration(0),
			// time.Nanosecond * 10,
			time.Microsecond * 10,
			time.Millisecond * 10,
			time.Second * 10,
			time.Minute * 10,
			time.Hour * 10,
		}

		for i, input := range inputs {
			expectedOutput := expectedOutputs[i]

			actualOutput, err := ParseDuration(input)
			require.NoError(t, err)
			require.Equal(t, expectedOutput, actualOutput)
		}
	})
}