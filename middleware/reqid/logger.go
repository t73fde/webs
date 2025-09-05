//-----------------------------------------------------------------------------
// Copyright (c) 2025-present Detlef Stern
//
// This file is part of webs.
//
// webs is licensed under the latest version of the EUPL (European Union Public
// License. Please see file LICENSE.txt for your rights and obligations under
// this license.
//
// SPDX-License-Identifier: EUPL-1.2
// SPDX-FileCopyrightText: 2025-present Detlef Stern
//-----------------------------------------------------------------------------

package reqid

import (
	"context"
	"log/slog"
)

// WithLogger enhances a base logger to add logging the value of the
// request identifier, stored in a Context.
func (c *Config) WithLogger(base *slog.Logger) *slog.Logger {
	if !c.WithContext || base == nil {
		return base
	}
	key := c.LoggingKey
	if key == "" {
		key = DefaultLoggingKey
	}
	return slog.New(&reqidLogHandler{h: base.Handler(), key: key})
}

// reqidLogHandler ist ein Wrapper um einen bestehenden slog.Handler.
type reqidLogHandler struct {
	h   slog.Handler
	key string
}

func (rid *reqidLogHandler) Handle(ctx context.Context, r slog.Record) error {
	if ctx != nil {
		reqID, ok := getReqID(ctx)
		if ok {
			r.AddAttrs(slog.Any(rid.key, reqID))
		}
	}
	return rid.h.Handle(ctx, r)
}

func (rid *reqidLogHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return rid.h.Enabled(ctx, level)
}

func (rid *reqidLogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &reqidLogHandler{h: rid.h.WithAttrs(attrs), key: rid.key}
}

func (rid *reqidLogHandler) WithGroup(name string) slog.Handler {
	return &reqidLogHandler{h: rid.h.WithGroup(name), key: rid.key}
}
