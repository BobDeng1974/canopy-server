// Copyright 2014 SimpleThings, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package endpoints

import (
    "canopy/rest/adapter"
    "canopy/rest/rest_errors"
    //"canopy/mail"
    "net/http"
)

func POST_create_account(w http.ResponseWriter, r *http.Request, info adapter.CanopyRestInfo) (map[string]interface{}, rest_errors.CanopyRestError) {
    username, ok := info.BodyObj["username"].(string)
    if !ok {
        return nil, rest_errors.NewBadInputError("String \"username\" expected")
    }

    email, ok := info.BodyObj["email"].(string)
    if !ok {
        return nil, rest_errors.NewBadInputError("String \"email\" expected")
    }

    password, ok := info.BodyObj["password"].(string)
    if !ok {
        return nil, rest_errors.NewBadInputError("String \"password\" expected")
    }

    account, err := info.Conn.LookupAccount(username)
    if err == nil {
        return nil, rest_errors.NewUsernameTakenError()
    }

    account, err = info.Conn.LookupAccount(email)
    if err == nil {
        return nil, rest_errors.NewEmailTakenError()
    }

    account, err = info.Conn.CreateAccount(username, email, password)
    if err != nil {
        return nil, rest_errors.NewInternalServerError("Problem Creating Account")
    }

    info.Session.Values["logged_in_username"] = username
    err = info.Session.Save(r, w)
    if err != nil {
        return nil, rest_errors.NewInternalServerError("Problem saving session")
    }

    // TODO send welcome email
    /*canolog.Trace("Sending email")
    mailer, err := mail.NewDefaultMailClient()
    if (err != nil) {
        canolog.Error(err)
        fmt.Fprintf(w, "{\"error\" : \"initializing_mail_client\"}")
        return
    }

    msg := mailer.NewMail();
    msg.AddTo(account.Email(), account.Username())
    msg.SetFrom("no-reply@canopy.link", "Canopy Cloud Service")
    msg.SetReplyTo("no-reply@canopy.link")
    msg.SetSubject("Welcome to Canopy")
    msg.SetHTML("Thank you for creating a Canopy account!")
    err = mailer.Send(msg)
    if (err != nil) {
        fmt.Fprintf(w, "{\"error\" : \"sending_email\"}")
        return
    }*/

    out := map[string]interface{} {
        "result" : "ok",
        "username" : account.Username(),
        "email" : account.Email(),
    }
    return out, nil
}

