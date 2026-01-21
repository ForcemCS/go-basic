TCP 的状态看似多，其实就分为三个阶段：**打招呼（建立连接）、干活（传输数据）、告别（断开连接）**。

### 第一阶段：打招呼（建立连接 - 三次握手）

在这个阶段，双方确认对方都在，且能听到说话。

1. **LISTEN（监听）**
   - **角色**：仓库（服务端）。
   - **动作**：仓库管理员打开门，坐在柜台前等着电话响。
   - **含义**：服务端启动了，正在等待连接。
2. **SYN_SENT（请求发送）**
   - **角色**：餐厅服务员（客户端）。
   - **动作**：服务员拨通电话，喊了一句：“喂？有人吗？”（发送 SYN 包）。
   - **含义**：客户端刚发起了连接请求，正在等对方回应。
   - *注：如果此时断网，客户端就会卡在这个状态。*
3. **SYN_RCVD（收到请求）**
   - **角色**：仓库（服务端）。
   - **动作**：管理员听到了电话，回复说：“我在！你能听到我吗？”（回复 ACK + SYN）。
   - **含义**：服务端收到了请求，并确认了自己的发送能力，正在等客户端最后的确认。
   - *注：这种状态通常很短。如果大量出现，可能是遭受了 SYN Flood 攻击（伪造 IP 狂发请求）。*
4. **ESTABLISHED（已建立 - \**最幸福的状态\**）**
   - **角色**：双方。
   - **动作**：服务员说：“听到了！咱们开始干活吧。”
   - **含义**：连接彻底打通，大家可以愉快地传输数据了。**这是正常运行时的主要状态。**

------

### 第二阶段：告别（断开连接 - 四次挥手）

这是最容易出问题的阶段。TCP 允许任何一方先挂电话。
在 HTTP 短连接中，通常是**客户端（餐厅）**先说“再见”。

#### 一、 主动分手方（比如客户端）的状态：

1. **FIN_WAIT_1**
   - **动作**：服务员说：“活干完了，我要挂了。”（发送 FIN）。
   - **含义**：我已经想走了，正在等对方回复“知道了”。
2. **FIN_WAIT_2**
   - **动作**：听到了仓库说“知道了”。
   - **含义**：对方已经同意我走，但我还得等对方把手里剩下的活处理完，等对方说最后一句“我也要挂了”。
   - *注：如果对方迟迟不说最后一句“我也挂了”，你就会卡在这里。*
3. **TIME_WAIT（这就是刚才讲的那个）**
   - **动作**：终于听到仓库说“我也挂了”。服务员回复“好，拜拜”，然后**原地站立等待**。
   - **含义**：确保最后那句“好，拜拜”送到了，防止网络延迟的旧消息捣乱。

#### 二、 被动分手方（比如服务端）的状态：

这里有个非常关键的状态叫 **CLOSE_WAIT**，它比 TIME_WAIT 更常见，也更像代码 Bug。

1. **CLOSE_WAIT（等待关闭 - \**重点中的重点\**）**
   - **场景**：仓库听到了服务员说“我要挂了”。
   - **动作**：仓库回复“知道了”（ACK）。然后，仓库转头对内部喊：“哎，对方要挂电话了，咱们这边还有数据没发完吗？赶紧处理！”
   - **含义**：**对方已经想关连接了，但我这边还没准备好关（或者代码忘了关）。**
   - **故障特征**：如果你的服务器上堆积了成千上万个 `CLOSE_WAIT`，说明**你的代码有 Bug**。通常是你打开了连接，对方都断了，你却**忘记调用 `close()`**，导致连接一直半死不活地吊着，占用资源。
2. **LAST_ACK（最后确认）**
   - **动作**：仓库终于处理完了，对服务员说：“我也没事了，我也挂了。”（发送 FIN）。
   - **含义**：发送了最后的告别，等待对方回复。
3. **CLOSED（关闭）**
   - **动作**：双方都挂断，回家睡觉。
   - **含义**：连接彻底不存在了。

------

### 总结图（谁先挂电话，谁就进左边）

假设是 **客户端** 先发起断开（大部分 HTTP 请求都是这样）：

