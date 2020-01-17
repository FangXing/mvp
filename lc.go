package main

import (
	"fmt"
	"encoding/json"
	"strconv"
	"bytes"
	"time"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"

	"encoding/pem"
	"crypto/x509"
)

type cc1 struct {

}
type Type int
const(
	averageCapitalPlusInterest  = iota 	//0	等额本息
	averageCapital						//1	等额本金
)

type Bank struct{
	Cname	string	`json:"Cname"`
	Ename	string	`json:"Ename"`
	Logo	string	`json:"logo"`
	Tyshxydm	string	`json:"tyshxydm"`
}

type product struct{
	Id	string	`json:"productId"`
	ProductName	string	`json:"productName"`
	MaxLimit	float64	`json:"MaxLimit"`
	InterestRate	float64	`json:"interestRate"`
	Term	string	`json:"Term"`
	PaymentMethod	string	`json:"PaymentMethod"`
	Bank	Bank	`json:"bank"`
}

type Application struct{
	Id string	`json:"id"`
	Source	[]byte	`json:"Source"`
	Product	product `json:"product"`
	Sponsor	string	`json:"sponsor"`
	Approver	string	`json:"approver"`
	Status	string	`json:"status"`
	Credential	string	`json:"credential"`
	CreateTime	string	`json:"createTime"`
	ApproveTime	string	`json:"approveTime"`
	ActualCredit	float64	`json:"ActualCredit"`
	ActualTerm	string	`json:"ActualTerm"`
	ActualInteresRate	float64	`json:"ActualInteresRate"`
	ActualPaymentMethod	string	`json:"ActualPaymentMethod"`

}

type LetterCredit struct{
	Id string `json:"Id"`
	OwnerId	string	`json:"OwnerId"`
	ApplicationId string `json:"ApplicationId"`
	Balance  float64 `json:"Balance"`
	CreateTime	string	`json:"CreateTime"`
	DividedCount	int	`json:"DividedCount"`
	Symbol  string `json:"symbol"`
	Application Application `json:"Application"`
}

type TrasnferBill struct{
	Id string `json:"Id"`
	LcId string `json:"LcId"`
	ToAcct string `json:"ToAcct"`
	ToAcctName  string `json:"ToAcctName"`
	ToAcctPhone  string `json:"ToAcctPhone"`
	FromAcct   string `json:"FromAcct"`
	Amount  float64 `json:"Amount"`
	CreateTime   string `json:"CreateTime"`
	Description   string `json:"Description"`
}

type FinancingBill struct{
	Id string `json:"Id"`
	LcId string `json:"LcId"`
	Term	string	`json:"Term"`
	Status string	`json:"Status"`
	Amount  float64 `json:"Amount"`
	CreateTime   string `json:"CreateTime"`
	ApproveTime	string	`json:"approveTime"`
}

type LetterCreditFB struct{
	LetterCredit LetterCredit `json:"LetterCredit"`
	FinancingBills  []FinancingBill `json:"FinancingBills"`
}

type Msg struct{
	Status  bool    `json:"Status"`
	Code    int     `json:"Code"`
	Message string  `json:"Message"`
}

func (t *cc1) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

func (t *cc1) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("ex02 Invoke")
	function, args := stub.GetFunctionAndParameters()
	if function == "insert" {
		return t.insert(stub, args)
	}else if function == "queryPrefix" {
		return t.queryPrefix(stub, args)
	}else if function == "delete" {
		return t.delete(stub, args)
	}else if function == "deleteByPrefix" {
		return t.deleteByPrefix(stub, args)
	}else if function == "query" {
		return t.query(stub, args)
	}else if function == "createApplication" {
		return t.createApplication(stub, args)
	}else if function == "CreateProduct" {
		return t.CreateProduct(stub, args)
	}else if function == "createLC" {
		return t.createLC(stub, args)
	}else if function == "transferLC" {
		return t.transferLC(stub, args)
	}else if function == "updateApplication" {
		return t.updateApplication(stub, args)
	}else if function == "queryPrefixS" {
		return t.queryPrefixS(stub, args)
	}else if function == "financingLC" {
		return t.financingLC(stub, args)
	}else if function == "queryWithLoan" {
		return t.queryWithLoan(stub, args)
	}else if function == "updatafinancingLC" {
		return t.updatafinancingLC(stub, args)
	}else if function == "deleteCreat" {
		return t.deleteCreat(stub, args)
	}


	return shim.Error("Invalid Smart Contract function name.")
}

