HOST ?= localhost:8080
PROTOSET ?= descriptor.pb

.PHONY: descriptor unary unrary-h

# Unary call
unary:
	grpcurl -plaintext -protoset $(PROTOSET) -d '$(DATA)' $(HOST) $(SERVICE)/$(METHOD)

# make unary DATA='{"username":"kruyneg", "email":"kruyneg@gmail.com", "password":"PassW*rd228", "gender":true, "country":"russia", "age":20}' SERVICE=auth.AuthService METHOD=RegisterUser

# make unary DATA='{"username":"kruyneg", "password":"PassW*rd228"}' SERVICE=auth.AuthService METHOD=Login

# make unary DATA='{"username":"artist", "password":"PassW*rd228", "author":"Singer", "producer":"Super Corp.", "country":"russia", "description":"This is artist"}' SERVICE=auth.AuthService METHOD=RegisterArtist

# make unary DATA='{"username":"artist", "password":"PassW*rd228"}' SERVICE=auth.AuthService METHOD=Login

# Unary with headers
unary-h:
	grpcurl -plaintext -protoset $(PROTOSET) -rpc-header '$(HEADERS)' -d '$(DATA)' $(HOST) $(SERVICE)/$(METHOD)

# make unary-h HEADERS="Authorization: Bearer $(cat token.txt)" DATA='{"id":""}' SERVICE=musical.MusicalService METHOD=GetArtist

# make unary-h HEADERS="Authorization: Bearer $(cat token.txt)" DATA='{"artist_id":"389f3196-9c38-4ca2-bd77-cb0d35386953", "title":"dream"}' SERVICE=musical.MusicalService METHOD=CreateAlbum

# make unary-h HEADERS="Authorization: Bearer $(cat token.txt)" DATA='{"id":"bf13d0c8-7ac7-4806-a4c2-fb547e8e29de"}' SERVICE=musical.MusicalService METHOD=GetAlbum

# make unary-h HEADERS="Authorization: Bearer $(cat token.txt)" DATA='{"refresh_token": "$(cat refresh.txt)"}' SERVICE=auth.AuthService METHOD=RefreshToken

descriptor:
	protoc \
	--proto_path=../musical/api \
	--proto_path=../auth/api \
	--include_imports \
	--include_source_info \
	--descriptor_set_out=descriptor.pb \
	../musical/api/album.proto \
	../musical/api/artist.proto \
	../musical/api/like.proto \
	../musical/api/musical.proto \
	../musical/api/track.proto \
	../auth/api/auth_proto.proto