| 客户端（主动方）      | 状态流转 | 服务端（被动方）         | 解释                                    |
| --------------------- | -------- | ------------------------ | --------------------------------------- |
| **ESTABLISHED**       | 正在通话 | **ESTABLISHED**          | 正常数据传输                            |
| 发送 `FIN` (我要走了) | ->       | 收到 `FIN`               | 客户端发起关闭                          |
| **FIN_WAIT_1**        |          | **CLOSE_WAIT**           | **Server 此时需要代码显式调用 Close()** |
| 收到 `ACK` (好的)     | <-       | 发送 `ACK`               | 服务端确认收到关闭请求                  |
| **FIN_WAIT_2**        |          | (处理剩余工作...)        | 客户端等服务端忙完                      |
|                       |          | 调用 `Close()`, 发 `FIN` | **服务端忙完了，终于说我也要走了**      |
| 收到 `FIN` (我也走了) | <-       | **LAST_ACK**             | 服务端等待最后的确认                    |
| **TIME_WAIT**         | ->       | 收到 `ACK`               | **客户端进入 60秒 等待**                |
| (等 2MSL 时间)        |          | **CLOSED**               | 服务端收到确认，彻底关闭                |
| **CLOSED**            |          |                          | 客户端等待结束，彻底关闭                |

### 重点：

#### 一、 通俗解释：什么是 TIME_WAIT？

想象你在经营一家**非常火爆的餐厅（客户端）**，你要派服务员去对面的**仓库（服务器）**取食材。

1. **建立连接（三次握手）**：服务员跑过去，敲门，仓库开门，两人确认身份。

2. **传输数据**：服务员拿到了食材。

3. 关闭连接（四次挥手）

   ：

   - 服务员说：“我拿完了，我要走了。”（主动关闭）
   - 仓库说：“好的，稍等，我确认下没东西落下了……好，你可以走了。”
   - 服务员说：“再见。”

**重点来了！**

在说完“再见”之后，服务员**不能立马回到餐厅干别的事**，他必须在仓库门口**原地站立等待 1~2 分钟**。

这就是 **TIME_WAIT** 状态。

#### 为什么要傻站着等（TIME_WAIT）？

这是 TCP 协议为了保险起见设计的：

1. **防止“再见”没听到**：万一仓库没听到最后那句“再见”，仓库会重新问“你走了没？”。如果服务员已经跑了，仓库就会一直困惑。所以服务员要多站会儿，确保对方收到了告别。
2. **防止旧数据混淆**：防止上一次请求的迟到的数据包，混入到下一次新的请求中。

------

#### 二、 为什么高并发下这是个灾难？

回到计算机世界：

- **服务员** = 你的客户端的一个端口（Port）。
- **原地站立** = 端口处于 `TIME_WAIT` 状态，无法被使用。
- **等待时间** = 通常是 60 秒（Linux 默认）。

##### 算一笔账（端口耗尽）

一台电脑的端口数量是有限的（最多 65535 个，除去系统保留，能用的也就几万个）。

如果你写了一个循环，**每秒钟发送 1000 个请求**，并且**每次都新建连接、用完就关**：

1. 第 1 秒：消耗 1000 个端口，它们全部进入 `TIME_WAIT`（要等 60 秒才能释放）。
2. 第 2 秒：又消耗 1000 个。
3. ...
4. 第 30 秒：你已经有 30,000 个端口在“傻站着”（TIME_WAIT），无法用来发新请求。
5. 第 60 秒：大概 60,000 个端口都在 TIME_WAIT。**你的电脑没端口可用了！**

这时候，你的程序就会报错：`dial tcp: bind: address already in use`（地址/端口已被占用），请求发不出去了。

------

#### 三、 Go 的默认坑：MaxIdleConnsPerHost

Go 的 HTTP 客户端默认支持**长连接（Keep-Alive）**。也就是服务员拿完食材**不走**，站在那等下一次任务，这样就不需要“握手”和“挥手”，也不会有 `TIME_WAIT`。

**但是**，Go 的默认配置有一个坑：

- `MaxIdleConnsPerHost` 默认值是 **2**。

这意味着：默认情况下，针对同一个服务器（比如 `localhost:23456`），连接池里**最多只保留 2 个**空闲连接供下次复用。

