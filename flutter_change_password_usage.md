# Flutter 更改密码API使用说明

## 1. API信息

- **接口**: `POST /api/users/change-password`
- **认证**: Bearer Token
- **请求体**:
```json
{
  "old_password": "当前密码",
  "new_password": "新密码"
}
```

## 2. 简单使用示例

### 2.1 基本API调用

```dart
import 'dart:convert';
import 'package:http/http.dart' as http;

// 更改密码
Future<Map<String, dynamic>> changePassword({
  required String token,
  required String oldPassword,
  required String newPassword,
}) async {
  final response = await http.post(
    Uri.parse('http://localhost:8008/api/users/change-password'),
    headers: {
      'Content-Type': 'application/json',
      'Authorization': 'Bearer $token',
    },
    body: jsonEncode({
      'old_password': oldPassword,
      'new_password': newPassword,
    }),
  );

  final data = jsonDecode(response.body);
  
  if (response.statusCode == 200) {
    return {'success': true, 'message': data['message']};
  } else {
    return {'success': false, 'message': data['error']};
  }
}
```

### 2.2 使用示例

```dart
// 调用API
final result = await changePassword(
  token: 'your-jwt-token',
  oldPassword: '123456',
  newPassword: 'newpassword123',
);

// 处理结果
if (result['success']) {
  print('密码修改成功');
} else {
  print('密码修改失败: ${result['message']}');
}
```

## 3. 常见错误

- **400**: 旧密码错误、新密码太短
- **401**: 未授权，token无效
- **500**: 服务器错误

## 4. 注意事项

1. 新密码长度至少6位
2. 需要有效的JWT token
3. 使用HTTPS进行安全传输 