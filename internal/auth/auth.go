package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/user"
	"path/filepath"

	"golang.org/x/oauth2"
)

type Auth struct {
	cfg      *oauth2.Config
	authFile string
}

type Data struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Scopes       []string
	AuthURL      string
	TokenURL     string
}

func NewAuth(data Data) *Auth {
	usr, _ := user.Current()
	authFile := filepath.Join(usr.HomeDir, ".y2spot_token.json")

	return &Auth{
		cfg: &oauth2.Config{
			ClientID:     data.ClientID,
			ClientSecret: data.ClientSecret,
			RedirectURL:  data.RedirectURL,
			Scopes:       data.Scopes,
			Endpoint: oauth2.Endpoint{
				AuthURL:  data.AuthURL,
				TokenURL: data.TokenURL,
			},
		},
		authFile: authFile,
	}
}

func (a *Auth) Config() *oauth2.Config {
	return a.cfg
}

func (a *Auth) IsAuthenticated() bool {
	_, err := os.Stat(a.authFile)
	return err == nil
}

func (a *Auth) Authenticate() error {
	codeCh := make(chan string)

	url := a.cfg.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Println("🌐 Please log in via your browser:")
	fmt.Println(url)

	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		fmt.Fprintln(w, "✅ Authentication successful! You may close this window.")
		codeCh <- code
	})

	server := &http.Server{Addr: ":8080"}
	go func() {
		_ = server.ListenAndServe()
	}()

	code := <-codeCh
	_ = server.Shutdown(context.Background())

	token, err := a.cfg.Exchange(context.Background(), code)
	if err != nil {
		return fmt.Errorf("failed to exchange token: %w", err)
	}

	return a.SaveToken(token)
}

func (a *Auth) LoadToken() (*oauth2.Token, error) {
	f, err := os.Open(a.authFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var token oauth2.Token
	if err := json.NewDecoder(f).Decode(&token); err != nil {
		return nil, err
	}
	return &token, err
}

func (a *Auth) SaveToken(token *oauth2.Token) error {
	f, err := os.Create(a.authFile)
	if err != nil {
		return err
	}
	defer f.Close()

	return json.NewEncoder(f).Encode(token)
}
