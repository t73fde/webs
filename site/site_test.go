//-----------------------------------------------------------------------------
// Copyright (c) 2026-present Detlef Stern
//
// This file is part of webs.
//
// webs is licensed under the latest version of the EUPL (European Union Public
// License. Please see file LICENSE.txt for your rights and obligations under
// this license.
//
// SPDX-License-Identifier: EUPL-1.2
// SPDX-FileCopyrightText: 2026-present Detlef Stern
//-----------------------------------------------------------------------------

package site_test

import (
	"strconv"
	"testing"

	"t73f.de/r/webs/site"
)

func TestBuilderFor(t *testing.T) {
	st := &site.Site{
		Name: "TEST-SITE",
		Root: site.Node{
			ID: "root", Title: "Start",
			Children: []*site.Node{
				{
					ID: "terms", Nodepath: "terms", Title: "Terms",
				},
				{
					ID: "courses", Nodepath: "courses", Title: "Courses",
					Children: []*site.Node{
						{
							ID: "course-show", Nodepath: "{course}", Title: "Course Show",
							Children: []*site.Node{
								{
									ID: "topic-show", Nodepath: "{topic}", Title: "Topic Show",
									Children: []*site.Node{
										{
											ID: "topic-edit", Nodepath: "edit", Title: "Topic Edit",
										},
									},
								},
								{
									ID: "teams", Nodepath: "teams", Title: "Teams",
								},
							},
						},
					},
				},
				{
					ID: "users", Nodepath: "users", Title: "Users",
				},
			},
		},
	}
	if err := st.Bake(); err != nil {
		t.Error(err)
		return
	}

	testcases := []struct {
		nodeID string
		args   []string
		exp    string
	}{
		{"terms", nil, "/terms/"},
		{"terms", []string{"e"}, "/terms/e/"},
		{"courses", nil, "/courses/"},
		{"courses", []string{"e"}, "/courses/e/"},
		{"course-show", nil, "/courses/missing-arg-0/"},
		{"course-show", []string{"c"}, "/courses/c/"},
		{"course-show", []string{"c", "e"}, "/courses/c/e/"},
		{"topic-show", nil, "/courses/missing-arg-0/missing-arg-1/"},
		{"topic-show", []string{"c"}, "/courses/c/missing-arg-1/"},
		{"topic-show", []string{"c", "t"}, "/courses/c/t/"},
		{"topic-show", []string{"c", "t", "e"}, "/courses/c/t/e/"},
		{"topic-edit", nil, "/courses/missing-arg-0/missing-arg-1/edit/"},
		{"topic-edit", []string{"c"}, "/courses/c/missing-arg-1/edit/"},
		{"topic-edit", []string{"c", "t"}, "/courses/c/t/edit/"},
		{"topic-edit", []string{"c", "t", "e"}, "/courses/c/t/edit/e/"},
		{"teams", nil, "/courses/missing-arg-0/teams/"},
		{"teams", []string{"c"}, "/courses/c/teams/"},
		{"teams", []string{"c", "e"}, "/courses/c/teams/e/"},
		{"users", nil, "/users/"},
		{"users", []string{"e"}, "/users/e/"},
	}
	for i, tc := range testcases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			args := make([]any, 0, len(tc.args))
			for _, arg := range tc.args {
				args = append(args, arg)
			}
			b := st.BuilderFor(tc.nodeID, args...)
			if got := b.String(); tc.exp != got {
				t.Errorf("Expected %q, but got %q", tc.exp, got)
			}
		})
	}
}
