package integration_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/TeoPlow/online-music-service/src/musical/internal/models"
	pb "github.com/TeoPlow/online-music-service/src/musical/pkg/musicalpb"
	"github.com/TeoPlow/online-music-service/src/musical/tests/testutils"
)

type LikeServiceTest struct {
	testutils.BaseIntegrationSuite
}

func (ts *LikeServiceTest) SetupTest() {
	artistsData := []models.Artist{
		{ID: uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"),
			Name: "Test Artist 1"},
		{ID: uuid.MustParse("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"),
			Name: "Test Artist 2"},
	}

	albumsData := []models.Album{
		{
			ID:          uuid.MustParse("11111111-1111-1111-1111-111111111111"),
			Title:       "Test Album 1",
			ArtistID:    artistsData[0].ID,
			ReleaseDate: time.Now(),
		},
	}

	tracksData := []models.Track{
		{
			ID:         uuid.MustParse("7cb40172-a203-49e1-9c94-61f4d472b2a4"),
			Title:      "Test Track 1",
			AlbumID:    albumsData[0].ID,
			Genre:      "pop",
			Duration:   3*time.Minute + 32*time.Second,
			Lyrics:     testutils.ToPtr("Test lyrics 1"),
			IsExplicit: false,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			ID:         uuid.MustParse("1c2e394f-f036-4971-8b0a-4cbead2d80d4"),
			Title:      "Test Track 2",
			AlbumID:    albumsData[0].ID,
			Genre:      "rock",
			Duration:   3*time.Minute + 7*time.Second,
			IsExplicit: true,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
	}

	userID := uuid.MustParse("cccccccc-cccc-cccc-cccc-cccccccccccc")

	testutils.TruncateTables(ts.TxM.GetDatabase(),
		"liked_tracks",
		"liked_artists",
		"tracks",
		"albums",
		"artists")

	for _, artist := range artistsData {
		if err := ts.ArtistRepo.Add(ts.T().Context(), artist); err != nil {
			panic(err)
		}
	}

	for _, album := range albumsData {
		if err := ts.AlbumRepo.Add(ts.T().Context(), album); err != nil {
			panic(err)
		}
	}

	for _, track := range tracksData {
		if err := ts.TrackRepo.Add(ts.T().Context(), track); err != nil {
			panic(err)
		}
	}

	if err := ts.LikeRepo.LikeTrack(ts.T().Context(),
		userID,
		tracksData[0].ID,
	); err != nil {
		panic(err)
	}

	if err := ts.LikeRepo.LikeArtist(ts.T().Context(),
		userID,
		artistsData[0].ID,
	); err != nil {
		panic(err)
	}
}

func (ts *LikeServiceTest) TearDownTest() {
	testutils.TruncateTables(ts.TxM.GetDatabase(), "liked_tracks", "liked_artists")
}

func (ts *LikeServiceTest) TestLikeTrack() {
	userID := uuid.MustParse("cccccccc-cccc-cccc-cccc-cccccccccccc")

	tests := []struct {
		name     string
		req      *pb.LikeTrackRequest
		wantCode codes.Code
	}{
		{
			name: "success",
			req: &pb.LikeTrackRequest{
				TrackId: "1c2e394f-f036-4971-8b0a-4cbead2d80d4",
			},
			wantCode: codes.OK,
		},
		{
			name: "already liked",
			req: &pb.LikeTrackRequest{
				TrackId: "7cb40172-a203-49e1-9c94-61f4d472b2a4",
			},
			wantCode: codes.OK,
		},
		{
			name: "track not found",
			req: &pb.LikeTrackRequest{
				TrackId: "dddddddd-dddd-dddd-dddd-dddddddddddd",
			},
			wantCode: codes.Internal,
		},
		{
			name: "invalid track id",
			req: &pb.LikeTrackRequest{
				TrackId: "invalid-uuid",
			},
			wantCode: codes.InvalidArgument,
		},
	}

	for _, tt := range tests {
		ts.T().Run(tt.name, func(t *testing.T) {
			ctx := testutils.ContextWithUserID(ts.T().Context(), userID)

			res, err := ts.Client.LikeTrack(ctx, tt.req)
			if err != nil {
				st, ok := status.FromError(err)
				ts.Require().True(ok)
				ts.Assert().Equal(tt.wantCode, st.Code(), st.String())
				return
			}

			ts.Require().IsType(&emptypb.Empty{}, res)

			if tt.wantCode == codes.OK {
				trackID, _ := uuid.Parse(tt.req.TrackId)
				isLiked, err := ts.LikeRepo.IsTrackLiked(ctx, userID, trackID)
				ts.Require().NoError(err)
				ts.Assert().True(isLiked)
			}
		})
	}
}

func (ts *LikeServiceTest) TestUnlikeTrack() {
	userID := uuid.MustParse("cccccccc-cccc-cccc-cccc-cccccccccccc")

	tests := []struct {
		name     string
		req      *pb.UnlikeTrackRequest
		wantCode codes.Code
	}{
		{
			name: "success",
			req: &pb.UnlikeTrackRequest{
				TrackId: "7cb40172-a203-49e1-9c94-61f4d472b2a4",
			},
			wantCode: codes.OK,
		},
		{
			name: "not liked",
			req: &pb.UnlikeTrackRequest{
				TrackId: "1c2e394f-f036-4971-8b0a-4cbead2d80d4",
			},
			wantCode: codes.NotFound,
		},
		{
			name: "invalid track id",
			req: &pb.UnlikeTrackRequest{
				TrackId: "invalid-uuid",
			},
			wantCode: codes.InvalidArgument,
		},
	}

	for _, tt := range tests {
		ts.T().Run(tt.name, func(t *testing.T) {
			ctx := testutils.ContextWithUserID(ts.T().Context(), userID)

			res, err := ts.Client.UnlikeTrack(ctx, tt.req)
			if err != nil {
				st, ok := status.FromError(err)
				ts.Require().True(ok)
				ts.Assert().Equal(tt.wantCode, st.Code(), st.String())
				return
			}

			ts.Require().IsType(&emptypb.Empty{}, res)

			if tt.wantCode == codes.OK {
				trackID, _ := uuid.Parse(tt.req.TrackId)
				isLiked, err := ts.LikeRepo.IsTrackLiked(ctx, userID, trackID)
				ts.Require().NoError(err)
				ts.Assert().False(isLiked)
			}
		})
	}
}

func (ts *LikeServiceTest) TestLikeArtist() {
	userID := uuid.MustParse("cccccccc-cccc-cccc-cccc-cccccccccccc")

	tests := []struct {
		name     string
		req      *pb.LikeArtistRequest
		wantCode codes.Code
	}{
		{
			name: "success",
			req: &pb.LikeArtistRequest{
				ArtistId: "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb",
			},
			wantCode: codes.OK,
		},
		{
			name: "already liked",
			req: &pb.LikeArtistRequest{
				ArtistId: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
			},
			wantCode: codes.OK,
		},
		{
			name: "artist not found",
			req: &pb.LikeArtistRequest{
				ArtistId: "eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee",
			},
			wantCode: codes.Internal,
		},
		{
			name: "invalid artist id",
			req: &pb.LikeArtistRequest{
				ArtistId: "invalid-uuid",
			},
			wantCode: codes.InvalidArgument,
		},
	}

	for _, tt := range tests {
		ts.T().Run(tt.name, func(t *testing.T) {
			ctx := testutils.ContextWithUserID(ts.T().Context(), userID)

			res, err := ts.Client.LikeArtist(ctx, tt.req)
			if err != nil {
				st, ok := status.FromError(err)
				ts.Require().True(ok)
				ts.Assert().Equal(tt.wantCode, st.Code(), st.String())
				return
			}

			ts.Require().IsType(&emptypb.Empty{}, res)

			if tt.wantCode == codes.OK {
				artistID, _ := uuid.Parse(tt.req.ArtistId)
				isLiked, err := ts.LikeRepo.IsArtistLiked(ctx, userID, artistID)
				ts.Require().NoError(err)
				ts.Assert().True(isLiked)
			}
		})
	}
}

func (ts *LikeServiceTest) TestUnlikeArtist() {
	userID := uuid.MustParse("cccccccc-cccc-cccc-cccc-cccccccccccc")

	tests := []struct {
		name     string
		req      *pb.UnlikeArtistRequest
		wantCode codes.Code
	}{
		{
			name: "success",
			req: &pb.UnlikeArtistRequest{
				ArtistId: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
			},
			wantCode: codes.OK,
		},
		{
			name: "not liked",
			req: &pb.UnlikeArtistRequest{
				ArtistId: "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb",
			},
			wantCode: codes.NotFound,
		},
		{
			name: "invalid artist id",
			req: &pb.UnlikeArtistRequest{
				ArtistId: "invalid-uuid",
			},
			wantCode: codes.InvalidArgument,
		},
	}

	for _, tt := range tests {
		ts.T().Run(tt.name, func(t *testing.T) {
			ctx := testutils.ContextWithUserID(ts.T().Context(), userID)

			res, err := ts.Client.UnlikeArtist(ctx, tt.req)
			if err != nil {
				st, ok := status.FromError(err)
				ts.Require().True(ok)
				ts.Assert().Equal(tt.wantCode, st.Code(), st.String())
				return
			}

			ts.Require().IsType(&emptypb.Empty{}, res)

			if tt.wantCode == codes.OK {
				artistID, _ := uuid.Parse(tt.req.ArtistId)
				isLiked, err := ts.LikeRepo.IsArtistLiked(ctx, userID, artistID)
				ts.Require().NoError(err)
				ts.Assert().False(isLiked)
			}
		})
	}
}

func (ts *LikeServiceTest) TestGetLikedTracks() {
	userID := uuid.MustParse("cccccccc-cccc-cccc-cccc-cccccccccccc")

	tests := []struct {
		name      string
		req       *pb.GetLikedTracksRequest
		userID    uuid.UUID
		wantCode  codes.Code
		wantCount int
	}{
		{
			name: "success",
			req: &pb.GetLikedTracksRequest{
				Page:     1,
				PageSize: 10,
			},
			userID:    userID,
			wantCode:  codes.OK,
			wantCount: 1,
		},
		{
			name: "empty result",
			req: &pb.GetLikedTracksRequest{
				Page:     1,
				PageSize: 10,
			},
			userID:    uuid.New(),
			wantCode:  codes.OK,
			wantCount: 0,
		},
		{
			name: "invalid parameters",
			req: &pb.GetLikedTracksRequest{
				Page:     0,
				PageSize: 0,
			},
			userID:    userID,
			wantCode:  codes.InvalidArgument,
			wantCount: 0,
		},
	}

	for _, tt := range tests {
		ts.T().Run(tt.name, func(t *testing.T) {
			ctx := testutils.ContextWithUserID(ts.T().Context(), tt.userID)

			res, err := ts.Client.GetLikedTracks(ctx, tt.req)
			if err != nil {
				st, ok := status.FromError(err)
				ts.Require().True(ok)
				ts.Assert().Equal(tt.wantCode, st.Code(), st.String())
				return
			}

			ts.Require().NoError(err)
			ts.Assert().Equal(tt.wantCount, len(res.GetTracks()))

			if tt.wantCount > 0 {
				track := res.GetTracks()[0]
				ts.Assert().Equal("7cb40172-a203-49e1-9c94-61f4d472b2a4", track.GetId())
				ts.Assert().Equal("Test Track 1", track.GetTitle())
				ts.Assert().Equal("11111111-1111-1111-1111-111111111111", track.GetAlbumId())
				ts.Assert().Equal("pop", track.GetGenre())
				ts.Assert().Equal(int32(3*60+32), track.GetDuration())
				ts.Assert().Equal("Test lyrics 1", track.GetLyrics())
				ts.Assert().Equal(false, track.GetIsExplicit())
			}
		})
	}
}

func (ts *LikeServiceTest) TestGetLikedArtists() {
	userID := uuid.MustParse("cccccccc-cccc-cccc-cccc-cccccccccccc")

	tests := []struct {
		name      string
		req       *pb.GetLikedArtistsRequest
		userID    uuid.UUID
		wantCode  codes.Code
		wantCount int
	}{
		{
			name: "success",
			req: &pb.GetLikedArtistsRequest{
				Page:     1,
				PageSize: 10,
			},
			userID:    userID,
			wantCode:  codes.OK,
			wantCount: 1,
		},
		{
			name: "empty result",
			req: &pb.GetLikedArtistsRequest{
				Page:     1,
				PageSize: 10,
			},
			userID:    uuid.New(),
			wantCode:  codes.OK,
			wantCount: 0,
		},
		{
			name: "invalid parameters",
			req: &pb.GetLikedArtistsRequest{
				Page:     0,
				PageSize: 0,
			},
			userID:    userID,
			wantCode:  codes.InvalidArgument,
			wantCount: 0,
		},
	}

	for _, tt := range tests {
		ts.T().Run(tt.name, func(t *testing.T) {
			ctx := testutils.ContextWithUserID(ts.T().Context(), tt.userID)

			res, err := ts.Client.GetLikedArtists(ctx, tt.req)
			if err != nil {
				st, ok := status.FromError(err)
				ts.Require().True(ok)
				ts.Assert().Equal(tt.wantCode, st.Code(), st.String())
				return
			}

			ts.Require().NoError(err)
			ts.Assert().Equal(tt.wantCount, len(res.GetArtists()))

			if tt.wantCount > 0 {
				artist := res.GetArtists()[0]
				ts.Assert().Equal("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa", artist.GetId())
				ts.Assert().Equal("Test Artist 1", artist.GetName())
			}
		})
	}
}

func TestLikeServiceGRPC(t *testing.T) {
	suite.Run(t, new(LikeServiceTest))
}
