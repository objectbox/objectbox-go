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
#cgo LDFLAGS: -lobjectbox
#include <stdlib.h>
#include "objectbox.h"
*/
import "C"

// Internal: won't be publicly exposed in a future version!
type Transaction struct {
	txn       *C.OBX_txn
	objectBox *ObjectBox
}

func (txn *Transaction) Close() error {
	rc := C.obx_txn_close(txn.txn)
	txn.txn = nil
	if rc != 0 {
		return createError()
	}
	return nil
}

func (txn *Transaction) Abort() error {
	rc := C.obx_txn_abort(txn.txn)
	if rc != 0 {
		return createError()
	}
	return nil
}

func (txn *Transaction) Commit() error {
	rc := C.obx_txn_commit(txn.txn)
	if rc != 0 {
		return createError()
	}
	return nil
}
