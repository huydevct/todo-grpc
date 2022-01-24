# Todo-grpc

## gRPC server
Entry file: cmd/server/main.go

## câu 1
Thêm trường status có kiểu INTEGER với 3 giá trị 1,2,3 đại diện cho 3 trạng thái TODO, DOING, DONE.
Khi call api sẽ truyền 1 param status có value là todo, doing hoặc done. Từ đó, ta có biến todo=1, query vào DB `SELECT * FROM list_task WHERE status = 1` 
Tương tự với doing, done.
## câu 2
Trong DB ta tạo thêm 1 bảng gồm title_id và tags. Khi call api, thêm 1 param tags gồm nhiều giá trị ngăn cách bằng dấu cách. Backend nhận chuỗi, tách chuỗi lấy các tag và select vào DB để lấy ra các title.

## câu 3.a 
Khi call api ta nhận 1 param mới là page. Backend query vào DB: `SELECT * FROM todo LIMIT 10 OFFSET ?` truyền page vào ? trong query.

## câu 3.b
Khi call api, ta nhận 2 param với chuỗi string là status và tags. Backend nhận vào chuỗi, tách chuỗi, sử dụng kết quả sau khi tách chuỗi để query vào DB như sau:
`SELECT * FROM list_task WHERE status = ? AND (tags = ? OR tags = ?)` 

## câu 4
Mỗi bảng ghi có trường index tự tăng, khi đảo các bản ghi thì swap index. Khi backend query ta order theo index.

### from huydevct