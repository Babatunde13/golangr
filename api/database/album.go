package database

import (
	"bkoiki950/go-store/api/utils"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Album struct {
	ID        string `json:"id" bson:"_id,omitempty"`
	Title     string `json:"title" bson:"title"`
	Artist	  string `json:"artist" bson:"artist"`
	Price	 float64 `json:"price" bson:"price"`
}

type IAlbumDB interface {
	GetAllAlbums (filter Album) ([]Album, error)
	GetAlbumByID(id string) (Album, error)
	CreateAlbum(data Album) (Album, error)
	UpdateAlbum(id string, data Album) (Album, error)
	DeleteAlbum (id string) (Album, error)
}

type AlbumDB struct {}

var albumCtx = context.TODO()

func AlbumCollection() *AlbumDB {
	var db IAlbumDB = &AlbumDB{}
	return db.(*AlbumDB)
}

func getAlbumFilterQuery (data Album) (bson.M, error){
	filter := bson.M{}
	if data.ID != "" {
		objectId, err := utils.ConvertStringToObjectId(data.ID); if err != nil {
			return filter, err
		}

		filter["_id"] = objectId
	}
	if data.Artist != "" {
		filter["artist"] = data.Artist
	}
	if data.Title != "" {
		filter["title"] = data.Title
	}

	if data.Price != 0 {
		filter["price"] = data.Price
	}

	return filter, nil
}

func getAlbumUpdateQuery(data Album) bson.M {
	update := bson.M{}
	if data.Artist != "" {
		update["artist"] = data.Artist
	}
	if data.Title != "" {
		update["title"] = data.Title
	}
	if data.Price != 0 {
		update["price"] = data.Price
	}

	return bson.M{ "$set": update }
}

func (c *AlbumDB) GetAllAlbums(filter Album) ([]Album, error) {
	var albums []Album
	var albumColl, err = GetCollection("albums")
	if albumColl == nil {
		return albums, err
	}
	query, err := getAlbumFilterQuery(filter); if err != nil {
		return albums, err
	}
	cursor, err := albumColl.Find(albumCtx, query)
	if err != nil {
		fmt.Println(err)
		return albums, err
	}

	defer cursor.Close(albumCtx)
	
	if err = cursor.All(ctx, &albums); err != nil {
		return albums, err
	}

	if len(albums) == 0 {
		return []Album{}, nil
	}

	return albums, nil
}

func (c *AlbumDB) GetAlbumByID(id string) (Album, error) {
	var album Album
	var albumColl, err = GetCollection("albums")
	if albumColl == nil {
		return album, err
	}
	filter, err := getAlbumFilterQuery(Album{ID: id}); if err != nil {
		return Album{}, err
	}
	err = albumColl.FindOne(albumCtx, filter).Decode(&album); if err != nil {
		return Album{}, utils.HandleError(err, "Album not found")
	}

	fmt.Println(album)

	return album, nil
}

func (c *AlbumDB) CreateAlbum(data Album) (Album, error) {
	var albumColl, err = GetCollection("albums")
	if albumColl == nil {
		return Album{}, err
	}
	res, err := albumColl.InsertOne(albumCtx, data); if err != nil {
		fmt.Println(err)
		return Album{}, err
	}

	var album Album
	err = albumColl.FindOne(albumCtx, bson.M{"_id": res.InsertedID}).Decode(&album); if err != nil {
		return Album{}, err
	}
	idstring := utils.ConvertObjectIdToString(res.InsertedID.(primitive.ObjectID))
	album.ID = idstring

	return album, nil
}

func (c *AlbumDB) UpdateAlbum(id string, data Album) (Album, error) {
	var updatedAlbum Album
	var albumColl, err = GetCollection("albums")
	if albumColl == nil {
		return updatedAlbum, err
	}
	after := options.After
	opts := options.FindOneAndUpdateOptions{ReturnDocument: &after }
	filter, err := getAlbumFilterQuery(Album{ID: id}); if err != nil {
		return Album{}, err
	}
	update := getAlbumUpdateQuery(data)
	err = albumColl.FindOneAndUpdate(albumCtx, filter, update, &opts).Decode(&updatedAlbum); if err != nil {
		return Album{}, utils.HandleError(err, "Album not found")
	}

	return updatedAlbum, nil
}

func (c *AlbumDB) DeleteAlbum(id string) (Album, error) {
	var albumColl, err = GetCollection("albums")
	if albumColl == nil {
		return Album{}, err
	}

	filter, err := getAlbumFilterQuery(Album{ID: id}); if err != nil {
		return Album{}, err
	}

	var deletedAlbum Album
	err = albumColl.FindOneAndDelete(albumCtx, filter).Decode(&deletedAlbum); if err != nil {
		return Album{}, utils.HandleError(err, "Album not found")
	}

	return deletedAlbum, nil
}
