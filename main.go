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
	"math/rand"
    "strings"
	"errors"
	"time"
)

type cc1 struct {
}

func (t *cc1) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

func (t *cc1) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	//queryWSQ

	function, args := stub.GetFunctionAndParameters()
	if function == "submit" {
		return t.submit(stub, args)
	}else if function == "query" {
		return t.query(stub, args)
	} else if function == "rangeQuery" {
		return t.rangeQuery(stub, args)
	}else if function == "unGrant"{
		return t.unGrant(stub, args)
	}else if function == "accountUserInfo"{
		return t.accountUserInfo(stub, args)
	}else if function == "queryByPrefix"{
		return t.queryByPrefix(stub, args)
	}else if function == "reportCreate"{
		return t.reportCreate(stub, args)
	}else if function == "grant"{
		return t.grant(stub, args)
	}else if function == "showPriv"{
		return t.showPriv(stub, args)
	}else if function == "revoke"{
		return t.revoke(stub, args)
	}else if function == "reportList"{
		return t.reportList(stub, args)
	}else if function == "reportDetail"{
		return t.reportDetail(stub, args)
	}
	return shim.Error("Invalid Smart Contract function name.")
}

func (t *cc1) submit(stub shim.ChaincodeStubInterface, args []string) pb.Response{
	fmt.Println("进入submit方法")
	datas ,_:= base64.StdEncoding.DecodeString(args[0])
	var s []map[string]interface {}
	var a map[string]interface {}
	var err error
	err =json.Unmarshal(datas, &s)

	if err!= nil {
		fmt.Println(err)
		return shim.Error(err.Error())
	}
	var gfKey string
	var xfKey string
	for i:=0;i<len(s);i++{
		fmt.Println("遍历数组")
		a = s[i]["fpxx"].(map[string]interface {})
		gfKey = fmt.Sprintf("fpdata:gffp:%s:%s",a["gfsbh"],a["kprq"])
		datastr,err := json.Marshal(s[i])

		if err != nil {
			fmt.Println(err)
			return shim.Error(err.Error())
		}

		err = stub.PutState(gfKey,[]byte(string(datastr)))
		if err!= nil {
			fmt.Println(err)
			return shim.Error(err.Error())
		}

		xfKey = fmt.Sprintf("fpdata:xffp:%s:%s",a["xfsbh"],a["kprq"])
		err  = stub.PutState(xfKey,[]byte(string(datastr)))
		if err!= nil {
			fmt.Println(err)
			return shim.Error(err.Error())
		}
		//创建组合键
		gfmc := fmt.Sprintf("%s",a["gfmc"])
		gfsbh := fmt.Sprintf("%s",a["gfsbh"])
		xfmc := fmt.Sprintf("%s",a["xfmc"])
		xfsbh := fmt.Sprintf("%s",a["xfsbh"])
		key1,_:=stub.CreateCompositeKey("gfmc",[]string{gfmc,gfsbh})
		key2,_:=stub.CreateCompositeKey("xfmc",[]string{xfmc,xfsbh})
		fmt.Println("key1",key1)
		fmt.Println("key2",key2)
		stub.PutState(key1, []byte(fmt.Sprintf(`{"gsmc":"%s","gfsbh":"%s"}`,gfmc,gfsbh)))
		stub.PutState(key2, []byte(fmt.Sprintf(`{"xfmc":"%s","xfsbh":"%s"}`,xfmc,xfsbh)))
		fmt.Println("遍历数组结算")
	}
	fmt.Println("发票上传成功，submit方法结束")
	return shim.Success([]byte("submit success"))
}

func (t *cc1) accountNumber (stub shim.ChaincodeStubInterface, args []string) pb.Response{
	if args[0]=="bank" {
		fmt.Println("basnk")
		key1, _ := stub.CreateCompositeKey("account", []string{ args[0], args[1]})
		fmt.Println(key1)
		stub.PutState(key1, []byte(fmt.Sprintf("{name:%s}", args[2])))
	}else if args[0]=="zx"{
		fmt.Println("zx")
		key1, _ := stub.CreateCompositeKey("account", []string{ args[0], args[1]})
		fmt.Println(key1)
		stub.PutState(key1, []byte(fmt.Sprintf("{name:%s}", args[2])))
	}else if args[0]=="qiye"{
		fmt.Println("qiye")
		key1, _ := stub.CreateCompositeKey("account", []string{ args[0], args[1]})
		fmt.Println(key1)
		stub.PutState(key1, []byte(fmt.Sprintf("{name:%s}", args[2])))
	}
	return shim.Success([]byte("submit success"))
}
func (t *cc1) accountUserInfo (stub shim.ChaincodeStubInterface, args []string) pb.Response{
	userinfo,_:=stub.CreateCompositeKey(fmt.Sprintf("account%s",args[0]),[]string{args[1],args[2]})
	stub.PutState(userinfo, []byte(fmt.Sprintf(`{"mc":"%s","taxno":"%s"}`,args[1],args[2])))

	return shim.Success([]byte("submit success"))
}



