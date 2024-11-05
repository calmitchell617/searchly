// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.27.3
// source: points_service.proto

package qdrant

import (
	context "context"
	grpc "google.golang.org/grpc"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	Points_Upsert_FullMethodName              = "/qdrant.Points/Upsert"
	Points_Delete_FullMethodName              = "/qdrant.Points/Delete"
	Points_Get_FullMethodName                 = "/qdrant.Points/Get"
	Points_UpdateVectors_FullMethodName       = "/qdrant.Points/UpdateVectors"
	Points_DeleteVectors_FullMethodName       = "/qdrant.Points/DeleteVectors"
	Points_SetPayload_FullMethodName          = "/qdrant.Points/SetPayload"
	Points_OverwritePayload_FullMethodName    = "/qdrant.Points/OverwritePayload"
	Points_DeletePayload_FullMethodName       = "/qdrant.Points/DeletePayload"
	Points_ClearPayload_FullMethodName        = "/qdrant.Points/ClearPayload"
	Points_CreateFieldIndex_FullMethodName    = "/qdrant.Points/CreateFieldIndex"
	Points_DeleteFieldIndex_FullMethodName    = "/qdrant.Points/DeleteFieldIndex"
	Points_Search_FullMethodName              = "/qdrant.Points/Search"
	Points_SearchBatch_FullMethodName         = "/qdrant.Points/SearchBatch"
	Points_SearchGroups_FullMethodName        = "/qdrant.Points/SearchGroups"
	Points_Scroll_FullMethodName              = "/qdrant.Points/Scroll"
	Points_Recommend_FullMethodName           = "/qdrant.Points/Recommend"
	Points_RecommendBatch_FullMethodName      = "/qdrant.Points/RecommendBatch"
	Points_RecommendGroups_FullMethodName     = "/qdrant.Points/RecommendGroups"
	Points_Discover_FullMethodName            = "/qdrant.Points/Discover"
	Points_DiscoverBatch_FullMethodName       = "/qdrant.Points/DiscoverBatch"
	Points_Count_FullMethodName               = "/qdrant.Points/Count"
	Points_UpdateBatch_FullMethodName         = "/qdrant.Points/UpdateBatch"
	Points_Query_FullMethodName               = "/qdrant.Points/Query"
	Points_QueryBatch_FullMethodName          = "/qdrant.Points/QueryBatch"
	Points_QueryGroups_FullMethodName         = "/qdrant.Points/QueryGroups"
	Points_Facet_FullMethodName               = "/qdrant.Points/Facet"
	Points_SearchMatrixPairs_FullMethodName   = "/qdrant.Points/SearchMatrixPairs"
	Points_SearchMatrixOffsets_FullMethodName = "/qdrant.Points/SearchMatrixOffsets"
)

