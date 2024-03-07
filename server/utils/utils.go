package utils

import (
	"strings"

	"github.com/google/uuid"
)

func Int32Ptr(i int32) *int32 { return &i }

var precedence = map[string]int{
	"(": 0,
	")": 0,
	"+": 1,
	"-": 1,
	"*": 2,
	"/": 2,
}

type ImageData struct {
	ImageName string
	FuncName  string
	Id        int32
	Port      int32
	NodePort  int32
}

var ImagesList []ImageData = []ImageData{
	{
		ImageName: "rajeshreddyt/grpcserveradd:latest",
		FuncName:  "add",
		Id:        1,
		Port:      3000,
		NodePort:  30001,
	},
	{
		ImageName: "rajeshreddyt/grpcserversubtract:latest",
		FuncName:  "subtract",
		Id:        2,
		Port:      3000,
		NodePort:  30002,
	},
	{
		ImageName: "rajeshreddyt/grpcservermultiply:latest",
		FuncName:  "multiply",
		Id:        3,
		Port:      3000,
		NodePort:  30003,
	},
	{
		ImageName: "rajeshreddyt/grpcserverdivision:latest",
		FuncName:  "division",
		Id:        4,
		Port:      3000,
		NodePort:  30004,
	},
}
var PbftImagesList []ImageData = []ImageData{
	{
		ImageName: "rajeshreddyt/grpcserveradd",
		FuncName:  "add",
		Id:        1,
		Port:      3000,
		NodePort:  30001,
	},
	// {
	// 	ImageName: "rajeshreddyt/grpcserversubtract",
	// 	FuncName:  "subtract",
	// 	Id:        2,
	// 	Port:      3000,
	// 	NodePort:  30002,
	// },
	// {
	// 	ImageName: "rajeshreddyt/grpcservermultiply",
	// 	FuncName:  "multiply",
	// 	Id:        3,
	// 	Port:      3000,
	// 	NodePort:  30003,
	// },
	// {
	// 	ImageName: "rajeshreddyt/grpcserverdivision",
	// 	FuncName:  "division",
	// 	Id:        4,
	// 	Port:      3000,
	// 	NodePort:  30004,
	// },
}

func GetImageNameByFuncationality(functionality string) ImageData {

	for _, image := range ImagesList {
		if strings.Contains(image.FuncName, functionality) {
			return image
		}
	}
	return ImageData{}

}
func GetPbftImageNameByFuncationality(functionality string) ImageData {

	for _, image := range ImagesList {
		if strings.Contains(image.FuncName, functionality) {
			return image
		}
	}
	return ImageData{}

}

func GenerateUUID() string {
	uuid := uuid.New()
	return uuid.String()
}
