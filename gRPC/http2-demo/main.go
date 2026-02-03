package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"golang.org/x/net/http2"
	"io"
	"log"
	"math/big"
	"net/http"
	"time"
)

// 服务器地址
const serverAddr = "localhost:8443"

func main() {
	// --- 为服务器生成自签名证书 ---
	certPEM, keyPEM, err := generateSelfSignedCert()
	if err != nil {
		log.Fatalf("无法生成证书: %v", err)
	}
	serverCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		log.Fatalf("无法加载密钥对: %v", err)
	}

	// --- 配置并启动 HTTP/2 服务器 ---
	go func() {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("[服务器] 收到来自 %s 的请求", r.RemoteAddr)
			// 检查请求协议并响应
			fmt.Fprintf(w, "你好，世界！你的请求协议是: %s", r.Proto)
		})

		server := &http.Server{
			Addr:    serverAddr,
			Handler: handler,
			// 依然提供我们生成的证书
			TLSConfig: &tls.Config{
				Certificates: []tls.Certificate{serverCert},
			},
		}

		// 2. 明确为服务器配置 HTTP/2
		// 这是最关键的修正！
		// ConfigureServer 会修改 server.TLSConfig 以确保 NextProtos 包含 "h2"
		if err := http2.ConfigureServer(server, &http2.Server{}); err != nil {
			log.Fatalf("无法配置 HTTP/2 服务器: %v", err)
		}

		log.Printf("[服务器] 正在 %s 上启动 HTTPS 服务器 (已明确启用 HTTP/2)...\n", serverAddr)
		if err := server.ListenAndServeTLS("", ""); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[服务器] 启动失败: %v", err)
		}
	}()

	// 等待服务器启动
	time.Sleep(1 * time.Second)

	// --- 客户端部分保持不变，因为它默认就会尝试 HTTP/2 ---
	certPool := x509.NewCertPool()
	if ok := certPool.AppendCertsFromPEM(certPEM); !ok {
		log.Fatal("无法将服务器证书添加到客户端证书池")
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: certPool,
			},
		},
	}

	log.Println("[客户端] 正在向服务器发送请求...")
	resp, err := client.Get("https://" + serverAddr)
	if err != nil {
		log.Fatalf("[客户端] 请求失败: %v", err)
	}
	defer resp.Body.Close()

	log.Printf("[客户端] 收到响应状态码: %s", resp.Status)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("[客户端] 读取响应体失败: %v", err)
	}

	fmt.Println("---------------------------------")
	fmt.Printf("服务器响应内容:\n%s\n", string(body))
	fmt.Println("---------------------------------")
}

// generateSelfSignedCert 函数保持不变
func generateSelfSignedCert() (certPEM, keyPEM []byte, err error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject:      pkix.Name{Organization: []string{"我的测试公司"}},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(time.Hour * 24 * 365),
		KeyUsage:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		DNSNames:     []string{"localhost"},
	}
	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return nil, nil, err
	}
	certPEM = pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	keyPEM = pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)})
	return certPEM, keyPEM, nil
}
