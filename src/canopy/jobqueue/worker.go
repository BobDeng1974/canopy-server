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
    "fmt"
    "net"
    "net/rpc"
    "net/http"
)

type PigeonWorker struct {
    sys *PigeonSystem
    hostname string
    listeners map[string]PigeonListener
}

type PigeonRequest struct {
    body map[string]interface{}
}

type PigeonResponse struct {
    err error
    body map[string]interface{}
}

type PigeonListener struct {
    worker *PigeonWorker
    requestChan chan<- Request
    responseChan <-chan Response
}
func (worker *PigeonWorker) HandleRequest(request map[string]interface{}, response *PigeonResponse) error {
    // Get job key from request
    key, ok := request["key"].(string)
    if !ok {
        return fmt.Errorf("Pigeon Worker: Expected string \"key\" in request")
    }

    // Lookup the listener for that job type
    listener, ok := worker.listeners[key]
    if !ok {
        // NOT FOUND
        return fmt.Errorf("Pigeon Worker: No handler for job key %s on worker %s", key, worker.hostname)
    }

    // Post to the request to the listener's channel
    req := &PigeonRequest{
        body: request,
    }

    listener.requestChan <- req

    // Wait for response
    response, ok = (<-listener.responseChan).(*PigeonResponse)
    if !ok {
        return fmt.Errorf("Pigeon Worker: Expected PigeonResponse from response handler)")
    }

    return nil
}

func (worker *PigeonWorker) serveRPC() error {
    PIGEON_RPC_PORT := ":1888"
    rpc.Register(worker)
    rpc.HandleHTTP()
    l, err := net.Listen("tcp", PIGEON_RPC_PORT)
    if err != nil {
        return err
    }
    go http.Serve(l, nil)
    return nil
}

func (worker *PigeonWorker) Listen(key string, requestChan chan<- Request, responseChan <-chan Response) error {
    err := worker.sys.dl.RegisterListener(worker.hostname, key)

    listener := PigeonListener{
        worker: worker,
        requestChan: requestChan,
        responseChan: responseChan,
    }

    worker.listeners[key] = listener

    return err
}

func (worker *PigeonWorker) Start() error {
    err := worker.sys.dl.RegisterWorker(worker.hostname)
    if err != nil {
        return err
    }

    err = worker.serveRPC()
    if err != nil {
        // TODO: unregister?
        return err
    }

    return nil
}

func (worker *PigeonWorker) Status() error {
    return fmt.Errorf("Not implemented")
}

func (worker *PigeonWorker) Stop() error {
    return fmt.Errorf("Not implemented")
}

func (worker *PigeonWorker) StopListening(key string) error {
    return fmt.Errorf("Not implemented")
}
