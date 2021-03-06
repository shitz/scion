// Copyright 2020 ETH Zurich, Anapaya Systems
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package colibri_mgmt

import (
	"github.com/scionproto/scion/go/proto"
)

type Response struct {
	Which                    proto.Response_Which
	SegmentSetup             *SegmentSetupRes
	SegmentRenewal           *SegmentSetupRes
	SegmentTeardown          *SegmentTeardownRes
	SegmentIndexConfirmation *SegmentIndexConfirmationRes
	SegmentCleanup           *SegmentCleanupRes
	E2ESetup                 *E2ESetupRes   `capnp:"e2eSetup"`
	E2ERenewal               *E2ESetupRes   `capnp:"e2eRenewal"`
	E2ECleanup               *E2ECleanupRes `capnp:"e2eCleanup"`
	Accepted                 bool
	FailedHop                uint8 // only relevant if Accepted == false
}

func (r *Response) ProtoId() proto.ProtoIdType {
	return proto.Response_TypeID
}
