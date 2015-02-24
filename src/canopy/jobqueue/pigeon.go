// Copyright 2015 Gregory Prisament
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

package jobqueue

import (
    "canopy/datalayer"
    "fmt"
)

type PigeonSystem struct {
    dl datalayer.PigeonSystem
}
func (pigeon *PigeonSystem) NewClient() Client {
    return &PigeonClient{
        sys: pigeon,
        timeoutms: -1,
    }
}

func (pigeon *PigeonSystem) NewResponse() Response {
    return &PigeonResponse{}
}

func (pigeon *PigeonSystem) StartServer(hostname string) (Server, error) {
    server := &PigeonServer{
        sys : pigeon,
        hostname: hostname,
        handlers : map[string]HandlerFunc{},
    }

    err := server.Start()
    if err != nil {
        return nil, err
    }

    return server, nil
}

func (pigeon *PigeonSystem) Server(hostname string) (Server, error) {
    return nil, fmt.Errorf("Not Implemented")
}

func (pigeon *PigeonSystem) Servers() ([]Server, error) {
    return nil, fmt.Errorf("Not Implemented")
}

func (resp *PigeonResponse) Body() map[string]interface{} {
    return resp.RespBody
}

func (resp *PigeonResponse) Error() error {
    return resp.RespErr
}

func (resp *PigeonResponse) SetBody(body map[string]interface{}) {
    resp.RespBody = body
}

func (resp *PigeonResponse) SetError(err error) {
    resp.RespErr = err
}

func (req *PigeonRequest) Body() map[string]interface{} {
    return req.ReqBody
}

