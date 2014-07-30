// Package flowdock contains helpful structs and
// methods for dealing with Flowdock's RESTful API's.
// Structs are based on the message types defined here: https://www.flowdock.com/api/message-types
package flowdock

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// flowdockGET is a convenience function for performing
// GET requests against the Flowdock API
func flowdockGET(apiKey, url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(apiKey, "BATMAN")
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func pushMessage(flowAPIKey, message, sender string, threadID int64) error {
	v := url.Values{}
	v.Set("content", message)
	v.Set("external_user_name", sender)
	if threadID != 0 {
		v.Set("message_id", string(threadID))
	}

	pushURL := fmt.Sprintf("https://api.flowdock.com/v1/messages/chat/%s", flowAPIKey)
	resp, err := http.PostForm(pushURL, v)
	defer resp.Body.Close()

	if err != nil {
		return err
	}
	return nil
}

// PushMessageToFlowWithKey can uses the Flowdock "Push" API to start
// a new thread in a flow using any pseudonym the client wishes. Useful
// for e.g implementing bots.
func PushMessageToFlowWithKey(flowAPIKey, message, sender string) error {
	return pushMessage(flowAPIKey, message, sender, 0)
}

// ReplyToThreadInFlowWithKey is similar to PushMessageToFlowWithKey
// except that it is used for replies rather than starting a new thread.
func ReplyToThreadInFlowWithKey(flowAPIKey, message, sender string, threadID int64) error {
	return pushMessage(flowAPIKey, message, sender, threadID)
}
