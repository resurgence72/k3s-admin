package service

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"k3s-admin/config"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
	"net/http"
	"time"
)

// 定义终端和容器的shell交互的格式
type TerminalMessage struct {
	Operation string `json:"operation" desc:"操作类型"`
	Data      string `json:"data" desc:"数据内容"`
	Rows      uint16 `json:"rows" desc:"终端宽"`
	Cols      uint16 `json:"cols" desc:"终端高"`
}

// 初始化一个websocket.upgrader对象,用于将http协议升级为websocket协议
var upgrader = func() websocket.Upgrader {
	upgrader := websocket.Upgrader{}
	upgrader.HandshakeTimeout = time.Second * 2
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	return upgrader
}()

// 定义结构体，实现ptyHandler接口
type TerminalSession struct {
	wsConn   *websocket.Conn
	sizeChan chan remotecommand.TerminalSize
	doneChan chan struct{}
}

var Terminal terminal

type terminal struct{}

// 定义ws的handler方法
func (t *terminal) WsHandler(w http.ResponseWriter, r *http.Request) {
	// 加载k8s配置
	conf, err := clientcmd.BuildConfigFromFlags("", config.Kubeconfig)
	if err != nil {
		fmt.Println("wshandler build config failed ", err)
		return
	}

	// 解析form入参 namespace podName container 参数
	if err := r.ParseForm(); err != nil {
		fmt.Println("parseform failed ", err)
		return
	}

	namespace, podName, containerName := r.Form.Get("namespace"), r.Form.Get("podName"), r.Form.Get("container")
	fmt.Println(namespace, podName, containerName)

	// new一个 terminalsession类型的pty实例
	pty, err := NewTerminalSession(w, r, nil)
	if err != nil {
		fmt.Println("new terminal failed ", err)
		return
	}
	defer func() {
		pty.Close()
	}()

	req := K8s.ClientSet.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(podName).
		Namespace(namespace).
		SubResource("exec").
		VersionedParams(&v1.PodExecOptions{
			Stdin:     true,
			Stdout:    true,
			Stderr:    true,
			TTY:       true,
			Container: containerName,
			Command:   []string{"sh"},
		}, scheme.ParameterCodec)

	fmt.Println(req.URL())

	// remotecommand 主要实现了http 转 spdy 并添加相关header发送请求
	executor, err := remotecommand.NewSPDYExecutor(conf, "POST", req.URL())
	if err != nil {
		fmt.Println("new spdy executor failed ", err)
		return
	}

	// 建立连接后从请求stream中发送 接受数据
	err = executor.Stream(remotecommand.StreamOptions{
		Stdin:             pty,
		Stdout:            pty,
		Stderr:            pty,
		Tty:               true,
		TerminalSizeQueue: pty,
	})
	if err != nil {
		fmt.Println("executor stream failed ", err)
		// 返回报错
		pty.Write([]byte(err.Error()))
		// 标记退出stream流
		pty.Done()
		return
	}
}

// 升级http至ws,并返回一个 terminalsession对象
func NewTerminalSession(w http.ResponseWriter, r *http.Request, respHeader http.Header) (*TerminalSession, error) {
	conn, err := upgrader.Upgrade(w, r, respHeader)
	if err != nil {
		return nil, err
	}

	session := &TerminalSession{
		wsConn:   conn,
		sizeChan: make(chan remotecommand.TerminalSize),
		doneChan: make(chan struct{}),
	}
	return session, nil
}

// 关闭doneChan 关闭触发退出终端
func (t *TerminalSession) Done() {
	close(t.doneChan)
}

// 获取web端是否resize 以及是否退出终端
func (t *TerminalSession) Next() *remotecommand.TerminalSize {
	select {
	case size := <-t.sizeChan:
		return &size
	case <-t.doneChan:
		return nil
	}
}

// 读取web端输入，接受web段输入的指令类容
func (t *TerminalSession) Read(p []byte) (int, error) {
	_, message, err := t.wsConn.ReadMessage()
	if err != nil {
		return 0, err
	}

	var msg TerminalMessage
	err = json.Unmarshal(message, &msg)
	if err != nil {
		return 0, err
	}

	switch msg.Operation {
	case "stdin":
		return copy(p, msg.Data), nil
	case "resize":
		t.sizeChan <- remotecommand.TerminalSize{Width: msg.Cols, Height: msg.Rows}
		return 0, nil
	case "ping":
		return 0, nil
	default:
		return 0, nil
	}
}

// 向web端输出，接受web指令执行后，将结果返回出去
func (t *TerminalSession) Write(p []byte) (int, error) {
	msg, err := json.Marshal(TerminalMessage{
		Operation: "stdout",
		Data:      string(p),
	})
	if err != nil {
		return 0, err
	}

	err = t.wsConn.WriteMessage(websocket.TextMessage, msg)
	if err != nil {
		return 0, err
	}
	return len(p), nil
}

// 关闭websocket
func (t *TerminalSession) Close() error {
	return t.wsConn.Close()
}
