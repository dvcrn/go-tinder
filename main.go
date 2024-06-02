package tinder

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	pb "github.com/dvcrn/go-tinder/pb"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type FilterOpt func(q *url.Values)
type PageOpt func(q *url.Values)

func WithPageToken(timestamp int) PageOpt {
	// check if timestamp is in milliseconds
	if timestamp < 1000000000000 {
		timestamp = timestamp * 1000
	}
	token := base64.StdEncoding.EncodeToString([]byte(strconv.FormatInt(int64(timestamp), 10)))

	return func(q *url.Values) {
		q.Set("page_token", token)
	}
}

type ClientOpt func(c *Client)

func WithDebug() ClientOpt {
	return func(c *Client) {
		c.debug = true
	}
}

type Client struct {
	token              string
	cookiejar          *cookiejar.Jar
	debug              bool
	logger             *log.Logger
	persistentDeviceID string
}

func NewClient(token string, persistentDeviceID string, opts ...ClientOpt) *Client {
	jar, _ := cookiejar.New(nil)
	log := log.Default()
	log.SetPrefix("[go-tinder] ")

	c := &Client{
		token:              token,
		cookiejar:          jar,
		debug:              false,
		logger:             log,
		persistentDeviceID: persistentDeviceID,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func (c *Client) log(msg string, args map[string]string) {
	if c.debug {
		argStrs := []string{}
		for k, v := range args {
			argStrs = append(argStrs, fmt.Sprintf("%s=%s", k, v))
		}

		c.logger.Print(msg, " ", strings.Join(argStrs, ", "))
	}
}

func (c *Client) createRequest(method string, url string, body []byte) (*http.Request, error) {
	bodyReader := bytes.NewReader(body)
	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept-Language", "en,en-US")
	req.Header.Set("app-version", "1051800")

	req.Header.Set("platform", "web")
	req.Header.Set("support-short-video", "1")
	req.Header.Set("tinder-version", "5.18.0")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.3.1 Safari/605.1.15")
	req.Header.Set("user-session-time-elapsed", "1704")
	req.Header.Set("X-Auth-Token", c.token)
	req.Header.Set("x-supported-image-formats", "jpeg")
	req.Header.Set("persistent-device-id", c.persistentDeviceID)

	// these don't seem needed
	// req.Header.Set("app-session-id", "xxx")
	// req.Header.Set("user-session-id", "xxx")
	// req.Header.Set("app-session-time-elapsed", "32828842")

	return req, nil
}

func (c *Client) doRequestAndUnmarshalProto(req *http.Request, target protoreflect.ProtoMessage) error {
	client := http.DefaultClient
	client.Jar = c.cookiejar

	c.log("sending request", map[string]string{
		"url": req.URL.String(),
	})

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = proto.Unmarshal(body, target)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) doRequestAndUnmarshal(req *http.Request, target any) error {
	client := http.DefaultClient
	client.Jar = c.cookiejar

	c.log("sending request", map[string]string{
		"url": req.URL.String(),
	})

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, target)
	if err != nil {
		return err
	}

	return nil
}

type GetMatchesFilterOpt FilterOpt

// WithhasMessagesFilter adds a filter to GetMatches to only return matches that either have messages or don't have messages
func WithHasMessagesFilter(hasMessages bool) GetMatchesFilterOpt {
	return func(q *url.Values) {
		if hasMessages {
			q.Set("message", "1")
		} else {
			q.Set("message", "0")
		}
	}
}

// GetMatches returns `count` amount of matches
func (c *Client) GetMatches(count int, opts ...GetMatchesFilterOpt) ([]*Match, error) {
	u, err := url.Parse("https://api.gotinder.com/v2/matches")
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Set("locale", "en")
	q.Set("count", strconv.FormatInt(int64(count), 10))
	q.Set("is_tinder_u", "false")

	for _, opt := range opts {
		opt(&q)
	}

	u.RawQuery = q.Encode()

	req, err := c.createRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	var responseObject Response
	if err := c.doRequestAndUnmarshal(req, &responseObject); err != nil {
		return nil, err
	}

	if responseObject.Data.Matches == nil {
		return nil, errors.New("invalid response returned")
	}

	return *responseObject.Data.Matches, nil
}

// GetMessages returns `count` amount of messages for the given `matchID`
// Specify `opts` for pagination by providing a page token
func (c *Client) GetMessages(matchID string, count int, opts ...PageOpt) ([]*Message, error) {
	u, err := url.Parse("https://api.gotinder.com/v2/matches/" + matchID + "/messages")
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Set("locale", "en")
	q.Set("count", strconv.FormatInt(int64(count), 10))
	for _, opt := range opts {
		opt(&q)
	}

	u.RawQuery = q.Encode()

	req, err := c.createRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	var responseObject Response
	if err := c.doRequestAndUnmarshal(req, &responseObject); err != nil {
		return nil, err
	}

	if responseObject.Data.Messages == nil {
		return nil, errors.New("invalid response returned")
	}

	return *responseObject.Data.Messages, nil
}

// GetWsToken returns a new token to be used to establish a websocket connection
func (c *Client) GetWsToken() (string, error) {
	u, err := url.Parse("https://api.gotinder.com/ws/generate?locale=en")
	if err != nil {
		return "", err
	}

	req, err := c.createRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return "", err
	}

	var responseObject struct {
		Token string `json:"token"`
	}

	if err := c.doRequestAndUnmarshal(req, &responseObject); err != nil {
		return "", err
	}

	return responseObject.Token, nil
}

