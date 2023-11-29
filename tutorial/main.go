package main

import (
	"fmt"
	"go-grpc-app/pb"
	"log"

	"github.com/golang/protobuf/jsonpb"
)

func main(){
	employee := &pb.Person{
		Name: "Suzuki",
		Age: 25,
		Occupation: pb.Occupation_ENGINEER,
		PhoneNumber: []string{"090-1234-5678"},
		Project: map[string]*pb.Company_Project{"Project": &pb.Company_Project{}},
		Profile: &pb.Person_Text{
			Text: "My name is Suzuki",
		},
		Birthday: &pb.Date{
			Year: 2000,
			Month: 1,
			Day: 1,
		},
	}

	// birtDate, err := proto.Marshal(employee)
	// if err != nil {
	// 	log.Fatalln("Cant serialize", err)
	// }

	// if err := os.WriteFile("test.bin", birtDate, 0666); err != nil{
	// 	log.Fatal("Cant write", err)
	// }

	// in, err := os.ReadFile("test.bin")
	// if err != nil {
	// 	log.Fatalln("cant read file", err)
	// }

	// readEmployee := &pb.Person{}

	// err = proto.Unmarshal(in, readEmployee)
	// if err != nil {
	// 	log.Fatalln("cat deserialize",err)
	// }

	// fmt.Println(readEmployee)
	

	// jsonに変換する。
	m := jsonpb.Marshaler{}
	out, err := m.MarshalToString(employee)
	if err != nil{
		log.Fatalln("cant marshal to json", err)
	}

	// fmt.Println(out)

	// jsonから構造対にもどする。
	readEmployee := &pb.Person{}
	if err := jsonpb.UnmarshalString(out, readEmployee); err != nil{
		log.Fatalln("catn unmarshal frin json", err)
	}

	fmt.Println(readEmployee)
}