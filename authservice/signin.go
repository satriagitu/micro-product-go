// Package Authentication of Product API
//
// # Documentation for Authentication of Product API
//
// Schemes : http
// BasePath : /auth
// Version : 1.0.0
//
// Consumes:
//   - application/json
//
// Produces:
//   - application/json
//
// swagger:meta
package authservice

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/satriagitu/micro-product-go/authservice/data"
	"github.com/satriagitu/micro-product-go/authservice/jwt"
	"go.uber.org/zap"
)

var (
	signinRequests = promauto.NewCounter(prometheus.CounterOpts{
		Name: "signin_total",
		Help: "Total number of signup requests",
	})
	signinSuccess = promauto.NewCounter(prometheus.CounterOpts{
		Name: "signin_success",
		Help: "Successful signup requests",
	})
	signinFail = promauto.NewCounter(prometheus.CounterOpts{
		Name: "signin_fail",
		Help: "Failed signup requests",
	})
	signinError = promauto.NewCounter(prometheus.CounterOpts{
		Name: "signin_error",
		Help: "Erroneous signup requests",
	})
)

// SigninController is the Signin route handler
type SigninController struct {
	logger            *zap.Logger
	promSigninTotal   prometheus.Counter
	promSigninSuccess prometheus.Counter
	promSigninFail    prometheus.Counter
	promSigninError   prometheus.Counter
}

// NewSigninController returns a frsh Signin controller
func NewSigninController(logger *zap.Logger) *SigninController {
	return &SigninController{
		logger:            logger,
		promSigninTotal:   signinRequests,
		promSigninSuccess: signinSuccess,
		promSigninFail:    signinFail,
		promSigninError:   signinError,
	}
}

// we need this function to be private
func getSignedToken() (string, error) {
	// we make a JWT Token here with signing method of ES256 and claims.
	// claims are attributes.
	// Aud - audience
	// Iss - issuer
	// Exp - expiration of the Token
	// token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
	// 	"Aud": "frontend.knowsearch.ml",
	// 	"Iss": "knowsearch.ml",
	// 	"Exp": string(time.Now().Add(time.Minute * 1).Unix()),
	// })
	claimsMap := jwt.ClaimsMap{
		Aud: "frontend.knowsearch.ml",
		Iss: "knowsearch.ml",
		Exp: fmt.Sprint(time.Now().Add(time.Minute * 1).Unix()),
	}

	secret := jwt.GetSecret()
	if secret == "" {
		return "", errors.New("empty JWT secret")
	}

	header := "HS256"
	tokenString, err := jwt.GenerateToken(header, claimsMap, secret)
	if err != nil {
		return tokenString, err
	}
	return tokenString, nil
}

// searches the user in the database.
func validateUser(email string, passwordHash string) (bool, error) {
	usr, exists := data.GetUserObject(email)
	if !exists {
		return false, errors.New("user does not exist")
	}
	passwordCheck := usr.ValidatePasswordHash(passwordHash)

	if !passwordCheck {
		return false, nil
	}
	return true, nil
}

// This will be supplied to the MUX router. It will be called when signin request is sent
// if user not found or not validates, returns the Unauthorized error
// if found, returns the JWT back. [How to return this?]
func (ctrl *SigninController) SigninHandler(rw http.ResponseWriter, r *http.Request) {
	// increment total singin requests
	ctrl.promSigninTotal.Inc()

	// validate the request first.
	if _, ok := r.Header["Email"]; !ok {
		ctrl.logger.Warn("Email was not found in the header")
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Email Missing"))
		ctrl.promSigninFail.Inc()
		return
	}
	if _, ok := r.Header["Passwordhash"]; !ok {
		ctrl.logger.Warn("Passwordhash was not found in the header")
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Passwordhash Missing"))
		ctrl.promSigninFail.Inc()
		return
	}
	// lets see if the user exists
	valid, err := validateUser(r.Header["Email"][0], r.Header["Passwordhash"][0])
	if err != nil {
		// this means either the user does not exist
		ctrl.logger.Warn("User does not exist", zap.String("email", r.Header["Email"][0]))
		rw.WriteHeader(http.StatusUnauthorized)
		rw.Write([]byte("User Does not Exist"))
		ctrl.promSigninFail.Inc()
		return
	}

	if !valid {
		// this means the password is wrong
		ctrl.logger.Warn("Password is wrong", zap.String("email", r.Header["Email"][0]))
		rw.WriteHeader(http.StatusUnauthorized)
		rw.Write([]byte("Incorrect Password"))
		ctrl.promSigninFail.Inc()
		return
	}
	tokenString, err := getSignedToken()
	if err != nil {
		ctrl.logger.Error("unable to sign the token", zap.Error(err))
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Internal Server Error"))
		ctrl.promSigninError.Inc()
		return
	}
	ctrl.logger.Info("Token sign", zap.String("token", tokenString), zap.String("email", r.Header["Email"][0]))

	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte(tokenString))
	ctrl.promSigninSuccess.Inc()
}
