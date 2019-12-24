package main

import (
	"fmt"
	"strconv"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"encoding/base64"
	"encoding/json"
	"bytes"
	"encoding/pem"
	"crypto/x509"
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
	}else if function == "grant"{
		return t.grant(stub, args)
	}else if function == "showPriv"{
		return t.showPriv(stub, args)
	}else if function == "revoke"{
		return t.revoke(stub, args)
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

func string_to_map(s string)(map[string]interface {}){
    var fp map[string]interface{}
    // var fp map[string]string
    // 将字符串反解析为字典
    json.Unmarshal([]byte(s), &fp)
    // fmt.Println(fp)
    return fp 
}

func (t *cc1) report(stub shim.ChaincodeStubInterface, args []string) pb.Response{

	name,err := getCreator(stub)
	if err != nil{
			return shim.Error(err.Error())
	}	

	gxf := args[0]
	sh := args[1]
	start := args[2]
	end := args[3]
	startKey:=fmt.Sprintf("%s:%s:%s",gxf,sh,start)
	endKey:=fmt.Sprintf("%s:%s:%s",gxf,sh,end)
	info,err := stub.GetStateByRange(startKey,endKey)

	key := fmt.Sprintf("%s:%s", sh, name)

	priv_type_bytes,err := stub.GetState(key)

	//没有授权，直接返回
	if priv_type_bytes == nil {
		jsonResp := fmt.Sprintf("{\"Error\":\"[%s]没有得到[%s]的授权，不能生成报告 \"}",name,sh)
		return shim.Error(jsonResp)
	}

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


//发票表
//fp:gffp:91420300571533687D:20181205
//fp:xffp:91420300571533687D:20181205

//授权表
//91420300571533687D:aisino  -> 1/2...



func  (t *cc1)  grant(stub shim.ChaincodeStubInterface, args []string) pb.Response{

	if len(args) != 3{
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	var owner,taker,priv_type string

	owner = args[0]
	taker = args[1]
	priv_type = args[2]

	key := fmt.Sprintf("%s:%s", owner,taker)
	stub.PutState(key,[]byte(priv_type))

	fmt.Printf("%s grant priv[%s] to %s",owner,taker,priv_type)
	return shim.Success(nil)
}


func  (t *cc1)  revoke(stub shim.ChaincodeStubInterface, args []string) pb.Response{

	if len(args) != 3{
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	var owner,taker,priv_type string

	owner = args[0]
	taker = args[1]
	priv_type = args[2]

	key := fmt.Sprintf("%s:%s", owner,taker)

	err := stub.DelState(key)
	if err != nil {
	   return shim.Error(fmt.Sprintf("Faild, [%s] revoke priv[%s] from [%s]",owner,priv_type,taker))
	}

	fmt.Printf("success [%s] revoke priv[%s] from [%s]",owner,priv_type,taker)
	return shim.Success([]byte(fmt.Sprintf("[%s] revoke priv[%s] from [%s]",owner,priv_type,taker)))
}



func  (t *cc1) showPriv(stub shim.ChaincodeStubInterface, args []string) pb.Response{

	if len(args) != 2{
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	var owner,taker,priv_type string

	owner = args[0]
	taker = args[1]

	key := fmt.Sprintf("%s:%s", owner, taker)

	priv_type_bytes,err := stub.GetState(key)



	if err != nil {
		jsonResp := fmt.Sprintf("{\"Error\":\"Fail to get priv for %s-%s \"}",owner, taker)
		return shim.Error(jsonResp)
	}

	


	if priv_type_bytes == nil {
		jsonResp := fmt.Sprintf("{\"Error\":\"Null priv for %s-%s \"}",owner, taker)
		return shim.Error(jsonResp)
	}

	priv_type = string(priv_type_bytes)

	jsonResp := fmt.Sprintf("{\" [%s] have priv [%s] from  [%s']\"}",owner,priv_type,taker)

	fmt.Printf("Query Response:%s",jsonResp)
	// return shim.Success(priv_type_bytes)
	return shim.Success( []byte(jsonResp))
}


// 获取操作成员
func getCreator(stub shim.ChaincodeStubInterface) (string, error) {
	creatorByte, _ := stub.GetCreator()
	certStart := bytes.IndexAny(creatorByte, "-----BEGIN")
	if certStart == -1 {
		fmt.Errorf("No certificate found")
	}
	certText := creatorByte[certStart:]
	bl, _ := pem.Decode(certText)
	if bl == nil {
		fmt.Errorf("Could not decode the PEM structure")
	}

	cert, err := x509.ParseCertificate(bl.Bytes)
	if err != nil {
		fmt.Errorf("ParseCertificate failed")
	}
	uname := cert.Subject.CommonName
	return uname, nil
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
