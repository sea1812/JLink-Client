package app

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/util/guid"
	_ "golang.org/x/tools/go/ssa"
	"time"
)

//加载所有参数设置
func LoadAllConfig() g.Map {
	m_broker_topic := g.Cfg().GetString("mqtt.broker_topic")
	m_broker_ip := g.Cfg().GetString("mqtt.broker_ip")
	m_broker_port := g.Cfg().GetInt("mqtt.broker_port")
	m_broker_user := g.Cfg().GetString("mqtt.broker_user")
	m_broker_pwd := g.Cfg().GetString("mqtt.broker_pwd")
	return g.Map{
		"broker_topic": m_broker_topic,
		"broker_ip":    m_broker_ip,
		"broker_port":  m_broker_port,
		"broker_user":  m_broker_user,
		"broker_pwd":   m_broker_pwd,
	}
}

var mc mqtt.Client

//初始化MQTT
func InitMQTT(server_ip string, server_port int, server_user string, server_pwd string, server_topic string) {
	var broker = server_ip
	var port = server_port
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("mqtt://%s:%d", broker, port))
	opts.SetKeepAlive(time.Second * 30)
	opts.SetPingTimeout(time.Second * 30)
	opts.SetClientID(guid.S())
	opts.SetUsername(server_user)
	opts.SetPassword(server_pwd)
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	mc = mqtt.NewClient(opts)
	if token := mc.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	//sub(mc, server_topic) //不订阅任何频道
}

//处理消息
var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("MQTT Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

//订阅主题
func sub(client mqtt.Client, topic string) {
	token := client.Subscribe(topic, 1, nil)
	token.Wait()
	fmt.Printf("MQTT Subscribed to topic %s", topic)
}

//发布消息
//func publish(client mqtt.Client, command TMqttFirstBill) string {
//	//command.Msg_Id = guid.S()
//	command.SendTime = time.Now()
//	j, err := json.Marshal(command)
//	if err == nil {
//		token := client.Publish(command.Topic, 0, false, string(j))
//		token.Wait()
//		//time.Sleep(time.Second)
//		return ""
//	} else {
//		return err.Error()
//	}
//}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("MQTT Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("MQTT Connect lost: %v", err)
	token := mc.Connect()
	token.Wait()
	sub(mc, g.Cfg().GetString("mqtt.broker_topic"))
}
