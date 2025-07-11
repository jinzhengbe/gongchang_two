# Flutter 工厂信息API用法

## API端点

- `GET /api/factories/profile` - 获取当前用户的工厂详细信息
- `PUT /api/factories/profile` - 更新当前用户的工厂详细信息

## 数据模型

### 工厂信息字段

```dart
class FactoryProfile {
  final int id;
  final String userID;
  final String companyName;      // 公司名称
  final String address;          // 地址
  final int capacity;            // 产能
  final String equipment;        // 设备
  final String certificates;     // 证书
  final List<String> photos;     // 工厂照片URL数组
  final List<String> videos;     // 工厂视频URL数组
  final int employeeCount;       // 员工数量
  final double rating;           // 评分
  final int status;              // 状态
  final DateTime createdAt;
  final DateTime updatedAt;
}
```

### 更新请求字段

```dart
class UpdateFactoryProfileRequest {
  final String? companyName;     // 公司名称
  final String? address;         // 地址
  final int? capacity;           // 产能
  final String? equipment;       // 设备
  final String? certificates;    // 证书
  final List<String>? photos;    // 照片URL数组
  final List<String>? videos;    // 视频URL数组
  final int? employeeCount;      // 员工数量
}
```

## API调用方法

### 1. 获取工厂信息

```dart
Future<FactoryProfile> getFactoryProfile() async {
  final response = await http.get(
    Uri.parse('$baseUrl/factories/profile'),
    headers: {
      'Content-Type': 'application/json',
      'Authorization': 'Bearer $token',
    },
  );

  if (response.statusCode == 200) {
    final data = jsonDecode(response.body);
    if (data['code'] == 0) {
      return FactoryProfile.fromJson(data['data']);
    }
  }
  throw Exception('获取工厂信息失败');
}
```

### 2. 更新工厂信息

```dart
Future<FactoryProfile> updateFactoryProfile(UpdateFactoryProfileRequest request) async {
  final response = await http.put(
    Uri.parse('$baseUrl/factories/profile'),
    headers: {
      'Content-Type': 'application/json',
      'Authorization': 'Bearer $token',
    },
    body: jsonEncode(request.toJson()),
  );

  if (response.statusCode == 200) {
    final data = jsonDecode(response.body);
    if (data['code'] == 0) {
      return FactoryProfile.fromJson(data['data']);
    }
  }
  throw Exception('更新工厂信息失败');
}
```

## 使用示例

### 获取工厂信息
```dart
try {
  final profile = await getFactoryProfile();
  print('公司名称: ${profile.companyName}');
  print('员工数量: ${profile.employeeCount}');
  print('照片数量: ${profile.photos.length}');
  print('视频数量: ${profile.videos.length}');
} catch (e) {
  print('获取失败: $e');
}
```

### 更新工厂信息
```dart
try {
  final request = UpdateFactoryProfileRequest(
    companyName: '新工厂名称',
    address: '新地址',
    capacity: 1000,
    equipment: '新设备',
    certificates: '新证书',
    photos: ['https://example.com/photo1.jpg'],
    videos: ['https://example.com/video1.mp4'],
    employeeCount: 50,
  );
  
  final updatedProfile = await updateFactoryProfile(request);
  print('更新成功');
} catch (e) {
  print('更新失败: $e');
}
```

## 注意事项

1. **认证**: 需要在请求头中包含 `Authorization: Bearer $token`
2. **照片/视频**: 必须是完整的URL地址
3. **可选字段**: 更新时只传递需要修改的字段
4. **错误处理**: 所有API调用都应该包含try-catch错误处理
5. **数据格式**: 照片和视频在数据库中存储为JSON字符串格式 