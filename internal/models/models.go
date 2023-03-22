package models

import "time"

type FeedFull struct {
	Feed []struct {
		Items []struct {
			AppID     int    `json:"appID"`
			ID        string `json:"id"`
			Developer struct {
				ID            int    `json:"id"`
				Name          string `json:"name"`
				ResidenceInfo struct {
					Eu bool `json:"EU"`
				} `json:"residence_info"`
			} `json:"developer"`
			Rating          float64 `json:"rating"`
			RatingCount     int     `json:"ratingCount"`
			Retention       int     `json:"retention"`
			Title           string  `json:"title"`
			URL             string  `json:"url"`
			CategoryIDs     []int   `json:"categoryIDs"`
			CategoriesNames any     `json:"categoriesNames"`
			Features        struct {
				SupportsSw  bool `json:"supports-sw"`
				NativeCache bool `json:"native_cache"`
			} `json:"features"`
			TagIDs []int `json:"tagIDs"`
			Media  struct {
				Cover struct {
					PrefixURL string `json:"prefix-url"`
					MainColor string `json:"mainColor"`
				} `json:"cover"`
				Icon struct {
					PrefixURL string `json:"prefix-url"`
					MainColor string `json:"mainColor"`
				} `json:"icon"`
				Videos []struct {
					EmbedURL           string `json:"embedUrl"`
					ThumbnailURL       string `json:"thumbnailUrl"`
					ThumbnailURLPrefix string `json:"thumbnailUrlPrefix"`
					StreamURL          string `json:"streamUrl"`
					PreviewURL         string `json:"previewUrl"`
					Mp4StreamURL       string `json:"mp4StreamUrl"`
					Height             int    `json:"height"`
					Width              int    `json:"width"`
				} `json:"videos"`
			} `json:"media"`
			PlayersCount int    `json:"playersCount"`
			Type         string `json:"type"`
			Column       int    `json:"column"`
			Row          int    `json:"row"`
		} `json:"items"`
		Type         string `json:"type"`
		Category     string `json:"category"`
		Title        string `json:"title"`
		HasTitle     bool   `json:"hasTitle"`
		Size         string `json:"size"`
		PageNumber   int    `json:"pageNumber"`
		HasMoreItems bool   `json:"hasMoreItems"`
	} `json:"feed"`
	PageInfo struct {
		IsFirstPage bool   `json:"isFirstPage"`
		HasNextPage bool   `json:"hasNextPage"`
		NextPageID  string `json:"nextPageId"`
		RtxReqID    string `json:"rtxReqId"`
	} `json:"pageInfo"`
	GamesWithPromos  int    `json:"gamesWithPromos"`
	UsedDJExperiment string `json:"usedDJExperiment"`
	LastPlayedTS     int    `json:"lastPlayedTS"`
}

type FeedDevelopers struct {
	Feed []struct {
		Items []struct {
			Developer struct {
				Name string `json:"name"`
			} `json:"developer"`
		} `json:"items"`
	} `json:"feed"`
}

type FeedGames struct {
	Feed []struct {
		Items []struct {
			Developer struct {
				Name string `json:"name"`
			} `json:"developer"`
			PlayersCount int    `json:"playersCount"`
			Title        string `json:"title"`
			URL          string `json:"url"`
			Retention    int    `json:"retention"`
			Media        struct {
				Cover struct {
					PrefixURL string `json:"prefix-url"`
					MainColor string `json:"mainColor"`
				} `json:"cover"`
				Icon struct {
					PrefixURL string `json:"prefix-url"`
					MainColor string `json:"mainColor"`
				} `json:"icon"`
				Videos []struct {
					EmbedURL           string `json:"embedUrl"`
					ThumbnailURL       string `json:"thumbnailUrl"`
					ThumbnailURLPrefix string `json:"thumbnailUrlPrefix"`
					StreamURL          string `json:"streamUrl"`
					PreviewURL         string `json:"previewUrl"`
					Mp4StreamURL       string `json:"mp4StreamUrl"`
					Height             int    `json:"height"`
					Width              int    `json:"width"`
				} `json:"videos"`
			} `json:"media"`
		} `json:"items"`
	} `json:"feed"`
}

type Developer struct {
	ID   int
	Name string
}

type Game struct {
	ID        int
	Name      string
	Img       string
	Players   int64
	Trend     int64
	Developer string
	CreatedAt time.Time
	URL       string
}