// PointsClient is the client API for Points service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PointsClient interface {
	// Perform insert + updates on points. If a point with a given ID already exists - it will be overwritten.
	Upsert(ctx context.Context, in *UpsertPoints, opts ...grpc.CallOption) (*PointsOperationResponse, error)
	// Delete points
	Delete(ctx context.Context, in *DeletePoints, opts ...grpc.CallOption) (*PointsOperationResponse, error)
	// Retrieve points
	Get(ctx context.Context, in *GetPoints, opts ...grpc.CallOption) (*GetResponse, error)
	// Update named vectors for point
	UpdateVectors(ctx context.Context, in *UpdatePointVectors, opts ...grpc.CallOption) (*PointsOperationResponse, error)
	// Delete named vectors for points
	DeleteVectors(ctx context.Context, in *DeletePointVectors, opts ...grpc.CallOption) (*PointsOperationResponse, error)
	// Set payload for points
	SetPayload(ctx context.Context, in *SetPayloadPoints, opts ...grpc.CallOption) (*PointsOperationResponse, error)
	// Overwrite payload for points
	OverwritePayload(ctx context.Context, in *SetPayloadPoints, opts ...grpc.CallOption) (*PointsOperationResponse, error)
	// Delete specified key payload for points
	DeletePayload(ctx context.Context, in *DeletePayloadPoints, opts ...grpc.CallOption) (*PointsOperationResponse, error)
	// Remove all payload for specified points
	ClearPayload(ctx context.Context, in *ClearPayloadPoints, opts ...grpc.CallOption) (*PointsOperationResponse, error)
	// Create index for field in collection
	CreateFieldIndex(ctx context.Context, in *CreateFieldIndexCollection, opts ...grpc.CallOption) (*PointsOperationResponse, error)
	// Delete field index for collection
	DeleteFieldIndex(ctx context.Context, in *DeleteFieldIndexCollection, opts ...grpc.CallOption) (*PointsOperationResponse, error)
	// Retrieve closest points based on vector similarity and given filtering conditions
	Search(ctx context.Context, in *SearchPoints, opts ...grpc.CallOption) (*SearchResponse, error)
	// Retrieve closest points based on vector similarity and given filtering conditions
	SearchBatch(ctx context.Context, in *SearchBatchPoints, opts ...grpc.CallOption) (*SearchBatchResponse, error)
	// Retrieve closest points based on vector similarity and given filtering conditions, grouped by a given field
	SearchGroups(ctx context.Context, in *SearchPointGroups, opts ...grpc.CallOption) (*SearchGroupsResponse, error)
	// Iterate over all or filtered points
	Scroll(ctx context.Context, in *ScrollPoints, opts ...grpc.CallOption) (*ScrollResponse, error)
	// Look for the points which are closer to stored positive examples and at the same time further to negative examples.
	Recommend(ctx context.Context, in *RecommendPoints, opts ...grpc.CallOption) (*RecommendResponse, error)
	// Look for the points which are closer to stored positive examples and at the same time further to negative examples.
	RecommendBatch(ctx context.Context, in *RecommendBatchPoints, opts ...grpc.CallOption) (*RecommendBatchResponse, error)
	// Look for the points which are closer to stored positive examples and at the same time further to negative examples, grouped by a given field
	RecommendGroups(ctx context.Context, in *RecommendPointGroups, opts ...grpc.CallOption) (*RecommendGroupsResponse, error)
	// Use context and a target to find the most similar points to the target, constrained by the context.
	//
	// When using only the context (without a target), a special search - called context search - is performed where
	// pairs of points are used to generate a loss that guides the search towards the zone where
	// most positive examples overlap. This means that the score minimizes the scenario of
	// finding a point closer to a negative than to a positive part of a pair.
	//
	// Since the score of a context relates to loss, the maximum score a point can get is 0.0,
	// and it becomes normal that many points can have a score of 0.0.
	//
	// When using target (with or without context), the score behaves a little different: The
	// integer part of the score represents the rank with respect to the context, while the
	// decimal part of the score relates to the distance to the target. The context part of the score for
	// each pair is calculated +1 if the point is closer to a positive than to a negative part of a pair,
	// and -1 otherwise.
	Discover(ctx context.Context, in *DiscoverPoints, opts ...grpc.CallOption) (*DiscoverResponse, error)
	// Batch request points based on { positive, negative } pairs of examples, and/or a target
	DiscoverBatch(ctx context.Context, in *DiscoverBatchPoints, opts ...grpc.CallOption) (*DiscoverBatchResponse, error)
	// Count points in collection with given filtering conditions
	Count(ctx context.Context, in *CountPoints, opts ...grpc.CallOption) (*CountResponse, error)
	// Perform multiple update operations in one request
	UpdateBatch(ctx context.Context, in *UpdateBatchPoints, opts ...grpc.CallOption) (*UpdateBatchResponse, error)
	// Universally query points. This endpoint covers all capabilities of search, recommend, discover, filters. But also enables hybrid and multi-stage queries.
	Query(ctx context.Context, in *QueryPoints, opts ...grpc.CallOption) (*QueryResponse, error)
	// Universally query points in a batch fashion. This endpoint covers all capabilities of search, recommend, discover, filters. But also enables hybrid and multi-stage queries.
	QueryBatch(ctx context.Context, in *QueryBatchPoints, opts ...grpc.CallOption) (*QueryBatchResponse, error)
	// Universally query points in a group fashion. This endpoint covers all capabilities of search, recommend, discover, filters. But also enables hybrid and multi-stage queries.
	QueryGroups(ctx context.Context, in *QueryPointGroups, opts ...grpc.CallOption) (*QueryGroupsResponse, error)
	// Perform facet counts. For each value in the field, count the number of points that have this value and match the conditions.
	Facet(ctx context.Context, in *FacetCounts, opts ...grpc.CallOption) (*FacetResponse, error)
	// Compute distance matrix for sampled points with a pair based output format
	SearchMatrixPairs(ctx context.Context, in *SearchMatrixPoints, opts ...grpc.CallOption) (*SearchMatrixPairsResponse, error)
	// Compute distance matrix for sampled points with an offset based output format
	SearchMatrixOffsets(ctx context.Context, in *SearchMatrixPoints, opts ...grpc.CallOption) (*SearchMatrixOffsetsResponse, error)
}

