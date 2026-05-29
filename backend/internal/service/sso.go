package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
	"gorm.io/gorm"

	"github.com/cashier-config/server/internal/config"
	"github.com/cashier-config/server/internal/model"
)

type SSOService struct {
	db       *gorm.DB
	cfg      *config.SSOConfig
	provider *oidc.Provider
	verifier *oidc.IDTokenVerifier
	oauthCfg *oauth2.Config
}

func NewSSOService(db *gorm.DB, cfg *config.SSOConfig) (*SSOService, error) {
	provider, err := oidc.NewProvider(context.Background(), cfg.ProviderURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create OIDC provider: %w", err)
	}

	oauthCfg := &oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		RedirectURL:  cfg.RedirectURL,
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}

	verifier := provider.Verifier(&oidc.Config{ClientID: cfg.ClientID})

	return &SSOService{
		db:       db,
		cfg:      cfg,
		provider: provider,
		verifier: verifier,
		oauthCfg: oauthCfg,
	}, nil
}

func (s *SSOService) AuthURL() (string, string, error) {
	state := randStr(32)
	nonce := randStr(32)

	s.db.Where("expires_at < ?", time.Now()).Delete(&model.SSOState{})
	s.db.Create(&model.SSOState{
		State:     state,
		Nonce:     nonce,
		ExpiresAt: time.Now().Add(10 * time.Minute),
	})

	return s.oauthCfg.AuthCodeURL(state, oidc.Nonce(nonce)), state, nil
}

func (s *SSOService) Callback(code, state string) (*model.User, error) {
	var stored model.SSOState
	if err := s.db.Where("state = ? AND expires_at > ?", state, time.Now()).First(&stored).Error; err != nil {
		return nil, errors.New("invalid or expired state")
	}
	s.db.Delete(&stored)

	ctx := context.Background()
	token, err := s.oauthCfg.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("token exchange failed: %w", err)
	}

	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		return nil, errors.New("no id_token in response")
	}

	idToken, err := s.verifier.Verify(ctx, rawIDToken)
	if err != nil {
		return nil, fmt.Errorf("id_token verification failed: %w", err)
	}

	var claims struct {
		Sub           string `json:"sub"`
		Name          string `json:"name"`
		Email         string `json:"email"`
		PreferredName string `json:"preferred_username"`
	}
	if err := idToken.Claims(&claims); err != nil {
		return nil, fmt.Errorf("failed to parse claims: %w", err)
	}

	username := claims.PreferredName
	if username == "" {
		username = claims.Sub
	}
	name := claims.Name
	if name == "" {
		name = username
	}

	var user model.User
	err = s.db.Preload("Role").Where("username = ?", username).First(&user).Error
	if err == nil {
		if user.Status != 1 {
			return nil, errors.New("account disabled")
		}
		return &user, nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	user = model.User{
		Username:     username,
		PasswordHash: "sso-" + claims.Sub,
		Name:         name,
		Email:        claims.Email,
		RoleID:       s.cfg.DefaultRoleID,
		AuthSource:   "sso",
		Status:       1,
	}
	if err := s.db.Create(&user).Error; err != nil {
		return nil, err
	}
	s.db.Preload("Role").First(&user, user.ID)
	return &user, nil
}

func (s *SSOService) OIDCConfig() (string, string, error) {
	return s.oauthCfg.Endpoint.AuthURL, s.cfg.ClientID, nil
}

func randStr(length int) string {
	b := make([]byte, length)
	rand.Read(b)
	return hex.EncodeToString(b)
}
