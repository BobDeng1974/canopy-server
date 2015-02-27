// Copyright 2015 Canopy Services, Inc.
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

package jobs

import (
    "canopy/config"
    "canopy/mail"
    "canopy/datalayer/cassandra_datalayer"
    "canopy/jobqueue"
    "canopy/jobs/rest"
)

func InitJobServer(cfg config.Config) error {
    pigeon, err := jobqueue.NewPigeonSystem(cfg)
    if err != nil {
        return err
    }

    mailer, err := mail.NewMailClient(cfg)
    if err != nil {
        return err
    }

    userCtx := map[string]interface{}{
        "cfg" : cfg,
        "mailer" : mailer,
    }

    server, err := pigeon.StartServer("localhost") // TODO use configured hostname
    if err != nil {
        return err
    }

    dl := cassandra_datalayer.NewDatalayer(cfg)
    conn, err := dl.Connect("canopy")
    if err != nil {
        return err
    }
    userCtx["db-conn"] = conn

    routes := map[string]jobqueue.HandlerFunc{
        "api/activate": rest.RestJobWrapper(rest.ApiActivateHandler),
        "api/create_account": rest.RestJobWrapper(rest.ApiCreateAccountHandler),
        "api/create_devices": rest.RestJobWrapper(rest.ApiCreateDevicesHandler),
        "api/info": rest.ApiInfoHandler,
        "api/me": rest.RestJobWrapper(rest.ApiMeHandler),
    }

    // Register handlers
    for jobKey, handler := range routes {
        err = server.Handle(jobKey, handler, userCtx)
        if err != nil {
            return err
        }
    }

    return nil
}