func (t *cc1) queryByPrefix(stub shim.ChaincodeStubInterface, args []string) pb.Response{
	fmt.Println("进入queryByPrefix方法")
	var result string
	var resulList []string
	rs, err := stub.GetStateByPartialCompositeKey(args[0], []string{})
	if err != nil{
		fmt.Println(err)
		return  shim.Error(err.Error())
	}
	defer rs.Close()

	for rs.HasNext(){
		fmt.Println("开始遍历")
		responseRange, err := rs.Next()
		if err != nil{
			fmt.Println(err)
		}
		fmt.Println(responseRange.Key)
		fmt.Println(string(responseRange.Value))
		result = string(responseRange.Value)
		resulList = append(resulList,result)
		fmt.Println("遍历结束")
	}
	var buffer bytes.Buffer
	buffer.WriteString("[")
	for i:=0;i<len(resulList) ;i++  {
		buffer.WriteString(resulList[i])
		if i != len(resulList)-1 {
			buffer.WriteString(",")
		}
	}
	buffer.WriteString("]")
	fmt.Println("查询成功，queryByPrefix方法结束")
	return shim.Success(buffer.Bytes())
}

func (t *cc1) queryBy(stub shim.ChaincodeStubInterface, args []string) pb.Response{
	fmt.Println("进入queryByPrefix方法")
	var result string
	var resulList []string
	rs, err := stub.GetStateByPartialCompositeKey(args[0], []string{})
	if err != nil{
		fmt.Println(err)
		return  shim.Error(err.Error())
	}
	defer rs.Close()

	for rs.HasNext(){
		fmt.Println("开始遍历")
		responseRange, err := rs.Next()
		if err != nil{
			fmt.Println(err)
		}
		fmt.Println(responseRange.Key)
		fmt.Println(string(responseRange.Value))
		result = string(responseRange.Value)
		resulList = append(resulList,result)
		fmt.Println("遍历结束")
	}
	var buffer bytes.Buffer
	buffer.WriteString("[")
	for i:=0;i<len(resulList) ;i++  {
		buffer.WriteString(resulList[i])
		if i != len(resulList)-1 {
			buffer.WriteString(",")
		}
	}
	buffer.WriteString("]")
	fmt.Println("查询成功，queryByPrefix方法结束")
	return shim.Success(buffer.Bytes())
}

func (t *cc1) query(stub shim.ChaincodeStubInterface, args []string) pb.Response{
	fmt.Println("进入query方法")
	key := fmt.Sprintf("fpdata:%s:%s:%s",args[0],args[1],args[2])
	fp,err := stub.GetState(key)

	if err != nil {
		return shim.Error(err.Error())
	}
	if fp == nil {
		return shim.Error("Entity not found")
	}
	fmt.Println("查询成功，query方法结束")
	return shim.Success(fp)

}
func (t *cc1)  rangeQuery(stub shim.ChaincodeStubInterface, args []string) pb.Response{
	fmt.Println("进入rangeQuery方法")
	var result string
	startKey:=fmt.Sprintf("fpdata:%s:%s:%s",args[0],args[1],args[2])
	endKey:=fmt.Sprintf("fpdata:%s:%s:%s",args[0],args[1],args[3])
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
	fmt.Println("查询成功，rangeQuery方法结束")
	return shim.Success([]byte(result))

}
func string_to_map(s string)(map[string]interface {}){
    var fp map[string]interface{}
    // 将字符串反解析为字典
    json.Unmarshal([]byte(s), &fp)
    return fp 
}

func randString(length int) string {
    rand.Seed(time.Now().UnixNano())
    rs := make([]string, length)
    for start := 0; start < length; start++ {
        t := rand.Intn(3)
        if t == 0 {
            rs = append(rs, strconv.Itoa(rand.Intn(10)))
        } else if t == 1 {
            rs = append(rs, string(rand.Intn(26)+65))
        } else {
            rs = append(rs, string(rand.Intn(26)+97))
        }
    }
    return strings.Join(rs, "")
}


