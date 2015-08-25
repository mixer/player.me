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

type Direction int

const (
    Forward Direction = iota
    Backward
)

func (d Direction) Next(pager Pager) int {
		var next int
        switch d {
        case Forward:
                next = pager.Limit + pager.From
				if next >= pager.Total {
					return -1
				}

                return next
        case Backward:
                if pager.Limit > pager.From {
					next = 0
				} else {
					next = pager.From - pager.Limit
				}
				if pager.From == 0 {
					return -1
				}
				return next
        }

        return -1
}


// Returns a GameList of games from the API.
func (c *Client) GameList(queries ...Query) (*GameListResult, error) {
	result := &GameListResult{}
	err := c.Request(GET, "/api/v1/games", grabQuery(queries)).Run(result)

	return result, err
}

// Returns a channel of games (and errors) which can be iterated
// over transparently, without manually dealing with pagination.
func (c *Client) GameListIterate(d Direction, queries ...Query) *GameListIterator {
	query := grabPaginateQuery(queries)
	it := &GameListIterator{
		Iterator: newIterator(),
		Results:  make(chan Game),
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
			next := d.Next(result.Pager)

			if next == -1 {
				close(it.Results)
				return
			}
			

			// Update the _from in the pagination query to get
			// the next page.
			query["_from"] = next
		}
	}()

	return it
}
