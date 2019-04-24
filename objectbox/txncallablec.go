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

// This file implements externs defined in txncallable.go.
// It needs to be separate or it would cause duplicate symbol errors during linking.
// See https://golang.org/cmd/cgo/#hdr-C_references_to_Go for more details.

/*
#include <stdbool.h>
#include <stdint.h>
#include "objectbox.h"
*/
import "C"

//export txnCallableDispatch
// txnCallableDispatch is called from C.txn_callable_read|write
func txnCallableDispatch(id C.uint, cTx *C.OBX_txn) C.bool {
	var fn = txnCallableLookup(uint32(id))
	var tx = &Transaction{txn: cTx}
	return C.bool(fn(tx))
}
