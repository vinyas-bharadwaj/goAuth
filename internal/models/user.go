package models

type User struct {
	Id       string `protobuf:"id,omitempty" bson:"_id,omitempty"`
	Username string `protobuf:"username,omitempty" bson:"username,omitempty"`
	Email    string `protobuf:"email,omitempty" bson:"email,omitempty"`
	Password string `protobuf:"password,omitempty" bson:"password,omitempty"`
	Role     string `protobuf:"role,omitempty" bson:"role,omitempty"`
	GoogleId string `protobuf:"google_id,omitempty" bson:"google_id,omitempty"`
	Picture  string `protobuf:"picture,omitempty" bson:"picture,omitempty"`
}
