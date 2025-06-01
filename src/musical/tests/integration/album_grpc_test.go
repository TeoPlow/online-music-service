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

type AlbumServiceTest struct {
	testutils.BaseIntegrationSuite
}

func (ts *AlbumServiceTest) SetupTest() {
	albumsData := []models.Album{
		{
			ID:          uuid.MustParse("8bd250c1-d0a8-45d5-9407-deae081e74b2"),
			Title:       "test_album 1",
			ArtistID:    uuid.MustParse("88ba1267-768e-410e-91bd-adacbcbeb6d5"),
			ReleaseDate: time.Now(),
		},
		{
			ID:          uuid.MustParse("d4b79799-4a3c-4abd-8323-6814d234fd8d"),
			Title:       "test_album 2",
			ArtistID:    uuid.MustParse("88ba1267-768e-410e-91bd-adacbcbeb6d5"),
			ReleaseDate: time.Now(),
		},
		{
			ID:          uuid.MustParse("c20129d7-8a2a-42bb-8210-9becf2dcece5"),
			Title:       "test_album 3",
			ArtistID:    uuid.MustParse("31a52596-5c16-4200-b270-1f51c2025093"),
			ReleaseDate: time.Now(),
		},
		{
			ID:          uuid.MustParse("f83202c3-e6aa-4156-b7c8-d29f6e39196b"),
			Title:       "test_album 4",
			ArtistID:    uuid.MustParse("250ab1b1-f157-4b14-8fd8-9e362eac3cb1"),
			ReleaseDate: time.Now(),
		},
		{
			ID:          uuid.MustParse("a7c6fdc8-8e58-40d3-abda-6440f8760de2"),
			Title:       "test_album 5",
			ArtistID:    uuid.MustParse("6aa52977-b1e5-4e30-b2ae-874616766c7d"),
			ReleaseDate: time.Now(),
		},
	}
	artistsData := []models.Artist{
		{ID: uuid.MustParse("88ba1267-768e-410e-91bd-adacbcbeb6d5")},
		{ID: uuid.MustParse("31a52596-5c16-4200-b270-1f51c2025093")},
		{ID: uuid.MustParse("250ab1b1-f157-4b14-8fd8-9e362eac3cb1")},
		{ID: uuid.MustParse("6aa52977-b1e5-4e30-b2ae-874616766c7d")},
	}

	testutils.TruncateTables(ts.TxM.GetDatabase(), "albums", "artists")
	for _, artist := range artistsData {
		err := ts.ArtistRepo.Add(ts.T().Context(), artist)
		if err != nil {
			panic(err)
		}
	}
	for _, album := range albumsData {
		err := ts.AlbumRepo.Add(ts.T().Context(), album)
		if err != nil {
			panic(err)
		}
	}
}

func (ts *AlbumServiceTest) TearDownTest() {
	testutils.TruncateTables(ts.TxM.GetDatabase(), "albums", "artists")
}

func (ts *AlbumServiceTest) TestCreateAlbum() {
	tests := []struct {
		name     string
		req      *pb.CreateAlbumRequest
		wantCode codes.Code
	}{
		{
			name: "success",
			req: &pb.CreateAlbumRequest{
				Title:    "test_album 6",
				ArtistId: "6aa52977-b1e5-4e30-b2ae-874616766c7d",
			},
			wantCode: codes.OK,
		},
		{
			name: "already exists",
			req: &pb.CreateAlbumRequest{
				Title:    "test_album 1",
				ArtistId: "88ba1267-768e-410e-91bd-adacbcbeb6d5",
			},
			wantCode: codes.AlreadyExists,
		},
		{
			name: "no artist",
			req: &pb.CreateAlbumRequest{
				Title:    "test_album 7",
				ArtistId: "30a1415e-1287-44bd-b856-78bb21b1375c",
			},
			wantCode: codes.NotFound,
		},
		{
			name: "invalid uuid",
			req: &pb.CreateAlbumRequest{
				Title:    "test_ablum 8",
				ArtistId: "i'm uuid",
			},
			wantCode: codes.InvalidArgument,
		},
	}

	for _, tt := range tests {
		ts.T().Run(tt.name, func(t *testing.T) {
			res, err := ts.Client.CreateAlbum(ts.T().Context(), tt.req)
			if err != nil {
				st, ok := status.FromError(err)
				ts.Require().True(ok)
				ts.Assert().Equal(tt.wantCode, st.Code(), st.String())
				return
			}

			actual, err := ts.AlbumRepo.GetByID(ts.T().Context(), uuid.MustParse(res.GetId()))
			ts.Require().NoError(err)

			ts.Assert().Equal(res.GetId(), actual.ID.String())
			ts.Assert().Equal(res.GetTitle(), actual.Title)
			ts.Assert().Equal(res.GetArtistId(), actual.ArtistID.String())
			ts.Assert().Equal(testutils.DateOnly(res.GetReleaseDate().AsTime()), testutils.DateOnly(actual.ReleaseDate))

			ts.Assert().Equal(tt.req.GetTitle(), actual.Title)
			ts.Assert().Equal(tt.req.GetArtistId(), actual.ArtistID.String())
		})
	}
}

