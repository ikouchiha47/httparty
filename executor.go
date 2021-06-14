package httparty

import (
	"encoding/json"
	"httparty/engine"
	"httparty/httpclient"
	"io/ioutil"
	"net/http"
)

type Executor interface {
	RunIt(step *engine.Step) (engine.Response, error)
}

type HttpExecutor struct {
	client httpclient.HttpClientI
}

func DefaultHttpExecutor() *HttpExecutor {
	return &HttpExecutor{client: httpclient.NewHttpClient(int64(1000))}
}

func (ht *HttpExecutor) RunIt(step *engine.Step) (engine.Response, error) {
	req, err := ht.prepareRequest(step.Request)
	if err != nil {
		return engine.Response{}, err
	}

	return ht.doReq(req, step.Name)
}

func (ht *HttpExecutor) prepareRequest(request engine.Request) (*http.Request, error) {
	body, err := engine.GetBodyString(request)
	if err != nil {
		return nil, err
	}

	var req *http.Request

	if request.ShouldMakeBody() {
		req, err = request.MakePost(body)
	} else {
		req, err = request.MakeGet(body)
	}

	if err != nil {
		return nil, err
	}

	for key, value := range request.Headers {
		req.Header.Set(key, value)
	}

	req.Header.Set(engine.ContentTypeKey, request.ReqType())
	req.Header.Set(engine.Accept, request.AcceptType())

	return req, nil
}

func (ht *HttpExecutor) doReq(req *http.Request, name string) (engine.Response, error) {
	httpresp := engine.Response{}

	resp, err := ht.client.Do(req)
	if err != nil {
		return httpresp, err
	}

	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return httpresp, err
	}

	if req.Header.Get(engine.Accept) != engine.JSON {
		return engine.Response{Body: engine.Body{name: string(b)}}, nil
	}
	var resbody engine.Body

	if err = json.Unmarshal(b, &resbody); err != nil {
		return httpresp, err
	}

	return engine.Response{Body: resbody}, nil
}
