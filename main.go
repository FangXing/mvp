package main

import (
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"encoding/base64"
	"encoding/json"
)

type cc1 struct {
}

func (t *cc1) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

func (t *cc1) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("ex02 Invoke")
	function, args := stub.GetFunctionAndParameters()
	if function == "submit" {
		return t.submit(stub, args)
	}else if function == "query" {
		return t.query(stub, args)
	}


	return shim.Error("Invalid Smart Contract function name.")
}

func (t *cc1) submit(stub shim.ChaincodeStubInterface, args []string) pb.Response{
	for i:=0;i<len(args);i++ {
		arrays:=args[i]
		bytes ,_:= base64.StdEncoding.DecodeString(arrays)
		m := make(map[string]interface{})
		err := json.Unmarshal([]byte(bytes), &m)
		if err!= nil {
			fmt.Println(err)
		} else {
			gffp := "gffp"
			var gfsbh  interface{} = m["fpxx"].(map[string]interface{})["gfsbh"]
			var gfkprq  interface{} = m["fpxx"].(map[string]interface{})["kprq"]
			gfKey:=fmt.Sprintf("%s:%s:%s",gffp,gfsbh,gfkprq)
			gferr:=stub.PutState(gfKey,[]byte(string(bytes)))
			if gferr!= nil {
				fmt.Println(gferr)
			}
			xffp :="xffp"
			var xfsbh  interface{} = m["fpxx"].(map[string]interface{})["xfsbh"]
			var xfkprq  interface{} = m["fpxx"].(map[string]interface{})["kprq"]
			xfKey:=fmt.Sprintf("%s:%s:%s",xffp,xfsbh,xfkprq)
			xferr:=stub.PutState(xfKey,[]byte(string(bytes)))
			if xferr!= nil {
				fmt.Println(xferr)
			}
		}
	}
	return shim.Success([]byte("submit success"))
}

func (t *cc1) query(stub shim.ChaincodeStubInterface, args []string) pb.Response{

	gxf := args[0]
	sh := args[1]
	start := args[2]
	end := args[3]
	startKey:=fmt.Sprintf("%s:%s:%s",gxf,sh,start)
	endKey:=fmt.Sprintf("%s:%s:%s",gxf,sh,end)
	info,err := stub.GetStateByRange(startKey,endKey)

	rsp := make(map[string]string)
	for info.HasNext(){
		response, interErr := info.Next()
		if interErr != nil{
			return shim.Error(interErr.Error())
		}
		rsp[response.Key] = string(response.Value)
		fmt.Println(response.Key, string(response.Value))
	}

	if err != nil {
		return shim.Error(err.Error())
	}
	if info == nil {
		return shim.Error("Entity not found")
	}
	jsonRsp, err := json.Marshal(rsp)
	if err != nil{
		return shim.Error(err.Error())
	}
	return shim.Success(jsonRsp)


}
func base64Decode(src []byte) ([]byte, error) {
	return base64.StdEncoding.DecodeString(string(src))
}

func main()  {


	err := shim.Start(new(cc1))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}

}
