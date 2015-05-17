package player

import (
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
)

func TestGrabsPaginated(t *testing.T) {
	assert.Equal(t,
		grabPaginateQuery([]Query{Query{}}),
		Query{"_limit": 100, "_from": 0},
		"Failed to grab from empty",
	)

	assert.Equal(t,
		grabPaginateQuery([]Query{}),
		Query{"_limit": 100, "_from": 0},
		"Failed to grab from none",
	)

	assert.Equal(t,
		grabPaginateQuery([]Query{Query{"_page": 3}}),
		Query{"_limit": 100, "_from": 30},
		"Failed to grab from page without limit",
	)

	assert.Equal(t,
		grabPaginateQuery([]Query{Query{"_page": 3, "_limit": 25}}),
		Query{"_limit": 100, "_from": 75},
		"Failed to grab from page with limit",
	)

	assert.Equal(t,
		grabPaginateQuery([]Query{Query{"_from": 300}}),
		Query{"_limit": 100, "_from": 300},
		"Gets with default from",
	)

	assert.Equal(t,
		grabPaginateQuery([]Query{Query{"_from": 300, "_extra": "foo"}}),
		Query{"_limit": 100, "_from": 300, "_extra": "foo"},
		"Passes extra",
	)
}

func TestAddsValues(t *testing.T) {
	var v *url.Values

	v = &url.Values{}
	assert.Nil(t, addToValues(v, "key", "value"))
	assert.Equal(t, url.Values{"key": []string{"value"}}, *v, "adds string to values")

	v = &url.Values{}
	assert.Nil(t, addToValues(v, "key", 42))
	assert.Equal(t, url.Values{"key": []string{"42"}}, *v, "adds int to values")

	v = &url.Values{}
	assert.Nil(t, addToValues(v, "key", true))
	assert.Equal(t, url.Values{"key": []string{"true"}}, *v, "adds bool to values")

	assert.NotNil(t, addToValues(v, "key", []int{}), "errors on unsupported")
}