func (t *cc1) insert(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args)!=4{
		return shim.Error("参数有误，请检查参数,参数个数为4个")
	}
	bankInfo :=Bank{
		Cname:args[0],//银行名称
		Ename:args[1],//银行英文代码
		Logo:args[2],//logo图片
		Tyshxydm:args[3]}//统一社会信用代码
	/*accountAsBytes, _ := json.Marshal(account)
	err = stub.PutState(key, accountAsBytes)*/
	key,_:=stub.CreateCompositeKey("bank",[]string{args[3]})
	//key1 := fmt.Sprintf("bank%s",args[3])
	//result := fmt.Sprintf("%+v",bankInfo)

	creator :=  initiator(stub)
	Cname, _ := json.Marshal(bankInfo.Cname)
	stub.PutState(creator, Cname)
	bankInfoAsBytes, _ := json.Marshal(bankInfo)
	fmt.Println(string(bankInfoAsBytes))
	stub.PutState(key, bankInfoAsBytes)
	//stub.PutState(key1, bankInfoAsBytes)
	return shim.Success(nil)
}

func (t *cc1) queryPrefix(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	fmt.Println("进入queryByPrefix方法")
	var result string
	var resulList []string
	var buffer bytes.Buffer

		rs, err := stub.GetStateByPartialCompositeKey(args[0], []string{})

		if err != nil {
			fmt.Println(err)
			return shim.Error(err.Error())
		}
		defer rs.Close()

		for rs.HasNext() {
			fmt.Println("开始遍历")
			responseRange, err := rs.Next()
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(responseRange.Key)
			fmt.Println(string(responseRange.Value))
			result = string(responseRange.Value)
			resulList = append(resulList, result)
			fmt.Println("遍历结束")
		}

		buffer.WriteString("[")
		for i := 0; i < len(resulList); i++ {
			buffer.WriteString(resulList[i])
			if i != len(resulList)-1 {
				buffer.WriteString(",")
			}
		}
		buffer.WriteString("]")
		fmt.Println("查询成功，queryByPrefix方法结束")
	return shim.Success(buffer.Bytes())
}

func (t *cc1) queryPrefixS(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	fmt.Println("进入queryByPrefix方法")
	var result string
	var resulList []string
	var buffer bytes.Buffer

	rs, err := stub.GetStateByPartialCompositeKey(args[0], []string{args[1]})

	if err != nil {
		fmt.Println(err)
		return shim.Error(err.Error())
	}
	defer rs.Close()

	for rs.HasNext() {
		fmt.Println("开始遍历")
		responseRange, err := rs.Next()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(responseRange.Key)
		fmt.Println(string(responseRange.Value))
		result = string(responseRange.Value)
		resulList = append(resulList, result)
		fmt.Println("遍历结束")
	}

	buffer.WriteString("[")
	for i := 0; i < len(resulList); i++ {
		buffer.WriteString(resulList[i])
		if i != len(resulList)-1 {
			buffer.WriteString(",")
		}
	}
	buffer.WriteString("]")
	fmt.Println("查询成功，queryByPrefix方法结束")
	return shim.Success(buffer.Bytes())
}

func (t *cc1) delete(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args)!=2{
		return shim.Error("参数有误")
	}
	key,_:=stub.CreateCompositeKey(args[0],[]string{args[1]})
	stub.DelState(key)
	return shim.Success(nil)
}

func (t *cc1) deleteCreat(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args)!=3{
		return shim.Error("参数有误")
	}
	key,_:=stub.CreateCompositeKey(args[0],[]string{args[1],args[2]})
	stub.DelState(key)
	return shim.Success(nil)
}

func (t *cc1) deleteByPrefix(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("enter deleteByPrefix")

	keyArr := make([]string, len(args) - 1)

	for i:=1;i<len(args);i++{
		keyArr[i-1] = args[i]
	}


	 key,_ :=stub.CreateCompositeKey(args[0],keyArr)
	stub.DelState(key)
	 return shim.Success(nil)

	rs, err := stub.GetStateByPartialCompositeKey(args[0],keyArr)

	if err != nil {
		fmt.Println(err)
		return shim.Error(err.Error())
	}
	defer rs.Close()
	fmt.Println("1开始遍历")
	for rs.HasNext() {

		responseRange, err := rs.Next()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(responseRange.Key)

		fmt.Println(string(responseRange.Value))
		/*result = string(responseRange.Value)
		resulList = append(resulList, result)*/

		stub.DelState(responseRange.Key)
	}
	fmt.Println("1遍历结束")
	return shim.Success(nil)
}