**场景复现：**
如果你并发发了 100 个请求：

1. Go 会瞬间建立 100 个连接。
2. 请求处理完，准备放回连接池。
3. 连接池说：“我只能装 2 个，剩下的 98 个你们走吧（关闭连接）。”
4. 这 **98 个连接被强制关闭**，于是全部进入 **TIME_WAIT** 状态。

在高并发下，这和没用连接池几乎没区别，瞬间就会把端口耗尽。

1. **CLOSE_WAIT**:
   - **是谁？** 被通知分手的那个人（通常是服务端）。
   - **好坏？** **坏的！** 大量出现通常意味着代码 Bug。
   - **原因？** 对方已经把连接断了，你的程序却还在傻傻地持有连接，需要关闭获取到的数据源

### 示例代码

+ Server

  ```go
  package main
  
  import (
  	"context"
  	"encoding/json"
  	"log"
  	"net/http"
  	"os"
  	"os/signal"
  	"strings"
  	"syscall"
  	"time"
  )
  
  type server struct{}
  
  type messageData struct {
  	Message string `json:"message"`
  }
  
  func (srv *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  	// 0. 养成好习惯，虽然 http.Server 会自动关闭它，但显式写出来更清晰
      // 谁实现了Reader就关闭谁
  	defer r.Body.Close()
  
  	// 1. 限制请求方法，仅允许 POST
  	if r.Method != http.MethodPost {
  		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
  		return
  	}
  
  	// 2. 限制读取大小（例如最大 1MB），防止恶意的大包撑爆内存
  	//    注意：MaxBytesReader 会在读取超限时自动返回错误
  	r.Body = http.MaxBytesReader(w, r.Body, 1048576)
  
  	var message messageData
  
  	// 3. 解码请求体
  	if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
  		// 区分是数据太大还是格式错误
  		if err.Error() == "http: request body too large" {
  			http.Error(w, "Request body too large", http.StatusRequestEntityTooLarge)
  		} else {
  			http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
  		}
  		return
  	}
  
  	// 模拟业务日志
  	log.Printf("Received: %s", message.Message)
  
  	// 业务处理
  	message.Message = strings.ToUpper(message.Message)
  
  	// 4. 设置 Header
  	w.Header().Set("Content-Type", "application/json")
  
  	// 5. 写入响应
  	if err := json.NewEncoder(w).Encode(message); err != nil {
  		log.Printf("Response encoding failed: %v", err)
  	}
  }
  
  func main() {
  	// 创建 Handler
  	myHandler := &server{}
  
  	// 6. 自定义 http.Server，而不是直接用 http.ListenAndServe
  	//    这是生产环境最重要的配置！防止 Slowloris 攻击。
  	httpServer := &http.Server{
  		Addr:    ":23456",
  		Handler: myHandler,
  		// 读超时：从连接建立到读完 Request Body 的时间
  		ReadTimeout: 5 * time.Second,
  		// 写超时：从读完 Header 到响应写完的时间
  		WriteTimeout: 10 * time.Second,
  		// 空闲超时：Keep-Alive 连接在两次请求之间的最大等待时间
  		IdleTimeout: 120 * time.Second,
  		// Header 读取最大时间，防止客户端发超大 Header 耗尽资源
  		ReadHeaderTimeout: 2 * time.Second,
  	}
  
  	// 启动服务器（在 Goroutine 中启动，这样主线程可以等待信号）
  	go func() {
  		log.Println("Server starting on :23456")
  		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
  			log.Fatalf("ListenAndServe error: %v", err)
  		}
  	}()
  
  	// 7. 优雅退出 (Graceful Shutdown)
  	//    监听系统信号 (Ctrl+C 或 kill 命令)
  	quit := make(chan os.Signal, 1)
  	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
  
  	// 阻塞等待信号
  	<-quit
  	log.Println("Shutting down server...")
  
  	// 创建一个 5 秒的超时上下文，给正在处理的请求一点收尾时间
  	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
  	defer cancel()
  
  	// 8. 执行关机
  	//    Shutdown 会关闭监听端口，不再接收新请求，
  	//    但会等待活跃的连接处理完（直到 ctx 超时）
  	if err := httpServer.Shutdown(ctx); err != nil {
  		log.Fatal("Server forced to shutdown:", err)
  	}
  
  	log.Println("Server exiting")
  }
  ```

