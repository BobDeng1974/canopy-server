package main

import (
    "fmt"
    "net/http"
    "code.google.com/p/go.net/websocket"
    "github.com/gocql/gocql"
    "github.com/gorilla/sessions"
    "github.com/gorilla/context"
    "github.com/gorilla/mux"
    "canopy/datalayer"
    "canopy/mail"
    "canopy/pigeon"
    "encoding/json"
    "time"
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
    _, err = dl.LookupAccountVerifyPassword(username, password)
    if err == nil {
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

func createAccountHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Access-Control-Allow-Origin", "http://canopy.link")
    w.Header().Set("Access-Control-Allow-Credentials", "true")

    var data map[string]interface{}
    decoder := json.NewDecoder(r.Body)
    err := decoder.Decode(&data)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest);
        fmt.Fprintf(w, "{\"error\" : \"json_decode_failed\"}")
        return
    }

    username, ok := data["username"].(string)
    if !ok {
        w.WriteHeader(http.StatusBadRequest);
        fmt.Fprintf(w, "{\"error\" : \"string_username_expected\"}")
        return
    }

    email, ok := data["username"].(string)
    if !ok {
        w.WriteHeader(http.StatusBadRequest);
        fmt.Fprintf(w, "{\"error\" : \"string_email_expected\"}")
        return
    }

    password, ok := data["password"].(string)
    if !ok {
        w.WriteHeader(http.StatusBadRequest);
        fmt.Fprintf(w, "{\"error\" : \"string_password_expected\"}")
        return
    }

    password_confirm, ok := data["password_confirm"].(string)
    if !ok {
        w.WriteHeader(http.StatusBadRequest);
        fmt.Fprintf(w, "{\"error\" : \"string_password_confirm_expected\"}")
        return
    }

    if (password != password_confirm) {
        w.WriteHeader(http.StatusBadRequest);
        fmt.Fprintf(w, "{\"error\" : \"passwords_dont_match\"}")
        return
    }

    dl := datalayer.NewCassandraDatalayer()
    dl.Connect("canopy")

    dl.CreateAccount(username, email, password);
    session, _ := store.Get(r, "canopy-login-session")
    session.Values["logged_in_username"] = username
    err = session.Save(r, w)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError);
        fmt.Fprintf(w, "{\"error\" : \"saving_session\"}")
        return
    }
    fmt.Fprintf(w, "{\"success\" : true}")
    return
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
            fmt.Fprintf(w, "{\"error\" : \"not_logged_in\"}");
            return
        }
    } else {
        w.WriteHeader(http.StatusUnauthorized);
        fmt.Fprintf(w, "{\"error\" : \"not_logged_in\"}");
        return
    }
}

func devicesHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Access-Control-Allow-Origin", "http://canopy.link")
    w.Header().Set("Access-Control-Allow-Credentials", "true")
    session, _ := store.Get(r, "canopy-login-session")
    
    var username_string string
    username, ok := session.Values["logged_in_username"]
    if ok {
        username_string, ok = username.(string)
        if !(ok && username_string != "") {
            w.WriteHeader(http.StatusUnauthorized);
            fmt.Fprintf(w, "{\"error\" : \"not_logged_in\"");
            return
        }
    } else {
        w.WriteHeader(http.StatusUnauthorized);
        fmt.Fprintf(w, "{\"error\" : \"not_logged_in\"");
        return
    }
    
    dl := datalayer.NewCassandraDatalayer()
    dl.Connect("canopy")
    account, err := dl.LookupAccount(username_string)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError);
        fmt.Fprintf(w, "{\"error\" : \"account_lookup_failed\"}");
        return
    }

    devices, err := account.GetDevices()
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError);
        fmt.Fprintf(w, "{\"error\" : \"device_lookup_failed\"}");
        return
    }
    out, err := devicesToJson(devices)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError);
        fmt.Fprintf(w, "{\"error\" : \"generating_json\"}");
        return
    }
    fmt.Fprintf(w, out);
  
    return 
}

func sensorDataHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    deviceIdString := vars["id"]
    sensorName := vars["sensor"]

    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Access-Control-Allow-Origin", "http://canopy.link")
    w.Header().Set("Access-Control-Allow-Credentials", "true")
    session, _ := store.Get(r, "canopy-login-session")
    
    var username_string string
    username, ok := session.Values["logged_in_username"]
    if ok {
        username_string, ok = username.(string)
        if !(ok && username_string != "") {
            w.WriteHeader(http.StatusUnauthorized);
            fmt.Fprintf(w, "{\"error\" : \"not_logged_in1\"}");
            return
        }
    } else {
        w.WriteHeader(http.StatusUnauthorized);
        fmt.Fprintf(w, "{\"error\" : \"not_logged_in2\"}");
        return
    }
    
    dl := datalayer.NewCassandraDatalayer()
    dl.Connect("canopy")
    account, err := dl.LookupAccount(username_string)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError);
        fmt.Fprintf(w, "{\"error\" : \"account_lookup_failed\"}");
        return
    }

    uuid, err := gocql.ParseUUID(deviceIdString)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest);
        fmt.Fprintf(w, "{\"error\" : \"Device UUID expected\"}");
        return
    }

    device, err := account.GetDeviceById(uuid)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest);
        fmt.Fprintf(w, "{\"error\" : \"Could not find or access device\"}");
        return
    }

    samples, err := device.GetSensorData(sensorName, time.Now(), time.Now())
    if err != nil {
        fmt.Println(err)
        w.WriteHeader(http.StatusInternalServerError);
        fmt.Fprintf(w, "{\"error\" : \"Could not obtain sample data\"}");
        return
    }

    out, err := samplesToJson(samples)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError);
        fmt.Fprintf(w, "{\"error\" : \"generating_json\"} : ", err);
        return
    }

    fmt.Fprintf(w, out);
    return 
}

func controlHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    deviceIdString := vars["id"]
    //controlName := vars["control"]

    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Access-Control-Allow-Origin", "http://canopy.link")
    w.Header().Set("Access-Control-Allow-Credentials", "true")
    session, _ := store.Get(r, "canopy-login-session")
    
    var username_string string
    username, ok := session.Values["logged_in_username"]
    if ok {
        username_string, ok = username.(string)
        if !(ok && username_string != "") {
            w.WriteHeader(http.StatusUnauthorized);
            fmt.Fprintf(w, "{\"error\" : \"not_logged_in1\"}");
            return
        }
    } else {
        w.WriteHeader(http.StatusUnauthorized);
        fmt.Fprintf(w, "{\"error\" : \"not_logged_in2\"}");
        return
    }
    
    dl := datalayer.NewCassandraDatalayer()
    dl.Connect("canopy")
    account, err := dl.LookupAccount(username_string)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError);
        fmt.Fprintf(w, "{\"error\" : \"account_lookup_failed\"}");
        return
    }

    uuid, err := gocql.ParseUUID(deviceIdString)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest);
        fmt.Fprintf(w, "{\"error\" : \"Device UUID expected\"}");
        return
    }

    _, err = account.GetDeviceById(uuid)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest);
        fmt.Fprintf(w, "{\"error\" : \"Could not find or access device\"}");
        return
    }

    /* Parse input as json and just forward it along using pigeon */
    var data map[string]interface{}
    decoder := json.NewDecoder(r.Body)
    err = decoder.Decode(&data)
    if err != nil {
        fmt.Fprintf(w, "{\"error\" : \"json_decode_failed\"}")
        return
    }

    msg := &pigeon.PigeonMessage { 
        Data : data,
    }
    err = gPigeon.SendMessage(deviceIdString, msg, time.Duration(100*time.Millisecond))
    if err != nil {
        fmt.Fprintf(w, "{\"error\" : \"SendMessage failed\"}");
    }

    fmt.Fprintf(w, "{\"result\" : \"ok\"}");
    return 
}

func shareHandler(w http.ResponseWriter, r *http.Request) {
    /*
     *  POST
     *  {
     *      "device_id" : <DEVICE_ID>,
     *      "access_level" : <ACCESS_LEVEL>,
     *      "sharing_level" : <SHARING_LEVEL>,
     *      "email_address" : <EMAIL_ADDRESS>,
     *  }
     *
     * TODO: Add to REST API documentation
     */
    var data map[string]interface{}
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Access-Control-Allow-Origin", "http://canopy.link")
    w.Header().Set("Access-Control-Allow-Credentials", "true")
    session, _ := store.Get(r, "canopy-login-session")

    decoder := json.NewDecoder(r.Body)
    err := decoder.Decode(&data)
    if err != nil {
        fmt.Fprintf(w, "{\"error\" : \"json_decode_failed\"}")
        return
    }

    deviceId, ok := data["device_id"].(string)
    if !ok {
        fmt.Fprintf(w, "{\"error\" : \"device_id expected\"}")
        return
    }

    //accessLevel, ok := data["access_level"].(int)
    /*_, ok = data["access_level"].(float)
    if !ok {
        fmt.Fprintf(w, "{\"error\" : \"access_level expected\"}")
        return
    }*/

    //sharingLevel, ok := data["sharing_level"].(int)
    /*_, ok = data["sharing_level"].(float)
    if !ok {
        fmt.Fprintf(w, "{\"error\" : \"sharing_level expected\"}")
        return
    }*/

    email, ok := data["email"].(string)
    if !ok {
        fmt.Fprintf(w, "{\"error\" : \"email expected\"}")
        return
    }
    var username_string string
    username, ok := session.Values["logged_in_username"]
    if ok {
        username_string, ok = username.(string)
        if !(ok && username_string != "") {
            w.WriteHeader(http.StatusUnauthorized);
            fmt.Fprintf(w, "{\"error\" : \"not_logged_in\"");
            return
        }
    } else {
        w.WriteHeader(http.StatusUnauthorized);
        fmt.Fprintf(w, "{\"error\" : \"not_logged_in\"");
        return
    }

    dl := datalayer.NewCassandraDatalayer()
    dl.Connect("canopy")
    account, err := dl.LookupAccount(username_string)
    if account == nil || err != nil {
        w.WriteHeader(http.StatusInternalServerError);
        fmt.Fprintf(w, "{\"error\" : \"account_lookup_failed\"}");
        return
    }

    mailer, err := mail.NewDefaultMailClient()
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError);
        fmt.Fprintf(w, "{\"error\" : \"Failed to initialize mail client\"}")
        return
    }
    mail := mailer.NewMail();
    err = mail.AddTo(email, "")
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError);
        fmt.Fprintf(w, "{\"error\" : \"Invalid email recipient\"}")
        return
    }
    mail.SetSubject("Greg's Smart fan")
    mail.SetHTML(`
<img src="http://canopy.link/canopy_logo.jpg"></img>
<h2>I've shared a device with you.</h2>
<a href="http://canopy.link/canopy-app/index_nodes.html?share_device=` + deviceId + `">Greg's Smart Fan</a>
<h2>What is Canopy?</h2>
<b>Canopy</b> is a secure platform for monitoring and controlling physical
devices.  Learn more at <a href=http://canopy.link>http://canopy.link</a>
`)
    mail.SetFrom("greg@canopy.link", "greg (via Canopy)")
    err = mailer.Send(mail)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError);
        fmt.Fprintf(w, "{\"error\" : \"Error sending email\"}")
        return
    }

    fmt.Fprintf(w, "{\"result\" : \"ok\"}");
    return 
}

