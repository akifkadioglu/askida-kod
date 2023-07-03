package auth

import (
	"net/http"
	"time"

	"github.com/akifkadioglu/askida-kod/database"
	"github.com/akifkadioglu/askida-kod/ent"
	"github.com/akifkadioglu/askida-kod/ent/user"
	"github.com/akifkadioglu/askida-kod/models"
	"github.com/akifkadioglu/askida-kod/resources/emails"
	"github.com/akifkadioglu/askida-kod/utils"
	"github.com/akifkadioglu/askida-kod/variables"

	"github.com/dghubble/gologin/v2/google"
	"github.com/go-chi/render"
)

func GoogleLogin(w http.ResponseWriter, r *http.Request) {
	var JWTModel models.JwtModel
	db := database.DBManager()

	userFromGoogle, err := google.UserFromContext(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var userClient *ent.User
	userClient, err = db.User.Query().
		Where(user.Email(userFromGoogle.Email)).
		First(r.Context())

	if ent.IsNotFound(err) {

		userClient, err = db.User.
			Create().
			SetEmail(userFromGoogle.Email).
			SetPicture(userFromGoogle.Picture).
			SetName(userFromGoogle.Name).
			SetIsActive(true).
			Save(r.Context())

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}

	JWTModel.ID = userClient.ID.String()
	JWTModel.Email = userClient.Email
	JWTModel.Name = userClient.Name
	JWTModel.Picture = &userClient.Picture
	JWTModel.Time = time.Now().String()

	tokenAsString, _ := utils.GenerateToken(JWTModel)

	w.Header().Add("Authorization", "Bearer "+tokenAsString)
	http.Redirect(w, r, variables.CLIENT+"/token/"+tokenAsString, http.StatusMovedPermanently)
}

func Register(w http.ResponseWriter, r *http.Request) {
	var input BodyRegister
	db := database.DBManager()
	if err := render.DecodeJSON(r.Body, &input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, utils.JSONMessage("Name, E-mail and Password are required"))
		return
	}

	_, err := db.User.
		Query().
		Where(user.Email(input.Email)).
		First(r.Context())

	if ent.IsNotFound(err) {
		hashedPassword, err := utils.Hash(input.Password)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		newUser, err := db.User.
			Create().
			SetEmail(input.Email).
			SetName(input.Name).
			SetPassword(hashedPassword).
			Save(r.Context())

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, utils.JSONMessage("The user did not create"))
			return
		}

		utils.SendEmail(newUser.Email, "Register", emails.Register(newUser.Name, newUser.ActivationID.String()))

		render.JSON(w, r, utils.JSONMessage("The user created"))

	} else {
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, utils.JSONMessage("This E-mail already in use by someone"))
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	var input BodyLogin
	var JWTModel models.JwtModel
	db := database.DBManager()

	if err := render.DecodeJSON(r.Body, &input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, utils.JSONMessage("E-mail and Password are required"))
		return
	}
	user, err := db.User.
		Query().
		Where(user.Email(input.Email)).
		First(r.Context())

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, utils.JSONMessage("Something went wrong"))
		return
	}

	if utils.CompareHash(user.Password, input.Password) {
		JWTModel.Email = user.Email
		JWTModel.ID = user.ID.String()
		JWTModel.Name = user.Name
		JWTModel.Picture = &user.Picture
		JWTModel.Time = time.Now().String()

		token, err := utils.GenerateToken(JWTModel)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, utils.JSONMessage("Something went wrong"))
			return
		}

		render.JSON(w, r, map[string]string{
			"token": token,
		})
		return

	} else {
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, utils.JSONMessage("E-mail and password dont match"))
		return
	}
}

func SendEmailAgain(w http.ResponseWriter, r *http.Request) {
	var input BodySendEmailAgain
	db := database.DBManager()

	if err := render.DecodeJSON(r.Body, &input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, utils.JSONMessage("E-mail is required"))
		return
	}

	user, _ := db.User.
		Query().
		Where(user.Email(input.Email)).
		First(r.Context())

	err := utils.SendEmail(user.Email, "Register", emails.Register(user.Name, user.ActivationID.String()))

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, utils.JSONMessage("Something went wrong"))
		return

	} else {
		render.JSON(w, r, utils.JSONMessage("sent"))
		return
	}
}
