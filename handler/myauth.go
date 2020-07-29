package handler

import (
	"context"
	"encoding/json"
	"time"

	log "github.com/micro/go-micro/v2/logger"

	"myauth/lib/token"

	api "github.com/micro/go-micro/v2/api/proto"
)

type Myauth struct{}

func extractValue(pair *api.Pair) string {
	if pair == nil {
		return ""
	}
	if len(pair.Values) == 0 {
		return ""
	}
	return pair.Values[0]
}

type rspMsg struct {
	Code int
	Err  string
	Msg  map[string]interface{}
}

// Myauth.Call is called by the API as /myauth/call with post body {"name": "foo"}
func (e *Myauth) Call(ctx context.Context, req *api.Request, rsp *api.Response) error {
	log.Info("Received Myauth.Call request")

	// // extract the client from the context
	// myauthClient, ok := client.MyauthFromContext(ctx)
	// if !ok {
	// 	return errors.InternalServerError("go.micro.api.myauth.myauth.call", "myauth client not found")
	// }

	// // make request
	// response, err := myauthClient.Call(ctx, &myauth.Request{
	// 	Name: extractValue(req.Post["name"]),
	// })
	// if err != nil {
	// 	return errors.InternalServerError("go.micro.api.myauth.myauth.call", err.Error())
	// }
	response := "myauth.call ok"
	b, _ := json.Marshal(response)

	rsp.StatusCode = 200
	rsp.Body = string(b)

	return nil
}

// Myauth.Call is called by the API as /myauth/call with post body {"name": "foo"}
func (e *Myauth) GetJwt(ctx context.Context, req *api.Request, rsp *api.Response) error {
	log.Info("Received Myauth.GetJwt request")

	getmap := req.GetGet()
	log.Info("Received getmap %+v\n", getmap)
	postmap := req.GetPost()
	log.Info("Received postmap %+v\n", postmap)

	// 内部服务调用底层service，通过jwt验证
	// 定义服务名和key，通过这2个参数获取jwt access token
	// serviceName := "order"
	// serviceKey := "123456"
	// 对比serviceName\serviceKey 也可以是用户名密码等，这里的示例为了方便硬编码在代码中
	// 实际项目中应该从数据库或文件读取
	serviceName := extractValue(postmap["service"])
	serviceKey := extractValue(postmap["key"])
	log.Info("serviceName %+v\n", serviceName, serviceKey)
	if serviceName != "order" || serviceKey != "123456" {
		Rsp(rsp, 403, "服务名称或key错误", nil)
		return nil
	}

	//生成jwt
	// expireTime := time.Now().Add(time.Hour * 24 * 3).Unix()
	expireTime := time.Now().Add(time.Second * 60 * 60).Unix()
	token := &token.Token{}
	token.Init([]byte("key123456")) //实际项目需从数据库或文件读取
	jwtstring, err := token.Encode("auth jwt", serviceName, expireTime)
	if err != nil {
		Rsp(rsp, 403, "jwt 生成错误", nil)
		return nil
	}

	msg := make(map[string]interface{})
	msg["jwt"] = jwtstring
	Rsp(rsp, 200, "ok", msg)
	return nil
}

// 验证jwt
func (e *Myauth) InspectJwt(ctx context.Context, req *api.Request, rsp *api.Response) error {
	log.Info("Received Myauth.InspectJwt request")
	// getmap := req.GetGet()
	// log.Info("Received getmap %+v\n", getmap)
	postmap := req.GetPost()
	// log.Info("Received postmap %+v\n", postmap)

	jwtString := extractValue(postmap["jwt"])
	log.Info("jwtString %+v\n", jwtString)
	if len(jwtString) == 0 {
		Rsp(rsp, 403, "jwt参数错误", nil)
		return nil
	}

	//解析jwt
	token := &token.Token{}
	token.Init([]byte("key123456"))
	info, err := token.Decode(jwtString)
	if err != nil {
		Rsp(rsp, 403, "jwt 解析错误", nil) //过期或jwt有问题
		return nil
	}

	t := make(map[string]interface{})
	t["data"] = info
	Rsp(rsp, 200, "ok", t)
	return nil
}

// 返回func
func Rsp(rsp *api.Response, code int, err string, msg map[string]interface{}) error {
	if msg == nil {
		msg = make(map[string]interface{})
	}
	r := &rspMsg{
		Code: code,
		Err:  err,
		Msg:  msg,
	}

	b, err2 := json.Marshal(r)
	if err2 != nil {
		log.Info("json.Marshal err", err2)
	}
	rsp.StatusCode = int32(code)
	rsp.Body = string(b)
	return nil
}
