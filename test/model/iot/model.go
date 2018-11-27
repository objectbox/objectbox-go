/*
 * Copyright 2018 ObjectBox Ltd. All rights reserved.
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

package iot

//go:generate objectbox-gogen

type Event struct {
	Id     uint64 `id`
	Uid    string `unique`
	Device string
	Date   int64 `date`
}

type Reading struct {
	Id   uint64 `id`
	Date int64  `date`

	/// to-one relation
	EventId uint64

	ValueName string

	/// Device sensor data value
	ValueString string

	/// Device sensor data value
	ValueInteger int64

	/// Device sensor data value
	ValueFloating float64
}
