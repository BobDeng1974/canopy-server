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

package config

import (
    "encoding/json"
    "flag"
    "fmt"
    "io/ioutil"
    "os"
    "strconv"
)

type CanopyConfig struct {
    allowAnonDevices bool
    allowOrigin string
    forwardOtherHosts string
    hostname string
    httpPort int16
    logFile string
    webManagerPath string
    javascriptClientPath string
}

func (config *CanopyConfig) DumpToString() string {
    return fmt.Sprint(`SERVER CONFIG SETTINGS:
allow-anon-devices:  `, config.allowAnonDevices, `
allow-origin:        `, config.allowOrigin, `
forward-other-hosts: `, config.forwardOtherHosts, `
hostname:            `, config.hostname, `
http-port:           `, config.httpPort, `
js-client-path:      `, config.javascriptClientPath, `
log-file:            `, config.logFile, `
web-manager-path:    `, config.webManagerPath)
}

func (config *CanopyConfig) LoadConfig() error {
    err := config.LoadConfigFile("/etc/canopy/server.conf")
    if os.IsNotExist(err) {
        // If file doesn't exist, just move on to the next one.
    } else if err != nil {
        return err
    }

    homeDir := os.Getenv("HOME")
    if homeDir != "" {
        err = config.LoadConfigFile(homeDir + "/.canopy/server.conf")
        if os.IsNotExist(err) {
            // If file doesn't exist, just move on to the next one.
        } else if err != nil {
            return err
        }
    }

    confFile := os.Getenv("CANOPY_SERVER_CONFIG_FILE")
    if confFile != "" {
        err = config.LoadConfigFile(confFile)
        if err != nil {
            // If config file is specified explicitely, it must be readable
            return err
        }
    }

    err = config.LoadConfigEnv()
    if err != nil {
        return err
    }

    err = config.LoadConfigCLI()
    if err != nil {
        return err
    }

    return nil
}

func (config *CanopyConfig) LoadConfigEnv() error {
    allowAnonDevices := os.Getenv("CCS_ALLOW_ANON_DEVICES")
    if allowAnonDevices == "1" || allowAnonDevices == "true" {
        config.allowAnonDevices = true
    } else if allowAnonDevices == "0" || allowAnonDevices == "false" {
        config.allowAnonDevices = false
    } else if allowAnonDevices != "" {
        return fmt.Errorf("Invalid value for CCS_ALLOW_ANON_DEVICES: %s",  allowAnonDevices)
    }

    allowOrigin := os.Getenv("CCS_ALLOW_ORIGIN")
    if allowOrigin != "" {
        config.allowOrigin = allowOrigin
    }

    forwardOtherHosts := os.Getenv("CCS_FORWARD_OTHER_HOSTS")
    if forwardOtherHosts != "" {
        config.forwardOtherHosts = forwardOtherHosts
    }

    hostname := os.Getenv("CCS_HOSTNAME")
    if hostname != "" {
        config.hostname = hostname
    }

    httpPort := os.Getenv("CCS_HTTP_PORT")
    if httpPort != "" {
        port, err := strconv.ParseInt(httpPort, 0, 16)
        if err != nil {
            return fmt.Errorf("Invalid value for CCS_HTTP_PORT: %s",  httpPort)
        }
        config.httpPort = int16(port)
    }

    jsClientPath := os.Getenv("CCS_JS_CLIENT_PATH")
    if jsClientPath != "" {
        config.javascriptClientPath = jsClientPath
    }

    logFile := os.Getenv("CCS_LOG_FILE")
    if logFile != "" {
        config.logFile = logFile
    }

    webMgrPath := os.Getenv("CCS_WEB_MANAGER_PATH")
    if webMgrPath != "" {
        config.webManagerPath = webMgrPath
    }

    return nil
}

