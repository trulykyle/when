package common_test

import (
	"testing"
	"time"

	"github.com/olebedev/when"
	"github.com/olebedev/when/rules"
	"github.com/olebedev/when/rules/common"
	"github.com/stretchr/testify/require"
)

func TestDateTimeRule(t *testing.T) {
	refTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	fixt := []Fixture{
		// Timezone formats
		{"The Deadline is 2025-04-01T13:00:00+08", 16, "2025-04-01T13:00:00+08", (90*24*time.Hour + 5*time.Hour)},
		{"The Deadline is 2025-04-01T13:00:00-05", 16, "2025-04-01T13:00:00-05", (90*24*time.Hour + 18*time.Hour)},
		{"The Deadline is 2025-04-01T13:00:00Z", 16, "2025-04-01T13:00:00Z", (90*24*time.Hour + 13*time.Hour)},

		// Timezone-naive formats
		{"The Deadline is 2025-04-28 15:00", 16, "2025-04-28 15:00", (117*24*time.Hour + 15*time.Hour)},
		{"The Deadline is 2025-04-28 3pm", 16, "2025-04-28 3pm", (117*24*time.Hour + 15*time.Hour)},
		{"The Deadline is 2025-04-28 09:30:45", 16, "2025-04-28 09:30:45", (117*24*time.Hour + 9*time.Hour + 30*time.Minute + 45*time.Second)},
	}

	nilFixt := []Fixture{
		{"The Deadline is 2025-04-28 25:00", 16, "invalid hour", 0},
		{"The Deadline is 2025-04-28 12:60", 16, "invalid minute", 0},
		{"The Deadline is 2025-04-28 12:00:99", 16, "invalid second", 0},
	}

	w := when.New(nil)
	w.Add(common.DateTimeRule(rules.Override))

	t.Run("Valid", func(t *testing.T) {
		for i, f := range fixt {
			res, err := w.Parse(f.Text, refTime)
			t.Logf("%v vs %v\n", f.Text, refTime)
			require.Nil(t, err, "[%d] error", i)
			require.NotNil(t, res, "[%d] result", i)
			require.Equal(t, f.Diff, res.Time.Sub(refTime), "[%d] duration mismatch", i)
		}
	})

	t.Run("Invalid", func(t *testing.T) {
		for i, f := range nilFixt {
			res, err := w.Parse(f.Text, refTime)
			require.Nil(t, err, "[%d] error", i)
			require.Nil(t, res, "[%d] should be nil", i)
		}
	})
}
