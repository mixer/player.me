package player

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
)

// The Query can be provided to pass values to the API.
type Query map[string]interface{}

// Converts the query values to
func (q Query) toParams() (url.Values, error) {
	values := url.Values{}

	for key, value := range q {
		if err := addToValues(&values, key, value); err != nil {
			return nil, err
		}
	}

	return values, nil
}

// Adds a key/value pair to the url values, returning an error if it's
// an unexpected type.
func addToValues(values *url.Values, key string, value interface{}) error {
	switch v := value.(type) {
	case string:
		values.Set(key, v)
	case bool:
		values.Set(key, fmt.Sprintf("%t", v))
	case uint8, uint16, uint32, uint64, int8, int16, int32, int64, int, uint:
		values.Set(key, fmt.Sprintf("%d", v))
	default:
		return errors.New("Unknown type for key " + key)
	}

	return nil
}

// Returns the first query in the slice, or nil. Useful for allowing
// optional queries.
func grabQuery(queries []Query) Query {
	var query Query
	if len(queries) > 0 {
		query = queries[0]
	}

	return query
}

// Returns the first query in the slice, and ensures the
// _from is set.
func grabPaginateQuery(queries []Query) Query {
	query := grabQuery(queries)
	if len(queries) == 0 {
		return Query{
			"_from":  0,
			"_limit": 100,
		}
	}

	// Convert pages to using `from`
	if page, ok := query["_page"]; ok {
		if limit, ok := query["_limit"]; ok {
			query["_from"] = page.(int) * limit.(int)
		} else {
			query["_from"] = page.(int) * 10
		}
		delete(query, "_page")
	}
	if _, ok := query["_from"]; !ok {
		query["_from"] = 0
	}

	query["_limit"] = 100

	return query
}

type Request struct {
	base      string
	path      string
	method    string
	params    Query
	fulfiller Fulfiller
}

func (r *Request) Run(target interface{}) error {
	result, err := r.fulfiller.Request(r.method, r.base+r.path, r.params)
	if err != nil {
		return err
	}

	return json.Unmarshal(result, &target)
}
