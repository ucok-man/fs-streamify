package main

import (
	"errors"
	"net/http"
	"time"

	dto "github.com/ucok-man/streamify-api/cmd/dtos"
	response "github.com/ucok-man/streamify-api/cmd/responses"
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

	// TODO: create user in streamio

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
		Expires:  expiration,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Secure:   app.config.Env == "production",
	})

	err = app.writeJSON(w, http.StatusAccepted, envelope{"user": response.SignupResponse{
		ID:          user.ID,
		FullName:    user.FullName,
		Email:       user.Email,
		Bio:         user.Bio,
		ProfilePic:  user.ProfilePic,
		NativeLng:   user.NativeLng,
		LearningLng: user.LearningLng,
		Location:    user.Location,
		IsOnboarded: user.IsOnboarded,
		FriendIDs:   user.FriendIDs,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}}, nil)
	if err != nil {
		app.errInternalServer(w, r, err)
	}
}

func (app *application) signin(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World"))
}

func (app *application) signout(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World"))
}
