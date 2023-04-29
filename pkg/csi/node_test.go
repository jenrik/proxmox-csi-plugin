/*
Copyright 2023 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package csi_test

import (
	"context"
	"fmt"
	"testing"

	proto "github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/stretchr/testify/assert"

	"github.com/sergelogvinov/proxmox-csi-plugin/pkg/csi"
)

var _ proto.NodeServer = (*csi.NodeService)(nil)

type nodeServiceTestEnv struct {
	service *csi.NodeService
}

func newNodeServerTestEnv() nodeServiceTestEnv {
	return nodeServiceTestEnv{
		service: csi.NewNodeService("fake-proxmox-node", nil),
	}
}

func TestNodeStageVolumeErrors(t *testing.T) {
	t.Parallel()

	env := newNodeServerTestEnv()
	volcap := &proto.VolumeCapability{
		AccessMode: &proto.VolumeCapability_AccessMode{
			Mode: proto.VolumeCapability_AccessMode_SINGLE_NODE_WRITER,
		},
		AccessType: &proto.VolumeCapability_Mount{
			Mount: &proto.VolumeCapability_MountVolume{
				FsType: "ext4",
			},
		},
	}

	params := map[string]string{
		"DevicePath": "/dev/zero",
	}

	tests := []struct {
		msg           string
		request       *proto.NodeStageVolumeRequest
		expectedError error
	}{
		{
			msg: "VolumePath",
			request: &proto.NodeStageVolumeRequest{
				StagingTargetPath: "/staging",
				VolumeCapability:  volcap,
				PublishContext:    params,
			},
			expectedError: fmt.Errorf("VolumeID must be provided"),
		},
		{
			msg: "StagingTargetPath",
			request: &proto.NodeStageVolumeRequest{
				VolumeId:         "pvc-1",
				VolumeCapability: volcap,
				PublishContext:   params,
			},
			expectedError: fmt.Errorf("StagingTargetPath must be provided"),
		},
		{
			msg: "VolumeCapability",
			request: &proto.NodeStageVolumeRequest{
				VolumeId:          "pvc-1",
				StagingTargetPath: "/staging",
				PublishContext:    params,
			},
			expectedError: fmt.Errorf("VolumeCapability must be provided"),
		},
		{
			msg: "BlockVolume",
			request: &proto.NodeStageVolumeRequest{
				VolumeId:          "pvc-1",
				StagingTargetPath: "/staging",
				VolumeCapability: &proto.VolumeCapability{
					AccessMode: &proto.VolumeCapability_AccessMode{
						Mode: proto.VolumeCapability_AccessMode_SINGLE_NODE_WRITER,
					},
					AccessType: &proto.VolumeCapability_Block{
						Block: &proto.VolumeCapability_BlockVolume{},
					},
				},
				PublishContext: params,
			},
			expectedError: nil,
		},

		{
			msg: "DevicePath",
			request: &proto.NodeStageVolumeRequest{
				VolumeId:          "pvc-1",
				StagingTargetPath: "/staging",
				VolumeCapability:  volcap,
				PublishContext:    map[string]string{},
			},
			expectedError: fmt.Errorf("DevicePath must be provided"),
		},
	}

	for _, testCase := range tests {
		testCase := testCase

		t.Run(fmt.Sprint(testCase.msg), func(t *testing.T) {
			t.Parallel()

			resp, err := env.service.NodeStageVolume(context.Background(), testCase.request)

			if testCase.expectedError != nil {
				assert.NotNil(t, err)
				assert.Contains(t, err.Error(), testCase.expectedError.Error())
			} else {
				assert.Nil(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, *resp, proto.NodeStageVolumeResponse{})
			}
		})
	}
}

// nolint:dupl
func TestNodeUnstageVolumeErrors(t *testing.T) {
	t.Parallel()

	env := newNodeServerTestEnv()
	tests := []struct {
		msg           string
		request       *proto.NodeUnstageVolumeRequest
		expectedError error
	}{
		{
			msg: "VolumePath",
			request: &proto.NodeUnstageVolumeRequest{
				VolumeId: "pvc-1",
			},
			expectedError: fmt.Errorf("StagingTargetPath must be provided"),
		},
	}

	for _, testCase := range tests {
		testCase := testCase

		t.Run(fmt.Sprint(testCase.msg), func(t *testing.T) {
			t.Parallel()

			_, err := env.service.NodeUnstageVolume(context.Background(), testCase.request)

			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), testCase.expectedError.Error())
		})
	}
}

func TestNodeServiceNodePublishVolumeErrors(t *testing.T) {
	t.Parallel()

	env := newNodeServerTestEnv()
	volcap := &proto.VolumeCapability{
		AccessMode: &proto.VolumeCapability_AccessMode{
			Mode: proto.VolumeCapability_AccessMode_SINGLE_NODE_WRITER,
		},
		AccessType: &proto.VolumeCapability_Mount{
			Mount: &proto.VolumeCapability_MountVolume{
				FsType: "ext4",
			},
		},
	}

	params := map[string]string{
		"DevicePath": "/dev/zero",
	}

	tests := []struct {
		msg           string
		request       *proto.NodePublishVolumeRequest
		expectedError error
	}{
		{
			msg: "StagingTargetPath",
			request: &proto.NodePublishVolumeRequest{
				VolumeId:         "pvc-1",
				TargetPath:       "/target",
				VolumeCapability: volcap,
				PublishContext:   params,
			},
			expectedError: fmt.Errorf("StagingTargetPath must be provided"),
		},
		{
			msg: "TargetPath",
			request: &proto.NodePublishVolumeRequest{
				VolumeId:          "pvc-1",
				StagingTargetPath: "/staging",
				VolumeCapability:  volcap,
				PublishContext:    params,
			},
			expectedError: fmt.Errorf("TargetPath must be provided"),
		},
		{
			msg: "VolumeCapability",
			request: &proto.NodePublishVolumeRequest{
				VolumeId:          "pvc-1",
				StagingTargetPath: "/staging",
				TargetPath:        "/target",
				PublishContext:    params,
			},
			expectedError: fmt.Errorf("VolumeCapability must be provided"),
		},
		{
			msg: "VolumeCapability",
			request: &proto.NodePublishVolumeRequest{
				VolumeId:          "pvc-1",
				StagingTargetPath: "/staging",
				TargetPath:        "/target",
				VolumeCapability: &proto.VolumeCapability{
					AccessMode: &proto.VolumeCapability_AccessMode{
						Mode: proto.VolumeCapability_AccessMode_MULTI_NODE_MULTI_WRITER,
					},
				},
				PublishContext: params,
			},
			expectedError: fmt.Errorf("VolumeCapability not supported"),
		},
		{
			msg: "BlockVolume",
			request: &proto.NodePublishVolumeRequest{
				VolumeId:          "pvc-1",
				StagingTargetPath: "/staging",
				TargetPath:        "/target",
				VolumeCapability: &proto.VolumeCapability{
					AccessMode: &proto.VolumeCapability_AccessMode{
						Mode: proto.VolumeCapability_AccessMode_SINGLE_NODE_WRITER,
					},
					AccessType: &proto.VolumeCapability_Block{
						Block: &proto.VolumeCapability_BlockVolume{},
					},
				},
				PublishContext: params,
			},
			expectedError: fmt.Errorf("publish block volume is not supported"),
		},
		{
			msg: "VolumeCapability",
			request: &proto.NodePublishVolumeRequest{
				VolumeId:          "pvc-1",
				StagingTargetPath: "/staging",
				TargetPath:        "/target",
				VolumeCapability:  volcap,
				PublishContext:    map[string]string{},
			},
			expectedError: fmt.Errorf("DevicePath must be provided"),
		},
	}

	for _, testCase := range tests {
		testCase := testCase

		t.Run(fmt.Sprint(testCase.msg), func(t *testing.T) {
			t.Parallel()

			_, err := env.service.NodePublishVolume(context.Background(), testCase.request)

			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), testCase.expectedError.Error())
		})
	}
}

// nolint:dupl
func TestNodeUnpublishVolumeErrors(t *testing.T) {
	t.Parallel()

	env := newNodeServerTestEnv()
	tests := []struct {
		msg           string
		request       *proto.NodeUnpublishVolumeRequest
		expectedError error
	}{
		{
			msg:           "EmptyRequest",
			request:       &proto.NodeUnpublishVolumeRequest{},
			expectedError: fmt.Errorf("TargetPath must be provided"),
		},
		{
			msg: "TargetPath",
			request: &proto.NodeUnpublishVolumeRequest{
				VolumeId: "pvc-1",
			},
			expectedError: fmt.Errorf("TargetPath must be provided"),
		},
	}

	for _, testCase := range tests {
		testCase := testCase

		t.Run(fmt.Sprint(testCase.msg), func(t *testing.T) {
			t.Parallel()

			_, err := env.service.NodeUnpublishVolume(context.Background(), testCase.request)

			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), testCase.expectedError.Error())
		})
	}
}

// nolint:dupl
func TestNodeGetVolumeStatsErrors(t *testing.T) {
	t.Parallel()

	env := newNodeServerTestEnv()
	tests := []struct {
		msg           string
		request       *proto.NodeGetVolumeStatsRequest
		expectedError error
	}{
		{
			msg: "VolumePath",
			request: &proto.NodeGetVolumeStatsRequest{
				VolumeId: "pvc-1",
			},
			expectedError: fmt.Errorf("VolumePath must be provided"),
		},
	}

	for _, testCase := range tests {
		testCase := testCase

		t.Run(fmt.Sprint(testCase.msg), func(t *testing.T) {
			t.Parallel()

			_, err := env.service.NodeGetVolumeStats(context.Background(), testCase.request)

			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), testCase.expectedError.Error())
		})
	}
}

func TestNodeServiceNodeExpandVolumeErrors(t *testing.T) {
	t.Parallel()

	env := newNodeServerTestEnv()
	tests := []struct {
		msg           string
		request       *proto.NodeExpandVolumeRequest
		expectedError error
	}{
		{
			msg: "EmptyRequest",
			request: &proto.NodeExpandVolumeRequest{
				VolumePath: "/path",
			},
			expectedError: fmt.Errorf("VolumeID must be provided"),
		},
		{
			msg: "EmptyRequest",
			request: &proto.NodeExpandVolumeRequest{
				VolumeId: "pvc-1",
			},
			expectedError: fmt.Errorf("VolumePath must be provided"),
		},
	}

	for _, testCase := range tests {
		testCase := testCase

		t.Run(fmt.Sprint(testCase.msg), func(t *testing.T) {
			t.Parallel()

			_, err := env.service.NodeExpandVolume(context.Background(), testCase.request)

			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), testCase.expectedError.Error())
		})
	}
}

func TestNodeServiceNodeGetCapabilities(t *testing.T) {
	env := newNodeServerTestEnv()

	resp, err := env.service.NodeGetCapabilities(context.Background(), &proto.NodeGetCapabilitiesRequest{})
	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.GetCapabilities())

	for _, capability := range resp.GetCapabilities() {
		switch capability.GetRpc().Type { //nolint:exhaustive
		case proto.NodeServiceCapability_RPC_STAGE_UNSTAGE_VOLUME:
		case proto.NodeServiceCapability_RPC_EXPAND_VOLUME:
		case proto.NodeServiceCapability_RPC_GET_VOLUME_STATS:
		default:
			t.Fatalf("Unknown capability: %v", capability.Type)
		}
	}
}

// func TestNodeServiceNodeGetInfo(t *testing.T) {
// 	env := newNodeServerTestEnv()

// 	resp, err := env.service.NodeGetInfo(context.Background(), &proto.NodeGetInfoRequest{})
// 	assert.Nil(t, err)
// 	assert.NotNil(t, resp)

// 	assert.Equal(t, resp.NodeId, "fake-proxmox-node")
// 	assert.Equal(t, resp.MaxVolumesPerNode, csi.MaxVolumesPerNode)
// }