func (t *cc1) CreateProduct(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	key,_:=stub.CreateCompositeKey("bank",[]string{args[6]})
	result,_:=stub.GetState(key)
	coinbase := &Bank{}
	json.Unmarshal(result, &coinbase)
	result2 := fmt.Sprintf("%+v",coinbase)
	fmt.Println(result2)

	maxLimit,_ :=strconv.ParseFloat(args[2],64)
	interestRate,_:=strconv.ParseFloat(args[3],64)
	period := args[4]
	product :=product{
		Id:args[0],//产品id TODO自动生成
		ProductName:args[1],//产品名称
		MaxLimit:maxLimit,//最高额度
		InterestRate:interestRate,//利率
		Term:period,//授信期限
		PaymentMethod:args[5],//还款方式
		Bank:*coinbase}//银行id

	result1 := fmt.Sprintf("%+v",product)
	fmt.Println(result1)

	//result := fmt.Sprintf("%+v",product)
	productAsBytes,_ := json.Marshal(product)
	productKey,err :=stub.CreateCompositeKey("product",[]string{args[0]})
	if err != nil {
		return shim.Error(err.Error())
	}
	stub.PutState(productKey, productAsBytes)
	key1 := fmt.Sprintf("product%s",args[0])
	stub.PutState(key1, productAsBytes)
	return shim.Success(nil)
}

func (t *cc1) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	result,_:=stub.GetState(args[0])
	return shim.Success(result)
}

func (t *cc1) createApplication(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	//var application Application
		key,_:=stub.CreateCompositeKey("product",[]string{args[0]})
		result,_:=stub.GetState(key)
		prod := &product{}
		json.Unmarshal(result, &prod)
		creator :=  initiator(stub)
		Source,_:=stub.GetState(creator)
		//ApproveTime := time.Now().Format("2006-01-02 15:04:05")
		//now := time.Now()
	now := time.Now()
	h, _ := time.ParseDuration("1h")
	h1 := now.Add(8 * h).Format("2006-01-02 15:04:05")
		billId := fmt.Sprintf("%02d%02d%02d%02d%02d%02d",
			now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())

		application := Application{

			Id:                  billId,                                   //序号
			Source:			     Source,								   //来源
			Product:             *prod,                                    //授信产品id
			Sponsor:             creator,                                  //申请机构
			Approver:            args[1],                                  //审批机构
			Status:              "0",                                      //状态
			Credential:          args[2],                                  //凭证
			CreateTime:          h1, //申请时间
			ApproveTime:         "",                                       //审批时间
			ActualCredit:        0,                                        //使用的授信额度
			ActualTerm:          "",                                       //授信期限
			ActualInteresRate:   0,                                        //利率
			ActualPaymentMethod: ""}
		key1,_:=stub.CreateCompositeKey("application",[]string{billId})
		bankInfoAsBytes, _ := json.Marshal(application)
	application_key := fmt.Sprintf("application%s",billId)
		stub.PutState(key1, bankInfoAsBytes)
	stub.PutState(application_key, bankInfoAsBytes)
	return shim.Success(nil)
}

