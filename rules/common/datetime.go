package common

import (
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/AlekSi/pointer"
	"github.com/olebedev/when/rules"
)

// DateTimeRule returns a custom rule for parsing full datetime strings.
func DateTimeRule(s rules.Strategy) rules.Rule {
	overwrite := s == rules.Override

	return &rules.F{
		RegExp: regexp.MustCompile(
			`(?i)\b` +
				`(\d{4}-\d{2}-\d{2})` + // Date part
				`[ T]` + // Separator
				`(` +
				`\d{1,2}(?::\d{2}){0,2}` + // Hours, optional minutes, optional seconds
				`(?: *(?:am|pm))?` + // Optional AM/PM
				`(?: *(Z|[+-]\d{2}(:?\d{2})?))?` + // Optional timezone
				`)` +
				`\b`,
		),
		Applier: func(m *rules.Match, c *rules.Context, o *rules.Options, ref time.Time) (bool, error) {
			if (c.Year != nil || c.Month != nil || c.Day != nil) && !overwrite {
				return false, nil
			}

			datePart := m.Captures[0]
			timePart := strings.TrimSpace(m.Captures[1])

			// Extract timezone
			tzRegex := regexp.MustCompile(`(?i)(Z|[+-]\d{2}(:?\d{2})?)$`)
			tzMatch := tzRegex.FindStringSubmatch(timePart)
			var tz string
			if len(tzMatch) > 0 {
				tz = tzMatch[0]
				timePart = strings.TrimSuffix(timePart, tz)
				timePart = strings.TrimSpace(timePart)
			}

			// Extract AM/PM
			ampmRegex := regexp.MustCompile(`(?i)(am|pm)$`)
			ampmMatch := ampmRegex.FindStringSubmatch(timePart)
			var ampm string
			if len(ampmMatch) > 0 {
				ampm = ampmMatch[0]
				timePart = strings.TrimSuffix(timePart, ampm)
				timePart = strings.TrimSpace(timePart)
			}

			// Split into time components
			timeComponents := strings.Split(timePart, ":")
			numComponents := len(timeComponents)

			// Determine time format based on components and AM/PM
			var timeFormat string
			switch numComponents {
			case 1: // hours only
				if ampm != "" {
					timeFormat = "3pm"
				} else {
					timeFormat = "15"
				}
			case 2: // hours and minutes
				if ampm != "" {
					timeFormat = "3:04pm"
				} else {
					timeFormat = "15:04"
				}
			case 3: // hours, minutes, seconds
				if ampm != "" {
					timeFormat = "3:04:05pm"
				} else {
					timeFormat = "15:04:05"
				}
			default:
				log.Printf("Invalid time components: %v", timeComponents)
				return false, nil
			}

			// Append timezone to layout if present
			if tz != "" {
				if strings.EqualFold(tz, "Z") {
					timeFormat += "Z"
				} else {
					// Determine timezone format
					if strings.Contains(tz, ":") {
						timeFormat += "Z07:00"
					} else if len(tz) == 5 { // e.g., +0800
						timeFormat += "Z0700"
					} else if len(tz) == 3 { // e.g., +08
						timeFormat += "Z07"
					} else {
						log.Printf("Unsupported timezone format: %s", tz)
						return false, nil
					}
				}
			}

			// Combine date and time formats
			layout := "2006-01-02T" + timeFormat
			datetimeStr := datePart + "T" + strings.TrimSpace(m.Captures[1]) // Use original timePart with proper spacing

			t, err := time.Parse(layout, datetimeStr)
			if err != nil {
				log.Printf("Parse error: %v", err)
				return false, nil
			}

			utcTime := t.UTC()

			c.Year = pointer.ToInt(utcTime.Year())
			c.Month = pointer.ToInt(int(utcTime.Month()))
			c.Day = pointer.ToInt(utcTime.Day())
			c.Hour = pointer.ToInt(utcTime.Hour())
			c.Minute = pointer.ToInt(utcTime.Minute())
			c.Second = pointer.ToInt(utcTime.Second())

			return true, nil
		},
	}
}