func (t *cc1) reportCreate(stub shim.ChaincodeStubInterface, args []string) pb.Response{

	name,err := getCreator(stub)
	if err != nil{
			return shim.Error(err.Error())
	}	

	fpdata:="fpdata"
	gxf := args[0]
	sh := args[1]
	start := args[2]
	end := args[3]
	yhsh := args[4]
	startKey:=fmt.Sprintf("%s:%s:%s:%s",fpdata,gxf,sh,start)
	endKey:=fmt.Sprintf("%s:%s:%s:%s",fpdata,gxf,sh,end)

	info,err := stub.GetStateByRange(startKey,endKey)


	certMsg,_ := GetMsgFromCert(stub)
	rs, err := stub.GetStateByPartialCompositeKey(fmt.Sprintf("priv1"), []string{certMsg["oid"],sh})
	//没有授权，直接返回
	if rs == nil {
		jsonResp := fmt.Sprintf("{\"Error\":\"[%s]没有得到[%s]的授权，不能生成报告 \"}",name,args[4])
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
	result["je"] = strconv.FormatFloat(je, 'f', -1, 64)
	jsonRsp, err := json.Marshal(result)
	if err != nil{
		return shim.Error(err.Error())
	}
	fmt.Println(jsonRsp)
	random:=randString(15)

	report_key1,_ := stub.CreateCompositeKey(fmt.Sprintf("report%s",certMsg["oid"]),[]string{args[1],random}) //征信
	report_key2,_ := stub.CreateCompositeKey(fmt.Sprintf("report%s",yhsh),[]string{args[1],random})

	reporterr1 := stub.PutState(report_key1,[]byte(string(jsonRsp)))
	reporterr2 := stub.PutState(report_key2,[]byte(string(jsonRsp)))

	if reporterr2 != nil {
		fmt.Println(reporterr2)
		return shim.Error(reporterr2.Error())
	}

	if reporterr1 != nil {
		fmt.Println(reporterr1)
		return shim.Error(reporterr1.Error())
	} else {
		fmt.Println("report key",reporterr1,reporterr2)
	}
	fmt.Println(string(jsonRsp))

	return shim.Success(jsonRsp)

}

func (t *cc1) reportList(stub shim.ChaincodeStubInterface, args []string) pb.Response{
	certMsg,_ := GetMsgFromCert(stub)

	reports,err:= stub.GetStateByPartialCompositeKey(fmt.Sprintf("report%s",certMsg["oid"]),[]string{})
	if err != nil {
        return shim.Error(err.Error())
    }
    defer reports.Close()
    
    var report []string

    for i := 0; reports.HasNext(); i++ {
        responseRange, responseerr := reports.Next()

		_, reportKeyParts, _  := stub.SplitCompositeKey(responseRange.Key)

		if responseerr != nil {
            return shim.Error(responseerr.Error())
        }

		report = append(report, fmt.Sprintf("%s_%s",reportKeyParts[0],reportKeyParts[1]))
		fmt.Println("reportKeyParts",reportKeyParts)

    }

	var buffer bytes.Buffer

	for i:=0;i<len(report) ;i++  {
		buffer.WriteString(report[i])
		if i != len(report)-1 {
			buffer.WriteString(",")
		}
	}

	return shim.Success(buffer.Bytes())
}

func (t *cc1) reportDetail(stub shim.ChaincodeStubInterface, args []string) pb.Response{

	certMsg,_ := GetMsgFromCert(stub)

	reports,err := stub.GetStateByPartialCompositeKey(fmt.Sprintf("report%s",certMsg["oid"]),args)

	fmt.Println("reports",reports)

	if err != nil {
		return shim.Error(err.Error())
	}
	defer reports.Close()

	var report []string
	for i := 0; reports.HasNext(); i++ {
		responseRange, responseerr := reports.Next()
		if responseerr != nil {
			return shim.Error(responseerr.Error())
		}

		report_value := responseRange.Value
		report = append(report, string(report_value))
	}

	reportjson,_ :=json.Marshal(report)

	return shim.Success(reportjson)
}


func  (t *cc1)  grant(stub shim.ChaincodeStubInterface, args []string) pb.Response{

	certMsg,_ := GetMsgFromCert(stub)

	key1,_:=stub.CreateCompositeKey(fmt.Sprintf("priv%s",args[0]),[]string{certMsg["oid"],args[1]})

	//key2,_:=stub.CreateCompositeKey(fmt.Sprintf("priv%s%s",args[0],certMsg["oid"]),[]string{args[1]})
	//反向
	/*key3,_:=stub.CreateCompositeKey(fmt.Sprintf("privreversal%s",args[0]),[]string{args[1],certMsg["oid"]})
	key4,_:=stub.CreateCompositeKey(fmt.Sprintf("privreversal%s%s",args[0],args[1]),[]string{certMsg["oid"]})*/

	stub.PutState(key1,[]byte(fmt.Sprintf(`{"gsmc":"%s","tyshxydm":"%s"}`,args[2],args[1])))
	//stub.PutState(key2,[]byte(fmt.Sprintf(`{"gsmc":"%s","tyshxydm":"%s"}`,args[2],args[1])))

	/*stub.PutState(key3,[]byte(fmt.Sprintf(`{"gsmc":"%s","tyshxydm":"%s"}`,args[2],args[1])))
	stub.PutState(key4,[]byte(fmt.Sprintf(`{"gsmc":"%s","tyshxydm":"%s"}`,args[2],args[1])))*/
	fmt.Println(key1,certMsg["oid"])
	//fmt.Println(key2,certMsg["oid"])
	//stub.DelState(args[3])
	return shim.Success(nil)
}
func  (t *cc1)  unGrant(stub shim.ChaincodeStubInterface, args []string) pb.Response{
	certMsg,_ := GetMsgFromCert(stub)
	key1,_:=stub.CreateCompositeKey(fmt.Sprintf("priv%s",args[0]),[]string{certMsg["oid"],args[1]})
	stub.DelState(key1)
	//key2,_:=stub.CreateCompositeKey(fmt.Sprintf("priv%s%s",args[0],certMsg["oid"]),[]string{args[1]})
	//stub.DelState(key2)
	//fmt.Println(fmt.Sprintf("priv%s%s%s",args[0],certMsg["oid"],args[1]))
	//stub.DelState(args[3])
	return shim.Success([]byte(string("删除成功")))
}



func GetMsgFromCert(stub shim.ChaincodeStubInterface) (map[string]string, error) {
	var role string
	var id string
	var application string
	var taxNum string
	var oid string
	certMsg := make(map[string]string)

	creator, _ := stub.GetCreator()
	certStart := bytes.IndexAny(creator, "-----BEGIN")
	if certStart == -1 {
		return nil, errors.New("8002")
	}

	certText := creator[certStart:]
	certDERBlock, _ := pem.Decode(certText)
	cert, err := x509.ParseCertificate(certDERBlock.Bytes)
	if err != nil {
		return nil, errors.New("8002")
	}

	for i := 0; i < len(cert.Extensions); i++ {
		ExtensionId := cert.Extensions[i].Id.String()
		fmt.Println(ExtensionId)
		fmt.Println(string(cert.Extensions[i].Value))
		if ExtensionId == "0.9.2342.19200300.100.1.1" {
			oid = string(cert.Extensions[i].Value)
		}
	}
	certMsg["role"] = role
	certMsg["id"] = id
	certMsg["application"] = application
	certMsg["taxNum"] = taxNum
	certMsg["oid"] = oid

	return certMsg, nil
}

func  (t *cc1)  revoke(stub shim.ChaincodeStubInterface, args []string) pb.Response{

	if len(args) != 3{
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	var owner,taker,priv_type string

	owner = args[0]
	taker = args[1]
	priv_type = args[2]

	key := fmt.Sprintf("privilege:%s:%s", owner,taker)

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

	key := fmt.Sprintf("privilege:%s:%s", owner, taker)

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
	return shim.Success( []byte(jsonResp))
}

func querySchoolIds(stub shim.ChaincodeStubInterface) []string {
	resultsIterator, err := stub.GetStateByPartialCompositeKey("School", []string{"school"})
	if err != nil {
		return nil
	}
	defer resultsIterator.Close()

	scIds := make([]string,0)
	for i := 0; resultsIterator.HasNext(); i++ {
		responseRange, err := resultsIterator.Next()
		if err != nil {
			return nil
		}
		_, compositeKeyParts, err := stub.SplitCompositeKey(responseRange.Key)
		if err != nil {
			return nil
		}
		returnedSchoolId := compositeKeyParts[1]
		scIds = append(scIds, returnedSchoolId)
	}
	return scIds
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
