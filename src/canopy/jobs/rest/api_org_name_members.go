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

package rest

import (
)

// Backend implementation /api/org/{name}/members endpoint
// 
func GET__api__org__name__members(info *RestRequestInfo, sideEffects *RestSideEffects) (map[string]interface{}, RestError) {
    if info.Account == nil {
        return nil, NotLoggedInError().Log()
    }

    // Lookup organization by name
    org, err := info.Conn.LookupOrganization(info.URLVars["name"])
    if err != nil {
        // TODO: determine actual cause of error
        return nil, URLNotFoundError().Log()
    }

    // Is authenticated user member of this organization?
    ok, err := org.IsMember(info.Account)
    if !ok || err != nil {
        return nil, UnauthorizedError("Must be member of org to see other members").Log()
    }

    // List members
    members, err := org.Members()
    if err != nil {
        return nil, InternalServerError(err.Error()).Log()
    }

    // Compose JSON response payload

    memberJsonList := []map[string]interface{} {}
    for _, member := range members {
        memberJsonList = append(memberJsonList, map[string]interface{} {
                "username" : member.Username(),
                "email" : member.Email(),
                // TODO: add teams
        })
    }

    out := map[string]interface{} {
        "members" : memberJsonList,
        "result" : "ok",
    }

    return out, nil
}

// Add/Remoe members
func POST__api__org__name__members(info *RestRequestInfo, sideEffects *RestSideEffects) (map[string]interface{}, RestError) {
    if info.Account == nil {
        return nil, NotLoggedInError().Log()
    }

    // Lookup organization by name
    org, err := info.Conn.LookupOrganization(info.URLVars["name"])
    if err != nil {
        // TODO: determine actual cause of error
        return nil, URLNotFoundError().Log()
    }

    // Is authenticated user member of this organization?
    // TODO: They need to be the owner.
    ok, err := org.IsOwner(info.Account)
    if !ok || err != nil {
        return nil, UnauthorizedError("Must be owner of org to add members").Log()
    }

    addMemberObj, ok := info.BodyObj["add_member"].(map[string]interface{})
    if ok {
        user, ok := addMemberObj["user"].(string)
        if !ok {
            return nil, BadInputError("Expected string \"user\"").Log()
        }
        acct, err := info.Conn.LookupAccount(user)
        if err != nil {
            // TODO: Proper error reporting
            return nil, InternalServerError(err.Error()).Log()
        }

        err = org.AddMember(acct)
        if err != nil {
            // TODO: Proper error reporting
            return nil, InternalServerError(err.Error()).Log()
        }
    }

    removeMember, ok := info.BodyObj["remove_member"].(string)
    if ok {
        acct, err := info.Conn.LookupAccount(removeMember)
        if err != nil {
            // TODO: Proper error reporting
            return nil, InternalServerError(err.Error()).Log()
        }

        err = org.RemoveMember(acct)
        if err != nil {
            // TODO: Proper error reporting
            return nil, InternalServerError(err.Error()).Log()
        }
    }

    out := map[string]interface{} {
        "result" : "ok",
    }

    return out, nil
}
