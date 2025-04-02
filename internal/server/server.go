package server

import (
	"context"
	"findx/config"
	"findx/internal/lockdb"
	"findx/internal/search"
	"findx/pkg/protogen"
	"fmt"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type SearchServer struct {
	protogen.UnimplementedSearchServiceServer
	keyManager   *search.ApiKeyManager
	googleClient *search.Client

	lockDb      lockdb.ILockDb
	rateLimiter lockdb.RateLimiter
}

func NewSearchServer(conf *config.Config, lockDb lockdb.ILockDb, rateLimiter lockdb.RateLimiter) *SearchServer {
	return &SearchServer{
		googleClient: search.NewClient(),
		keyManager:   search.NewApiKeyManager(conf.POSTGRES_DSN, lockDb, rateLimiter),
		lockDb:       lockDb,
		rateLimiter:  rateLimiter,
	}
}

func (s *SearchServer) Search(ctx context.Context, req *protogen.SearchRequest) (*protogen.SearchResponse, error) {
	if req.Query == "" {
		return nil, status.Error(codes.InvalidArgument, "query is required")
	}

	// Number of bucket can be configured
	bucketList, err := s.keyManager.GetKeyBucket(ctx, 5)
	if err != nil {
		return nil, status.Error(codes.ResourceExhausted, err.Error())
	}
	availableKey, err := s.keyManager.GetAvailableKey(ctx, bucketList)
	if err != nil {
		return nil, status.Errorf(codes.ResourceExhausted, "no available API keys")
	}
	defer func() {
		var (
			dateOnlyCurrentTime = time.Now().UTC().
						Add(-7 * time.Hour).Truncate(24 * time.Hour)
			dateOnlyResetedAt = availableKey.ResetedAt.UTC().
						Add(-7 * time.Hour).Truncate(24 * time.Hour)
		)
		if !dateOnlyCurrentTime.
			After(dateOnlyResetedAt) {
			return
		}
		goCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		go func(ctx context.Context) {
			//TODO: handle error
			ourLock, err := s.lockDb.LockSimple(ctx, "search:get_key:reset")
			if err != nil {
				return
			}
			defer ourLock.Unlock()
			s.keyManager.ResetDailyCounts(100, availableKey.ResetedAt)
		}(goCtx)
	}()

	var (
		weShouldDoSomething bool
	)
	if weShouldDoSomething = bucketList.Avg() < 10; weShouldDoSomething {
		fmt.Println("nooooooooo")
	}

	params := map[string]string{
		"lr":  fmt.Sprintf("lang_%s", req.Language),
		"cr":  req.Country,
		"num": fmt.Sprintf("%d", req.NumResults),
	}

	results, err := s.googleClient.Search(ctx, availableKey.ApiKey, availableKey.EngineId, req.Query, params)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "search failed: %v", err)
	}

	response := &protogen.SearchResponse{
		Results: make([]*protogen.SearchResult, 0, len(results)),
	}

	for _, r := range results {
		response.Results = append(response.Results, &protogen.SearchResult{
			Title:   r.Title,
			Link:    r.Link,
			Snippet: r.Snippet,
		})
	}

	return response, nil
}