func (config *CanopyConfig) LoadConfigCLI() error {
    allowAnonDevices := flag.String("allow-anon-devices", "", "")
    allowOrigin := flag.String("allow-origin", "", "")
    forwardOtherHosts := flag.String("forward-other-hosts", "", "")
    hostname := flag.String("hostname", "", "")
    httpPort := flag.String("http-port", "", "")
    jsClientPath := flag.String("js-client-path", "", "")
    logFile := flag.String("log-file", "", "")
    webMgrPath := flag.String("web-manager-path", "", "")

    flag.Parse()

    if *allowAnonDevices != "" {
        if *allowAnonDevices == "1" || *allowAnonDevices == "true" {
            config.allowAnonDevices = true
        } else if *allowAnonDevices == "0" || *allowAnonDevices == "false" {
            config.allowAnonDevices = false
        } else if *allowAnonDevices != "" {
            return fmt.Errorf("Invalid value for --allow-anon-devices: %s",  *allowAnonDevices)
        }
    }

    if *allowOrigin != "" {
        config.allowOrigin = *allowOrigin
    }

    if *forwardOtherHosts != "" {
        config.forwardOtherHosts = *forwardOtherHosts
    }

    if *hostname != "" {
        config.hostname = *hostname
    }

    if *httpPort != "" {
        port, err := strconv.ParseInt(*httpPort, 0, 16)
        if err != nil {
            return fmt.Errorf("Invalid value for CCS_HTTP_PORT: %s",  httpPort)
        }
        config.httpPort = int16(port)
    }

    if *jsClientPath != "" {
        config.javascriptClientPath = *jsClientPath
    }

    if *logFile != "" {
        config.logFile = *logFile
    }

    if *webMgrPath != "" {
        config.webManagerPath = *webMgrPath
    }

    return nil
}

func (config *CanopyConfig) LoadConfigFile(filename string) error {
    bytes, err := ioutil.ReadFile(filename)
    if err != nil {
        return err
    }

    s := string(bytes)

    return config.LoadConfigJsonString(s)
}

func (config *CanopyConfig) LoadConfigJsonString(jsonString string) error {
    var jsonObj map[string]interface{}

    err := json.Unmarshal([]byte(jsonString), &jsonObj)
    if err != nil {
        return err
    }

    return config.LoadConfigJson(jsonObj)

}

func (config *CanopyConfig) LoadConfigJson(jsonObj map[string]interface{}) error {
    for k, v := range jsonObj {
        ok := false
        switch k {
        case "allow-anon-devices":
            config.allowAnonDevices, ok = v.(bool)
        case "allow-origin":
            config.allowOrigin, ok = v.(string)
        case "forward-other-hosts": 
            config.forwardOtherHosts, ok = v.(string)
        case "hostname": 
            config.hostname, ok = v.(string)
        case "http-port": 
            port, ok := v.(int)
            if ok {
                config.httpPort = int16(port)
            }
        case "js-client-path": 
            config.javascriptClientPath, ok = v.(string)
        case "log-file": 
            config.logFile, ok = v.(string)
        case "web-manager-path": 
            config.webManagerPath, ok = v.(string)
        default:
            return fmt.Errorf("Unknown configuration option: %s", k)
        }

        if !ok {
            return fmt.Errorf("Incorrect JSON type for %s", k)
        }
    }
    return nil
}
func (config *CanopyConfig) OptAllowAnonDevices() bool {
    return config.allowAnonDevices
}

func (config *CanopyConfig) OptAllowOrigin() string {
    return config.allowOrigin
}

func (config *CanopyConfig) OptForwardOtherHosts() string {
    return config.forwardOtherHosts
}

func (config *CanopyConfig) OptHostname() string {
    return config.hostname
}

func (config *CanopyConfig) OptHTTPPort() int16 {
    return config.httpPort
}

func (config *CanopyConfig) OptJavascriptClientPath() string {
    return config.javascriptClientPath
}

func (config *CanopyConfig) OptLogFile() string {
    return config.logFile
}

func (config *CanopyConfig) OptWebManagerPath() string {
    return config.webManagerPath
}