type pointsClient struct {
	cc grpc.ClientConnInterface
}

func NewPointsClient(cc grpc.ClientConnInterface) PointsClient {
	return &pointsClient{cc}
}

func (c *pointsClient) Upsert(ctx context.Context, in *UpsertPoints, opts ...grpc.CallOption) (*PointsOperationResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(PointsOperationResponse)
	err := c.cc.Invoke(ctx, Points_Upsert_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pointsClient) Delete(ctx context.Context, in *DeletePoints, opts ...grpc.CallOption) (*PointsOperationResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(PointsOperationResponse)
	err := c.cc.Invoke(ctx, Points_Delete_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pointsClient) Get(ctx context.Context, in *GetPoints, opts ...grpc.CallOption) (*GetResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetResponse)
	err := c.cc.Invoke(ctx, Points_Get_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pointsClient) UpdateVectors(ctx context.Context, in *UpdatePointVectors, opts ...grpc.CallOption) (*PointsOperationResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(PointsOperationResponse)
	err := c.cc.Invoke(ctx, Points_UpdateVectors_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pointsClient) DeleteVectors(ctx context.Context, in *DeletePointVectors, opts ...grpc.CallOption) (*PointsOperationResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(PointsOperationResponse)
	err := c.cc.Invoke(ctx, Points_DeleteVectors_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pointsClient) SetPayload(ctx context.Context, in *SetPayloadPoints, opts ...grpc.CallOption) (*PointsOperationResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(PointsOperationResponse)
	err := c.cc.Invoke(ctx, Points_SetPayload_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pointsClient) OverwritePayload(ctx context.Context, in *SetPayloadPoints, opts ...grpc.CallOption) (*PointsOperationResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(PointsOperationResponse)
	err := c.cc.Invoke(ctx, Points_OverwritePayload_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pointsClient) DeletePayload(ctx context.Context, in *DeletePayloadPoints, opts ...grpc.CallOption) (*PointsOperationResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(PointsOperationResponse)
	err := c.cc.Invoke(ctx, Points_DeletePayload_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pointsClient) ClearPayload(ctx context.Context, in *ClearPayloadPoints, opts ...grpc.CallOption) (*PointsOperationResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(PointsOperationResponse)
	err := c.cc.Invoke(ctx, Points_ClearPayload_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pointsClient) CreateFieldIndex(ctx context.Context, in *CreateFieldIndexCollection, opts ...grpc.CallOption) (*PointsOperationResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(PointsOperationResponse)
	err := c.cc.Invoke(ctx, Points_CreateFieldIndex_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pointsClient) DeleteFieldIndex(ctx context.Context, in *DeleteFieldIndexCollection, opts ...grpc.CallOption) (*PointsOperationResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(PointsOperationResponse)
	err := c.cc.Invoke(ctx, Points_DeleteFieldIndex_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pointsClient) Search(ctx context.Context, in *SearchPoints, opts ...grpc.CallOption) (*SearchResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SearchResponse)
	err := c.cc.Invoke(ctx, Points_Search_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pointsClient) SearchBatch(ctx context.Context, in *SearchBatchPoints, opts ...grpc.CallOption) (*SearchBatchResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SearchBatchResponse)
	err := c.cc.Invoke(ctx, Points_SearchBatch_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pointsClient) SearchGroups(ctx context.Context, in *SearchPointGroups, opts ...grpc.CallOption) (*SearchGroupsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SearchGroupsResponse)
	err := c.cc.Invoke(ctx, Points_SearchGroups_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pointsClient) Scroll(ctx context.Context, in *ScrollPoints, opts ...grpc.CallOption) (*ScrollResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ScrollResponse)
	err := c.cc.Invoke(ctx, Points_Scroll_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pointsClient) Recommend(ctx context.Context, in *RecommendPoints, opts ...grpc.CallOption) (*RecommendResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RecommendResponse)
	err := c.cc.Invoke(ctx, Points_Recommend_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pointsClient) RecommendBatch(ctx context.Context, in *RecommendBatchPoints, opts ...grpc.CallOption) (*RecommendBatchResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RecommendBatchResponse)
	err := c.cc.Invoke(ctx, Points_RecommendBatch_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pointsClient) RecommendGroups(ctx context.Context, in *RecommendPointGroups, opts ...grpc.CallOption) (*RecommendGroupsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RecommendGroupsResponse)
	err := c.cc.Invoke(ctx, Points_RecommendGroups_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pointsClient) Discover(ctx context.Context, in *DiscoverPoints, opts ...grpc.CallOption) (*DiscoverResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DiscoverResponse)
	err := c.cc.Invoke(ctx, Points_Discover_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pointsClient) DiscoverBatch(ctx context.Context, in *DiscoverBatchPoints, opts ...grpc.CallOption) (*DiscoverBatchResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DiscoverBatchResponse)
	err := c.cc.Invoke(ctx, Points_DiscoverBatch_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pointsClient) Count(ctx context.Context, in *CountPoints, opts ...grpc.CallOption) (*CountResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CountResponse)
	err := c.cc.Invoke(ctx, Points_Count_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pointsClient) UpdateBatch(ctx context.Context, in *UpdateBatchPoints, opts ...grpc.CallOption) (*UpdateBatchResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UpdateBatchResponse)
	err := c.cc.Invoke(ctx, Points_UpdateBatch_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pointsClient) Query(ctx context.Context, in *QueryPoints, opts ...grpc.CallOption) (*QueryResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(QueryResponse)
	err := c.cc.Invoke(ctx, Points_Query_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pointsClient) QueryBatch(ctx context.Context, in *QueryBatchPoints, opts ...grpc.CallOption) (*QueryBatchResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(QueryBatchResponse)
	err := c.cc.Invoke(ctx, Points_QueryBatch_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pointsClient) QueryGroups(ctx context.Context, in *QueryPointGroups, opts ...grpc.CallOption) (*QueryGroupsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(QueryGroupsResponse)
	err := c.cc.Invoke(ctx, Points_QueryGroups_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pointsClient) Facet(ctx context.Context, in *FacetCounts, opts ...grpc.CallOption) (*FacetResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(FacetResponse)
	err := c.cc.Invoke(ctx, Points_Facet_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pointsClient) SearchMatrixPairs(ctx context.Context, in *SearchMatrixPoints, opts ...grpc.CallOption) (*SearchMatrixPairsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SearchMatrixPairsResponse)
	err := c.cc.Invoke(ctx, Points_SearchMatrixPairs_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pointsClient) SearchMatrixOffsets(ctx context.Context, in *SearchMatrixPoints, opts ...grpc.CallOption) (*SearchMatrixOffsetsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SearchMatrixOffsetsResponse)
	err := c.cc.Invoke(ctx, Points_SearchMatrixOffsets_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}