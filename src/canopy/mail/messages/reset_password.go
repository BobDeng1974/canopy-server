// Copright 2014-2015 Canopy Services, Inc.
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

func MailMessageResetPassword(msg mail.MailMessage, username, resetLink, manageLink, hostname string) {
    msg.SetSubject("Reset your Canopy password (on " + hostname + ")")

    msg.SetHTML(`<html>
    <body style='font-family: sans-serif'>
    <table align="center" width="600" border=0 cellspacing=0 cellpadding=0 style="border-collapse: collapse;">
        <tr>
            <td bgcolor=#204080 style='border:4px solid #204080; color:#ffffff; padding: 16px 16px 0px 16px;'>
                <p>
                    Hi <b>` + username + `</b>,
                </p>
                <p>
                    <font size=6><b>Canopy Password Reset</b></font>
                </p>
                <br>
            </td>
        </tr>
        <tr>
            <td bgcolor=#f0f0f0 style='border:4px solid #204080; color:#303030; padding: 16px 16px 16px 16px;'>
                <p>
                    <br><i>If you believe you have received this email in error
                    then simply disregard this message.</i>
                </p>
                <h3><br>Reset Password</h3>
                <p>
                    To reset your Canopy password, click the link below.  The
                    link will expire in 24 hours.
                </p>

                <p>
                    <a href=` + resetLink + `>Reset your password.</a>
                </p>
                <h3><br>Manage Your Devices</h3>
                After resetting your password, you can manage your
                Canopy-enabled devices by going here:
                <p>
                    <a href=` + manageLink + `>` + manageLink + `</a>
                </p>
                <br>
            </td>
        </tr>
        <tr>
            <td bgcolor=#ffff80 style='border:4px solid #204080; color:#303030; padding: 16px 16px 16px 16px;'>
                <b>Note</b>: This is only for
                <b>` + hostname + `</b>.  Other deployments of the Canopy
                Server have separate accounts.
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
