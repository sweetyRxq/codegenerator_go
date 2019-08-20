package model

type ResourceCar struct {
//GENERATE_START	 
	CarId string`json:"carId"`	 
	CarNo string`json:"carNo"`	 
	PlateNum string`json:"plateNum"`	 
	State int`json:"state"`	
	DataType string `json:"dataType"`
//GENERATE_END
}