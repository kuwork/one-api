package mokaai

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/songquanpeng/one-api/relay/adaptor"
	"github.com/songquanpeng/one-api/relay/meta"
	"github.com/songquanpeng/one-api/relay/model"
	"github.com/songquanpeng/one-api/relay/relaymode"
)

type Adaptor struct {
	meta *meta.Meta
}

// ConvertImageRequest implements adaptor.Adaptor.
func (*Adaptor) ConvertImageRequest(request *model.ImageRequest) (any, error) {
	return nil, errors.New("not implemented")
}

// ConvertImageRequest implements adaptor.Adaptor.

func (a *Adaptor) Init(meta *meta.Meta) {
	a.meta = meta
}


func (a *Adaptor) GetRequestURL(meta *meta.Meta) (string, error) {
	
	var urlPrefix = meta.BaseURL
	
	switch meta.Mode {
	case relaymode.ChatCompletions:
		return fmt.Sprintf("%s/chat/completions", urlPrefix), nil
	case relaymode.Embeddings:
		return fmt.Sprintf("%s/embeddings", urlPrefix), nil
	default:
		return fmt.Sprintf("%s/run/%s", urlPrefix, meta.ActualModelName), nil
	}
}

func (a *Adaptor) SetupRequestHeader(c *gin.Context, req *http.Request, meta *meta.Meta) error {
	adaptor.SetupCommonRequestHeader(c, req, meta)
	req.Header.Set("Authorization", "Bearer "+meta.APIKey)
	return nil
}

func (a *Adaptor) ConvertRequest(c *gin.Context, relayMode int, request *model.GeneralOpenAIRequest) (any, error) {
	if request == nil {
		return nil, errors.New("request is nil")
	}
	switch relayMode {
	case relaymode.ChatCompletions, relaymode.Completions:
		return nil, errors.New("not implemented")
	case  relaymode.Embeddings:
		// return ConvertCompletionsRequest(*request), nil
		return ConvertEmbeddingRequest(*request), nil
	default:
		return nil, errors.New("not implemented")
	}
}

func (a *Adaptor) DoRequest(c *gin.Context, meta *meta.Meta, requestBody io.Reader) (*http.Response, error) {
	return adaptor.DoRequestHelper(a, c, meta, requestBody)
}

func (a *Adaptor) DoResponse(c *gin.Context, resp *http.Response, meta *meta.Meta) (usage *model.Usage, err *model.ErrorWithStatusCode) {
	if meta.IsStream {
		err, usage = StreamHandler(c, resp, meta.PromptTokens, meta.ActualModelName)
	} else {
		err, usage = Handler(c, resp, meta.PromptTokens, meta.ActualModelName)
	}
	return
}

func (a *Adaptor) GetModelList() []string {
	return ModelList
}

func (a *Adaptor) GetChannelName() string {
	return "mokaai"
}
