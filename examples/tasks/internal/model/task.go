/*
 * Copyright 2018-2022 ObjectBox Ltd. All rights reserved.
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

package model

import (
	"time"
)

//go:generate go run github.com/objectbox/objectbox-go/cmd/objectbox-gogen

// Put this on a new line to enable sync: // `objectbox:"sync"`
type Task struct {
	Id          uint64
	Text        string
	DateCreated time.Time `objectbox:"date"`

	// DateFinished is initially set to unix epoch (value 0 in ObjectBox DB) to tag the task as "unfinished"
	DateFinished time.Time `objectbox:"date"`
}
