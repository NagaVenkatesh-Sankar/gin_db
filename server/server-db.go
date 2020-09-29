package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type smsSchema struct {
	To      string `json:"to" binding:"required"`
	Message string `json:"message"`
}
type deviceEntity struct {
	ID                   bson.ObjectId `json:",omitempty" binding:"-" bson:"_id"`
	Name                 string        `json:"name"`
	Os                   string        `json:"os"`
	PhoneNumber          string        `json:"phoneNumber"`
	CreatedTimestamp     time.Time     `json:",omitempty" binding:"-"`
	LastUpdatedTimestamp time.Time     `json:",omitempty" binding:"-"`
}

var dbSession *mgo.Session

const (
	database   string = "deviceDB"
	collection string = "listings"
)

func getSession() (*mgo.Session, error) {
	var err error
	if dbSession == nil {
		url := "mongodb://localhost:27017/"

		dbSession, err = mgo.Dial(url)

		fmt.Println("Connecting to: " + url)

		if err != nil {
			panic(err)
		}
	}
	return dbSession.Clone(), err
}

//MiddleDbSession session middleware
func MiddleDbSession() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Setup session, collection
		session, err := getSession()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer session.Close()
		deviceListCollection := session.DB(database).C(collection)
		c.Set("collection", deviceListCollection)
		c.Next()
	}
}
func main() {
	r := gin.Default()
	r.Use(MiddleDbSession())
	r.POST("/device/sms/:device_id/", sms)
	r.GET("/device/:device_id/", getDeviceByID)
	r.POST("/device/new/", newDevice)
	r.PUT("/device/:device_id/", updateDevice)
	r.DELETE("/device/:device_id/", deleteDevice)
	r.GET("/devices", getDevices)

	r.Run(":8090") // listen and serve on 0.0.0.0:8090
}

func newDevice(c *gin.Context) {
	// Setup session, collection
	deviceListCollection := c.Keys["collection"].(*mgo.Collection)

	jsonDeviceData := deviceEntity{}
	if err := c.ShouldBindJSON(&jsonDeviceData); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	jsonDeviceData.ID = bson.NewObjectId()
	jsonDeviceData.CreatedTimestamp = time.Now()
	jsonDeviceData.LastUpdatedTimestamp = time.Now()
	fmt.Println(jsonDeviceData)

	err := deviceListCollection.Insert(jsonDeviceData)
	if err != nil {
		log.Fatal(err)
	}
	c.JSON(http.StatusOK, gin.H{"device created id": jsonDeviceData.ID})
}

func getDeviceByID(c *gin.Context) {

	deviceListCollection := c.Keys["collection"].(*mgo.Collection)

	var device deviceEntity
	if bson.IsObjectIdHex(c.Param("device_id")) == false {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid id provided"})
		return
	}
	idParam := bson.ObjectIdHex(c.Param("device_id"))
	if err := deviceListCollection.Find(bson.M{"_id": idParam}).One(&device); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, device)
}
func updateDevice(c *gin.Context) {
	// Setup session, collection
	deviceListCollection := c.Keys["collection"].(*mgo.Collection)

	var device deviceEntity
	if bson.IsObjectIdHex(c.Param("device_id")) == false {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid id provided"})
		return
	}
	idParam := bson.ObjectIdHex(c.Param("device_id"))
	if err := deviceListCollection.Find(bson.M{"_id": idParam}).One(&device); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	jsonDeviceData := deviceEntity{}
	if err := c.ShouldBindJSON(&jsonDeviceData); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	deviceListCollection.Update(bson.M{"_id": idParam},
		bson.M{"$set": bson.M{
			"os":                   jsonDeviceData.Os,
			"name":                 jsonDeviceData.Name,
			"phonenumber":          jsonDeviceData.PhoneNumber,
			"lastupdatedtimestamp": time.Now()}})

	if err := deviceListCollection.Find(bson.M{"_id": idParam}).One(&device); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, device)
}
func deleteDevice(c *gin.Context) {
	// Setup session, collection
	deviceListCollection := c.Keys["collection"].(*mgo.Collection)

	if bson.IsObjectIdHex(c.Param("device_id")) == false {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid id provided"})
		return
	}
	idParam := bson.ObjectIdHex(c.Param("device_id"))

	if err := deviceListCollection.Remove(bson.M{"_id": idParam}); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "device deleted successfully"})
}
func getDevices(c *gin.Context) {
	// Setup session, collection
	deviceListCollection := c.Keys["collection"].(*mgo.Collection)

	var device []deviceEntity
	if err := deviceListCollection.Find(nil).All(&device); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, device)
}

func sms(c *gin.Context) {

	deviceListCollection := c.Keys["collection"].(*mgo.Collection)

	var device deviceEntity
	if bson.IsObjectIdHex(c.Param("device_id")) == false {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid id provided"})
		return
	}
	idParam := bson.ObjectIdHex(c.Param("device_id"))
	if err := deviceListCollection.Find(bson.M{"_id": idParam}).One(&device); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Create the device type
	currentDevice, err := CreateDevice(device.Os)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	jsonSmsData := smsSchema{}
	if err := c.ShouldBindJSON(&jsonSmsData); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	msg, err := currentDevice.SendSms(jsonSmsData)
	if err != nil {
		log.Printf("Issue with the SMS module, %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": msg})

}
