1. etcd在网关中是具体用来做什么的？为什么只取backends\[0\]作为machine？
2. dns srv记录具体是用来做什么的？本地如何添加srv记录？
3. consul是来监听bloomfilter rpc server 的吗？
4. rotate bloomfilter 比 普通的 bloomfilter好在哪里？为什么这么设计？有什么是rotate可以完成，而普通的不能完成的？
5. rpc input 为什么要用 [][]byte 而不是 []byte？