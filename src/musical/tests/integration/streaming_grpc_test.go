package integration_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"

	pb "github.com/TeoPlow/online-music-service/src/musical/pkg/musicalpb"
	"github.com/TeoPlow/online-music-service/src/musical/tests/testutils"
)

type StreamingServiceTest struct {
	testutils.BaseIntegrationSuite
}

func (ts *StreamingServiceTest) SetupTest() {
	ts.MinIOClient.Truncate(ts.T().Context())

	track := testutils.TrackFromFile("meow-meow")
	err := ts.MinIOClient.Upload(ts.T().Context(), "97e6da54-54d6-4e74-890c-970a47543554",
		track, int64(track.Len()))
	if err != nil {
		panic(err)
	}
}

func (ts *StreamingServiceTest) TearDownTest() {
	ts.MinIOClient.Truncate(ts.T().Context())
}

func (ts *StreamingServiceTest) TestDownloadTrack() {
	ctx := ts.T().Context()

	ctx = testutils.ContextWithUserID(ctx, uuid.New())
	original := testutils.TrackFromFile("meow-meow")
	id := "97e6da54-54d6-4e74-890c-970a47543554"

	stream, err := ts.Client.DownloadTrack(ctx, &pb.IDRequest{Id: id})
	ts.Require().NoError(err)

	var downloaded bytes.Buffer
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		ts.Require().NoError(err)
		_, err = downloaded.Write(resp.GetChunk())
		ts.Require().NoError(err)
	}

	ts.Assert().True(testutils.CompareTracks(original.Bytes(), downloaded.Bytes()),
		"Downloaded track does not match original")
}

func TestStreamingServiceGRPC(t *testing.T) {
	suite.Run(t, new(StreamingServiceTest))
}
