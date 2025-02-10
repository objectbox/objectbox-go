/*
 * Copyright 2018-2021 ObjectBox Ltd. All rights reserved.
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
#include "objectbox-sync.h"
*/
import "C"

// SyncIsAvailable returns true if the loaded ObjectBox native library supports Sync.
// [ObjectBox Sync](https://objectbox.io/sync/) makes data available and synchronized across devices, online and offline.
func SyncIsAvailable() bool {
	return bool(C.obx_has_feature(C.OBXFeature_Sync))
}

// SyncCredentials are used to authenticate a sync client against a server.
type SyncCredentials struct {
	cType C.OBXSyncCredentialsType
	data  []byte
	data2 []byte // used for passing username-password credentials
}

// SyncCredentialsNone - no credentials - usually only for development, with a server configured to accept all
// connections without authentication.
func SyncCredentialsNone() *SyncCredentials {
	return &SyncCredentials{
		cType: C.OBXSyncCredentialsType_NONE,
		data:  nil,
	}
}

// SyncCredentialsSharedSecret - shared secret authentication
func SyncCredentialsSharedSecret(data []byte) *SyncCredentials {
	return &SyncCredentials{
		cType: C.OBXSyncCredentialsType_SHARED_SECRET,
		data:  data,
	}
}

// SyncCredentialsGoogleAuth - Google authentication
func SyncCredentialsGoogleAuth(data []byte) *SyncCredentials {
	return &SyncCredentials{
		cType: C.OBXSyncCredentialsType_GOOGLE_AUTH,
		data:  data,
	}
}

// SyncCredentialsGoogleAuth - Google authentication
func SyncCredentialsUsernamePassword(username []byte, password []byte) *SyncCredentials {
	return &SyncCredentials{
		cType: C.OBXSyncCredentialsType_USER_PASSWORD,
		data:  username,
		data2: password,
	}
}

// SyncCredentialsJwtId - JWT authentication with an ID token
func SyncCredentialsJwtId(data []byte) *SyncCredentials {
	return &SyncCredentials{
		cType: C.OBXSyncCredentialsType_JWT_ID,
		data:  data,
	}
}

// SyncCredentialsJwtAccess - JWT authentication with an access token
func SyncCredentialsJwtAccess(data []byte) *SyncCredentials {
	return &SyncCredentials{
		cType: C.OBXSyncCredentialsType_JWT_ACCESS,
		data:  data,
	}
}

// SyncCredentialsJwtRefresh - JWT authentication with a refresh token
func SyncCredentialsJwtRefresh(data []byte) *SyncCredentials {
	return &SyncCredentials{
		cType: C.OBXSyncCredentialsType_JWT_REFRESH,
		data:  data,
	}
}

// SyncCredentialsJwtCustom - JWT authentication with a custom token
func SyncCredentialsJwtCustom(data []byte) *SyncCredentials {
	return &SyncCredentials{
		cType: C.OBXSyncCredentialsType_JWT_CUSTOM,
		data:  data,
	}
}
