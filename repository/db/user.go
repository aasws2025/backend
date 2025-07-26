package db

import (
	"api/config"
	"api/model"
	"errors"
	"github.com/kamagasaki/go-utils/mongo"
	"go.mongodb.org/mongo-driver/bson"
)

func InsertUserSata(requestData model.UserAccount) error {
	db := mongo.MongoConnect(DBATS)
	insertedID := mongo.InsertOneDoc(db, config.UserColl, requestData)
	if insertedID == nil {
		return errors.New("couldn't insert data")
	}
	return nil
}

func CheckEmailOrTelfonExists(email, telfon string) (int, error) {
	db := mongo.MongoConnect(DBATS)
	emailFilter := bson.M{"email": email}
	emailCount, err := mongo.CountDocuments(db, config.UserColl, emailFilter)
	if err != nil {
		return -1, err
	}

	telfonFilter := bson.M{"telfon": telfon}
	telfonCount, err := mongo.CountDocuments(db, config.UserColl, telfonFilter)
	if err != nil {
		return -1, err
	}
	switch {
	case emailCount > 0 && telfonCount > 0:
		return 3, nil
	case emailCount > 0:
		return 1, nil
	case telfonCount > 0:
		return 2, nil
	default:
		return 0, nil
	}
}

func GetUserFilter(filter bson.M) ([]model.UserAccount, error) {
	db := mongo.MongoConnect(DBATS)
	var datas []model.UserAccount
	err := mongo.GetAllDocByFilter[model.UserAccount](db, config.UserColl, filter, &datas)
	if err != nil {
		return nil, err
	}
	return datas, nil
}

func GetOneUserFilter(filter bson.M) (model.UserAccount, error) {
	db := mongo.MongoConnect(DBATS)
	var data model.UserAccount
	err := mongo.GetOneDocByFilter[model.UserAccount](db, config.UserColl, filter, &data)
	if err != nil {
		return model.UserAccount{}, err
	}
	return data, nil
}
