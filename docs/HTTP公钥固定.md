# HTTP公钥固定（Public Key Pinning）

## 什么是HPKP？

HTTP公钥固定（HPKP）是一种安全功能，它告诉Web客户端将特定加密公钥与某个Web服务器相关联，以降低使用伪造证书进行“MITM攻击（中间人攻击）”的风险。

为了确保TLS会话中使用的服务器公钥的真实性，此公钥将包装到X.509证书中，该证书通常由证书颁发机构（CA）签名。诸如浏览器之类的Web客户端信任许多这些CA，它们都可以为任意域名创建证书。如果攻击者能够破坏单个CA，则他们可以对各种TLS连接执行MITM攻击。 HPKP可以通过告知客户端哪个公钥属于某个Web服务器来规避HTTPS协议的威胁。

HPKP是首次使用信任（TOFU）技术。 Web服务器第一次通过特殊的HTTP头告诉客户端哪些公钥属于它，客户端将该信息存储在给定的时间段内。当客户端再次访问服务器时，它希望证书链中至少有一个证书包含一个公钥，其指纹已通过HPKP获知。如果服务器提供未知的公钥，则客户端应向用户发出警告。

## 如何启用HPKP？

要为您的站点启用此功能，您需要在通过HTTPS访问站点时返回Public-Key-Pins HTTP标头：
```shell script
Public-Key-Pins: pin-sha256="base64=="; max-age=expireTime [; includeSubDomains][; report-uri="reportURI"]
```
