package booking_common

import "fmt"

func parseFTSearchResult(result interface{}) ([]map[string]string, error) {
	resultMap, ok := result.(map[interface{}]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid result type: %T", result)
	}

	totalResults, _ := resultMap["total_results"].(int64)
	if totalResults == 0 {
		return []map[string]string{}, nil
	}

	resultsInterface, ok := resultMap["results"]
	if !ok {
		return nil, fmt.Errorf("no results field")
	}

	results, ok := resultsInterface.([]interface{})
	if !ok {
		return nil, fmt.Errorf("results is not array")
	}

	var captains []map[string]string

	for _, resInterface := range results {
		resMap, ok := resInterface.(map[interface{}]interface{})
		if !ok {
			continue
		}

		id, _ := resMap["id"].(string)

		extraAttrsInterface, ok := resMap["extra_attributes"]
		if !ok {
			continue
		}

		extraAttrs, ok := extraAttrsInterface.(map[interface{}]interface{})
		if !ok {
			continue
		}
		captain := make(map[string]string)
		captain["id"] = id
		for k, _ := range extraAttrs {
			if _, exists := captain[k.(string)]; !exists {
				captain[k.(string)] = getInterfaceMapValue(extraAttrs, k.(string))
			}
		}
		captains = append(captains, captain)
	}

	return captains, nil
}

func getInterfaceMapValue(m map[interface{}]interface{}, key string) string {
	if val, ok := m[key]; ok {
		switch v := val.(type) {
		case string:
			return v
		case int:
			return fmt.Sprintf("%d", v)
		case int64:
			return fmt.Sprintf("%d", v)
		case float64:
			return fmt.Sprintf("%v", v)
		default:
			return fmt.Sprintf("%v", v)
		}
	}
	return ""
}

func parseLocation(loc string) (float64, float64) {
	var longitude, latitude float64
	fmt.Sscanf(loc, "%f,%f", &longitude, &latitude)
	return longitude, latitude
}