func (ts *AlbumServiceTest) TestGetAlbum() {
	tests := []struct {
		name     string
		req      *pb.IDRequest
		wantCode codes.Code
	}{
		{
			name: "success",
			req: &pb.IDRequest{
				Id: "8bd250c1-d0a8-45d5-9407-deae081e74b2",
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
			res, err := ts.Client.GetAlbum(ts.T().Context(), tt.req)
			if err != nil {
				st, ok := status.FromError(err)
				ts.Require().True(ok)
				ts.Assert().Equal(tt.wantCode, st.Code(), st.String())
				return
			}

			expect := models.Album{
				ID:          uuid.MustParse("8bd250c1-d0a8-45d5-9407-deae081e74b2"),
				Title:       "test_album 1",
				ArtistID:    uuid.MustParse("88ba1267-768e-410e-91bd-adacbcbeb6d5"),
				ReleaseDate: time.Now(),
			}
			ts.Require().NoError(err)

			ts.Assert().Equal(tt.req.GetId(), res.GetId())

			ts.Assert().Equal(expect.ID.String(), res.GetId())
			ts.Assert().Equal(expect.Title, res.GetTitle())
			ts.Assert().Equal(expect.ArtistID.String(), res.GetArtistId())
			ts.Assert().Equal(testutils.DateOnly(expect.ReleaseDate),
				testutils.DateOnly(res.GetReleaseDate().AsTime()))
		})
	}
}

func (ts *AlbumServiceTest) TestListAlbum() {
	tests := []struct {
		name     string
		req      *pb.ListAlbumsRequest
		wantCode codes.Code
	}{
		{
			name: "full query",
			req: &pb.ListAlbumsRequest{
				Page:        1,
				PageSize:    2,
				ArtistId:    testutils.ToPtr("88ba1267-768e-410e-91bd-adacbcbeb6d5"),
				SearchQuery: testutils.ToPtr("2"),
			},
			wantCode: codes.OK,
		},
		{
			name: "only artist",
			req: &pb.ListAlbumsRequest{
				Page:     1,
				PageSize: 2,
				ArtistId: testutils.ToPtr("6aa52977-b1e5-4e30-b2ae-874616766c7d"),
			},
			wantCode: codes.OK,
		},
		{
			name: "only search",
			req: &pb.ListAlbumsRequest{
				Page:        1,
				PageSize:    2,
				SearchQuery: testutils.ToPtr("5"),
			},
			wantCode: codes.OK,
		},
		{
			name: "invalid uuid",
			req: &pb.ListAlbumsRequest{
				Page:        1,
				PageSize:    2,
				ArtistId:    testutils.ToPtr("I'm uuid"),
				SearchQuery: testutils.ToPtr("5"),
			},
			wantCode: codes.InvalidArgument,
		},
		{
			name: "empty result",
			req: &pb.ListAlbumsRequest{
				Page:     1,
				PageSize: 2,
				ArtistId: testutils.ToPtr("a592ffc7-fccb-4e29-8b90-647b80cbb758"),
			},
			wantCode: codes.OK,
		},
	}

	for _, tt := range tests {
		ts.T().Run(tt.name, func(t *testing.T) {
			res, err := ts.Client.ListAlbums(ts.T().Context(), tt.req)
			if err != nil {
				st, ok := status.FromError(err)
				ts.Require().True(ok)
				ts.Assert().Equal(tt.wantCode, st.Code(), st.String())
				return
			}

			ts.Require().NoError(err)

			ts.Assert().GreaterOrEqual(int(tt.req.GetPageSize()), len(res.GetAlbums()))

			for _, album := range res.GetAlbums() {
				if tt.req.ArtistId != nil {
					ts.Assert().Equal(tt.req.GetArtistId(), album.GetArtistId())
				}
			}
			if tt.req.SearchQuery != nil {
				ts.Assert().Contains(res.GetAlbums()[0].GetTitle(), tt.req.GetSearchQuery())
			}
		})
	}
}

func TestAlbumServiceGRPC(t *testing.T) {
	suite.Run(t, new(AlbumServiceTest))
}
