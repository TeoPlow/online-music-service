package integration_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/TeoPlow/online-music-service/src/musical/internal/models"
	pb "github.com/TeoPlow/online-music-service/src/musical/pkg/musicalpb"
	"github.com/TeoPlow/online-music-service/src/musical/tests/testutils"
)

type TrackServiceTest struct {
	testutils.BaseIntegrationSuite
}

func (ts *TrackServiceTest) SetupTest() {
	artistsData := []models.Artist{
		{ID: uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")},
	}

	albumsData := []models.Album{
		{ID: uuid.MustParse("11111111-1111-1111-1111-111111111111"), ArtistID: artistsData[0].ID},
		{ID: uuid.MustParse("22222222-2222-2222-2222-222222222222"), ArtistID: artistsData[0].ID},
		{ID: uuid.MustParse("33333333-3333-3333-3333-333333333333"), ArtistID: artistsData[0].ID},
	}

	// Данные треков
	tracksData := []models.Track{
		{
			ID:         uuid.MustParse("7cb40172-a203-49e1-9c94-61f4d472b2a4"),
			Title:      "Never Gonna Give You Up",
			AlbumID:    albumsData[0].ID,
			Genre:      "pop",
			Duration:   3*time.Minute + 32*time.Second,
			Lyrics:     testutils.ToPtr("We're no strangers to love"),
			IsExplicit: false,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			ID:         uuid.MustParse("1c2e394f-f036-4971-8b0a-4cbead2d80d4"),
			Title:      "Imagine",
			AlbumID:    albumsData[0].ID,
			Genre:      "rock",
			Duration:   3*time.Minute + 7*time.Second,
			Lyrics:     testutils.ToPtr("Imagine there's no heaven"),
			IsExplicit: false,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			ID:         uuid.MustParse("999873a9-5589-467c-b3cd-b2fa42b87fb6"),
			Title:      "Yellow Submarine",
			AlbumID:    albumsData[1].ID,
			Genre:      "rock",
			Duration:   2*time.Minute + 45*time.Second,
			Lyrics:     testutils.ToPtr("In the town where I was born"),
			IsExplicit: false,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			ID:         uuid.MustParse("896150e2-de0a-49fd-9abf-3b4fd159fb16"),
			Title:      "Smells Like Teen Spirit",
			AlbumID:    albumsData[2].ID,
			Genre:      "grunge",
			Duration:   5*time.Minute + 1*time.Second,
			Lyrics:     nil,
			IsExplicit: true,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			ID:         uuid.MustParse("e189c67c-4864-485b-9fcc-75f99bf35a10"),
			Title:      "Cat words",
			AlbumID:    albumsData[2].ID,
			Genre:      "rock",
			Duration:   6*time.Minute + 3*time.Second,
			Lyrics:     testutils.ToPtr("Meow, meow, meow..."),
			IsExplicit: false,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
	}

	testutils.TruncateTables(ts.TxM.GetDatabase(), "tracks", "albums", "artists")

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
}

func (ts *TrackServiceTest) TearDownTest() {
	testutils.TruncateTables(ts.TxM.GetDatabase(), "tracks", "albums")
}

func (ts *TrackServiceTest) TestCreateTrack() {
	tests := []struct {
		name     string
		req      *pb.TrackInfo
		wantCode codes.Code
	}{
		{
			name: "success",
			req: &pb.TrackInfo{
				Title:   "test_track 6",
				AlbumId: "11111111-1111-1111-1111-111111111111",
			},
			wantCode: codes.OK,
		},
		{
			name: "already exists",
			req: &pb.TrackInfo{
				Title:      "Never Gonna Give You Up",
				AlbumId:    "11111111-1111-1111-1111-111111111111",
				Genre:      "pop",
				Duration:   int32(time.Duration(3*time.Minute + 32*time.Second).Seconds()),
				Lyrics:     testutils.ToPtr("We're no strangers to love"),
				IsExplicit: false,
			},
			wantCode: codes.AlreadyExists,
		},
		{
			name: "no album",
			req: &pb.TrackInfo{
				Title:   "test_track 7",
				AlbumId: "30a1415e-1287-44bd-b856-78bb21b1375c",
				Genre:   "classic",
			},
			wantCode: codes.NotFound,
		},
		{
			name: "invalid uuid",
			req: &pb.TrackInfo{
				Title:   "test_ablum 8",
				AlbumId: "i'm uuid",
				Genre:   "classic",
			},
			wantCode: codes.InvalidArgument,
		},
	}

	for _, tt := range tests {
		ts.T().Run(tt.name, func(t *testing.T) {
			stream, err := ts.Client.CreateTrack(ts.T().Context())
			ts.Require().NoError(err)

			err = stream.Send(&pb.CreateTrackRequest{
				Data: &pb.CreateTrackRequest_Metadata{
					Metadata: tt.req,
				},
			})
			ts.Require().NoError(err)

			res, err := stream.CloseAndRecv()
			if err != nil {
				st, ok := status.FromError(err)
				ts.Require().True(ok)
				ts.Assert().Equal(tt.wantCode, st.Code(), st.String())
				return
			}

			actual, err := ts.TrackRepo.GetByID(ts.T().Context(), uuid.MustParse(res.GetId()))
			ts.Require().NoError(err)

			ts.Assert().Equal(res.GetId(), actual.ID.String())
			ts.Assert().Equal(res.GetTitle(), actual.Title)
			ts.Assert().Equal(res.GetAlbumId(), actual.AlbumID.String())
			ts.Assert().Equal(res.GetDuration(), int32(actual.Duration.Seconds()))
			ts.Assert().Equal(res.GetGenre(), actual.Genre)
			ts.Assert().Equal(res.GetIsExplicit(), actual.IsExplicit)
			if res.Lyrics != nil {
				ts.Assert().Equal(res.GetLyrics(), *actual.Lyrics)
			} else {
				ts.Assert().Nil(actual.Lyrics)
			}
			ts.Assert().Equal(testutils.DateOnly(res.GetCreatedAt().AsTime()),
				testutils.DateOnly(actual.CreatedAt))
			ts.Assert().Equal(testutils.DateOnly(res.GetUpdatedAt().AsTime()),
				testutils.DateOnly(actual.UpdatedAt))

			ts.Assert().Equal(tt.req.GetTitle(), res.GetTitle())
			ts.Assert().Equal(tt.req.GetAlbumId(), res.GetAlbumId())
			ts.Assert().Equal(tt.req.GetGenre(), res.GetGenre())
			ts.Assert().Equal(tt.req.GetDuration(), res.GetDuration())
			ts.Assert().Equal(tt.req.GetLyrics(), res.GetLyrics())
			ts.Assert().Equal(tt.req.GetIsExplicit(), res.GetIsExplicit())
		})
	}
}

func (ts *TrackServiceTest) TestGetTrack() {
	tests := []struct {
		name     string
		req      *pb.IDRequest
		wantCode codes.Code
	}{
		{
			name: "success",
			req: &pb.IDRequest{
				Id: "e189c67c-4864-485b-9fcc-75f99bf35a10",
			},
			wantCode: codes.OK,
		},
		{
			name: "not found",
			req: &pb.IDRequest{
				Id: "6ca52977-b1e5-4e30-b2ae-874616766c7d",
			},
			wantCode: codes.NotFound,
		},
		{
			name: "invalid uuid",
			req: &pb.IDRequest{
				Id: "I'm uuid",
			},
			wantCode: codes.InvalidArgument,
		},
	}

	for _, tt := range tests {
		ts.T().Run(tt.name, func(t *testing.T) {
			res, err := ts.Client.GetTrack(ts.T().Context(), tt.req)
			if err != nil {
				st, ok := status.FromError(err)
				ts.Require().True(ok)
				ts.Assert().Equal(tt.wantCode, st.Code(), st.String())
				return
			}

			expect := models.Track{
				ID:         uuid.MustParse("e189c67c-4864-485b-9fcc-75f99bf35a10"),
				Title:      "Cat words",
				AlbumID:    uuid.MustParse("33333333-3333-3333-3333-333333333333"),
				Genre:      "rock",
				Duration:   6*time.Minute + 3*time.Second,
				Lyrics:     testutils.ToPtr("Meow, meow, meow..."),
				IsExplicit: false,
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			}
			ts.Require().NoError(err)

			ts.Assert().Equal(tt.req.GetId(), res.GetId())

			ts.Assert().Equal(expect.ID.String(), res.GetId())
			ts.Assert().Equal(expect.Title, res.GetTitle())
			ts.Assert().Equal(expect.AlbumID.String(), res.GetAlbumId())
			ts.Assert().Equal(expect.Genre, res.Genre)
			ts.Assert().EqualValues(expect.Duration.Seconds(), res.GetDuration())
			ts.Assert().Equal(*expect.Lyrics, res.GetLyrics())
			ts.Assert().Equal(expect.IsExplicit, res.IsExplicit)
			ts.Assert().Equal(testutils.DateOnly(expect.UpdatedAt),
				testutils.DateOnly(res.GetUpdatedAt().AsTime()))
			ts.Assert().Equal(testutils.DateOnly(expect.CreatedAt),
				testutils.DateOnly(res.GetCreatedAt().AsTime()))
		})
	}
}

func (ts *TrackServiceTest) TestListTrack() {
	tests := []struct {
		name     string
		req      *pb.ListTracksRequest
		wantCode codes.Code
	}{
		{
			name: "full query",
			req: &pb.ListTracksRequest{
				Page:        1,
				PageSize:    2,
				Genre:       testutils.ToPtr("rock"),
				AlbumId:     testutils.ToPtr("33333333-3333-3333-3333-333333333333"),
				SearchQuery: testutils.ToPtr("Cat"),
			},
			wantCode: codes.OK,
		},
		{
			name: "only album",
			req: &pb.ListTracksRequest{
				Page:     1,
				PageSize: 2,
				AlbumId:  testutils.ToPtr("11111111-1111-1111-1111-111111111111"),
			},
			wantCode: codes.OK,
		},
		{
			name: "only search",
			req: &pb.ListTracksRequest{
				Page:        1,
				PageSize:    2,
				SearchQuery: testutils.ToPtr("Spirit"),
			},
			wantCode: codes.OK,
		},
		{
			name: "invalid uuid",
			req: &pb.ListTracksRequest{
				Page:        1,
				PageSize:    2,
				AlbumId:     testutils.ToPtr("I'm uuid"),
				SearchQuery: testutils.ToPtr("Cat"),
			},
			wantCode: codes.InvalidArgument,
		},
		{
			name: "empty result",
			req: &pb.ListTracksRequest{
				Page:     1,
				PageSize: 2,
				AlbumId:  testutils.ToPtr("a592ffc7-fccb-4e29-8b90-647b80cbb758"),
			},
			wantCode: codes.OK,
		},
	}

	for _, tt := range tests {
		ts.T().Run(tt.name, func(t *testing.T) {
			res, err := ts.Client.ListTracks(ts.T().Context(), tt.req)
			if err != nil {
				st, ok := status.FromError(err)
				ts.Require().True(ok)
				ts.Assert().Equal(tt.wantCode, st.Code(), st.String())
				return
			}

			ts.Require().NoError(err)

			ts.Assert().GreaterOrEqual(int(tt.req.GetPageSize()), len(res.GetTracks()))

			for _, track := range res.GetTracks() {
				if tt.req.AlbumId != nil {
					ts.Assert().Equal(tt.req.GetAlbumId(), track.GetAlbumId())
				}
				if tt.req.Genre != nil {
					ts.Assert().Equal(tt.req.GetGenre(), track.Genre)
				}
			}
			if tt.req.SearchQuery != nil {
				ts.Assert().Contains(res.GetTracks()[0].GetTitle(), tt.req.GetSearchQuery())
			}
		})
	}
}

func TestTrackServiceGRPC(t *testing.T) {
	suite.Run(t, new(TrackServiceTest))
}
