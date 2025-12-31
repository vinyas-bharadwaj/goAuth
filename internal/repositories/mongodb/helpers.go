package mongodb

import (
	"goAuth/internal/models"
	pb "goAuth/proto/gen"
	"reflect"
)

func MapModelToPb[M any, P any](model *M, newPb func() *P) *P {
	pbEntity := newPb()
	modelVal := reflect.ValueOf(model).Elem()
	pbVal := reflect.ValueOf(pbEntity).Elem()

	for i := 0; i < modelVal.NumField(); i++ {
		modelField := modelVal.Field(i)
		modelFieldType := modelVal.Type().Field(i)

		pbField := pbVal.FieldByName(modelFieldType.Name)
		if pbField.IsValid() && pbField.CanSet() {
			pbField.Set(modelField)
		}

	}
	return pbEntity
}

func MapModelUserToPbUser(userModel *models.User) *pb.User {
	return MapModelToPb(userModel, func() *pb.User { return &pb.User{} })
}

func MapPbToModel[P any, M any](pbStruct *P, newModel func() *M) *M {
	modelEntity := newModel()
	pbVal := reflect.ValueOf(pbStruct).Elem()
	modelVal := reflect.ValueOf(modelEntity).Elem()

	for j := 0; j < pbVal.NumField(); j++ {
		pbField := pbVal.Field(j)
		fieldName := pbVal.Type().Field(j).Name

		modelField := modelVal.FieldByName(fieldName)
		if modelField.IsValid() && modelField.CanSet() {
			modelField.Set(pbField)
		}
	}

	return modelEntity
}

func MapPbUserToModelUser(pbTeacher *pb.User) *models.User {
	return MapPbToModel(pbTeacher, func() *models.User { return &models.User{} })
}
