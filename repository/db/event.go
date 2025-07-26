package db

import (
	"api/config"
	"api/model"
	"errors"
	"fmt"
	"github.com/kamagasaki/go-utils/mongo"
	"go.mongodb.org/mongo-driver/bson"
)

func InsertDataEvent(requestData model.DataEvent) error {
	db := mongo.MongoConnect(DBATS)
	insertedID := mongo.InsertOneDoc(db, config.ColEvent, requestData)
	if insertedID == nil {
		return errors.New("couldn't insert data")
	}
	return nil
}

func GetDataEventFilter(filter bson.M) ([]model.DataEvent, error) {
	db := mongo.MongoConnect(DBATS)
	var datas []model.DataEvent
	err := mongo.GetAllDocByFilter[model.DataEvent](db, config.ColEvent, filter, &datas)
	if err != nil {
		return nil, err
	}
	return datas, nil
}

func GetOneDataEventFilter(filter bson.M) (model.DataEvent, error) {
	db := mongo.MongoConnect(DBATS)
	var data model.DataEvent
	err := mongo.GetOneDocByFilter[model.DataEvent](db, config.ColEvent, filter, &data)
	if err != nil {
		return model.DataEvent{}, err
	}
	return data, nil
}

func EditEvent(filter bson.M, updateData model.DataEvent) error {
	db := mongo.MongoConnect(DBATS)
	result := mongo.ReplaceOneDoc(db, config.ColEvent, filter, updateData)
	if result == nil || result.MatchedCount == 0 {
		return fmt.Errorf("no matching profiloe found for update")
	}
	return nil
}

func DeleteEvent(filter bson.M) (model.DataEvent, error) {
	db := mongo.MongoConnect(DBATS)
	var data model.DataEvent
	result := mongo.DeleteOneDoc(db, config.ColEvent, filter)

	if result == nil || result.DeletedCount == 0 {
		return model.DataEvent{}, fmt.Errorf("failed to delete document: no documents matched the filter")
	}

	return data, nil
}
