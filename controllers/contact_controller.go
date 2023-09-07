package controllers

import (
	"context"
	"mongo-crud/configs"
	"mongo-crud/models"
	"mongo-crud/responses"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var contactCollection *mongo.Collection = configs.GetCollection(configs.DB, "contacts")
var validate = validator.New()

func Store() gin.HandlerFunc {
    return func(c *gin.Context) {
        ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
        var Contact models.Contact
        defer cancel()

        //validate the request body
        if err := c.BindJSON(&Contact); err != nil {
            c.JSON(http.StatusBadRequest, responses.ContactResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
            return
        }

        //use the validator library to validate required fields
        if validationErr := validate.Struct(&Contact); validationErr != nil {
            c.JSON(http.StatusBadRequest, responses.ContactResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
            return
        }

        newContact := models.Contact{
            Id:       primitive.NewObjectID(),
            Name:     Contact.Name,
						Email:    Contact.Email,
        }

        result, err := contactCollection.InsertOne(ctx, newContact)
        if err != nil {
            c.JSON(http.StatusInternalServerError, responses.ContactResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
            return
        }

        c.JSON(http.StatusCreated, responses.ContactResponse{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": result}})
    }
}


func Show() gin.HandlerFunc {
    return func(c *gin.Context) {
        ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
        contactId := c.Param("contactId")
        var contact models.Contact
        defer cancel()

        objId, _ := primitive.ObjectIDFromHex(contactId)

        err := contactCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&contact)
        if err != nil {
            c.JSON(http.StatusInternalServerError, responses.ContactResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
            return
        }

        c.JSON(http.StatusOK, responses.ContactResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": contact}})
    }
}



func Update() gin.HandlerFunc {
    return func(c *gin.Context) {
        ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
        contactId := c.Param("contactId")
        var Contact models.Contact
        defer cancel()
        objId, _ := primitive.ObjectIDFromHex(contactId)

        //validate the request body
        if err := c.BindJSON(&Contact); err != nil {
            c.JSON(http.StatusBadRequest, responses.ContactResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
            return
        }

        //use the validator library to validate required fields
        if validationErr := validate.Struct(&Contact); validationErr != nil {
            c.JSON(http.StatusBadRequest, responses.ContactResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
            return
        }

        update := bson.M{"name": Contact.Name, "email": Contact.Email}
        result, err := contactCollection.UpdateOne(ctx, bson.M{"id": objId}, bson.M{"$set": update})
        if err != nil {
            c.JSON(http.StatusInternalServerError, responses.ContactResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
            return
        }

        //get updated Contact details
        var updatedContact models.Contact
        if result.MatchedCount == 1 {
            err := contactCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&updatedContact)
            if err != nil {
                c.JSON(http.StatusInternalServerError, responses.ContactResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
                return
            }
        }

        c.JSON(http.StatusOK, responses.ContactResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": updatedContact}})
    }
}


func Delete() gin.HandlerFunc {
    return func(c *gin.Context) {
        ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
        contactId := c.Param("contactId")
        defer cancel()

        objId, _ := primitive.ObjectIDFromHex(contactId)

        result, err := contactCollection.DeleteOne(ctx, bson.M{"id": objId})
        if err != nil {
            c.JSON(http.StatusInternalServerError, responses.ContactResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
            return
        }

        if result.DeletedCount < 1 {
            c.JSON(http.StatusNotFound,
                responses.ContactResponse{Status: http.StatusNotFound, Message: "error", Data: map[string]interface{}{"data": "User with specified ID not found!"}},
            )
            return
        }

        c.JSON(http.StatusOK,
            responses.ContactResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": "User successfully deleted!"}},
        )
    }
}




func Index() gin.HandlerFunc {
    return func(c *gin.Context) {
        ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
        var contacts []models.Contact
        defer cancel()

        results, err := contactCollection.Find(ctx, bson.M{})

        if err != nil {
            c.JSON(http.StatusInternalServerError, responses.ContactResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
            return
        }

        //reading from the db in an optimal way
        defer results.Close(ctx)
        for results.Next(ctx) {
            var singleContact models.Contact
            if err = results.Decode(&singleContact); err != nil {
                c.JSON(http.StatusInternalServerError, responses.ContactResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
            }

            contacts = append(contacts, singleContact)
        }

        c.JSON(http.StatusOK,
            responses.ContactResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": contacts}},
        )
    }
}