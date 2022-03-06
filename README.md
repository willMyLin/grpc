Client发送完成后需要手动调用Close()或者CloseSend()方法关闭stream，Server端则return nil就会自动 Close。

1）ServerStream

服务端处理完成后return nil代表响应完成
客户端通过 err == io.EOF判断服务端是否响应完成
2）ClientStream

客户端发送完毕通过`CloseAndRecv关闭stream 并接收服务端响应
服务端通过 err == io.EOF判断客户端是否发送完毕，完毕后使用SendAndClose关闭 stream并返回响应。
3）服务端和客户端都为stream

客户端服务端都通过stream向对方推送数据
客户端推送完成后通过CloseSend关闭流，通过err == io.EOF`判断服务端是否响应完成
服务端通过err == io.EOF判断客户端是否响应完成,通过return nil表示已经完成响应
通过err == io.EOF来判定是否把对方推送的数据全部获取到了。

客户端通过CloseAndRecv或者CloseSend关闭 Stream，服务端则通过SendAndClose或者直接 return nil来返回响应。

