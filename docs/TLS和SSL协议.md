go语言中的crypto/tls包 实现了TLS 1.2  https://www.php.cn/manual/view/35164.html

深层次的了解TLS的工作原理可以去RFC的官方网站：www.rfc-editor.org，搜索RFC2246即可找到RFC文档



### SSL (Secure Socket Layer)

SSL(Secure Socket Layer 安全套接层)是基于HTTPS下的一个**协议加密层**，最初是由**网景**公司（Netscape）研发，后被IETF（The Internet Engineering Task Force - 互联网工程任务组）标准化后写入（RFCRequest For Comments 请求注释），RFC里包含了很多互联网技术的规范！
**HTTPS与SSL的关系**

起初是因为**HTTP**在传输数据时使用的是**明文**（虽然说**POST**提交的数据时放在报体里**看不到**的，但是还是可以通过抓包工具**窃取**到）是**不安全**的，为了解决这一隐患网景公司推出了**SSL安全套接字协议层**，SSL是基于HTTP之下TCP之上的一个协议层，是基于**HTTP标准**并对**TCP**传输数据时进行**加密**，所以**HPPTS是HTTP+SSL/TCP的简称**。

**划重点**：HPPTS是HTTP+SSL/TCP的简称

![1582140506344](http://picture.zyuhn.top/myblog/promise/20200220032827-32616.png)

[SSL协议](https://www.baidu.com/s?wd=SSL%E5%8D%8F%E8%AE%AE&tn=SE_PcZhidaonwhc_ngpagmjz&rsv_dl=gh_pc_zhidao)位于[TCP/IP协议](https://www.baidu.com/s?wd=TCP%2FIP%E5%8D%8F%E8%AE%AE&tn=SE_PcZhidaonwhc_ngpagmjz&rsv_dl=gh_pc_zhidao)与各种应用层协议之间，为数据通讯提供安全支持。SSL协议可分为两层： SSL记录协议（SSL Record Protocol）：它建立在可靠的传输协议（如TCP）之上，为高层协议提供数据封装、压缩、加密等基本功能的支持。 SSL握手协议（SSL Handshake Protocol）：它建立在SSL记录协议之上，用于在实际的数据传输开始前，通讯双方进行身份认证、协商加密算法、交换加密密钥等。



**SSL协议提供的服务主要有：**
1）认证用户和服务器，确保数据发送到正确的客户机和服务器；
2）加密数据以防止数据中途被窃取；
3）维护数据的完整性，确保数据在传输过程中不被改变。

**SSL协议的工作流程：**

**服务器认证阶段：**

1）客户端向服务器发送一个开始信息“Hello”以便开始一个新的会话连接；

2）服务器根据客户的信息确定是否需要生成新的主密钥，如需要则服务器在响应客户的“Hello”信息时将包含生成主密钥所需的信息；

3）客户根据收到的服务器响应信息，产生一个主密钥，并用服务器的[公开密钥](https://www.baidu.com/s?wd=%E5%85%AC%E5%BC%80%E5%AF%86%E9%92%A5&tn=SE_PcZhidaonwhc_ngpagmjz&rsv_dl=gh_pc_zhidao)加密后传给服务器；

4）服务器恢复该主密钥，并返回给客户一个用主密钥认证的信息，以此让客户认证服务器。

**用户认证阶段：**在此之前，服务器已经通过了客户认证，这一阶段主要完成对客户的认证。经认证的服务器发送一个提问给客户，客户则返回（数字）签名后的提问和其[公开密钥](https://www.baidu.com/s?wd=%E5%85%AC%E5%BC%80%E5%AF%86%E9%92%A5&tn=SE_PcZhidaonwhc_ngpagmjz&rsv_dl=gh_pc_zhidao)，从而向服务器提供认证。



从SSL 协议所提供的服务及其工作流程可以看出，SSL协议运行的基础是商家对消费者信息保密的承诺，这就有利于商家而不利于消费者。在电子商务初级阶段，由于运作电子商务的企业大多是信誉较高的大公司，因此这问题还没有充分暴露出来。但随着电子商务的发展，各中小型公司也参与进来，这样在[电子支付](https://www.baidu.com/s?wd=%E7%94%B5%E5%AD%90%E6%94%AF%E4%BB%98&tn=SE_PcZhidaonwhc_ngpagmjz&rsv_dl=gh_pc_zhidao)过程中的单一认证问题就越来越突出。虽然在SSL3.0中通过数字签名和[数字证书](https://www.baidu.com/s?wd=%E6%95%B0%E5%AD%97%E8%AF%81%E4%B9%A6&tn=SE_PcZhidaonwhc_ngpagmjz&rsv_dl=gh_pc_zhidao)可实现浏览器和Web服务器双方的身份验证，但是SSL协议仍存在一些问题，比如，只能提供交易中客户与服务器间的双方认证，在涉及多方的电子交易中，SSL协议并不能协调各方间的安全传输和信任关系。在这种情况下，Visa和 MasterCard两大信用卡公组织制定了[SET协议](https://www.baidu.com/s?wd=SET%E5%8D%8F%E8%AE%AE&tn=SE_PcZhidaonwhc_ngpagmjz&rsv_dl=gh_pc_zhidao)，为网上信用卡支付提供了全球性的标准。



### TLS：安全传输层协议

由于**HTTPS**的推出受到了很多人的欢迎，在**SSL更新到3.0**时，IETF对SSL3.0进行了**标准化**，并添加了少数机制(但是几乎和SSL3.0无差异)，标准化后的IETF**更名**为**TLS1.0**(Transport Layer Security 安全传输层协议)，可以说**TLS**就是**SSL的新版本3.1**

**划重点：TLS就是SSL的新版本3.1**

安全传输层协议（TLS）用于在两个通信应用程序之间提供保密性和数据完整性。该协议由两层组成： **TLS 记录协议**（TLS Record）和 **TLS 握手协议**（TLS Handshake）。较低的层为 TLS 记录协议，位于某个可靠的传输协议（例如 TCP）上面。 

TLS 记录协议提供的连接安全性具有两个基本特性：

**私有**： 对称加密用以数据加密（DES 、RC4 等）。对称加密所产生的密钥对每个连接都是唯一的，且此密钥基于另一个协议（如握手协议）协商。记录协议也可以不加密使用。
**可靠**：[信息传输](https://www.baidu.com/s?wd=%E4%BF%A1%E6%81%AF%E4%BC%A0%E8%BE%93&tn=SE_PcZhidaonwhc_ngpagmjz&rsv_dl=gh_pc_zhidao)包括使用密钥的 MAC 进行信息完整性检查。安全哈希功能（ SHA、MD5 等）用于 MAC 计算。记录协议在没有 MAC 的情况下也能操作，但一般只能用于这种模式，即有另一个协议正在使用记录协议传输协商安全参数。

TLS 记录协议用于封装各种高层协议。作为这种封装协议之一的握手协议允许服务器与客户机在应用程序协议传输和接收其第一个数据字节前彼此之间相互认证，协商加密算法和加密密钥。 TLS 握手协议提供的连接安全具有三个基本属性：

可以使用**非对称**的，或**公共密钥**的密码术来**认证**对等方的**身份**。该认证是可选的，但至少需要一个结点方。
共享加密密钥的协商是安全的。对偷窃者来说协商加密是难以获得的。此外经过认证过的连接不能获得加密，即使是进入连接中间的攻击者也不能。

协商是可靠的。没有经过通信方成员的检测，任何攻击者都不能修改通信协商。
TLS 的**最大优势**就在于：TLS 是独立于应用协议。高层协议可以透明地分布在 TLS 协议上面。然而， TLS 标准并没有规定应用程序如何在 TLS 上增加安全性；它把如何启动 TLS 握手协议以及如何解释交换的认证证书的决定权留给协议的设计者和实施者来判断。

**协议结构**

TLS 协议包括两个协议组―― TLS 记录协议和 TLS 握手协议――每组具有很多不同格式的信息。在此文件中我们只列出协议摘要并不作具体解析。具体内容可参照相关文档。

TLS 记录协议是一种分层协议。每一层中的信息可能包含长度、描述和内容等字段。记录协议支持[信息传输](https://www.baidu.com/s?wd=%E4%BF%A1%E6%81%AF%E4%BC%A0%E8%BE%93&tn=SE_PcZhidaonwhc_ngpagmjz&rsv_dl=gh_pc_zhidao)、将数据分段到可处理块、压缩数据、应用 MAC 、加密以及传输结果等。对接收到的数据进行解密、校验、解压缩、重组等，然后将它们传送到高层客户机。

TLS 连接状态指的是 TLS 记录协议的操作环境。它规定了压缩算法、加密算法和 MAC 算法。

TLS 记录层从高层接收任意大小无空块的连续数据。密钥计算：记录协议通过算法从握手协议提供的安全参数中产生密钥、 IV 和 MAC 密钥。 TLS 握手协议由三个子协议组构成，允许对等双方在记录层的安全参数上达成一致、自我认证、例示协商安全参数、互相报告出错条件。  