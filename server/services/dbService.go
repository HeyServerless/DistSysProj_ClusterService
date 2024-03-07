package services

import (
	"fmt"
	"log"
	"time"

	db "github.com/clusterService/server/db"
	model "github.com/clusterService/server/model"
	utils "github.com/clusterService/server/utils"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func InsertRequestDataIntoDb(c *gin.Context, createRequestRecord *model.CreateRequestRecordModel) (string, error) {
	// generate uuid
	uuid := utils.GenerateUUID()

	stmt, err := db.Init().Prepare("insert into requests (id, request_body, response, created_at, updated_at) values(?,?,?,?,?);")

	if err != nil {
		fmt.Print(err.Error())
		c.JSON(500, gin.H{
			"message": "Internal Server Error: Error while preparing the statement to insert data into db",
		})
		return "", err
	}
	defer stmt.Close()

	_, err = stmt.Exec(uuid, createRequestRecord.RequestBody, "", time.Now(), time.Now())

	if err != nil {
		c.JSON(500, gin.H{
			"message": "Internal Server Error: Error while inserting data into db",
		})
		fmt.Print(err.Error())
		return "", err
	}

	return uuid, nil
}

func GetResponseFromDb(c *gin.Context, getResponseModel *model.GetResponseModel) (model.CalculationResponseModel, error) {
	// get the result from db based on the request id
	stmt, err := db.Init().Prepare("select * from requests where id = ?;")
	if err != nil {
		fmt.Print(err.Error())
		c.JSON(500, gin.H{
			"message": "Internal Server Error: Error while preparing the statement to get data from db",
		})
		return model.CalculationResponseModel{}, err
	}

	defer stmt.Close()

	var response model.CalculationResponseModel
	log.Println("request id: ", getResponseModel.RequestId)

	err = stmt.QueryRow(getResponseModel.RequestId).Scan(&response.RequestId, &response.RequestBody, &response.Response, &response.Created_at, &response.Updated_at)

	if err != nil {
		c.JSON(500, gin.H{
			"message": "Internal Server Error: Error while getting data from db",
		})
		fmt.Print(err.Error())
		return model.CalculationResponseModel{}, err
	}

	return response, nil

}
