package player

type Game struct {
	Alias string
	Box   struct {
		Cached           string
		Original         string
		OriginalFilename string `json:"original_filename"`
	}
	Description    string
	FavoritesCount int `json:"favourites_count"`
	HasFavorited   bool
	HasLiked       bool
	LikesCount     int `json:"likes_count"`
	Slug           string
	Title          string
	URL            string
}

type GameListResult struct {
	PagedResult
	Results []Game
}

type GameListIterator struct {
	Iterator
	Results chan Game
	reverse bool
}

// Returns a GameList of games from the API.
func (c *Client) GameList(queries ...Query) (*GameListResult, error) {
	result := &GameListResult{}
	err := c.Request(GET, "/api/v1/games", grabQuery(queries)).Run(result)

	return result, err
}

// Returns a channel of games (and errors) which can be iterated
// over transparently, without manually dealing with pagination.
func (c *Client) GameListIterate(reverse bool, queries ...Query) *GameListIterator {
	query := grabPaginateQuery(queries)
	it := &GameListIterator{
		Iterator: newIterator(),
		Results:  make(chan Game),
		reverse: reverse,
	}

	go func() {
		for {
			// Grab the result, and fail if we got an error.
			result, err := c.GameList(query)
			if err != nil {
				it.Errors <- err
				close(it.Results)
				return
			}

			// Send all the results back down the channel, or quit
			// if no onw is reading the channel.
			for _, r := range result.Results {
				select {
				case <-it.Quit:
					return
				case it.Results <- r:
					// ignore
				}
			}

			// Close the results and end if we're at the end.
			var next int
			if it.reverse {
				if result.Pager.Limit > result.Pager.From {
					next = 0
				} else {
					next = result.Pager.From - result.Pager.Limit
				}
				if result.Pager.From == 0 {
					close(it.Results)
					return
				}
			} else {
				next = result.Pager.Limit + result.Pager.From
				if next >= result.Pager.Total {
					close(it.Results)
					return
				}
			}
			

			// Update the _from in the pagination query to get
			// the next page.
			query["_from"] = next
		}
	}()

	return it
}
