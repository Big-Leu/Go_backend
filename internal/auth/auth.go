package auth

import (
	"fmt"
	"kubequntumblock/internal/util"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/google"
)
var djlsf *goth.User
const (
	Key    = "ldsfjlasdjflajflj"
	MaxAge = 86400 
	IsProd = true
)

func NewAuth() {
	googleClientId := os.Getenv("GOOGLE_CLIENT_ID")
	gooleClientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	githubClientId := os.Getenv("GITHUB_CLIENT_ID")
	githubClientSecret := os.Getenv("GITHUB_CLIENT_SECRET")

	store := sessions.NewCookieStore([]byte(Key))
    store.MaxAge(MaxAge)

	store.Options.Path="/"
	store.Options.HttpOnly = true
	store.Options.Domain = "localhost"
	store.Options.SameSite = http.SameSiteDefaultMode
	store.Options.Secure = IsProd

	gothic.Store = store
    goth.UseProviders(
	google.New(googleClientId, gooleClientSecret ,os.Getenv("REDIRECT_LOGIN_URL")),
	github.New(githubClientId, githubClientSecret ,os.Getenv("REDIRECT_GITHUB_URL")),
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
	cookie, err := c.Cookie("_gothic_session")
	if err!= nil {
		fmt.Println("Error is here")
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Request.URL.RawQuery = q.Encode()
	user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		fmt.Println("Error is here")
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
    userId,err := util.CreateUser(c,user.Email)
	if err != nil {
		fmt.Println("Error is here")
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	fmt.Println(userId)
	w1 := util.Token(c,userId,cookie)
	http.Redirect(w1.Writer,w1.Request,"http://localhost:3000/home",http.StatusTemporaryRedirect)
}
func GitHubAuth(c *gin.Context){
	q := c.Request.URL.Query()
	q.Add("provider", "github")
	c.Request.URL.RawQuery = q.Encode()
	gothic.BeginAuthHandler(c.Writer, c.Request)
}

func GitHubAuthCallbackFunction(c *gin.Context){
	q := c.Request.URL.Query()
	q.Add("provider", "github")
	cookie, err := c.Cookie("_gothic_session")
	if err!= nil {
		fmt.Println("Error is here")
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Request.URL.RawQuery = q.Encode()
	user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		fmt.Println("Error is here")
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
    userId,err := util.CreateUser(c,user.UserID)
	if err != nil {
		fmt.Println("Error is here")
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	util.Token(c,userId,cookie)
	http.Redirect(c.Writer,c.Request,"http://localhost:3000/home",http.StatusTemporaryRedirect)
}

func Logout(c *gin.Context) {
	gothic.Logout(c.Writer, c.Request)
	c.SetCookie("Authorization", "", -1, "/", "localhost", false, true)
	c.Redirect(http.StatusTemporaryRedirect, "http://localhost:3000")
}