package parser

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

func parseJSON(line string) (Entry, error) {
	var raw map[string]interface{}
	if err := json.Unmarshal([]byte(line), &raw); err != nil {
		return Entry{Raw: line}, fmt.Errorf("parser: invalid JSON: %w", err)
	}
	e := Entry{Raw: line, Fields: make(map[string]string)}
	for _, key := range []string{"msg", "message"} {
		if v, ok := raw[key]; ok {
			e.Message = fmt.Sprintf("%v", v)
			delete(raw, key)
			break
		}
	}
	for _, key := range []string{"level", "severity"} {
		if v, ok := raw[key]; ok {
			e.Level = strings.ToUpper(fmt.Sprintf("%v", v))
			delete(raw, key)
			break
		}
	}
	for _, key := range []string{"time", "ts", "timestamp"} {
		if v, ok := raw[key]; ok {
			if s, ok2 := v.(string); ok2 {
				if t, err := time.Parse(time.RFC3339, s); err == nil {
					e.Timestamp = t
				}
			}
			delete(raw, key)
			break
		}
	}
	for k, v := range raw {
		e.Fields[k] = fmt.Sprintf("%v", v)
	}
	return e, nil
}

func parseLogfmt(line string) (Entry, error) {
	e := Entry{Raw: line, Fields: make(map[string]string)}
	for _, pair := range strings.Fields(line) {
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) != 2 {
			continue
		}
		k, v := parts[0], strings.Trim(parts[1], `"`)
		switch strings.ToLower(k) {
		case "msg", "message":
			e.Message = v
		case "level", "severity":
			e.Level = strings.ToUpper(v)
		case "time", "ts":
			if t, err := time.Parse(time.RFC3339, v); err == nil {
				e.Timestamp = t
			}
		default:
			e.Fields[k] = v
		}
	}
	return e, nil
}

func parseCommon(line string) (Entry, error) {
	m := commonLogRe.FindStringSubmatch(line)
	if m == nil {
		return Entry{Raw: line}, fmt.Errorf("parser: line does not match common log format")
	}
	e := Entry{
		Raw:     line,
		Message: m[3],
		Fields:  map[string]string{"host": m[1], "status": m[4], "bytes": m[5]},
	}
	if t, err := time.Parse("02/Jan/2006:15:04:05 -0700", m[2]); err == nil {
		e.Timestamp = t
	}
	return e, nil
}

func parseAuto(line string) (Entry, error) {
	if strings.HasPrefix(strings.TrimSpace(line), "{") {
		return parseJSON(line)
	}
	if strings.Contains(line, "=") {
		return parseLogfmt(line)
	}
	return parseCommon(line)
}
