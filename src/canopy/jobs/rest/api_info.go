// Copyright 2014-2015 Canopy Services, Inc.
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

package rest

import (
    canotime "canopy/util/time"
    "time"
)

// Constructs the response body for the /api/info REST endpoint
func GET__api__info(info *RestRequestInfo, sideEffect *RestSideEffects) (map[string]interface{}, RestError) {
    t := time.Now().UTC()
    return map[string]interface{}{
        "config" : info.Config.ToJsonObject(),
        "result" : "ok",
        "clock_us" : canotime.EpochMicroseconds(t),
        "clock_utc" : canotime.RFC3339(t),
        "service-name" : "Canopy Cloud Service",
        "version" : "0.9.2-beta",
    }, nil
}
