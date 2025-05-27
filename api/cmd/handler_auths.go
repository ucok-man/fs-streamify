package main

import (
	"context"
	"errors"
	"net/http"
	"time"

	stream_chat "github.com/GetStream/stream-chat-go/v5"
	dto "github.com/ucok-man/streamify-api/cmd/dtos"
	"github.com/ucok-man/streamify-api/internal/models"
	"github.com/ucok-man/streamify-api/internal/validator"
)

func (app *application) signup(w http.ResponseWriter, r *http.Request) {
	var dto dto.SignupDTO
	err := app.readJSON(w, r, &dto)
	if err != nil {
		app.errBadRequest(w, r, err)
		return
	}

	errmap := validator.Schema().SignupDTO.Validate(&dto)
	if errmap != nil {
		app.errFailedValidation(w, r, validator.Sanitize(errmap))
		return
	}

	user := &models.User{
		FullName:   dto.Fullname,
		Email:      dto.Email,
		ProfilePic: app.getRandomPicturePlaceholder(),
	}

	if err := user.Password.Set(dto.Password); err != nil {
		app.errInternalServer(w, r, err)
		return
	}

	user, err = app.models.User.Insert(user)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrDuplicateEmail):
			app.errFailedValidation(w, r, map[string][]string{
				"email": {"User with this email already exist"},
			})
		default:
			app.errInternalServer(w, r, err)
		}
		return
	}

	// Create user in getstream.io
	_, err = app.stream.UpsertUser(context.Background(), &stream_chat.User{
		ID:    user.ID.Hex(),
		Name:  user.FullName,
		Image: user.ProfilePic,
	})
	if err != nil {
		app.errInternalServer(w, r, err)
		return
	}

	expiration := time.Now().Add(7 * 24 * time.Hour)
	claim := app.NewJWTClaim(user.ID.Hex(), expiration)
	token, err := app.GenerateJwtToken(claim, app.config.JWT.AuthSecret)
	if err != nil {
		app.errInternalServer(w, r, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "jwt-auth-token.streamify",
		Value:    token,
		Path:     "/",
		Expires:  expiration,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Secure:   app.config.Env == "production",
	})

	err = app.writeJSON(w, http.StatusCreated, envelope{"user": user}, nil)
	if err != nil {
		app.errInternalServer(w, r, err)
	}
}

func (app *application) signin(w http.ResponseWriter, r *http.Request) {
	var dto dto.SigninDTO
	err := app.readJSON(w, r, &dto)
	if err != nil {
		app.errBadRequest(w, r, err)
		return
	}

	errmap := validator.Schema().SigninDTO.Validate(&dto)
	if errmap != nil {
		app.errFailedValidation(w, r, validator.Sanitize(errmap))
		return
	}

	user, err := app.models.User.GetByEmail(dto.Email)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrRecordNotFound):
			app.errInvalidCredentials(w, r)
		default:
			app.errInternalServer(w, r, err)
		}
		return
	}

	match, err := user.Password.Matches(dto.Password)
	if err != nil {
		app.errInternalServer(w, r, err)
		return
	}
	if !match {
		app.errInvalidCredentials(w, r)
		return
	}

	expiration := time.Now().Add(7 * 24 * time.Hour)
	claim := app.NewJWTClaim(user.ID.Hex(), expiration)
	token, err := app.GenerateJwtToken(claim, app.config.JWT.AuthSecret)
	if err != nil {
		app.errInternalServer(w, r, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "jwt-auth-token.streamify",
		Value:    token,
		Path:     "/",
		Expires:  expiration,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Secure:   app.config.Env == "production",
	})

	err = app.writeJSON(w, http.StatusCreated, envelope{"user": user}, nil)
	if err != nil {
		app.errInternalServer(w, r, err)
	}
}

func (app *application) signout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "jwt-auth-token.streamify",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Secure:   app.config.Env == "production",
	})
	err := app.writeJSON(w, http.StatusOK, envelope{"message": "Session signout success!"}, nil)
	if err != nil {
		app.errInternalServer(w, r, err)
	}
}

func (app *application) onboarding(w http.ResponseWriter, r *http.Request) {
	var dto dto.OnboardingDTO
	err := app.readJSON(w, r, &dto)
	if err != nil {
		app.errBadRequest(w, r, err)
		return
	}

	errmap := validator.Schema().OnboardingDTO.Validate(&dto)
	if errmap != nil {
		app.errFailedValidation(w, r, validator.Sanitize(errmap))
		return
	}

	user := app.contextGetUser(r)
	user.Bio = dto.Bio
	user.FullName = dto.Fullname
	user.NativeLng = dto.NativeLng
	user.LearningLng = dto.LearningLng
	user.Location = dto.Location
	user.IsOnboarded = true

	user, err = app.models.User.Update(user)
	if err != nil {
		app.errInternalServer(w, r, err)
		return
	}

	// Update user in getstream.io
	_, err = app.stream.UpsertUser(context.Background(), &stream_chat.User{
		ID:    user.ID.Hex(),
		Name:  user.FullName,
		Image: user.ProfilePic,
	})
	if err != nil {
		app.errInternalServer(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"user": user}, nil)
	if err != nil {
		app.errInternalServer(w, r, err)
	}
}

func (app *application) whoami(w http.ResponseWriter, r *http.Request) {
	user := app.contextGetUser(r)
	err := app.writeJSON(w, http.StatusCreated, envelope{"user": user}, nil)
	if err != nil {
		app.errInternalServer(w, r, err)
	}
}
