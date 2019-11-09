// Copyright 2019 WhizUs GmbH. All rights reserved.
//
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package common

import "strings"

type status string

type statusList struct {
	Error      status
	Forbidden  status
	Created    status
	Ok         status
	NotFound   status
	BadRequest status
}

type StatusError struct {
	message string
}

var Status = &statusList{
	Error:      "error",
	Forbidden:  "FORBIDDEN",
	Created:    "CREATED",
	Ok:         "OK",
	NotFound:   "NOT_FOUND",
	BadRequest: "BAD_REQUEST",
}

type StatusResponse struct {
	Status      *string            `json:"status"`
	Message     *string            `json:"message"`
	Reason      *string            `json:"reason"`
	InvalidKeys *map[string]string `json:"invalid_keys"`
}

func NewStatusError(reason string, invalidKeys map[string]string) *StatusError {
	if &reason == nil {
		return &StatusError{
			message: "Unknown reason ¯\\_(ツ)_/¯",
		}
	}

	if invalidKeys != nil {
		reason = reason + " invalid keys: "

		keys := make([]string, 0, len(invalidKeys))
		for _, key := range invalidKeys {
			keys = append(keys, key)
		}
		reason = reason + strings.Join(keys, ", ")
	}

	return &StatusError{message: reason}
}

func (e StatusError) Error() string {
	return e.message
}
