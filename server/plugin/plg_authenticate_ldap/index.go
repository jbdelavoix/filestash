package plg_authenticate_ldap

import (
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/mickael-kerjean/filestash/server/common"
	"gopkg.in/ldap.v3"
)

type ldapParams struct {
	Hostname     string
	Port         string
	UseSSL       bool
	VerifySSL    bool
	BindDN       string
	BindPassword string
	BaseDN       string
	SearchFilter string
}

func init() {
	common.Hooks.Register.AuthenticationMiddleware("ldap", Ldap{})
}

type Ldap struct{}

func (this Ldap) Setup() common.Form {
	return common.Form{
		Elmnts: []common.FormElement{
			{
				Name:  "type",
				Type:  "hidden",
				Value: "ldap",
			},
			{
				Name:        "hostname",
				Type:        "text",
				Placeholder: "Hostname",
			},
			{
				Name:        "port",
				Type:        "text",
				Placeholder: "Port",
			},
			{
				Name:        "use_ssl",
				Type:        "boolean",
				Default:     true,
				Placeholder: "Use SSL",
			},
			{
				Name:        "verify_ssl",
				Type:        "boolean",
				Default:     true,
				Placeholder: "Verify SSL certificate",
			},
			{
				Name:        "bind_dn",
				Type:        "text",
				Placeholder: "Bind DN",
			},
			{
				Name:        "bind_password",
				Type:        "password",
				Placeholder: "Bind password",
			},
			{
				Name:        "base_dn",
				Type:        "text",
				Placeholder: "Base DN",
			},
			{
				Name: "search_filter",
				Type: "text",
			},
		},
	}
}

func (this Ldap) EntryPoint(idpParams map[string]string, req *http.Request, res http.ResponseWriter) error {
	getFlash := func() string {
		c, err := req.Cookie("flash")
		if err != nil {
			return ""
		}
		http.SetCookie(res, &http.Cookie{
			Name:   "flash",
			MaxAge: -1,
			Path:   "/",
		})
		return fmt.Sprintf(`<p class="flash">%s</p>`, c.Value)
	}
	res.Header().Set("Content-Type", "text/html; charset=utf-8")
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(common.Page(`
      <form action="/api/session/auth/" method="post" class="component_middleware">
        <label>
          <input type="text" name="user" value="" placeholder="User" />
        </label>
        <label>
          <input type="password" name="password" value="" placeholder="Password" />
        </label>
        <button>CONNECT</button>
        ` + getFlash() + `
        <style>
          .flash{ color: #f26d6d; font-weight: bold; }
          form { padding-top: 10vh; }
        </style>
      </form>`)))
	return nil
}

func (this Ldap) Callback(formData map[string]string, idpParams map[string]string, res http.ResponseWriter) (map[string]string, error) {
	if verifyLdapCreds(idpParams, formData["user"], formData["password"]) {
		return map[string]string{
			"user":     formData["user"],
			"password": formData["password"],
		}, nil
	}

	http.SetCookie(res, &http.Cookie{
		Name:   "flash",
		Value:  "Invalid username or password",
		MaxAge: 1,
		Path:   "/",
	})

	return nil, common.ErrAuthenticationFailed
}

func verifyLdapCreds(idpParams map[string]string, user string, password string) bool {
	params := ldapParams{
		Hostname:     idpParams["hostname"],
		Port:         idpParams["port"],
		UseSSL:       idpParams["use_ssl"] != "false",
		VerifySSL:    idpParams["verify_ssl"] != "false",
		BindDN:       idpParams["bind_dn"],
		BindPassword: idpParams["bind_password"],
		BaseDN:       idpParams["base_dn"],
		SearchFilter: idpParams["search_filter"],
	}

	dialAddr := func() string {
		if params.Port == "" {
			return params.Hostname
		}
		return fmt.Sprintf("%s:%s", params.Hostname, params.Port)
	}()

	var l *ldap.Conn
	var err error

	if params.UseSSL {
		tlsConfig := &tls.Config{InsecureSkipVerify: !params.VerifySSL}
		l, err = ldap.DialTLS("tcp", dialAddr, tlsConfig)
	} else {
		l, err = ldap.Dial("tcp", dialAddr)
	}
	if err != nil {
		common.Log.Warning("ldap: error connecting auth backend: %v", err)
		return false
	}
	defer l.Close()

	if err = l.Bind(params.BindDN, params.BindPassword); err != nil {
		common.Log.Warning("ldap: error binding auth backend: %v", err)
		return false
	}

	searchRequest := ldap.NewSearchRequest(
		params.BaseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf(params.SearchFilter, user), []string{}, nil,
	)

	sr, err := l.Search(searchRequest)
	if err != nil {
		common.Log.Warning("ldap: error searching in auth backend: %v", err)
		return false
	}

	if len(sr.Entries) != 1 {
		common.Log.Warning("ldap: user %v does not exist or too many entries returned", user)
		return false
	}

	userDN := sr.Entries[0].DN

	err = l.Bind(userDN, password)
	if err != nil {
		common.Log.Warning("ldap: provided password for user %v is invalid", user)
		return false
	}

	return true
}
