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

package messages

import (
    "canopy/mail"
)

func MailMessageEmailChangedAway(msg mail.MailMessage, username, activationLink, manageLink, hostname string) {
    msg.SetSubject("Your email address has changed (on " + hostname + ")")

    msg.SetHTML(`<html>
    <body style='font-family: sans-serif'>
    <table align="center" width="600" border=0 cellspacing=0 cellpadding=0 style="border-collapse: collapse;">
        <tr>
            <td bgcolor=#204080 style='border:4px solid #204080; color:#ffffff; padding: 16px 16px 0px 16px;'>
                <p>
                    Hi <b>` + username + `</b>,
                </p>
                <p>
                    <font size=6><b>Canopy</b></font>
                </p>
                <p>The open cloud for IoT.</p>
                <br>
            </td>
        </tr>
        <tr>
            <td bgcolor=#f0f0f0 style='border:4px solid #204080; color:#303030; padding: 16px 16px 16px 16px;'>
                <h3><br>Your Email Address Has Changed</h3>
                <p>
                    You will no longer recieve Canopy email at this address.
                    You should recieve a verification email at your new email
                    address.
                </p>

                <br>
            </td>
        </tr>
        <tr>
            <td bgcolor=#ffff80 style='border:4px solid #204080; color:#303030; padding: 16px 16px 16px 16px;'>
                <b>Note</b>: This email address change is only for
                <b>` + hostname + `</b>.  Other deployments of the Canopy
                Server may have separate accounts.
            </td>
        </tr>
        <tr>
            <td style='font-size:12px'>
                <br>
                <b>Web: </b><a href=http://canopy.link>canopy.link</a>
                <br><b>Twitter:</b><a href='http://twitter.com/CanopyIOT'>@CanopyIoT</a>
                <br><b>Github:</b><a href='http://github.com/canopy-project'>github.com/canopy-project</a>
                <br><b>Forum:</b><a href='http://canopy.lefora.com'>canopy.lefora.com</a>
            </td>
        </tr>
    </table>
    </body>
</html>`)
}
