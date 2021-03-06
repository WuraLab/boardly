package controllers_test

import (
	"bytes"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/wuraLab/boardly/src/backend/internal/controllers"
	"github.com/wuraLab/boardly/src/backend/internal/middlewares"
	"github.com/wuraLab/boardly/src/backend/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)
type LoginResult struct {
    Code      int             `json:"code"`
    Expire    string          `json:"expire"`
    Token     string           `json:"token"`
}

const (
	JWT_SECRET = "secret"
)

// SetupRouter setup routing here
func SetupRouter(DB *gorm.DB) *gin.Engine {
	//Start the default gin server
	r := gin.Default()
	gin.SetMode(gin.TestMode)
	authMiddleware := middlewares.JWTMiddleware(DB,JWT_SECRET,false,false)

	api := r.Group("/api/v1")
	{
		userController := controllers.User{
			DB: DB,
		}

		api.POST("/user/register", userController.Register)

		api.POST("/user/login", authMiddleware.LoginHandler)

	}
	auth := r.Group("/api/auth")
	auth.Use(authMiddleware.MiddlewareFunc())
	{
		auth.POST("/refresh_token", authMiddleware.RefreshHandler)
	}

	r.GET("/", controllers.Home)

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  404,
			"message": "Route Not Found",
		})
	})

	return r
}

/**
* TestRegister
*/
func TestRegister(t *testing.T) {

	testRouter := SetupRouter(DB)
		//create user first
	idealCase := models.User{
		FirstName: "registerfirstname",
		LastName:  "registerlastname",
		Email:     "register@test.com",
		Password:  "registerPassword123$",		
	}

	testCases := []struct{
		input          models.User
		expected int
	  }{
		//missing email
		{
		  input: models.User{
						FirstName: idealCase.FirstName,
						LastName:  idealCase.LastName,
						Password:  idealCase.Password,
		  },
		  expected: http.StatusUnprocessableEntity,
		},
		//missing password
		{
			input: models.User{
						FirstName: idealCase.FirstName,
						LastName:  idealCase.LastName,
						Email:     idealCase.Email,
			},
			expected: http.StatusUnprocessableEntity,
		},
		//missing firstname or lastname
		{
			input: models.User{
						LastName:  idealCase.LastName,
						Email:     idealCase.Email,
						Password:  idealCase.Password,
			},
			expected: http.StatusUnprocessableEntity,
		},
		//compltely filled out
		{
			input: models.User{
						  FirstName: idealCase.FirstName,
						  LastName:  idealCase.LastName,
						  Email:     idealCase.Email,
						  Password:  idealCase.Password,
			},
			expected: http.StatusOK,
		},
	}
   for _, testCase := range testCases {

		data, _ := json.Marshal(testCase.input)

		req, err := http.NewRequest("POST", "/api/v1/user/register", bytes.NewBufferString(string(data)))
		req.Header.Set("Content-Type", "application/json")

		if err != nil {
			log.Fatalln(err)
		}

		resp := httptest.NewRecorder()

		testRouter.ServeHTTP(resp, req)
		assert.Equal(t, testCase.expected, resp.Code)
   }
}


func TestLogin(t *testing.T) {

	testRouter := SetupRouter(DB)

	//create user first
	idealCase := models.User{
		FirstName: "loginfirstname",
		LastName:  "loginlastname",
		Email:     "login@test.com",
		Password:  "loginPassword123$",		
	}
	data, _ := json.Marshal(idealCase)

	req, err := http.NewRequest("POST", "/api/v1/user/register", bytes.NewBufferString(string(data)))
	req.Header.Set("Content-Type", "application/json")

	if err != nil {
		log.Fatalln(err)
	}
	resp := httptest.NewRecorder()
	testRouter.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)

	testCases := []struct{
		input          models.User
		expected int
	  }{
		//missing email
		{
		  input: models.User{
						Password:  idealCase.Password,
		  },
		  expected: http.StatusUnauthorized,
		},
		//missing password
		{
			input: models.User{
						  Email:     idealCase.Email,
			},
			expected: http.StatusUnauthorized,
		},
		//non existing email
		{
			input: models.User{
						  Email:     "wrong@test.com",
						  Password:  idealCase.Password,
			},
			expected: http.StatusUnauthorized,
		},
		//wrong password
		{
			input: models.User{
						  Email:     idealCase.Email,
						  Password:  "wrongpassword",
			},
			expected: http.StatusUnauthorized,
		},
		//completely filled out
		{
			input: models.User{
						  Email:     idealCase.Email,
						  Password:  idealCase.Password,
			},
			expected: http.StatusOK,
		},
	}
   for _, testCase := range testCases {

		data, _ := json.Marshal(testCase.input)

		req, err := http.NewRequest("POST", "/api/v1/user/login", bytes.NewBufferString(string(data)))
		req.Header.Set("Content-Type", "application/json")

		if err != nil {
			log.Fatalln(err)
		}

		resp := httptest.NewRecorder()

		testRouter.ServeHTTP(resp, req)
		assert.Equal(t, testCase.expected, resp.Code)
   }
}

func TestRefresh(t *testing.T) {
    var loginResult LoginResult
	testRouter := SetupRouter(DB)

	//create user first
	idealCase := models.User{
		FirstName: "refreshfirstname",
		LastName:  "refreshlastname",
		Email:     "refresh@test.com",
		Password:  "refreshPassword123$",		
	}

	//Register
	data, _ := json.Marshal(idealCase)
	req, err := http.NewRequest("POST", "/api/v1/user/register", bytes.NewBufferString(string(data)))
	req.Header.Set("Content-Type", "application/json")

	if err != nil {
		log.Fatalln(err)
	}
	resp := httptest.NewRecorder()
	testRouter.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)

	//login
	req, err = http.NewRequest("POST", "/api/v1/user/login", bytes.NewBufferString(string(data)))
	req.Header.Set("Content-Type", "application/json")

	if err != nil {
		log.Fatalln(err)
	}
	resp = httptest.NewRecorder()
	testRouter.ServeHTTP(resp, req)
	b, _ := ioutil.ReadAll(resp.Body)
	if err = json.Unmarshal(b,&loginResult); err != nil {
		log.Error(err)
	}
	log.Println(loginResult)
	assert.Equal(t, http.StatusOK, resp.Code)

	//refresh
	req, err = http.NewRequest("POST", "/api/auth/refresh_token", bytes.NewBufferString(string(data)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer " + loginResult.Token)

	if err != nil {
		log.Fatalln(err)
	}
	resp = httptest.NewRecorder()
	testRouter.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)
}



