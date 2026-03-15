package utils

// BuildProplist builds proplist argument for RouterOS API
// If proplist is empty, returns default proplist
// If proplist is "detail" or contains empty string, returns empty (all fields)
// Otherwise, joins the provided fields
func BuildProplist(defaultProplist string, proplist ...string) string {
	if len(proplist) == 0 {
		return defaultProplist
	}

	// Check if detail mode (empty string or "detail")
	for _, p := range proplist {
		if p == "" || p == "detail" {
			return "" // Return empty to get all fields
		}
	}

	// Join provided fields
	result := ""
	for i, p := range proplist {
		if i > 0 {
			result += ","
		}
		result += p
	}
	return result
}

// AppendProplist appends proplist to args if not empty
func AppendProplist(args []string, proplist string) []string {
	if proplist != "" {
		return append(args, "=.proplist="+proplist)
	}
	return args
}
