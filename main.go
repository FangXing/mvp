package main

import (
	"fmt"
	"strconv"
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
	}else if function == "report"{
		return t.report(stub, args)
	}


	return shim.Error("Invalid Smart Contract function name.")
}

func (t *cc1) submit(stub shim.ChaincodeStubInterface, args []string) pb.Response{
	datas ,_:= base64.StdEncoding.DecodeString(args[0])
	var s []map[string]interface {}
	var a map[string]interface {}
	err :=json.Unmarshal(datas, &s)
	if err!= nil {
		fmt.Println(err)
	}
	for i:=0;i<len(s);i++{
		data:=s[i]
		a = data["fpxx"].(map[string]interface {})
		gffp := "fpdata:gffp"
		var gfsbh  interface{} = a["gfsbh"]
		var gfkprq  interface{} =a["kprq"]
		gfKey:=fmt.Sprintf("%s:%s:%s",gffp,gfsbh,gfkprq)
		gferr:=stub.PutState(gfKey,datas)

		if gferr!= nil {
			fmt.Println(gferr)
		}
		xffp :="fpdata:xffp"
		var xfsbh  interface{} =a["xfsbh"]
		var xfkprq  interface{} =a["kprq"]
		xfKey:=fmt.Sprintf("%s:%s:%s",xffp,xfsbh,xfkprq)
		xferr:=stub.PutState(xfKey,datas)

		if xferr!= nil {
			fmt.Println(xferr)
		}

	}
	/*for i:=0;i<len(args);i++ {
		arrays:=args[i]
		*//*bytes ,_:= base64.StdEncoding.DecodeString(arrays)*//*
		m := make(map[string]interface{})
		err := json.Unmarshal([]byte(arrays), &m)
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
	}*/
	return shim.Success([]byte("submit success"))
}

func (t *cc1) query(stub shim.ChaincodeStubInterface, args []string) pb.Response{
	var result string
	/*var response string
	var interErr Error*/
	fpdata:="fpdata"
	gxf := args[0]
	sh := args[1]
	start := args[2]
	end := args[3]
	startKey:=fmt.Sprintf("%s:%s:%s:%s",fpdata,gxf,sh,start)
	endKey:=fmt.Sprintf("%s:%s:%s:%s",fpdata,gxf,sh,end)
	info,err := stub.GetStateByRange(startKey,endKey)

	rsp := make(map[string]string)
	for info.HasNext(){
		response, interErr := info.Next()
		if interErr != nil{
			return shim.Error(interErr.Error())
		}
		rsp[response.Key] = string(response.Value)
		fmt.Println(response.Key, string(response.Value))
		result = string(response.Value)
	}
	if err != nil {
		return shim.Error(err.Error())
	}
	if info == nil {
		return shim.Error("Entity not found")
	}

	return shim.Success([]byte(result))

}

func string_to_map(s string)(map[string]interface {}){
    var fp map[string]interface{}
    // var fp map[string]string
    // 将字符串反解析为字典
    json.Unmarshal([]byte(s), &fp)
    // fmt.Println(fp)
    return fp 
}

func (t *cc1) report(stub shim.ChaincodeStubInterface, args []string) pb.Response{

	fpdata:="fpdata"
	gxf := args[0]
	sh := args[1]
	start := args[2]
	end := args[3]
	startKey:=fmt.Sprintf("%s:%s:%s:%s",fpdata,gxf,sh,start)
	endKey:=fmt.Sprintf("%s:%s:%s:%s",fpdata,gxf,sh,end)

	info,err := stub.GetStateByRange(startKey,endKey)
	var num int = 0
	var je float64 = 0.0//float64 = 0.0
	// var temp float64 
	for info.HasNext(){
		response, interErr := info.Next()
		if interErr != nil{
			return shim.Error(interErr.Error())
		}
		
		fp := string_to_map(string(response.Value))
		var fpje interface{} = fp["fpxx"].(map[string]interface{})["je"]
		fpjestr:=fmt.Sprintf("%s",fpje)
		temp, err := strconv.ParseFloat(fpjestr, 64)
		if err != nil{
			return shim.Error(err.Error())
		}		
		je = je + temp
		num = num + 1
	}


	result := make(map[string]string)
	result["num"] = strconv.Itoa(num)
	result["je"] = strconv.FormatFloat(je, 'E', -1, 64)
	jsonRsp, err := json.Marshal(result)
	if err != nil{
		return shim.Error(err.Error())
	}
	fmt.Println(jsonRsp)

	report_key :=fmt.Sprintf("%s:%s:%s:%s",args[0],args[1],args[2],args[3])
	reporterr := stub.PutState(report_key,[]byte(string(jsonRsp)))
	if reporterr != nil {
		fmt.Println(reporterr)
	} else {
		fmt.Println("report key",report_key)
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
