# MIME sniffing（MIME 嗅探）

产生这个问题的主要原因是，http报文中的content-type头在早期版本的浏览器中(例如IE7/6），浏览器并不光会凭靠content-type头进行不同类型的解析，还会`对response内容进行自我解析`，例如text/plain将文件内容直接显示，jpeg格式进行jpeg的渲染，例如传输的content-type头格式为 text/plain ，将文件内容直接显示，response中内容为
```
<script> alert(1); </script> 
```
IE浏览器会将其自动嗅探，并认为是text/html类型，并执行相应的渲染逻辑，这就很容易引发XSS攻击。 攻击者可以将图片文件内容写入xss攻击语句，并上传到共享的网站上，当用户请求该文件时，老旧浏览器收到response进行违背content-type的html解析，从而触发MIME嗅探攻击，我所知的解决方法为资源服务器使用与主网站不同的域名，在同源政策下，阻止跨域请求并解析的行为，从而达到阻止这一类嗅探攻击的目的。 