func (t *cc1) updateApplication(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if args[5]=="1" {

		fmt.Println("进入UpdateApplication方法")
		key, _ := stub.CreateCompositeKey("application", []string{args[0]})
		result, _ := stub.GetState(key)
		prod := &Application{}
		json.Unmarshal(result, &prod)
		//result2 := fmt.Sprintf("%+v",prod)
		//fmt.Println(result2)
		prod.Status = "1"
		now := time.Now()
		h, _ := time.ParseDuration("1h")
		h1 := now.Add(8 * h).Format("2006-01-02 15:04:05")
		prod.ApproveTime = h1 //审批时间
		ActualCredit, _ := strconv.ParseFloat(args[1], 64)

		prod.ActualCredit = ActualCredit //使用的授信额度
		prod.ActualTerm = args[2]        //授信期限

		ActualInteresRate, _ := strconv.ParseFloat(args[3], 64)
		prod.ActualInteresRate = ActualInteresRate //利率
		prod.ActualPaymentMethod = args[4]         //还款方式
		bankInfoAsBytes, _ := json.Marshal(prod)
		//fmt.Println(bankInfoAsBytes)
		stub.PutState(key, bankInfoAsBytes)
		application_key := fmt.Sprintf("application%s", args[0])
		stub.PutState(application_key, bankInfoAsBytes)
		fmt.Println("UpdateApplication结束")
	}else if args[5]=="2" {
		fmt.Println("进入UpdateApplication方法")
		key, _ := stub.CreateCompositeKey("application", []string{args[0]})
		result, _ := stub.GetState(key)
		prod := &Application{}
		json.Unmarshal(result, &prod)
		//result2 := fmt.Sprintf("%+v",prod)
		//fmt.Println(result2)
		prod.Status = "2"
		bankInfoAsBytes, _ := json.Marshal(prod)
		//fmt.Println(bankInfoAsBytes)
		stub.PutState(key, bankInfoAsBytes)
		application_key := fmt.Sprintf("application%s", args[0])
		stub.PutState(application_key, bankInfoAsBytes)
		fmt.Println("UpdateApplication结束")
	}
	return shim.Success(nil)
}

//信证开立
func (t *cc1) createLC(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	application_id := args[0]
	//application_key,_:=stub.CreateCompositeKey("application",[]string{args[0]})
	application_key := fmt.Sprintf("application%s",args[0])

	applicationBytes,err := stub.GetState(application_key)
	if err != nil{
		return shim.Error(err.Error())
	}

	fmt.Printf("GetState(%s) %s \n", application_key, string(applicationBytes))
	if string(applicationBytes) == "" {
			msg := &Msg{Status: true, Code: 0, Message: "信证开立失败,没有申请记录"}
			rev, _ := json.Marshal(msg)
			return shim.Error(string(rev))
	}

	application := &Application{}
	json.Unmarshal(applicationBytes, &application)



	balance := args[1]
	bank_name := application.Product.Bank.Ename

	creator :=  initiator(stub)
	_name :=  fmt.Sprintf("%s-%s-%s",creator,bank_name,application_id)
	_symbol:= 	fmt.Sprintf("%s-%s-%s",creator,bank_name,application_id)

	fmt.Println(balance)
	fmt.Println(bank_name)
	fmt.Println(_name)
	fmt.Println(_symbol)
	fmt.Println(creator)
	//发币者账户初始额度是balance
	fmt.Println("111111111111111111111111111111")
	trans:=[][]byte{[]byte("initCurrency"),[]byte(_name),[]byte(_symbol),[]byte(balance),[]byte(creator)}
	stub.InvokeChaincode("AAA",trans,"chain1")
	/*trans:=[][]byte{[]byte(args[2]),[]byte(args[3])}
	response:= stub.InvokeChaincode(args[0],trans,args[1])*/
	fmt.Println("跨链码调用,initCurrency")
	/*if response.Status != shim.OK {
		errStr := fmt.Sprintf("Failed to invoke chaincode. Got error: %s", string(response.Payload))
		fmt.Printf(errStr)
		return shim.Error(errStr)
	}*/

	fmt.Println("22222222222222222222222222222222")

	//银行账户的初始值是0

	trans = [][]byte{[]byte("initCurrency"),[]byte(_name),[]byte(_symbol),[]byte("0"),[]byte(bank_name)}
	stub.InvokeChaincode("AAA",trans,"chain1")
	fmt.Println("跨链码调用,initCurrency")
	/*if response.Status != shim.OK {
		errStr := fmt.Sprintf("Failed to invoke chaincode. Got error: %s", string(response.Payload))
		fmt.Printf(errStr)
		return shim.Error(errStr)
	}*/
	fmt.Println("33333333333333333333333333333333")

	//保存信证
	now := time.Now()
	billId := fmt.Sprintf("%02d%02d%02d%02d%02d%02d",
	now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())

	h, _ := time.ParseDuration("1h")
	h1 := now.Add(8 * h).Format("2006-01-02 15:04:05")
	balanceValue,_  :=   strconv.ParseFloat(balance, 64)
	lc := LetterCredit{
		Id: billId,
		OwnerId:creator,
		ApplicationId: application_id,
		Balance:balanceValue,
		DividedCount:0,
		Symbol:_symbol,
		CreateTime:h1,
		Application:*application}
	fmt.Printf("LC %+v\n", lc)

	lcAsBytes, err := json.Marshal(lc)
	if err != nil {
		return shim.Error(err.Error())
	}

	LcAccKey, err := stub.CreateCompositeKey("LC", []string{creator,billId})
	err = stub.PutState(LcAccKey, lcAsBytes)	
	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Printf("created LC %s \n", string(lcAsBytes))
	
	msg := &Msg{Status: true, Code: 0, Message: "开立信证成功"}
	rev, _ := json.Marshal(msg)
	return shim.Success(rev)

}

