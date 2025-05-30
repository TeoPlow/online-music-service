package testutils

import (
	"context"
	"log"
	"net"
	"path/filepath"
	"runtime"

	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"

	"github.com/TeoPlow/online-music-service/src/musical/internal/config"
	grpcctrl "github.com/TeoPlow/online-music-service/src/musical/internal/controllers/grpc"
	"github.com/TeoPlow/online-music-service/src/musical/internal/db"
	"github.com/TeoPlow/online-music-service/src/musical/internal/domain"
	"github.com/TeoPlow/online-music-service/src/musical/internal/logger"
	"github.com/TeoPlow/online-music-service/src/musical/internal/storage"
	pb "github.com/TeoPlow/online-music-service/src/musical/pkg/musicalpb"
)

const bufSize = 1024 * 1024

var listener *bufconn.Listener

type BaseIntegrationSuite struct {
	suite.Suite

	ctx      context.Context
	cancel   context.CancelFunc
	grpcConn *grpc.ClientConn

	Client pb.MusicalServiceClient

	AlbumRepo  *storage.AlbumRepository
	ArtistRepo *storage.ArtistRepository
	TrackRepo  *storage.TrackRepository

	TxM         *db.TxManager
	MinIOClient *storage.MinIOClient
}

func (s *BaseIntegrationSuite) SetupSuite() {
	s.ctx, s.cancel = context.WithCancel(context.Background())

	err := config.Load(getTestConfigPath())
	s.Require().NoError(err)

	logger.InitLogger()

	// init DB, repos, services
	tmanager, err := db.NewTxManager(s.ctx)
	s.Require().NoError(err)

	artistRepo := storage.NewArtistRepo(tmanager.GetDatabase())
	artistService := domain.NewArtistService(artistRepo, tmanager)

	albumRepo := storage.NewAlbumRepo(tmanager.GetDatabase())
	albumService := domain.NewAlbumService(albumRepo, tmanager, artistService)

	minio, err := storage.NewMinIOClient(s.ctx, "test")
	s.Require().NoError(err)
	streamingService := domain.NewStreamingService(minio)

	trackRepo := storage.NewTrackRepo(tmanager.GetDatabase())
	trackService := domain.NewTrackService(trackRepo, tmanager, albumService, streamingService)

	// bufconn listener
	listener = bufconn.Listen(bufSize)
	server := grpc.NewServer()

	// grpc ctrl init + registration
	grpcSrv := grpcctrl.NewServer(artistService, albumService, trackService, streamingService)
	grpcSrv.Register(server)

	go func() {
		if err := server.Serve(listener); err != nil {
			log.Fatalf("server.Serve failed: %v", err)
		}
	}()

	// dial bufconn
	conn, err := grpc.NewClient("passthrough:///bufnet",
		grpc.WithContextDialer(bufDialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	s.Require().NoError(err)

	s.grpcConn = conn
	s.Client = pb.NewMusicalServiceClient(conn)
	s.AlbumRepo = albumRepo
	s.ArtistRepo = artistRepo
	s.TrackRepo = trackRepo
	s.TxM = tmanager
	s.MinIOClient = minio
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return listener.Dial()
}

func (s *BaseIntegrationSuite) TearDownSuite() {
	if err := s.grpcConn.Close(); err != nil {
		s.T().Logf("failed to close grpc connection: %v", err)
	}
	s.cancel()
}

// Возможно это не лучший способ задать путь к конфигу,
// но пусть будет так.
func getTestConfigPath() string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("cannot get current filename")
	}
	basePath := filepath.Dir(filename)
	return filepath.Join(basePath, "..", "..", "configs", "test.cfg.yml")
}
