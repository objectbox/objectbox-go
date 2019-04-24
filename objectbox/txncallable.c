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

#include "_cgo_export.h"
#include "objectbox.h"

// implements txn_callable forwarding, it's called from ObjectBox C-api (passed as the C.txn_callable_* defined bellow)
// it's necessary to keep in the separate file to avoid "multiple-definitions" error
// It forward the call to the Go function, which finds the correct Go callback and calls it
// *arg is a pointer to the visitorId associated with the Go callback

void txn_callable_read_(void* arg, OBX_txn* txn) {
    txnCallableDispatch(*((uint32_t*)arg), txn);
}

bool txn_callable_write_(void* arg, OBX_txn* txn) {
    return txnCallableDispatch(*((uint32_t*)arg), txn);
}

// C.txn_callable_read|write is used as an argument for C.obx_store_exec_read|write(), passing a pointer to "visitorId"
obx_txn_callable_read*  txn_callable_read  = &txn_callable_read_;
obx_txn_callable_write* txn_callable_write = &txn_callable_write_;