//信证流转
func (t *cc1) transferLC(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	toAccount := args[0]
	toAccountName := args[1]
	toAccountPhone := args[2]
	amount := args[3]
	description := args[4]
	lcId := args[5]
	lc := &LetterCredit{}
	var lcKey string

	creator :=  initiator(stub)

 	LcResultsIterator, err := stub.GetStateByPartialCompositeKey ("LC",[]string{creator,lcId}) 
	if err != nil {
		return shim.Error(err.Error())
	}

    defer LcResultsIterator.Close()
	var application Application
    for LcResultsIterator.HasNext(){
    	item, _ := LcResultsIterator.Next()
    	LCBytes,err :=stub.GetState(item.Key)
		if err != nil {
			return shim.Error(err.Error())
		}
		fmt.Printf("GetState(%s) %s \n", item.Key, string(LCBytes))
		if string(LCBytes) == "" {
				msg := &Msg{Status: true, Code: 0, Message: "信证流转失败,没有信证"}
				rev, _ := json.Marshal(msg)
				return shim.Error(string(rev))
		}

		lcKey =  item.Key
		json.Unmarshal(LCBytes, &lc)
		application = lc.Application
    }

	existAsBytes,err := stub.GetState(toAccount)
	fmt.Printf("GetState(%s) %s \n", lcKey, string(existAsBytes))
	if string(existAsBytes) == "" {
		fmt.Println("init account:%s",toAccount)
		// init_args  := []string{lc.Symbol,lc.Symbol,"0",toAccount}

		trans := [][]byte{[]byte("initCurrency"),[]byte(lc.Symbol),[]byte(lc.Symbol),[]byte("0"),[]byte(toAccount)}
		response  := stub.InvokeChaincode("AAA",trans,"chain1")
		fmt.Println("跨链码调用,initCurrency,%s",toAccount)
		if response.Status != shim.OK {
			errStr := fmt.Sprintf("Failed to invoke chaincode. Got error: %s", string(response.Payload))
			fmt.Printf(errStr)
			return shim.Error(errStr)
		}
	}


	trans := [][]byte{[]byte("transferToken"),[]byte(creator),[]byte(toAccount),[]byte(lc.Symbol),[]byte(amount)}
	response  := stub.InvokeChaincode("AAA",trans,"chain1")
	fmt.Println("跨链码调用,transferToken,%s --> %s",creator,toAccount)

	if response.Status != shim.OK {
		errStr := fmt.Sprintf("Failed to invoke chaincode. Got error: %s", string(response.Payload))
		fmt.Printf(errStr)
		return shim.Error(errStr)
	}else{
		now := time.Now()
		billId := fmt.Sprintf("%02d%02d%02d%02d%02d%02d",
		now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())

		/*timeStr2 := fmt.Sprintf("%02d-%02d-%02d %02d:%02d:%02d",
		now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())*/
		h, _ := time.ParseDuration("1h")
		h1 := now.Add(8 * h).Format("2006-01-02 15:04:05")
		//修改原来通证的余额
		amountVal,_ := strconv.ParseFloat(amount, 64)
		lc.Balance -= amountVal
		lc.DividedCount += 1
		lcAsBytes, _ := json.Marshal(lc)
		stub.PutState(lcKey,lcAsBytes)

		//新增子通证

		lcNewId, err := stub.CreateCompositeKey("LC", []string{toAccount,fmt.Sprintf("%s-%03d",lcId,lc.DividedCount)})
		fmt.Println(lcNewId)
		lcNew := LetterCredit{
			Id: fmt.Sprintf("%s-%03d",lcId,lc.DividedCount),
			OwnerId:creator,
			ApplicationId: lc.ApplicationId,
			Balance: amountVal,
			DividedCount:0,
			Symbol:lc.Symbol,
			CreateTime:h1,
			Application:application}

		lcNewAsBytes, _ := json.Marshal(lcNew)
		err = stub.PutState(lcNewId, lcNewAsBytes)	
		if err != nil {
			return shim.Error(err.Error())
		}

		//保存转账单

		transferBill := TrasnferBill{
			Id:billId,
			LcId:lcId,
			ToAcct:toAccount,
			ToAcctName:toAccountName,
			ToAcctPhone:toAccountPhone,
			FromAcct:creator,
			Amount:amountVal,
			CreateTime:time.Now().Format("2006-01-02 15:04:05"),
			Description:description}

		billkey, err := stub.CreateCompositeKey("tranferBill", []string{lcId,billId})

    	transferBillBytes, err := json.Marshal(transferBill)
		if err != nil {
			return shim.Error(err.Error())
		}
		err = stub.PutState(billkey, transferBillBytes)
		if err != nil {
			return shim.Error(err.Error())
		}else{
			fmt.Printf("transferBill %s \n", string(transferBillBytes))
		}

		msg := &Msg{Status: true, Code: 0, Message: "信证流转成功"}
		rev, _ := json.Marshal(msg)
		return shim.Success(rev)

	}

}

