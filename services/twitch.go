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
	"strings"
	"time"

	"github.com/pkg/browser"
)

//go:embed twitch_client_credentials.txt
var secrets string

var (
	clientID     string
	clientSecret string
	Token        string = ""
)

const (
	authUrl     = "https://id.twitch.tv/oauth2/authorize"
	tokenUrl    = "https://id.twitch.tv/oauth2/token"
	scopes      = "channel:bot channel:read:redemptions user:bot user:read:chat"
	redirectUrl = "http://localhost:3000"
)

func getEmoteUrl(id, format, theme, scale string) string {
	const emotesTemplateUrl string = "https://static-cdn.jtvnw.net/emoticons/v2/%s/%s/%s/%s"
	return fmt.Sprintf(emotesTemplateUrl, id, format, theme, scale)
}

type EmoteData struct {
	id     string
	text   string
	format map[string]bool //присутствуют ли стандартные static/animated
	scale  map[int]bool    //присутствуют ли стандартные 1.0 2.0 3.0
	theme  map[string]bool //присутствуют ли стандартные light/dark
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

var cachedEmojis map[string]EmoteData //id to emote

var cachedBadges map[string]BadgeSetData //id to set

func setupVars() {
	lines := strings.Split(strings.TrimSpace(secrets), "\n")

	// Проверяем, что в файле ровно 2 строки
	if len(lines) != 2 {
		panic("Ошибка: файл twitch_client_credentials.txt должен содержать ровно 2 строки")
	}
	clientID = lines[0]
	clientSecret = lines[1]

	if clientID[len(clientID)-1] == '\r' {
		clientID = clientID[:len(clientID)-1]
	}
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
	data.Set("client_id", clientID)
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

func TwitchLogin() {
	setupVars()
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
		authUrl, clientID, redirectUrl, scopesParam)

	err := browser.OpenURL(authURL)
	if err != nil {
		log.Fatalln("Ошибка в отрытии ссылки на авторизацию")
		return
	}
	var code string
	select {
	case code = <-codeChan:
	case <-time.After(90 * time.Second):
		log.Println("timeout")
		return
	}

	token, err := exchangeCodeForToken(code)
	if err != nil {
		log.Println("Error exchanging code:", err)
		return
	}

	log.Println("User access token:", token.AccessToken)

	database.CredentialsDB.UpdateENVValue("twitchRefreshToken", token.RefreshToken)
	Token = token.AccessToken
}

func TwitchNewToken() (string, error) {
	refreshToken, err := database.CredentialsDB.GetENVValue("twitchRefreshToken")
	if err != nil {
		log.Println(err)
	}

	data := url.Values{}
	data.Set("client_id", clientID)
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

	return tokenResp.AccessToken, nil
}
