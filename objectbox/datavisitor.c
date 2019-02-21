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

// implements data visitor forwarding, it's called from ObjectBox C-api (passed to it ac C.data_visitor
// it's necessary to keep in the separate file to avoid "multiple-definitions" error
bool data_visitor_(void* arg, const void* data, size_t size) {
    // forward the call to the Go function, which finds the correct Go callback and calls it
    // *arg is a pointer to the dataVisitorId associated with the Go callback
    // we need to get rid of the `const` because Go doesn't support that
    return dataVisitorDispatch(*((uint32_t*)arg), (void*) data, size);
}

// C.data_visitor is used as an argument, e. g. for C.obx_query_visit(), together with a pointer to "visitorId"
obx_data_visitor* data_visitor = &data_visitor_;