//信证融资
func (t *cc1) financingLC(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	lcId := args[0]
	amount := args[1]
	term := args[2]
	lc := &LetterCredit{}
	var lcKey string

	creator :=  initiator(stub)

	LcResultsIterator, err := stub.GetStateByPartialCompositeKey ("LC",[]string{creator,lcId})
	if err != nil {
		return shim.Error(err.Error())
	}

	defer LcResultsIterator.Close()

	for LcResultsIterator.HasNext(){
		item, _ := LcResultsIterator.Next()
		LCBytes,err :=stub.GetState(item.Key)
		if err != nil {
			return shim.Error(err.Error())
		}
		fmt.Printf("GetState(%s) %s \n", item.Key, string(LCBytes))
		if string(LCBytes) == "" {
			msg := &Msg{Status: true, Code: 0, Message: "信证融资失败,没有信证"}
			rev, _ := json.Marshal(msg)
			return shim.Error(string(rev))
		}
		lcKey =  item.Key
		json.Unmarshal(LCBytes, &lc)
	}

	application_key := fmt.Sprintf("application%s",lc.ApplicationId)

	applicationBytes,err := stub.GetState(application_key)
	if err != nil{
		return shim.Error(err.Error())
	}

	fmt.Printf("GetState(%s) %s \n", application_key, string(applicationBytes))
	if string(applicationBytes) == "" {
		msg := &Msg{Status: true, Code: 0, Message: "信证融资失败,没有申请记录"}
		rev, _ := json.Marshal(msg)
		return shim.Error(string(rev))
	}

	application := &Application{}
	json.Unmarshal(applicationBytes, &application)


	toAccount := application.Product.Bank.Ename


	trans := [][]byte{[]byte("transferToken"),[]byte(creator),[]byte(toAccount),[]byte(lc.Symbol),[]byte(amount)}
	response  := stub.InvokeChaincode("Currencys",trans,"chain1")
	fmt.Println(fmt.Sprintf("跨链码调用,transferToken,%s --> %s",creator,toAccount))
	if response.Status != shim.OK {
		errStr := fmt.Sprintf("Failed to invoke chaincode. Got error: %s", string(response.Payload))
		fmt.Printf(errStr)
		return shim.Error(errStr)
	}else{
		//修改原来通证的余额,并保存
		amountVal,_ := strconv.ParseFloat(amount, 64)
		lc.Balance -= amountVal
		lcAsBytes, _ := json.Marshal(lc)
		stub.PutState(lcKey,lcAsBytes)

		now := time.Now()
		billId := fmt.Sprintf("%02d%02d%02d%02d%02d%02d",
			now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())

		/*nowStr := fmt.Sprintf("%02d-%02d-%02d %02d:%02d:%02d",
			now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())*/
		h, _ := time.ParseDuration("1h")
		h1 := now.Add(8 * h).Format("2006-01-02 15:04:05")
		financingBill := FinancingBill{
			Id:billId,
			LcId:lc.Id,
			Status:"0",
			Term:term,
			Amount:amountVal,
			CreateTime:h1,
			ApproveTime:""}

		billkey, err := stub.CreateCompositeKey("financingBill", []string{lcId,billId})

		financingBillBytes, err := json.Marshal(financingBill)
		if err != nil {
			return shim.Error(err.Error())
		}
		err = stub.PutState(billkey, financingBillBytes)
		if err != nil {
			return shim.Error(err.Error())
		}else{
			fmt.Printf("financingBillBytes %s \n", string(financingBillBytes))
		}

		msg := &Msg{Status: true, Code: 0, Message: "信证融资成功"}
		rev, _ := json.Marshal(msg)
		return shim.Success(rev)

	}

}

