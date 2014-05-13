package main

import (
    "fmt"
    "net/http"
    "code.google.com/p/go.net/websocket"
    "github.com/gorilla/sessions"
    "github.com/gorilla/context"
    "canopy/datalayer"
    "encoding/json"
)

var store = sessions.NewCookieStore([]byte("my_production_secret"))

func loginHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Access-Control-Allow-Origin", "http://canopy.link")
    w.Header().Set("Access-Control-Allow-Credentials", "true")

    var data map[string]interface{}
    decoder := json.NewDecoder(r.Body)
    err := decoder.Decode(&data)
    if err != nil {
        fmt.Fprintf(w, "{\"error\" : \"json_decode_failed\"}")
        return
    }

    username, ok := data["username"].(string)
    if !ok {
        fmt.Fprintf(w, "{\"error\" : \"string_username_expected\"}")
        return
    }

    password, ok := data["password"].(string)
    if !ok {
        fmt.Fprintf(w, "{\"error\" : \"string_password_expected\"}")
        return
    }

    session, _ := store.Get(r, "canopy-login-session")
    dl := datalayer.NewCassandraDatalayer()
    dl.Connect("canopy")
    if dl.VerifyAccountPassword(username, password) {
        session.Values["logged_in_username"] = username
        err := session.Save(r, w)
        if err != nil {
            fmt.Fprintf(w, "{\"error\" : \"saving_session\"}")
            return
        }
        fmt.Fprintf(w, "{\"success\" : true}")
        return
    } else {
        fmt.Fprintf(w, "{\"error\" : \"incorrect_password\"}")
        return
    }
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Access-Control-Allow-Origin", "http://canopy.link")
    w.Header().Set("Access-Control-Allow-Credentials", "true")
    session, _ := store.Get(r, "canopy-login-session")
    session.Values["logged_in_username"] = ""
    err := session.Save(r, w)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError);
        fmt.Fprintf(w, "{ \"error\" : \"could_not_logout\"");
        return;
    }
    fmt.Fprintf(w, "{ \"success\" : true }")
}
func meHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Access-Control-Allow-Origin", "http://canopy.link")
    w.Header().Set("Access-Control-Allow-Credentials", "true")
    session, _ := store.Get(r, "canopy-login-session")
    
    username, ok := session.Values["logged_in_username"]
    if ok {
        username_string, ok := username.(string)
        if ok && username_string != "" {
            fmt.Fprintf(w, "{\"username\" : \"%s\"}", username_string);
            return
        } else {
            w.WriteHeader(http.StatusUnauthorized);
            fmt.Fprintf(w, "{\"error\" : \"not_logged_in\"");
            return
        }
    } else {
        w.WriteHeader(http.StatusUnauthorized);
        fmt.Fprintf(w, "{\"error\" : \"not_logged_in\"");
        return
    }
}

func main() {
    fmt.Println("starting server");
    http.Handle("/echo", websocket.Handler(CanopyWebsocketServer))
    http.HandleFunc("/login", loginHandler)
    http.HandleFunc("/logout", logoutHandler)
    http.HandleFunc("/me", meHandler)
    http.ListenAndServe(":8080", context.ClearHandler(http.DefaultServeMux))
}