// Nudge represents a tinder 'nudge'
// Instead of sending full content through the socket connection, Tinder has opted
// to instead send 'nudges' that something new has happened. The clients will then
// call GetUpdates() to figure out what's new.
// Nudges are using Protobuf and don't have a schema available in this library yet
type Nudge struct {
	Message []byte
	Time    time.Time `json:"time"`
}

// ConnectWebsocket establishes a new websocket connection with the given token
// When a nudge is received, it will get put into `nudgeChan`
//
// To stop a websocket connection, use `ctx` with cancel()
// The socket may periodically disconnect. To check for such an event, use `websocket.IsUnexpectedCloseError(err)`
//
// Nudges are very simple currently. We need to guess the protobuf schema to figure out
// what the nudges are actually for. This hasn't been done yet.
func (c *Client) ConnectWebsocket(ctx context.Context, token string, nudgeChan chan (*pb.ClientData)) error {
	wsURL := "wss://keepalive.gotinder.com/ws?token=" + token
	c.log("starting websocket connection", map[string]string{
		"url": wsURL,
	})
	headers := http.Header{}

	ctx, cancel := context.WithCancel(ctx)
	conn, _, err := websocket.DefaultDialer.DialContext(ctx, wsURL, headers)
	if err != nil {
		return err
	}
	defer conn.Close()

	pingTicker := time.NewTicker(time.Minute * 2)
	noMessagesTimeout := time.NewTimer(time.Minute * 15)
	defer cancel()
	defer noMessagesTimeout.Stop()

	conn.SetPingHandler(func(appData string) error {
		c.log("websocket received ping", nil)
		conn.SetReadDeadline(time.Now().Add(10 * time.Minute))
		conn.WriteMessage(websocket.PongMessage, nil)
		return nil
	})

	done := make(chan error)
	go func() {
		defer close(done)
		for {
			conn.SetReadDeadline(time.Now().Add(10 * time.Minute))
			c.log("waiting for websocket message", nil)
			_, message, err := conn.ReadMessage()
			if err != nil {
				c.log("websocket err: ", map[string]string{
					"err": err.Error(),
				})

				if errors.Is(err, os.ErrDeadlineExceeded) || os.IsTimeout(err) {
					c.log("websocket received io.Timeout, aborting", map[string]string{
						"err": err.Error(),
					})

					done <- &websocket.CloseError{
						Code: websocket.CloseAbnormalClosure,
						Text: "timeout reached without any new message. socket probably dead.",
					}
					return

				}

				done <- err
				return
			}
			noMessagesTimeout.Reset(15 * time.Minute)

			var clientData pb.ClientData
			if err := proto.Unmarshal(message, &clientData); err != nil {
				c.log("failed to unmarshal nudge", map[string]string{
					"err": err.Error(),
				})
				continue
			}

			c.log("received websocket message", map[string]string{
				"message":     clientData.String(),
				"byte_length": strconv.FormatInt(int64(len(message)), 10),
			})

			nudgeChan <- &clientData
		}

		c.log("event loop exited somehow, not sure what happend here??", nil)
	}()

	// control loop
	for {
		select {
		case <-ctx.Done():
			c.log("context.Done() received", nil)
			return nil

		case <-pingTicker.C:
			c.log("sending ping to websocket", nil)
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return err
			}

		case err := <-done:
			if err != nil {
				c.log("received err from socket loop", map[string]string{
					"err": err.Error(),
				})
			}
			return err

		case <-noMessagesTimeout.C:
			c.log("timeout reached without any new message. socket probably dead.", nil)
			return &websocket.CloseError{
				Code: websocket.CloseAbnormalClosure,
				Text: "timeout reached without any new message. socket probably dead.",
			}
		}
	}

	c.log("finished without errors, returning nil", nil)
	return nil
}