func (t *cc1) updatafinancingLC(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	now := time.Now()
	h, _ := time.ParseDuration("1h")
	h1 := now.Add(8 * h).Format("2006-01-02 15:04:05")
	lcId :=args[0]
	billId :=args[1]
	billkey, _ := stub.CreateCompositeKey("financingBill", []string{lcId,billId})
	result,_:=stub.GetState(billkey)
	prod := &FinancingBill{}
	json.Unmarshal(result, &prod)
	prod.Status="1";
	prod.ApproveTime = h1
	bankInfoAsBytes, _ := json.Marshal(prod)
	//fmt.Println(bankInfoAsBytes)
	stub.PutState(billkey, bankInfoAsBytes)
	return shim.Success(nil)
}

func (t *cc1) queryWithLoan(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	creator :=  initiator(stub)
	//var a map[string]string
	var result []string
	//var result []LetterCreditFB
	var FinancingBillList []FinancingBill
	var buffer bytes.Buffer
	rs, err := stub.GetStateByPartialCompositeKey("LC", []string{creator})
	if err != nil {
		fmt.Println(err)
		return shim.Error(err.Error())
	}
	defer rs.Close()

	for rs.HasNext() {
		fmt.Println("开始遍历")
		responseRange, err := rs.Next()
		if err != nil {
			fmt.Println(err)
		}
		LCBytes,err :=stub.GetState(responseRange.Key)

		LetterCredit := &LetterCredit{}
		json.Unmarshal(LCBytes, &LetterCredit)
		fmt.Println(LetterCredit.Id)
		rs1, err := stub.GetStateByPartialCompositeKey("financingBill", []string{LetterCredit.Id})

		if err != nil {
			fmt.Println(err)
			return shim.Error(err.Error())
		}
		defer rs1.Close()

		for rs1.HasNext() {
			fmt.Println("开始遍历")
			responseRange, err := rs1.Next()
			if err != nil {
				fmt.Println(err)
			}
			LCBytes,err :=stub.GetState(responseRange.Key)
			FinancingBill := &FinancingBill{}
			json.Unmarshal(LCBytes, &FinancingBill)
			fmt.Println(FinancingBill)
			FinancingBillList = append(FinancingBillList, *FinancingBill)
		}
		bankInfoAsBytes, _ := json.Marshal(FinancingBillList)
		LetterCreditAsBytes, _ := json.Marshal(LetterCredit)
		/*LetterCreditFB :=LetterCreditFB{
			LetterCredit:*LetterCredit,
			FinancingBills:FinancingBillList}*/
		result = append(result, string(LetterCreditAsBytes))
		result = append(result, string(bankInfoAsBytes))
	}
	buffer.WriteString("[")
	for i := 0; i < len(result); i++ {
		buffer.WriteString(result[i])
		if i != len(result)-1 {
			buffer.WriteString(",")
		}
	}
	buffer.WriteString("]")
	fmt.Println("查询成功，queryByPrefix方法结束")

	return shim.Success(buffer.Bytes())

	//return shim.Success(nil)


}

//交易发起人
func initiator(stub shim.ChaincodeStubInterface) string {
	//获取当前用户
	creatorByte, _ := stub.GetCreator()
	certStart := bytes.IndexAny(creatorByte, "-----BEGIN")
	if certStart == -1 {
		fmt.Println("No certificate found")
	}
	certText := creatorByte[certStart:]
	bl, _ := pem.Decode(certText)
	if bl == nil {
		fmt.Println("Could not decode the PEM structure")
	}

	cert, err := x509.ParseCertificate(bl.Bytes)
	if err != nil {
		fmt.Println("ParseCertificate failed")
	}
	name := cert.Subject.CommonName
	fmt.Println("initiator:" + name)
	return name
}


func main()  {
	err := shim.Start(new(cc1))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}

}
