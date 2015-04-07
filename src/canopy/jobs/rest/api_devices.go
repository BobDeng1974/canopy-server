/*
 * Copyright 2014-2015 Canopy Services, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package rest

import (
)

func GET__api__devices(info *RestRequestInfo, sideEffects *RestSideEffects) (map[string]interface{}, RestError) {
    if info.Account == nil {
        return nil, NotLoggedInError()
    }

    dq := info.Account.Devices()
    devices, err := dq.DeviceList()
    if err != nil {
        return nil, InternalServerError("Device lookup failed")
    }

    timestamps := info.Query["timestamps"]
    timestamp_type := "epoch_us"
    if timestamps != nil && timestamps[0] == "rfc3339" {
        timestamp_type = "rfc3339"
    }

    //out, err := devicesToJsonObj(info.PigeonSys, devices)
    // TODO: How do we tell ws connectivity status?
    out, err := devicesToJsonObj(devices, timestamp_type)
    if err != nil {
        return nil, InternalServerError("Generating JSON")
    }

    return out, nil
}
