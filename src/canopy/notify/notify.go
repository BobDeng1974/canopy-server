/*
 * Copyright 2014 SimpleThings, Inc.
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
package notify

import (
    "canopy/datalayer"
    "canopy/mail"
)

func ProcessNotification(device datalayer.Device, notifyType string, msg string) error {
    // TODO: Add to notification log

    // Send email
    if notifyType == "email" {
        mailClient, err := mail.NewDefaultMailClient()
        if err != nil {
            return err
        }

        mailMsg := mailClient.NewMail()
        mailMsg.SetSubject("Message from your device")
        mailMsg.SetText(msg)

        // TODO: hack!
        mailMsg.AddTo("greg@canopy.link", "Greg")
        err = mailClient.Send(mailMsg)
        if err != nil {
            return err
        }
    }

    return nil
}

