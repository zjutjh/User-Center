package oauth2_handler

import (
	"auth-server/config"
	"auth-server/database"
	"auth-server/util"
	"fmt"
	"github.com/go-oauth2/oauth2/v4/errors"
	"github.com/go-oauth2/oauth2/v4/generates"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-oauth2/oauth2/v4/store"
	"github.com/go-session/session"
	"log"
	"net/http"
	"strings"
)

var (
	manager *manage.Manager
	clientStore *store.ClientStore
	srv *server.Server
)

func importClient(clientStore *store.ClientStore) {
	clients := config.Get().Client
	length := len(clients)
	for i := 0; i < length; i++ {
		if err := clientStore.Set(clients[i].ID, &models.Client{
			ID: clients[i].Name,
			Secret: clients[i].Secret,
			Domain: clients[i].Domain,
		}); err != nil {
			log.Println("invalid client configuration: ", clients[i].ID)
		}
	}
}

type AuthRequiredError struct {}

func (a *AuthRequiredError) Error() string {
	return "auth_required"
}

func init() {
	manager = manage.NewDefaultManager()
	manager.SetAuthorizeCodeTokenCfg(manage.DefaultAuthorizeCodeTokenCfg)

	manager.MustTokenStorage(store.NewMemoryTokenStore())
	manager.MapAccessGenerate(generates.NewAccessGenerate())

	clientStore = store.NewClientStore()
	importClient(clientStore)
	manager.MapClientStorage(clientStore)

	srv = server.NewServer(server.NewConfig(), manager)
	srv.SetInternalErrorHandler(internalErrorHandler)
	srv.SetResponseErrorHandler(responseErrorHandler)
	srv.SetUserAuthorizationHandler(userAuthorizationHandler)
	srv.SetAuthorizeScopeHandler(authorizeScopeHandler)
}

func internalErrorHandler(err error) (re *errors.Response) {
	log.Println(err.Error())
	return
}

func responseErrorHandler(re *errors.Response) {
	log.Println(re.Error.Error())
}

func userAuthorizationHandler(w http.ResponseWriter, r *http.Request) (userID string, err error) {
	sess, err := session.Start(r.Context(), w, r)
	if err != nil {
		log.Println(err)
		return
	}
	uid, ok := sess.Get("uid")
	if !ok {
		if r.Form == nil {
			if err = r.ParseForm(); err != nil {
				util.RawHttpError(w, "invalid form", http.StatusNotAcceptable)
				log.Println(err)
				return
			}
		}
		sess.Set("raw_query_data", r.Form)
		sess.Set("passive_login", true)
		if err = sess.Save(); err != nil {
			log.Println(err)
		}
		w.Header().Set("Location", "/user/login")
		w.WriteHeader(http.StatusFound)
		return "", nil
	}
	userID = uid.(string)
	return
}

func authorizeScopeHandler(w http.ResponseWriter, r *http.Request) (scope string, err error) {
	sess, err := session.Start(r.Context(), w, r)
	if err != nil {
		log.Println(err)
		return
	}
	if r.Form == nil {
		if err = r.ParseForm(); err != nil {
			util.RawHttpError(w, "invalid form", http.StatusNotAcceptable)
			log.Println(err)
			return
		}
	}
	uid, _ := sess.Get("uid")
	userId := uid.(string)
	clientId := r.Form.Get("client_id")
	scope = r.Form.Get("scope")
	scopes := strings.Split(scope, ",")
	for i := 0; i < len(scopes); i++ {
		if !util.IsValidPermission(scopes[i]) || !util.IsAppliedPermission(scopes[i], clientId) {
			util.RawHttpError(w, "invalid scope", http.StatusNotAcceptable)
			return
		}
	}
	scopes = database.TrimAuthorizedScope(scopes, clientId, userId)
	scope = strings.Join(scopes, ",")
	if scope != "" {
		r.Form.Set("scope", scope)
		sess.Set("raw_query_data", r.Form)
		if err = sess.Save(); err != nil {
			log.Println(err)
		}
		w.Header().Set("Location", fmt.Sprintf("/oauth2/auth?text=%s", scope))
		w.WriteHeader(http.StatusFound)
		return "all", &AuthRequiredError{}
	}
	sess.Delete("raw_query_data")
	if err = sess.Save(); err != nil {
		log.Println(err)
	}
	return "all", nil
}
