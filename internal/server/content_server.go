package server

import (
	"context"
	"fmt"
	"io"
	"regexp"
	"strings"
	"sync"

	"findx/internal/libcmd"
	"findx/internal/liberror"
	"findx/pkg/contentsvc"

	libDocumentRead "github.com/go-shiori/go-readability"
	"golang.org/x/sync/errgroup"
)

var (
	newlineRegex = regexp.MustCompile(`[\n\r]{3,}`)          // 3 or more newlines
	spaceRegex   = regexp.MustCompile(`[ \t]{2,}`)           // 2+ spaces/tabs
	lineTrim     = regexp.MustCompile(`(?m)^[ \t]+|[ \t]+$`) // trim leading/trailing whitespace per line
)

type ContentServer struct {
	contentsvc.UnimplementedContentServiceServer

	curlCmd libcmd.ICustomCmd
}

func NewContentServer() *ContentServer {
	return &ContentServer{
		curlCmd: libcmd.NewCustomCmd("curl"),
	}
}

// TODO: support stream response later
// TODO: support store document in embedding db and mapping to postgres for future use
// TODO: consider replacing curl by go-colly
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

		mu        sync.Mutex
		eg, egCtx = errgroup.WithContext(ctx)
	)
	for _, link := range request.Links {
		link := link

		eg.Go(func() error {
			var content *contentsvc.ExtractedContent

			_, err := s.curlCmd.WithStreamReader(func(readCloser io.ReadCloser) error {
				defer readCloser.Close()

				article, err := libDocumentRead.FromReader(readCloser, nil)
				if err != nil {
					return err
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
				return nil

				// TODO: support dynamic selected user-agent
			}).Run(egCtx, "-sL", "--compressed", "-A", "Mozilla/5.0", link)
			if err != nil {
				return fmt.Errorf("failed to extract link [%s]: %w", link, err)
			}

			mu.Lock()
			defer mu.Unlock()
			response.Contents = append(response.Contents, content)

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, Error(liberror.WrapStack(err, "extract content: failed"))
	}
	return response, nil
}
