// Copyright 2026 Cognit. All rights reserved.
// License can be found in the LICENSE file.

package auth

import "errors"

var (
	ErrType        = errors.New("auth: unknown auth type")
	ErrConfig      = errors.New("auth: incomplete scheme config")
	ErrCredentials = errors.New("auth: invalid credentials")
)
