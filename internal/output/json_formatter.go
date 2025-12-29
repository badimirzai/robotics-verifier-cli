package output

import "encoding/json"

// FormatJSON renders the payload as JSON, optionally pretty printed.
func FormatJSON(payload any, pretty bool) ([]byte, error) {
	if pretty {
		return json.MarshalIndent(payload, "", "    ")
	}
	return json.Marshal(payload)
}
