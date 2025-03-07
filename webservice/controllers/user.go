package controllers

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"

	"shaphil.me/webservice/models"
)

type userController struct {
	userIDPattern *regexp.Regexp
}

func newUserController() *userController {
	return &userController{
		userIDPattern: regexp.MustCompile(`^/users/(\d+)/?`),
	}
}

func (uc userController) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if request.URL.Path == "/users" {
		switch request.Method {
		case http.MethodGet:
			uc.getAll(writer, request)
		case http.MethodPost:
			uc.post(writer, request)
		default:
			writer.WriteHeader(http.StatusNotImplemented)
		}
	} else {
		matches := uc.userIDPattern.FindStringSubmatch(request.URL.Path)
		if len(matches) == 0 {
			writer.WriteHeader(http.StatusNotFound)
		}

		id, err := strconv.Atoi(matches[1])
		if err != nil {
			writer.WriteHeader(http.StatusNotFound)
		}

		switch request.Method {
		case http.MethodGet:
			uc.get(id, writer)
		case http.MethodPut:
			uc.put(id, writer, request)
		case http.MethodDelete:
			uc.delete(id, writer)
		default:
			writer.WriteHeader(http.StatusNotImplemented)
		}
	}
}

func (uc *userController) getAll(w http.ResponseWriter, _ *http.Request) {
	encodeResponseAsJson(models.GetUsers(), w)
}

func (uc *userController) get(id int, w http.ResponseWriter) {
	user, err := models.GetUserById(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	encodeResponseAsJson(user, w)
}

func (uc *userController) post(w http.ResponseWriter, r *http.Request) {
	user, err := uc.parseRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Could not parse User object"))
		return
	}

	user, err = models.AddUser(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	encodeResponseAsJson(user, w)
}

func (uc *userController) put(id int, w http.ResponseWriter, r *http.Request) {
	user, err := uc.parseRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Could not parse User object"))
		return
	}

	if id != user.ID {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("ID of submitted user must match ID in the URL"))
		return
	}

	user, err = models.UpdateUser(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	encodeResponseAsJson(user, w)
}

func (uc *userController) delete(id int, w http.ResponseWriter) {
	err := models.RemoveUserById(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (uc *userController) parseRequest(request *http.Request) (models.User, error) {
	decoder := json.NewDecoder(request.Body)
	var user models.User
	err := decoder.Decode(&user)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}