+ Client

  ```go
  package main
  
  import (
  	"bytes"
  	"context"
  	"encoding/json"
  	"fmt"
  	"io"
  	"log"
  	"net"
  	"net/http"
  	"time"
  )
  
  // MessageData 定义数据结构
  type MessageData struct {
  	Message string `json:"message"`
  }
  
  // APIClient 封装 HTTP 客户端，避免使用全局 http.DefaultClient
  type APIClient struct {
  	baseURL    string
  	httpClient *http.Client
  }
  
  // NewAPIClient 创建一个新的客户端实例
  // 在这里配置连接池、超时等关键参数
  func NewAPIClient(baseURL string) *APIClient {
  	return &APIClient{
  		baseURL: baseURL,
  		httpClient: &http.Client{
  			// 1. 设置总超时时间（包括连接、写、读），防止请求永远挂起
  			Timeout: 10 * time.Second,
  
  			// 2. 自定义 Transport 以优化连接池
  			Transport: &http.Transport{
  				Proxy: http.ProxyFromEnvironment,
  				DialContext: (&net.Dialer{
  					Timeout:   5 * time.Second,  // 建连超时
  					KeepAlive: 30 * time.Second, // TCP KeepAlive
  				}).DialContext,
  				MaxIdleConns:          100,              // 最大空闲连接数
  				MaxIdleConnsPerHost:   10,               // 每个 Host 的最大空闲连接数
  				IdleConnTimeout:       90 * time.Second, // 空闲连接保持时间
  				TLSHandshakeTimeout:   10 * time.Second, // TLS 握手超时
  				ExpectContinueTimeout: 1 * time.Second,
  			},
  		},
  	}
  }
  
  // PostData 发送数据
  // 3. 引入 context.Context 用于控制取消和截止时间
  // 4. 返回 error 而不是直接 fatal，让调用者决定如何处理
  func (c *APIClient) PostData(ctx context.Context, msg MessageData) (*MessageData, error) {
  	// JSON 序列化
  	jsonBytes, err := json.Marshal(msg)
  	if err != nil {
  		return nil, fmt.Errorf("marshaling error: %w", err)
  	}
  
  	// 5. 使用 NewRequestWithContext 构建请求
  	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL, bytes.NewBuffer(jsonBytes))
  	if err != nil {
  		return nil, fmt.Errorf("create request error: %w", err)
  	}
  
  	// 设置 Header
  	req.Header.Set("Content-Type", "application/json")
  
  	// 发送请求
  	resp, err := c.httpClient.Do(req)
  	if err != nil {
  		return nil, fmt.Errorf("request failed: %w", err)
  	}
  
  	// 6. 确保 Body 关闭，并读取剩余内容以便 TCP 连接复用
  	defer func() {
  		// 丢弃 Body 剩余内容，这是复用连接的关键步骤
  		_, _ = io.Copy(io.Discard, resp.Body)
  		resp.Body.Close()
  	}()
  
  	// 7. 检查 HTTP 状态码，不要假设总是 200
  	if resp.StatusCode != http.StatusOK {
  		// 尝试读取服务器返回的错误信息（如果有）
  		bodyBytes, _ := io.ReadAll(resp.Body)
  		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(bodyBytes))
  	}
  
  	// 反序列化响应
  	var result MessageData
  	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
  		return nil, fmt.Errorf("decoding response error: %w", err)
  	}
  
  	return &result, nil
  }
  
  func main() {
  	// 初始化客户端（通常在应用启动时做一次，单例模式）
  	client := NewAPIClient("http://localhost:23456")
  
  	msg := MessageData{Message: "hi server!"}
  
  	// 8. 创建带超时的 Context，实现请求粒度的超时控制
  	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
  	defer cancel()
  
  	// 调用业务逻辑
  	respData, err := client.PostData(ctx, msg)
  	if err != nil {
  		// 生产环境通常记录日志而非 Fatal
  		log.Printf("Error processing request: %v", err)
  		return
  	}
  
  	fmt.Printf("Server responded: %s\n", respData.Message)
  }
  ```

  