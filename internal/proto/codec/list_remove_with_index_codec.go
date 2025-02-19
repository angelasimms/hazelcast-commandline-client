/*
* Copyright (c) 2008-2023, Hazelcast, Inc. All Rights Reserved.
*
* Licensed under the Apache License, Version 2.0 (the "License")
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
* http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
 */

package codec

import (
	iserialization "github.com/hazelcast/hazelcast-go-client"
	proto "github.com/hazelcast/hazelcast-go-client"
)

const (
	ListRemoveWithIndexCodecRequestMessageType  = int32(0x051200)
	ListRemoveWithIndexCodecResponseMessageType = int32(0x051201)

	ListRemoveWithIndexCodecRequestIndexOffset      = proto.PartitionIDOffset + proto.IntSizeInBytes
	ListRemoveWithIndexCodecRequestInitialFrameSize = ListRemoveWithIndexCodecRequestIndexOffset + proto.IntSizeInBytes
)

// Removes the element at the specified position in this list (optional operation). Shifts any subsequent elements
// to the left (subtracts one from their indices). Returns the element that was removed from the list.

func EncodeListRemoveWithIndexRequest(name string, index int32) *proto.ClientMessage {
	clientMessage := proto.NewClientMessageForEncode()
	clientMessage.SetRetryable(false)

	initialFrame := proto.NewFrameWith(make([]byte, ListRemoveWithIndexCodecRequestInitialFrameSize), proto.UnfragmentedMessage)
	EncodeInt(initialFrame.Content, ListRemoveWithIndexCodecRequestIndexOffset, index)
	clientMessage.AddFrame(initialFrame)
	clientMessage.SetMessageType(ListRemoveWithIndexCodecRequestMessageType)
	clientMessage.SetPartitionId(-1)

	EncodeString(clientMessage, name)

	return clientMessage
}

func DecodeListRemoveWithIndexResponse(clientMessage *proto.ClientMessage) iserialization.Data {
	frameIterator := clientMessage.FrameIterator()
	// empty initial frame
	frameIterator.Next()

	return DecodeNullableForData(frameIterator)
}
