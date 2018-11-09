package modelinfo

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
)

type IdUid string

func CreateIdUid(id id, uid uid) IdUid {
	return IdUid(strconv.FormatInt(int64(id), 10) + ":" + strconv.FormatInt(int64(uid), 10))
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

func (str IdUid) GetId() (id, error) {
	if i, err := str.getComponent(0, 32); err != nil {
		return 0, err
	} else {
		return id(i), nil
	}
}

func (str *IdUid) GetUid() (uid, error) {
	return str.getComponent(1, 64)
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

func (str IdUid) getIdSafe() id {
	i, _ := str.GetId()
	return i
}

func (str IdUid) getUidSafe() uid {
	i, _ := str.GetUid()
	return i
}

func generateUid(isUnique func(uid) bool) (result uid, err error) {
	result = 0

	for i := 0; i < 1000; i++ {
		t := uid(rand.Int63())
		if isUnique(t) {
			result = t
			break
		}
	}

	if result == 0 {
		err = fmt.Errorf("internal error = could not generate a unique UID")
	}

	return result, err
}
