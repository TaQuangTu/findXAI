package server

import (
	"bytes"
	"context"
	"regexp"
	"strings"
	"sync"
	"time"

	"findx/internal/liberror"
	"findx/internal/lockdb"
	"findx/pkg/contentsvc"

	libDocumentRead "github.com/go-shiori/go-readability"
	colly "github.com/gocolly/colly/v2"
	"golang.org/x/sync/errgroup"
)

var (
	newlineRegex = regexp.MustCompile(`[\n\r]{3,}`)          // 3 or more newlines
	spaceRegex   = regexp.MustCompile(`[ \t]{2,}`)           // 2+ spaces/tabs
	lineTrim     = regexp.MustCompile(`(?m)^[ \t]+|[ \t]+$`) // trim leading/trailing whitespace per line
)

type ContentServer struct {
	contentsvc.UnimplementedContentServiceServer
	lockDb lockdb.ILockDb
}

func NewContentServer(lockDb lockdb.ILockDb) *ContentServer {
	return &ContentServer{
		lockDb: lockDb,
	}
}

func (s *ContentServer) ExtractContentFromLinks(ctx context.Context, request *contentsvc.ExtractContentFromLinksRequest) (_ *contentsvc.ExtractContentFromLinksReponse, err error) {
	if len(request.Links) == 0 {
		err = Error(
			liberror.WrapStack(liberror.ErrorDataInvalid, "links is required"))
		return
	}
	var (
		response = &contentsvc.ExtractContentFromLinksReponse{
			Contents: make([]*contentsvc.ExtractedContent, 0),
		}

		mutex = sync.Mutex{}
		eg, _ = errgroup.WithContext(ctx)
	)
	for _, link := range request.Links {
		// Maximum allow 50 concurrent goroutines at a time
		queueLock, err := s.lockDb.AcquireSlot(ctx, "content:extract:concurrency:lock", 50, 10*time.Second, 1*time.Second)
		if err != nil {
			return nil, err
		}

		eg.Go(func() (err error) {
			defer queueLock.ReleaseSlot(ctx)
			var content *contentsvc.ExtractedContent
			var (
				htmlCollector = colly.NewCollector()
				callbackErr   error
			)
			htmlCollector.OnResponse(func(r *colly.Response) {
				var (
					article libDocumentRead.Article
				)
				article, callbackErr = libDocumentRead.FromReader(bytes.NewBuffer(r.Body), nil)
				if callbackErr != nil {
					return
				}
				cleanedContent := strings.TrimSpace(article.TextContent)
				cleanedContent = lineTrim.ReplaceAllString(cleanedContent, "")
				cleanedContent = newlineRegex.ReplaceAllString(cleanedContent, "\n\n")
				cleanedContent = spaceRegex.ReplaceAllString(cleanedContent, " ")
				content = &contentsvc.ExtractedContent{
					Link:    link,
					Content: cleanedContent,
					Title:   article.Title,
				}
			})
			err = htmlCollector.Visit(link)
			if err != nil {
				return err
			}
			if callbackErr != nil {
				return callbackErr
			}

			mutex.Lock()
			defer mutex.Unlock()
			response.Contents = append(response.Contents, content)
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, Error(liberror.WrapStack(err, "extract content: failed"))
	}
	return response, nil
}
