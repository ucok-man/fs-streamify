package main

import (
	"errors"
	"fmt"
	"net/http"
	"slices"

	"github.com/go-chi/chi/v5"
	"github.com/jinzhu/copier"
	response "github.com/ucok-man/streamify-api/cmd/responses"
	"github.com/ucok-man/streamify-api/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (app *application) recommended(w http.ResponseWriter, r *http.Request) {
	currentUser := app.contextGetUser(r)
	users, err := app.models.User.Recommended(currentUser)
	if err != nil {
		app.errInternalServer(w, r, err)
	}

	var payload response.UsersResponse
	copier.Copy(&payload, &users)

	err = app.writeJSON(w, http.StatusOK, envelope{"users": payload}, nil)
	if err != nil {
		app.errInternalServer(w, r, err)
	}
}

func (app *application) myfriend(w http.ResponseWriter, r *http.Request) {
	currentUser := app.contextGetUser(r)
	users, err := app.models.User.MyFriends(currentUser)
	if err != nil {
		app.errInternalServer(w, r, err)
	}

	var payload response.UsersResponse
	copier.Copy(&payload, &users)

	err = app.writeJSON(w, http.StatusOK, envelope{"users": payload}, nil)
	if err != nil {
		app.errInternalServer(w, r, err)
	}
}

func (app *application) requestFriend(w http.ResponseWriter, r *http.Request) {
	idparam := chi.URLParam(r, "recipientId")
	recipientId, err := bson.ObjectIDFromHex(idparam)
	if err != nil {
		app.errBadRequest(w, r, fmt.Errorf("invalid recipient id value"))
		return
	}

	currentUser := app.contextGetUser(r)

	recipient, err := app.models.User.GetById(recipientId)
	if err != nil {
		app.errBadRequest(w, r, fmt.Errorf("invalid recipient id value"))
		return
	}

	alreadyFriend := slices.Contains(currentUser.FriendIDs, recipient.ID)
	if alreadyFriend {
		app.errBadRequest(w, r, fmt.Errorf(`already friend with user %v`, idparam))
		return
	}

	exist, err := app.models.FriendRequest.CheckExisting(currentUser.ID, recipient.ID)
	if err != nil {
		app.errInternalServer(w, r, err)
		return
	}
	if exist {
		app.errBadRequest(w, r, fmt.Errorf(`friend request already exist between you and this user`))
		return
	}

	friendRequest := &models.FriendRequest{
		SenderID:    currentUser.ID,
		RecipientID: recipient.ID,
	}

	friendRequest, err = app.models.FriendRequest.Create(friendRequest)
	if err != nil {
		app.errInternalServer(w, r, err)
		return
	}

	var payload response.FriendRequestResponse
	copier.Copy(&payload, &friendRequest)

	err = app.writeJSON(w, http.StatusCreated, envelope{"friend_request": payload}, nil)
	if err != nil {
		app.errInternalServer(w, r, err)
	}
}

func (app *application) acceptFriend(w http.ResponseWriter, r *http.Request) {
	idparam := chi.URLParam(r, "friendRequestId")
	friendRequestId, err := bson.ObjectIDFromHex(idparam)
	if err != nil {
		app.errBadRequest(w, r, fmt.Errorf("invalid friend request id value"))
		return
	}

	currentUser := app.contextGetUser(r)

	friendRequest, err := app.models.FriendRequest.GetById(friendRequestId)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrRecordNotFound):
			app.errNotFound(w, r)
		default:
			app.errInternalServer(w, r, err)
		}
		return
	}

	if friendRequest.RecipientID.Hex() != currentUser.ID.Hex() {
		app.errNotPermitted(w, r)
		return
	}

	friendRequest.Status = models.FriendRequestStatusAccepted

	friendRequest, err = app.models.FriendRequest.Update(friendRequest)
	if err != nil {
		app.errInternalServer(w, r, err)
		return
	}

	if err := app.models.User.AddFriends(friendRequest.RecipientID, friendRequest.SenderID); err != nil {
		app.errInternalServer(w, r, err)
		return
	}

	if err := app.models.User.AddFriends(friendRequest.SenderID, friendRequest.RecipientID); err != nil {
		app.errInternalServer(w, r, err)
		return
	}

	var payload response.FriendRequestResponse
	copier.Copy(&payload, &friendRequest)

	err = app.writeJSON(w, http.StatusOK, envelope{"friend_request": payload}, nil)
	if err != nil {
		app.errInternalServer(w, r, err)
	}
}

func (app *application) getAllFromFriendRequest(w http.ResponseWriter, r *http.Request) {
	// currentUser := app.contextGetUser(r)

	// app.models.FriendRequest.GetAllFromFriendRequest()
}

func (app *application) getAllSendFriendRequest(w http.ResponseWriter, r *http.Request) {
}