func (s *SearchServer) SearchParams(ctx context.Context, req *protogen.GoogleSearchParamRequest) (*protogen.GoogleSearchParamResponse, error) {
	return &protogen.GoogleSearchParamResponse{
		Results: []*protogen.GoogleSearchParam{
			{
				Param: "c2coff",
				Type:  "string",
				Description: `
				The default value for this parameter is 0 (zero), meaning that the feature is enabled. Supported values are:
					1: Disabled
					0: Enabled (default)
				`,
			},
			{
				Param: "cr",
				Type:  "string",
				Description: `
					Restricts search results to documents originating in a particular country. You may use Boolean operators in the cr parameter's value.
					Google Search determines the country of a document by analyzing:
						the top-level domain (TLD) of the document's URL
						the geographic location of the Web server's IP address
					See the Country Parameter Values page for a list of valid values for this parameter.
				`,
			},
			{
				Param: "cx",
				Type:  "string",
				Description: `
					The Programmable Search Engine ID to use for this request.
				`,
			},
			{
				Param: "dateRestrict",
				Type:  "string",
				Description: `
				Restricts results to URLs based on date. Supported values include:
					d[number]: requests results from the specified number of past days.
					w[number]: requests results from the specified number of past weeks.
					m[number]: requests results from the specified number of past months.
					y[number]: requests results from the specified number of past years.
				`,
			},
			{
				Param: "exactTerms",
				Type:  "string",
				Description: `
					Identifies a phrase that all documents in the search results must contain.
				`,
			},
			{
				Param: "excludeTerms",
				Type:  "string",
				Description: `
					Identifies a word or phrase that should not appear in any documents in the search results.
				`,
			},
			{
				Param: "fileType",
				Type:  "string",
				Description: `
					Restricts results to files of a specified extension. A list of file types indexable by Google can be found in Search Console Help Center.
				`,
			},
			{
				Param: "filter",
				Type:  "string",
				Description: `
				Controls turning on or off the duplicate content filter.
				See Automatic Filtering for more information about Google's search results filters. Note that host crowding filtering applies only to multi-site searches.
				By default, Google applies filtering to all search results to improve the quality of those results.
				Acceptable values are:
					0: Turns off duplicate content filter.
					1: Turns on duplicate content filter.
				`,
			},
			{
				Param: "gl",
				Type:  "string",
				Description: `
				Geolocation of end user.
					The gl parameter value is a two-letter country code. The gl parameter boosts search results whose country of origin matches the parameter value. See the Country Codes page for a list of valid values.
					Specifying a gl parameter value should lead to more relevant results. This is particularly true for international customers and, even more specifically, for customers in English- speaking countries other than the United States.
				`,
			},
			{
				Param: "highRange",
				Type:  "string",
				Description: `
				Specifies the ending value for a search range.
					Use lowRange and highRange to append an inclusive search range of lowRange...highRange to the query.
				`,
			},
			{
				Param: "hl",
				Type:  "string",
				Description: `
				Sets the user interface language.
					Explicitly setting this parameter improves the performance and the quality of your search results.
					See the Interface Languages section of Internationalizing Queries and Results Presentation for more information, and Supported Interface Languages for a list of supported languages.
				`,
			},
			{
				Param: "hq",
				Type:  "string",
				Description: `
					Appends the specified query terms to the query, as if they were combined with a logical AND operator.
				`,
			},
			{
				Param: "imgColorType",
				Type:  "enum",
				Description: `
				Returns black and white, grayscale, transparent, or color images. Acceptable values are:
					"color"
					"gray"
					"mono": black and white
					"trans": transparent background
				`,
			},
			{
				Param: "imgDominantColor",
				Type:  "enum",
				Description: `
				Returns images of a specific dominant color. Acceptable values are:
					"black"
					"blue"
					"brown"
					"gray"
					"green"
					"orange"
					"pink"
					"purple"
					"red"
					"teal"
					"white"
					"yellow"
				`,
			},
			{
				Param: "imgSize",
				Type:  "enum",
				Description: `
				Returns images of a specified size. Acceptable values are:
					"huge"
					"icon"
					"large"
					"medium"
					"small"
					"xlarge"
					"xxlarge"
				`,
			},
			{
				Param: "imgType",
				Type:  "enum",
				Description: `
				Returns images of a type. Acceptable values are:
					"clipart"
					"face"
					"lineart"
					"stock"
					"photo"
					"animated"
				`,
			},
			{
				Param: "linkSite",
				Type:  "enum",
				Description: `
					Specifies that all search results should contain a link to a particular URL.
				`,
			},
			{
				Param: "lowRange",
				Type:  "string",
				Description: `
					Specifies the starting value for a search range. Use lowRange and highRange to append an inclusive search range of lowRange...highRange to the query.
				`,
			},
			{
				Param: "lr",
				Type:  "string",
				Description: `
				Restricts the search to documents written in a particular language (e.g., lr=lang_ja).
					Acceptable values are:
					"lang_ar": Arabic
					"lang_bg": Bulgarian
					"lang_ca": Catalan
					"lang_cs": Czech
					"lang_da": Danish
					"lang_de": German
					"lang_el": Greek
					"lang_en": English
					"lang_es": Spanish
					"lang_et": Estonian
					"lang_fi": Finnish
					"lang_fr": French
					"lang_hr": Croatian
					"lang_hu": Hungarian
					"lang_id": Indonesian
					"lang_is": Icelandic
					"lang_it": Italian
					"lang_iw": Hebrew
					"lang_ja": Japanese
					"lang_ko": Korean
					"lang_lt": Lithuanian
					"lang_lv": Latvian
					"lang_nl": Dutch
					"lang_no": Norwegian
					"lang_pl": Polish
					"lang_pt": Portuguese
					"lang_ro": Romanian
					"lang_ru": Russian
					"lang_sk": Slovak
					"lang_sl": Slovenian
					"lang_sr": Serbian
					"lang_sv": Swedish
					"lang_tr": Turkish
					"lang_zh-CN": Chinese (Simplified)
					"lang_zh-TW": Chinese (Traditional)
				`,
			},
			{
				Param: "num",
				Type:  "integer",
				Description: `
				Number of search results to return.
					Valid values are integers between 1 and 10, inclusive.
				`,
			},
			{
				Param: "orTerms",
				Type:  "string",
				Description: `
					Provides additional search terms to check for in a document, where each document in the search results must contain at least one of the additional search terms.
				`,
			},
			{
				Param:       "q",
				Type:        "string",
				Description: "query",
			},
			{
				Param: "rights",
				Type:  "string",
				Description: `
					Filters based on licensing. Supported values include: cc_publicdomain, cc_attribute, cc_sharealike, cc_noncommercial, cc_nonderived and combinations of these. See typical combinations.
				`,
			},
			{
				Param: "safe",
				Type:  "enum",
				Description: `
				Search safety level. Acceptable values are:
					"active": Enables SafeSearch filtering.
					"off": Disables SafeSearch filtering. (default)
				`,
			},
			{
				Param: "searchType",
				Type:  "enum",
				Description: `
				Specifies the search type: image. If unspecified, results are limited to webpages.
				Acceptable values are:
					"image": custom image search.
				`,
			},
			{
				Param: "siteSearch",
				Type:  "string",
				Description: `
					Specifies a given site which should always be included or excluded from results (see siteSearchFilter parameter, below).
				`,
			},
			{
				Param: "siteSearchFilter",
				Type:  "enum",
				Description: `
				Controls whether to include or exclude results from the site named in the siteSearch parameter.
				Acceptable values are:
					"e": exclude
					"i": include
				`,
			},
			{
				Param: "sort",
				Type:  "string",
				Description: `
				The sort expression to apply to the results. The sort parameter specifies that the results be sorted according to the specified expression i.e. sort by date. Example: sort=date.
				`,
			},
			{
				Param: "start",
				Type:  "integer (uint32 format)",
				Description: `
					The index of the first result to return. The default number of results per page is 10, so &start=11 would start at the top of the second page of results. Note: The JSON API will never return more than 100 results, even if more than 100 documents match the query, so setting the sum of start + num to a number greater than 100 will produce an error. Also note that the maximum value for num is 10.
				`,
			},
		},
	}, nil
}
