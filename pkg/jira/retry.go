package jira

import (
	"context"
	"io/ioutil"
	"net/http"

	"github.com/hashicorp/go-retryablehttp"
)

// retryPolicy implements CheckRetry interface to log more information about request fails
func (j *Jira) retryPolicy(ctx context.Context, resp *http.Response, err error) (bool, error) {
	shouldRetry, err := retryablehttp.DefaultRetryPolicy(ctx, resp, err)
	if shouldRetry {
		j.log.Warnf("HTTP request failed with code %d, retrying ...", resp.StatusCode)
		body, bodyErr := ioutil.ReadAll(resp.Body)
		if bodyErr != nil {
			return true, bodyErr
		}
		j.log.Debugf("HTTP request response body: %s", body)
	}

	return shouldRetry, err
}
