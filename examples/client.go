package main

import (
	"context"
	"log"
	"time"

	"findx/pkg/protogen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := protogen.NewSearchServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	resp, err := c.Search(ctx, &protogen.SearchRequest{
		Q:                "Crypto market sentiment today",
		Language:         "en",
		Num:              10,
		DateRestrict:     "d1", //10 days ago, can be d, w, m, y also
		ExactTerms:       "",
		ExcludeTerms:     "",
		Gl:               "en",      // from country code list: https://developers.google.com/custom-search/docs/json_api_reference#countryCodes
		Hl:               "en",      // from interface language list: https://developers.google.com/custom-search/docs/json_api_reference#interfaceLanguages
		Hq:               "",        //extra terms to append to the main query with a logical AND operator
		LinkSite:         "",        // search results should be links to this site
		Lr:               "lang_en", // strict search documents written in this language
		OrTerms:          "",        //extra terms to append to the main query with a logical OR operator
		Safe:             "off",     //off or active
		SiteSearch:       "",        //site should be included or excluded (see siteSearchFilter)
		SiteSearchFilter: "e",       //e or i
		Sort:             "",

		// StartDate:  "2023-01-01",
		// EndDate:    "2023-12-31",
	})
	if err != nil {
		log.Fatalf("could not search: %v", err)
	}

	log.Printf("Results:")
	for _, result := range resp.Results {
		log.Printf("- %s (%s)", result.Title, result.Link)
		log.Printf("  %s", result.Snippet)
	}
}
