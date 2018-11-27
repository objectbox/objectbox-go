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

package modelinfo

import (
	"fmt"
	"strconv"
	"strings"
)

type IdUid string

func CreateIdUid(id Id, uid Uid) IdUid {
	return IdUid(strconv.FormatInt(int64(id), 10) + ":" + strconv.FormatUint(uid, 10))
}

// performs initial validation of loaded data so that it doesn't have to be checked in each function
func (str *IdUid) Validate() error {
	if _, err := str.GetUid(); err != nil {
		return fmt.Errorf("uid: %s", err)
	}

	if _, err := str.GetId(); err != nil {
		return fmt.Errorf("id: %s", err)
	}

	if len(strings.Split(string(*str), ":")) != 2 {
		return fmt.Errorf("id invalid format - too many colons")
	}

	return nil
}

func (str IdUid) GetId() (Id, error) {
	if i, err := str.getComponent(0, 32); err != nil {
		return 0, err
	} else {
		return Id(i), nil
	}
}

func (str *IdUid) GetUid() (Uid, error) {
	return str.getComponent(1, 64)
}

func (str *IdUid) Get() (Id, Uid, error) {
	if id, err := str.GetId(); err != nil {
		return 0, 0, err
	} else if uid, err := str.GetUid(); err != nil {
		return 0, 0, err
	} else {
		return id, uid, nil
	}
}

func (str IdUid) getComponent(n, bitsize int) (uint64, error) {
	if len(str) == 0 {
		return 0, fmt.Errorf("is undefined")
	}

	idStr := strings.Split(string(str), ":")[n]
	if component, err := strconv.ParseUint(idStr, 10, bitsize); err != nil {
		return 0, fmt.Errorf("can't parse '%s' as unsigned int: %s", idStr, err)
	} else if component == 0 {
		return 0, fmt.Errorf("equals to zero")
	} else {
		return component, nil
	}
}

func (str IdUid) getIdSafe() Id {
	i, _ := str.GetId()
	return i
}

func (str IdUid) getUidSafe() Uid {
	i, _ := str.GetUid()
	return i
}
