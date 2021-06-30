package main

import (
	"context"
	"github.com/go-vgo/robotgo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strings"
	"time"
)

type LatLng struct {
	LAT string `json:"lat"`
	LNG string `json:"lng"`
}
type Pharmacy struct {
	Id                 primitive.ObjectID `bson:"_id" json:"_id,omitempty"`
	PharmacyName        string `json:"pharmacyName"`
	PharmacyAddress     string `json:"pharmacyAddress"`
	PharmacyPhoneNumber string `json:"pharmacyPhoneNumber"`
	PharmacyProvince    string `json:"pharmacyProvince"`
	PharmacyDistrict    string `json:"pharmacyDistrict"`
	PharmacyLatLng      string `json:"pharmacyLatLng"`
}

var listPharmacy []Pharmacy
func main() {

	getAllDocument()
	if len(listPharmacy)> 0 {
		for i := 0; i < len(listPharmacy); i++ {
			latitude , longitude := getLATLNG(listPharmacy[i].PharmacyName + " "+ listPharmacy[i].PharmacyDistrict+" "+listPharmacy[i].PharmacyProvince)
			if latitude == "" || longitude == "" {
				continue
			}
			updateMongoDbLatLng(listPharmacy[0].Id.Hex(),latitude,longitude)
		}
	}
}
func getAllDocument() {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://34.69.194.57:5783"))
	if err != nil {
		println("Hata Var")
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
	result , errorAll := client.Database("Eczane").Collection("Nobetci|08-Jun-21").Find(ctx,bson.M{})
	if errorAll != nil {
		println("Gelmedi")
	}

	writeError := result.All(ctx, &listPharmacy)
	if writeError != nil {
		println("Write Error")
	}

}
func updateMongoDbLatLng(_id string,lat string,lng string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://34.69.194.57:5783"))
	if err != nil {
		println("Hata Var")
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
	hex, idError := primitive.ObjectIDFromHex(_id)
	if idError != nil {
		println("İd error")

	}
	result , errorUpdate := client.Database("Eczane").Collection("Nobetci|08-Jun-21", options.Collection()).UpdateOne(ctx,bson.M{"_id": bson.M{"$eq": hex}},bson.M{"$set": bson.M{"pharmacyLatLng": lat+","+lng}})
	if errorUpdate != nil {
		println("Hata Var")
		println(errorUpdate.Error())
	}
	println(result.ModifiedCount)
}

func getLATLNG(pharmacyAddress string) (string, string) {

	robotgo.Move(190, 140)
	robotgo.MouseClick()
	//robotgo.KeyTap("backspace","")
	robotgo.TypeStr(pharmacyAddress)
	robotgo.KeyTap("enter")
	robotgo.Sleep(4)
	robotgo.Move(800, 100)
	robotgo.MouseClick("right")
	robotgo.Sleep(1)
	robotgo.Move(840, 192)
	robotgo.MouseClick()

	all, err := robotgo.ReadAll()
	if err != nil {
		print("Hata Var")
	}
	splitOne := strings.Split(all, "@")
	var splitTwo []string
	if len(splitOne) >= 2 {
		splitTwo = strings.Split(splitOne[1], ",")
	}
	//robotgo.Sleep(1)
	robotgo.Move(450, 140)
	robotgo.MouseClick()
	//robotgo.Sleep(1)
	robotgo.Move(190, 140)
	//robotgo.MouseClick()
	println(splitTwo[0],splitTwo[1])

	return splitTwo[0],splitTwo[1]
}
