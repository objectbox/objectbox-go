/*
 * Copyright 2019 ObjectBox Ltd. All rights reserved.
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

package objectbox

/*
#include <stdlib.h>
#include "objectbox.h"
*/
import "C"

type transaction struct {
	cTxn      *C.OBX_txn
	objectBox *ObjectBox
}

func (txn *transaction) Close() error {
	rc := C.obx_txn_close(txn.cTxn)
	txn.cTxn = nil
	if rc != 0 {
		return createError()
	}
	return nil
}

func (txn *transaction) Abort() error {
	rc := C.obx_txn_abort(txn.cTxn)
	if rc != 0 {
		return createError()
	}
	return nil
}

func (txn *transaction) Commit() error {
	rc := C.obx_txn_commit(txn.cTxn)
	if rc != 0 {
		return createError()
	}
	return nil
}