// GetUpdates returns the last updates that happened since `since`
// `since` **has** to be in UTC, so when in doubt, use t.UTC()
// This method is often used after receiving a nudge
func (c *Client) GetUpdates(since time.Time) (*UpdateResponse, error) {
	u, err := url.Parse("https://api.gotinder.com/updates")
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Set("locale", "en")
	u.RawQuery = q.Encode()

	sinceInUTC := since.UTC()

	body := struct {
		Nudge            bool   `json:"nudge"`
		LastActivityDate string `json:"last_activity_date"`
	}{
		Nudge:            true,
		LastActivityDate: sinceInUTC.Format("2006-01-02T15:04:05.000Z"),
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := c.createRequest(http.MethodPost, u.String(), jsonBody)
	if err != nil {
		return nil, err
	}

	var result UpdateResponse
	if err := c.doRequestAndUnmarshal(req, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetOwnUser tries to fetch the users own profile
func (c *Client) GetOwnUser() (*UserProfile, error) {
	u, err := url.Parse("https://api.gotinder.com/v2/profile")
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Set("locale", "en")
	q.Set("include", "user")
	u.RawQuery = q.Encode()

	req, err := c.createRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	var result Response
	if err := c.doRequestAndUnmarshal(req, &result); err != nil {
		return nil, err
	}

	if result.Data.User == nil {
		return nil, errors.New("invalid response returned")
	}

	return result.Data.User, nil
}

// GetUser fetches the user profile for the given `userID`
// userID is the id of the **user**, not of a match. To get to the user ID from a match,
// check match.Person.ID
func (c *Client) GetUser(userID string) (*UserProfileResponse, error) {
	u, err := url.Parse("https://api.gotinder.com/user/" + userID + "?locale=en")
	if err != nil {
		return nil, err
	}

	req, err := c.createRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	var result UserProfileResponse
	if err := c.doRequestAndUnmarshal(req, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// SendMessage sends a message from `userID` to `otherID` within match (`matchID`)
func (c *Client) SendMessage(userID, otherID, matchID, sessionID, message string) (*Message, error) {
	u, err := url.Parse("https://api.gotinder.com/user/matches/" + matchID + "?locale=en")
	if err != nil {
		return nil, err
	}

	body := struct {
		UserID    string `json:"userId"`
		OtherID   string `json:"otherId"`
		MatchID   string `json:"matchId"`
		SessionID string `json:"sessionId"`
		Message   string `json:"message"`
	}{
		UserID:    userID,
		OtherID:   otherID,
		MatchID:   matchID,
		SessionID: sessionID,
		Message:   message,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := c.createRequest(http.MethodPost, u.String(), jsonBody)
	if err != nil {
		return nil, err
	}

	var result Message
	if err := c.doRequestAndUnmarshal(req, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Client) RefreshToken(refreshToken string) (*pb.LoginResult, error) {
	u, err := url.Parse("https://api.gotinder.com/v3/auth/login?locale=en")
	if err != nil {
		return nil, err
	}

	factor := &pb.AuthGatewayRequest_RefreshAuth{
		RefreshAuth: &pb.RefreshAuth{
			RefreshToken: refreshToken,
		},
	}

	gatewayReq := &pb.AuthGatewayRequest{
		Factor: factor,
	}

	byteData, err := proto.Marshal(gatewayReq)
	if err != nil {
		log.Fatalf("Failed to marshal request: %v", err)
	}

	req, err := c.createRequest(http.MethodPost, u.String(), byteData)
	if err != nil {
		return nil, err
	}

	// NOTE TO SELF: maybe persistent-device-id is generated on first login
	// req.Header.Set("X-Auth-Token", "d789b8a3-c53d-4aa0-92a8-6f953c409bb1")
	// req.Header.Set("persistent-device-id", randomDeviceID)
	// req.Header.Set("persistent-device-id", "e02a2b63-9e91-496f-a9e6-d634b721fe07")
	req.Header.Set("is-created-as-guest", "false")
	req.Header.Set("Content-Type", "application/x-google-protobuf")

	var result pb.AuthGatewayResponse
	if err := c.doRequestAndUnmarshalProto(req, &result); err != nil {
		return nil, err
	}

	loginResult := result.GetLoginResult()
	if loginResult != nil {
		c.token = loginResult.AuthToken
	}

	return loginResult, nil
}
