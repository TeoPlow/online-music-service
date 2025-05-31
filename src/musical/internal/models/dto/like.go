package dto

type LikeArtistRequest struct {
	ArtistID string
}

type UnlikeArtistRequest struct {
	ArtistID string
}

type LikeTrackRequest struct {
	TrackID string
}

type UnlikeTrackRequest struct {
	TrackID string
}

type GetLikedArtistsRequest struct {
	Page     int
	PageSize int
}

type GetLikedTracksRequest struct {
	Page     int
	PageSize int
}