func finishShareTransactionHandler(w http.ResponseWriter, r *http.Request) {
    /*
     *  POST
     *  {
     *      "device_id" : <DEVICE_ID>,
     *  }
     *
     * TODO: Add to REST API documentation
     * TODO: Highly insecure!!!
     */
    var data map[string]interface{}
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Access-Control-Allow-Origin", "http://canopy.link")
    w.Header().Set("Access-Control-Allow-Credentials", "true")
    session, _ := store.Get(r, "canopy-login-session")

    decoder := json.NewDecoder(r.Body)
    err := decoder.Decode(&data)
    if err != nil {
        fmt.Fprintf(w, "{\"error\" : \"json_decode_failed\"}")
        return
    }

    deviceId, ok := data["device_id"].(string)
    if !ok {
        fmt.Fprintf(w, "{\"error\" : \"device_id expected\"}")
        return
    }

    var username_string string
    username, ok := session.Values["logged_in_username"]
    if ok {
        username_string, ok = username.(string)
        if !(ok && username_string != "") {
            w.WriteHeader(http.StatusUnauthorized);
            fmt.Fprintf(w, "{\"error\" : \"not_logged_in\"");
            return
        }
    } else {
        w.WriteHeader(http.StatusUnauthorized);
        fmt.Fprintf(w, "{\"error\" : \"not_logged_in\"");
        return
    }

    dl := datalayer.NewCassandraDatalayer()
    dl.Connect("canopy")
    account, err := dl.LookupAccount(username_string)
    if account == nil || err != nil {
        w.WriteHeader(http.StatusInternalServerError);
        fmt.Fprintf(w, "{\"error\" : \"account_lookup_failed\"}");
        return
    }

    device, err := dl.LookupDeviceByStringId(deviceId)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest);
        fmt.Fprintf(w, "{\"error\" : \"device_lookup_failed\"}");
        return
    }

    /* Grant permissions to the user to access the device */
    err = device.SetAccountAccess(account, datalayer.ReadWriteShareAccess)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError);
        fmt.Fprintf(w, "{\"error\" : \"could_not_grant_access\"}");
        return
    }

    fmt.Fprintf(w, "{\"result\" : \"ok\", \"device_friendly_name\" : \"%s\" }", device.GetFriendlyName());
    return 
}

var gPigeon = pigeon.InitPigeonSystem()

func main() {
    fmt.Println("starting server");

    r := mux.NewRouter()
    r.HandleFunc("/create_account", createAccountHandler)
    /*r.HandleFunc("/device/{id}", getDeviceInfoHandler).Methods("GET");*/
    r.HandleFunc("/device/{id}", controlHandler).Methods("POST");
    r.HandleFunc("/device/{id}/{sensor}", sensorDataHandler).Methods("GET");
    r.HandleFunc("/devices", devicesHandler)
    r.HandleFunc("/share", shareHandler)
    r.HandleFunc("/finish_share_transaction", finishShareTransactionHandler)
    r.HandleFunc("/login", loginHandler);
    r.HandleFunc("/logout", logoutHandler);
    r.HandleFunc("/me", meHandler);

    http.Handle("/echo", websocket.Handler(CanopyWebsocketServer))
    http.Handle("/", r)
    http.ListenAndServe(":8080", context.ClearHandler(http.DefaultServeMux))
}
