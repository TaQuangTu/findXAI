package server

import (
	"bytes"
	"context"
	"html"
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
	specialCodes = regexp.MustCompile(`&(?:#\d+|#x[0-9A-Fa-f]+|[A-Za-z][A-Za-z0-9]+);`)
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

func cleanArticleContent(content string) string {
	// Split into paragraphs for better processing
	paragraphs := strings.Split(content, "\n")
	var cleanParagraphs []string

	for _, para := range paragraphs {
		para = strings.TrimSpace(para)
		para = lineTrim.ReplaceAllString(para, "")
		para = newlineRegex.ReplaceAllString(para, "\n\n")
		para = spaceRegex.ReplaceAllString(para, " ")
		para = specialCodes.ReplaceAllString(para, "")
		if para == "" || len(para) < 50 {
			continue
		}

		// use regexp.Compile instead of MatchString
		matcher, _ := regexp.Compile(`^(\s*-\s*)?Ảnh:`)
		if matched := matcher.MatchString(para); matched {
			continue
		}

		// Skip short paragraphs that are likely captions or metadata
		if len(para) < 40 && (strings.Contains(para, ":") || strings.HasPrefix(para, "-")) {
			continue
		}

		src_matcher, _ := regexp.Compile(`^(Source|Nguồn):`)
		// Skip paragraphs that look like source citations
		if matched := src_matcher.MatchString(para); matched {
			continue
		}

		cleanParagraphs = append(cleanParagraphs, para)
	}

	return strings.Join(cleanParagraphs, "\n")
}

func (s *ContentServer) ExtractContentFromLinks(ctx context.Context, request *contentsvc.ExtractContentFromLinksRequest) (_ *contentsvc.ExtractContentFromLinksReponse, err error) {
	if len(request.Links) == 0 {
		err = Error(
			liberror.WrapStack(liberror.ErrorDataInvalid, "links is required"))
		return
	}
	var (
		response = &contentsvc.ExtractContentFromLinksReponse{
			Contents: make([]*contentsvc.ExtractedContent, len(request.Links)),
		}

		mutex = sync.Mutex{}
		eg, _ = errgroup.WithContext(ctx)
	)
	for index, link := range request.Links {
		// Maximum allow 50 concurrent goroutines at a time
		queueLock, err := s.lockDb.AcquireSlot(ctx, "content:extract:concurrency:lock", 50, 10*time.Second, 1*time.Second)
		if err != nil {
			return nil, err
		}
		index := index
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
				article.TextContent = cleanArticleContent(article.TextContent)
				article.Title = html.UnescapeString(article.Title)
				content = &contentsvc.ExtractedContent{
					Link:    link,
					Content: article.TextContent,
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
			response.Contents[index] = content
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, Error(liberror.WrapStack(err, "extract content: failed"))
	}
	return response, nil
}
