package database

// import "go.mongodb.org/mongo-driver/bson"

// func GetMongoQuery (filter interface{}) (bson.M, error) {
// 	query := bson.M{}
// 	if filter != nil {
// 		query = filter.(bson.M)
// 	}

// 	// check if there's an ID field in the query
// 	if _, ok := query["_id"]; ok {
// 		id, err := helpers.ConvertStringToObjectId(query["_id"].(string))
// 		if err != nil {
// 			return nil, err
// 		}
// 		query["_id"] = id
// 	}

// 	return query, nil
// }