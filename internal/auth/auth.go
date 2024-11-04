package auth

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

const (
	Key    = "ldsfjlasdjflajflj"
	MaxAge = 86400 * 30
	IsProd = true
)

func NewAuth() {
	googleClientId := os.Getenv("GOOGLE_CLIENT_ID")
	gooleClientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")

	store := sessions.NewCookieStore([]byte(Key))
    store.MaxAge(MaxAge)

	store.Options.Path="/"
	store.Options.HttpOnly = true
	store.Options.Domain = "localhost"
	store.Options.SameSite = http.SameSiteLaxMode
	store.Options.Secure = IsProd

	gothic.Store = store
   goth.UseProviders(
	google.New(googleClientId, gooleClientSecret ,os.Getenv("REDIRECT_LOGIN_URL")),
   )

   fmt.Println("The Auth client Created.")
}

func GoogleAuth(c *gin.Context){
	q := c.Request.URL.Query()
	q.Add("provider", "google")
	c.Request.URL.RawQuery = q.Encode()
	gothic.BeginAuthHandler(c.Writer, c.Request)
}

func GoogleAuthCallbackFunction(c *gin.Context){
	q := c.Request.URL.Query()
	q.Add("provider", "google")
	c.Request.URL.RawQuery = q.Encode()
	user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		fmt.Println("Error is here")
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	fmt.Println(user)
	http.Redirect(c.Writer,c.Request,"http://localhost:3000/home",http.StatusTemporaryRedirect)
}

func Logout(c *gin.Context) {
	gothic.Logout(c.Writer, c.Request)
	c.Redirect(http.StatusTemporaryRedirect, "/")
}