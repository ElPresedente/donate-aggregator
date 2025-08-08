package services

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"go-back/database"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/browser"
)

//go:embed twitch_client_credentials.txt
var secrets string

var (
	ClientID      string
	clientSecret  string
	Token         string = ""
	BroadcasterId int    = 0
)

const (
	authUrl     = "https://id.twitch.tv/oauth2/authorize"
	tokenUrl    = "https://id.twitch.tv/oauth2/token"
	scopes      = "channel:bot channel:read:redemptions user:bot user:read:chat"
	redirectUrl = "http://localhost:3000"
)

type EmoteData struct {
	Id     string
	Text   string
	Format map[string]bool //присутствуют ли стандартные static/animated
	Scale  map[string]bool //присутствуют ли стандартные 1.0 2.0 3.0
	Theme  map[string]bool //присутствуют ли стандартные light/dark
}

type BadgeData struct {
	id          string
	image_url   map[int]string //1/2/3
	title       string
	description string
}

type BadgeSetData struct {
	id     string
	badges map[string]BadgeData // id to badge
}

var cachedEmotes map[string]EmoteData //id to emote

var cachedBadges map[string]BadgeSetData //id to set

func TwitchHasAuth() (bool, error) {
	return database.CredentialsDB.CheckENVExists("twitchRefreshToken")
}

func TwitchNewToken() (string, error) {
	setupVars()
	if res, err := database.CredentialsDB.CheckENVExists("twitchRefreshToken"); err == nil && !res {
		err = twitchLogin()
		if err != nil {
			return Token, nil
		}
		return "", err
	}
	refreshToken, err := database.CredentialsDB.GetENVValue("twitchRefreshToken")
	if err != nil {
		log.Println(err)
	}

	data := url.Values{}
	data.Set("client_id", ClientID)
	data.Set("client_secret", clientSecret)
	data.Set("refresh_token", refreshToken)
	data.Set("grant_type", "refresh_token")

	resp, err := http.Post(tokenUrl, "application/x-www-form-urlencoded", bytes.NewBufferString(data.Encode()))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("twitch API error: %s", string(body))
	}

	var tokenResp TokenResponse
	err = json.Unmarshal(body, &tokenResp)
	if err != nil {
		return "", err
	}

	Token = tokenResp.AccessToken

	if len(cachedEmotes) == 0 {
		requestEmotes()
	}
	if len(cachedBadges) == 0 {
		requestBadges()
	}

	return tokenResp.AccessToken, nil
}

func TwitchSubscribeChatMessages(sessionId string) error {
	return twitchSubscribe(sessionId, "channel.chat.message", map[string]string{
		"broadcaster_user_id": strconv.Itoa(BroadcasterId),
		"user_id":             strconv.Itoa(BroadcasterId),
	})
}

func TwitchSubscribeRewardRedemption(sessionId string) error {
	return twitchSubscribe(sessionId, "channel.channel_points_custom_reward_redemption.add", map[string]string{
		"broadcaster_user_id": strconv.Itoa(BroadcasterId),
	})
}

//===============implementation===============

type eventSubRequest struct {
	Type      string            `json:"type"`
	Version   string            `json:"version"`
	Condition map[string]string `json:"condition"`
	Transport map[string]string `json:"transport"`
}

func twitchSubscribe(sessionId, eventType string, condition map[string]string) error {
	if Token == "" {
		return fmt.Errorf("токен доступа отсутствует")
	}

	reqBody := eventSubRequest{
		Type:      eventType,
		Version:   "1",
		Condition: condition,
		Transport: map[string]string{
			"method":     "websocket",
			"session_id": sessionId,
		},
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	const url = "https://api.twitch.tv/helix/eventsub/subscriptions"

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", Token))
	req.Header.Set("Client-Id", ClientID)
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("twitch subscribe error: %s", string(respBody))
	}

	return nil
}

func setupVars() {
	lines := strings.Split(strings.TrimSpace(secrets), "\n")

	// Проверяем, что в файле ровно 2 строки
	if len(lines) != 2 {
		panic("Ошибка: файл twitch_client_credentials.txt должен содержать ровно 2 строки")
	}
	ClientID = lines[0]
	clientSecret = lines[1]

	if ClientID[len(ClientID)-1] == '\r' {
		ClientID = ClientID[:len(ClientID)-1]
	}
}

