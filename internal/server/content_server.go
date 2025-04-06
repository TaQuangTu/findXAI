package server

import (
	"context"
	"fmt"
	"io"
	"regexp"
	"strings"
	"time"

	"findx/internal/libcmd"
	"findx/internal/liberror"
	"findx/internal/lockdb"
	"findx/pkg/contentsvc"

	libDocumentRead "github.com/go-shiori/go-readability"
	"github.com/google/uuid"
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

	curlCmd libcmd.ICustomCmd
}

func NewContentServer(lockDb lockdb.ILockDb) *ContentServer {
	return &ContentServer{
		lockDb:  lockDb,
		curlCmd: libcmd.NewCustomCmd("curl"),
	}
}

func (s *ContentServer) ExtractContentFromLinks(ctx context.Context, request *contentsvc.ExtractContentFromLinksRequest) (_ *contentsvc.ExtractContentFromLinksReponse, err error) {
	requestid := uuid.NewString()
	if len(request.Links) == 0 {
		err = Error(
			liberror.WrapStack(liberror.ErrorDataInvalid, "links is required"))
		return
	}
	queueLock, err := s.lockDb.AcquireSlot(ctx, "content:extract:concurrency:lock", 50, 10*time.Second, 1*time.Second)
	if err != nil {
		return nil, err
	}
	defer queueLock.ReleaseSlot(ctx)
	var (
		response = &contentsvc.ExtractContentFromLinksReponse{
			Contents: make([]*contentsvc.ExtractedContent, 0),
		}

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

			mutexLock, err := s.lockDb.LockSimple(ctx, "content:extract:request:%s", requestid)
			if err != nil {
				return err
			}
			defer mutexLock.Unlock(ctx)
			response.Contents = append(response.Contents, content)

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, Error(liberror.WrapStack(err, "extract content: failed"))
	}
	return response, nil
}