func requestBroadcasterId() error {
	const (
		userUrl = "https://api.twitch.tv/helix/users?login=moonseere"
	)
	if Token == "" {
		return fmt.Errorf("токен доступа отсутствует")
	}

	req, err := http.NewRequest("GET", userUrl, nil)

	if err != nil {
		return err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", Token))
	req.Header.Set("Client-Id", ClientID)

	client := http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	type User struct {
		ID              string `json:"id"`
		Login           string `json:"login"`
		DisplayName     string `json:"display_name"`
		Type            string `json:"type"`
		BroadcasterType string `json:"broadcaster_type"`
		Description     string `json:"description"`
		ProfileImageURL string `json:"profile_image_url"`
		OfflineImageURL string `json:"offline_image_url"`
		ViewCount       int    `json:"view_count"`
		CreatedAt       string `json:"created_at"`
	}

	type Response struct {
		Data []User `json:"data"`
	}

	var result Response
	body, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(body, &result)
	if err != nil {
		return err
	}

	BroadcasterId, err = strconv.Atoi(result.Data[0].ID)
	return err
}

func newEmote() EmoteData {
	var data EmoteData
	data.Format = make(map[string]bool)
	data.Format["static"] = false
	data.Format["animated"] = false

	data.Scale = make(map[string]bool)
	data.Scale["1.0"] = false
	data.Scale["2.0"] = false
	data.Scale["3.0"] = false

	data.Theme = make(map[string]bool)
	data.Theme["light"] = false
	data.Theme["dark"] = false

	return data
}

func newBadgeSet() BadgeSetData {
	var data BadgeSetData
	data.badges = make(map[string]BadgeData)

	return data
}

func newBadge() BadgeData {
	var data BadgeData
	data.image_url = make(map[int]string)

	return data
}

func requestEmotes() error {
	const (
		globalUrl  = "https://api.twitch.tv/helix/chat/emotes/global"
		channelUrl = "https://api.twitch.tv/helix/chat/emotes?broadcaster_id=%d"
	)
	if Token == "" {
		return fmt.Errorf("токен доступа отсутствует")
	}

	if BroadcasterId == 0 {
		err := requestBroadcasterId()
		if err != nil {
			return err
		}
	}

	type Images struct {
		URL1x string `json:"url_1x"`
		URL2x string `json:"url_2x"`
		URL4x string `json:"url_4x"`
	}

	type Emote struct {
		ID        string   `json:"id"`
		Name      string   `json:"name"`
		Images    Images   `json:"images"`
		Format    []string `json:"format"`
		Scale     []string `json:"scale"`
		ThemeMode []string `json:"theme_mode"`
	}

	type Response struct {
		Data []Emote `json:"data"`
	}

	req, err := http.NewRequest("GET", fmt.Sprintf(channelUrl, BroadcasterId), nil)

	if err != nil {
		return err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", Token))
	req.Header.Set("Client-Id", ClientID)

	client := http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	var result Response

	body, _ := io.ReadAll(resp.Body)
	log.Println(string(body))
	err = json.Unmarshal(body, &result)
	if err != nil {
		return err
	}

	cachedEmotes = make(map[string]EmoteData)

	for _, data := range result.Data {
		newEmote := newEmote()
		newEmote.Id = data.ID
		newEmote.Text = data.Name

		for _, format := range data.Format {
			newEmote.Format[format] = true
		}
		for _, scale := range data.Scale {
			newEmote.Scale[scale] = true
		}
		for _, theme := range data.ThemeMode {
			newEmote.Theme[theme] = true
		}
		cachedEmotes[data.ID] = newEmote
	}

	req, err = http.NewRequest("GET", fmt.Sprintf(channelUrl, BroadcasterId), nil)

	if err != nil {
		return err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", Token))
	req.Header.Set("Client-Id", ClientID)

	resp, err = client.Do(req)
	if err != nil {
		return err
	}

	var result2 Response

	body, _ = io.ReadAll(resp.Body)
	err = json.Unmarshal(body, &result2)
	if err != nil {
		return err
	}

	for _, data := range result2.Data {
		newEmote := newEmote()
		newEmote.Id = data.ID
		newEmote.Text = data.Name

		for _, format := range data.Format {
			newEmote.Format[format] = true
		}
		for _, scale := range data.Scale {
			newEmote.Scale[scale] = true
		}
		for _, theme := range data.ThemeMode {
			newEmote.Theme[theme] = true
		}
		cachedEmotes[data.ID] = newEmote
	}
	return nil
}

func requestBadges() error {
	const (
		globalUrl  = "https://api.twitch.tv/helix/chat/badges/global"
		channelUrl = "https://api.twitch.tv/helix/chat/badges?broadcaster_id=%d"
	)
	if Token == "" {
		return fmt.Errorf("токен доступа отсутствует")
	}

	if BroadcasterId == 0 {
		err := requestBroadcasterId()
		if err != nil {
			return err
		}
	}

	type VersionJson struct {
		ID          string  `json:"id"`
		ImageURL1x  string  `json:"image_url_1x"`
		ImageURL2x  string  `json:"image_url_2x"`
		ImageURL4x  string  `json:"image_url_4x"`
		Title       string  `json:"title"`
		Description string  `json:"description"`
		ClickAction string  `json:"click_action"`
		ClickURL    *string `json:"click_url"` // Use pointer for nullable field
	}

	type BadgeSetJson struct {
		SetID    string        `json:"set_id"`
		Versions []VersionJson `json:"versions"`
	}

	type BadgeDataJson struct {
		Data []BadgeSetJson `json:"data"`
	}

	req, err := http.NewRequest("GET", globalUrl, nil)

	if err != nil {
		return err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", Token))
	req.Header.Set("Client-Id", ClientID)

	client := http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	var result BadgeDataJson

	body, _ := io.ReadAll(resp.Body)
	log.Println(string(body))
	err = json.Unmarshal(body, &result)
	if err != nil {
		return err
	}

	cachedBadges = make(map[string]BadgeSetData)

	for _, data := range result.Data {
		newBadgeSet := newBadgeSet()
		newBadgeSet.id = data.SetID

		for _, badge := range data.Versions {
			newBadge := newBadge()
			newBadge.id = badge.ID
			newBadge.description = badge.Description
			newBadge.title = badge.Title
			newBadge.image_url[1] = badge.ImageURL1x
			newBadge.image_url[2] = badge.ImageURL2x
			newBadge.image_url[3] = badge.ImageURL4x
		}
		cachedBadges[data.SetID] = newBadgeSet
	}

	req, err = http.NewRequest("GET", fmt.Sprintf(channelUrl, BroadcasterId), nil)

	if err != nil {
		return err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", Token))
	req.Header.Set("Client-Id", ClientID)

	resp, err = client.Do(req)
	if err != nil {
		return err
	}

	var result2 BadgeDataJson

	body, _ = io.ReadAll(resp.Body)
	err = json.Unmarshal(body, &result2)
	if err != nil {
		return err
	}

	for _, data := range result2.Data {
		newBadgeSet := newBadgeSet()
		newBadgeSet.id = data.SetID

		for _, badge := range data.Versions {
			newBadge := newBadge()
			newBadge.id = badge.ID
			newBadge.description = badge.Description
		}
		cachedBadges[data.SetID] = newBadgeSet
	}
	return nil
}

// Структура ответа Twitch
type TokenResponse struct {
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
	TokenType    string   `json:"token_type"`
	Scope        []string `json:"scope"`
}

func exchangeCodeForToken(code string) (*TokenResponse, error) {
	data := url.Values{}
	data.Set("client_id", ClientID)
	data.Set("client_secret", clientSecret)
	data.Set("code", code)
	data.Set("grant_type", "authorization_code")
	data.Set("redirect_uri", redirectUrl)

	resp, err := http.Post(tokenUrl, "application/x-www-form-urlencoded", bytes.NewBufferString(data.Encode()))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("twitch API error: %s", string(body))
	}

	var tokenResp TokenResponse
	err = json.Unmarshal(body, &tokenResp)
	if err != nil {
		return nil, err
	}

	return &tokenResp, nil
}

func twitchLogin() error {
	codeChan := make(chan string)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code != "" {
			fmt.Fprintf(w, "Code received. You can close this window.")
			codeChan <- code
			return
		}
		http.Error(w, "No code found", http.StatusBadRequest)
	})

	httpServer := http.Server{
		Addr:    ":3000",
		Handler: mux,
	}
	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Ошибка запуска сервера: %v", err)
		}
	}()
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	defer httpServer.Shutdown(ctx)
	scopesParam := strings.ReplaceAll(scopes, " ", "+")
	authURL := fmt.Sprintf("%s?client_id=%s&redirect_uri=%s&response_type=code&scope=%s",
		authUrl, ClientID, redirectUrl, scopesParam)

	err := browser.OpenURL(authURL)
	if err != nil {
		return fmt.Errorf("ошибка в отрытии ссылки на авторизацию")
	}
	var code string
	select {
	case code = <-codeChan:
	case <-time.After(90 * time.Second):
		return fmt.Errorf("timeout")
	}

	token, err := exchangeCodeForToken(code)
	if err != nil {
		return fmt.Errorf("error exchanging code: %s", err)
	}

	log.Println("User access token:", token.AccessToken)

	database.CredentialsDB.UpdateENVValue("twitchRefreshToken", token.RefreshToken)
	Token = token.AccessToken
	return nil
}

func getEmoteUrl(id, format, theme, scale string) string {
	const emotesTemplateUrl string = "https://static-cdn.jtvnw.net/emoticons/v2/%s/%s/%s/%s"
	return fmt.Sprintf(emotesTemplateUrl, id, format, theme, scale)
}
