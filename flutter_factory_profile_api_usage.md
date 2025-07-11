# Flutter 工厂信息API使用指南

## 概述

本文档介绍如何在Flutter前端中使用工厂信息管理API，包括获取和更新工厂详细信息。

## API端点

- `GET /api/factories/profile` - 获取当前用户的工厂详细信息
- `PUT /api/factories/profile` - 更新当前用户的工厂详细信息

## 数据模型

### 工厂信息模型 (FactoryProfile)

```dart
class FactoryProfile {
  final int id;
  final String userID;
  final String companyName;
  final String address;
  final int capacity;
  final String equipment;
  final String certificates;
  final List<String> photos;      // 工厂照片URL数组
  final List<String> videos;      // 工厂视频URL数组
  final int employeeCount;
  final double rating;
  final int status;
  final DateTime createdAt;
  final DateTime updatedAt;

  FactoryProfile({
    required this.id,
    required this.userID,
    required this.companyName,
    required this.address,
    required this.capacity,
    required this.equipment,
    required this.certificates,
    required this.photos,
    required this.videos,
    required this.employeeCount,
    required this.rating,
    required this.status,
    required this.createdAt,
    required this.updatedAt,
  });

  factory FactoryProfile.fromJson(Map<String, dynamic> json) {
    return FactoryProfile(
      id: json['ID'] ?? 0,
      userID: json['UserID'] ?? '',
      companyName: json['CompanyName'] ?? '',
      address: json['Address'] ?? '',
      capacity: json['Capacity'] ?? 0,
      equipment: json['Equipment'] ?? '',
      certificates: json['Certificates'] ?? '',
      photos: _parseJsonArray(json['Photos']),
      videos: _parseJsonArray(json['Videos']),
      employeeCount: json['EmployeeCount'] ?? 0,
      rating: (json['Rating'] ?? 0).toDouble(),
      status: json['Status'] ?? 1,
      createdAt: DateTime.parse(json['CreatedAt']),
      updatedAt: DateTime.parse(json['UpdatedAt']),
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'company_name': companyName,
      'address': address,
      'capacity': capacity,
      'equipment': equipment,
      'certificates': certificates,
      'photos': photos,
      'videos': videos,
      'employee_count': employeeCount,
    };
  }

  static List<String> _parseJsonArray(dynamic json) {
    if (json == null || json == '') return [];
    if (json is String) {
      try {
        List<dynamic> list = jsonDecode(json);
        return list.map((e) => e.toString()).toList();
      } catch (e) {
        return [];
      }
    }
    if (json is List) {
      return json.map((e) => e.toString()).toList();
    }
    return [];
  }
}
```

### 更新请求模型

```dart
class UpdateFactoryProfileRequest {
  final String? companyName;
  final String? address;
  final int? capacity;
  final String? equipment;
  final String? certificates;
  final List<String>? photos;
  final List<String>? videos;
  final int? employeeCount;

  UpdateFactoryProfileRequest({
    this.companyName,
    this.address,
    this.capacity,
    this.equipment,
    this.certificates,
    this.photos,
    this.videos,
    this.employeeCount,
  });

  Map<String, dynamic> toJson() {
    final Map<String, dynamic> data = {};
    if (companyName != null) data['company_name'] = companyName;
    if (address != null) data['address'] = address;
    if (capacity != null) data['capacity'] = capacity;
    if (equipment != null) data['equipment'] = equipment;
    if (certificates != null) data['certificates'] = certificates;
    if (photos != null) data['photos'] = photos;
    if (videos != null) data['videos'] = videos;
    if (employeeCount != null) data['employee_count'] = employeeCount;
    return data;
  }
}
```

## API服务类

```dart
import 'dart:convert';
import 'package:http/http.dart' as http;

class FactoryProfileService {
  static const String baseUrl = 'http://your-api-domain:8008/api';
  static String? _authToken;

  // 设置认证token
  static void setAuthToken(String token) {
    _authToken = token;
  }

  // 获取工厂详细信息
  static Future<FactoryProfile> getFactoryProfile() async {
    try {
      final response = await http.get(
        Uri.parse('$baseUrl/factories/profile'),
        headers: {
          'Content-Type': 'application/json',
          'Authorization': 'Bearer $_authToken',
        },
      );

      if (response.statusCode == 200) {
        final Map<String, dynamic> data = jsonDecode(response.body);
        if (data['code'] == 0) {
          return FactoryProfile.fromJson(data['data']);
        } else {
          throw Exception(data['msg'] ?? '获取工厂信息失败');
        }
      } else {
        throw Exception('HTTP错误: ${response.statusCode}');
      }
    } catch (e) {
      throw Exception('获取工厂信息失败: $e');
    }
  }

  // 更新工厂详细信息
  static Future<FactoryProfile> updateFactoryProfile(
    UpdateFactoryProfileRequest request,
  ) async {
    try {
      final response = await http.put(
        Uri.parse('$baseUrl/factories/profile'),
        headers: {
          'Content-Type': 'application/json',
          'Authorization': 'Bearer $_authToken',
        },
        body: jsonEncode(request.toJson()),
      );

      if (response.statusCode == 200) {
        final Map<String, dynamic> data = jsonDecode(response.body);
        if (data['code'] == 0) {
          return FactoryProfile.fromJson(data['data']);
        } else {
          throw Exception(data['msg'] ?? '更新工厂信息失败');
        }
      } else {
        throw Exception('HTTP错误: ${response.statusCode}');
      }
    } catch (e) {
      throw Exception('更新工厂信息失败: $e');
    }
  }
}
```

## 使用示例

### 1. 获取工厂信息

```dart
class FactoryProfilePage extends StatefulWidget {
  @override
  _FactoryProfilePageState createState() => _FactoryProfilePageState();
}

class _FactoryProfilePageState extends State<FactoryProfilePage> {
  FactoryProfile? factoryProfile;
  bool isLoading = true;
  String? errorMessage;

  @override
  void initState() {
    super.initState();
    _loadFactoryProfile();
  }

  Future<void> _loadFactoryProfile() async {
    try {
      setState(() {
        isLoading = true;
        errorMessage = null;
      });

      final profile = await FactoryProfileService.getFactoryProfile();
      setState(() {
        factoryProfile = profile;
        isLoading = false;
      });
    } catch (e) {
      setState(() {
        errorMessage = e.toString();
        isLoading = false;
      });
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text('工厂信息'),
        actions: [
          IconButton(
            icon: Icon(Icons.edit),
            onPressed: () {
              // 导航到编辑页面
              Navigator.push(
                context,
                MaterialPageRoute(
                  builder: (context) => EditFactoryProfilePage(factoryProfile),
                ),
              );
            },
          ),
        ],
      ),
      body: _buildBody(),
    );
  }

  Widget _buildBody() {
    if (isLoading) {
      return Center(child: CircularProgressIndicator());
    }

    if (errorMessage != null) {
      return Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Text('加载失败: $errorMessage'),
            ElevatedButton(
              onPressed: _loadFactoryProfile,
              child: Text('重试'),
            ),
          ],
        ),
      );
    }

    if (factoryProfile == null) {
      return Center(child: Text('暂无工厂信息'));
    }

    return SingleChildScrollView(
      padding: EdgeInsets.all(16),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          _buildInfoCard('公司名称', factoryProfile!.companyName),
          _buildInfoCard('地址', factoryProfile!.address),
          _buildInfoCard('产能', '${factoryProfile!.capacity} 件/月'),
          _buildInfoCard('设备', factoryProfile!.equipment),
          _buildInfoCard('证书', factoryProfile!.certificates),
          _buildInfoCard('员工数量', '${factoryProfile!.employeeCount} 人'),
          _buildInfoCard('评分', '${factoryProfile!.rating}'),
          _buildPhotosSection(),
          _buildVideosSection(),
        ],
      ),
    );
  }

  Widget _buildInfoCard(String title, String content) {
    return Card(
      margin: EdgeInsets.only(bottom: 8),
      child: Padding(
        padding: EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              title,
              style: TextStyle(
                fontWeight: FontWeight.bold,
                fontSize: 16,
              ),
            ),
            SizedBox(height: 8),
            Text(content.isEmpty ? '暂无信息' : content),
          ],
        ),
      ),
    );
  }

  Widget _buildPhotosSection() {
    if (factoryProfile!.photos.isEmpty) {
      return SizedBox.shrink();
    }

    return Card(
      margin: EdgeInsets.only(bottom: 8),
      child: Padding(
        padding: EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              '工厂照片',
              style: TextStyle(
                fontWeight: FontWeight.bold,
                fontSize: 16,
              ),
            ),
            SizedBox(height: 8),
            SizedBox(
              height: 120,
              child: ListView.builder(
                scrollDirection: Axis.horizontal,
                itemCount: factoryProfile!.photos.length,
                itemBuilder: (context, index) {
                  return Container(
                    margin: EdgeInsets.only(right: 8),
                    child: ClipRRect(
                      borderRadius: BorderRadius.circular(8),
                      child: Image.network(
                        factoryProfile!.photos[index],
                        width: 120,
                        height: 120,
                        fit: BoxFit.cover,
                        errorBuilder: (context, error, stackTrace) {
                          return Container(
                            width: 120,
                            height: 120,
                            color: Colors.grey[300],
                            child: Icon(Icons.image_not_supported),
                          );
                        },
                      ),
                    ),
                  );
                },
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildVideosSection() {
    if (factoryProfile!.videos.isEmpty) {
      return SizedBox.shrink();
    }

    return Card(
      margin: EdgeInsets.only(bottom: 8),
      child: Padding(
        padding: EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              '工厂视频',
              style: TextStyle(
                fontWeight: FontWeight.bold,
                fontSize: 16,
              ),
            ),
            SizedBox(height: 8),
            SizedBox(
              height: 120,
              child: ListView.builder(
                scrollDirection: Axis.horizontal,
                itemCount: factoryProfile!.videos.length,
                itemBuilder: (context, index) {
                  return Container(
                    margin: EdgeInsets.only(right: 8),
                    child: ClipRRect(
                      borderRadius: BorderRadius.circular(8),
                      child: Container(
                        width: 120,
                        height: 120,
                        color: Colors.black,
                        child: Center(
                          child: Icon(
                            Icons.play_arrow,
                            color: Colors.white,
                            size: 48,
                          ),
                        ),
                      ),
                    ),
                  );
                },
              ),
            ),
          ],
        ),
      ),
    );
  }
}
```

### 2. 编辑工厂信息

```dart
class EditFactoryProfilePage extends StatefulWidget {
  final FactoryProfile? factoryProfile;

  EditFactoryProfilePage(this.factoryProfile);

  @override
  _EditFactoryProfilePageState createState() => _EditFactoryProfilePageState();
}

class _EditFactoryProfilePageState extends State<EditFactoryProfilePage> {
  final _formKey = GlobalKey<FormState>();
  late TextEditingController _companyNameController;
  late TextEditingController _addressController;
  late TextEditingController _capacityController;
  late TextEditingController _equipmentController;
  late TextEditingController _certificatesController;
  late TextEditingController _employeeCountController;
  List<String> _photos = [];
  List<String> _videos = [];
  bool _isLoading = false;

  @override
  void initState() {
    super.initState();
    _initializeControllers();
  }

  void _initializeControllers() {
    final profile = widget.factoryProfile;
    _companyNameController = TextEditingController(text: profile?.companyName ?? '');
    _addressController = TextEditingController(text: profile?.address ?? '');
    _capacityController = TextEditingController(text: profile?.capacity.toString() ?? '');
    _equipmentController = TextEditingController(text: profile?.equipment ?? '');
    _certificatesController = TextEditingController(text: profile?.certificates ?? '');
    _employeeCountController = TextEditingController(text: profile?.employeeCount.toString() ?? '');
    _photos = List.from(profile?.photos ?? []);
    _videos = List.from(profile?.videos ?? []);
  }

  @override
  void dispose() {
    _companyNameController.dispose();
    _addressController.dispose();
    _capacityController.dispose();
    _equipmentController.dispose();
    _certificatesController.dispose();
    _employeeCountController.dispose();
    super.dispose();
  }

  Future<void> _saveProfile() async {
    if (!_formKey.currentState!.validate()) return;

    setState(() {
      _isLoading = true;
    });

    try {
      final request = UpdateFactoryProfileRequest(
        companyName: _companyNameController.text,
        address: _addressController.text,
        capacity: int.tryParse(_capacityController.text),
        equipment: _equipmentController.text,
        certificates: _certificatesController.text,
        photos: _photos,
        videos: _videos,
        employeeCount: int.tryParse(_employeeCountController.text),
      );

      await FactoryProfileService.updateFactoryProfile(request);

      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('工厂信息更新成功')),
      );

      Navigator.pop(context, true);
    } catch (e) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('更新失败: $e')),
      );
    } finally {
      setState(() {
        _isLoading = false;
      });
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text('编辑工厂信息'),
        actions: [
          if (_isLoading)
            Padding(
              padding: EdgeInsets.all(16),
              child: SizedBox(
                width: 20,
                height: 20,
                child: CircularProgressIndicator(strokeWidth: 2),
              ),
            )
          else
            TextButton(
              onPressed: _saveProfile,
              child: Text('保存'),
            ),
        ],
      ),
      body: Form(
        key: _formKey,
        child: SingleChildScrollView(
          padding: EdgeInsets.all(16),
          child: Column(
            children: [
              TextFormField(
                controller: _companyNameController,
                decoration: InputDecoration(
                  labelText: '公司名称',
                  border: OutlineInputBorder(),
                ),
                validator: (value) {
                  if (value?.isEmpty ?? true) {
                    return '请输入公司名称';
                  }
                  return null;
                },
              ),
              SizedBox(height: 16),
              TextFormField(
                controller: _addressController,
                decoration: InputDecoration(
                  labelText: '地址',
                  border: OutlineInputBorder(),
                ),
                validator: (value) {
                  if (value?.isEmpty ?? true) {
                    return '请输入地址';
                  }
                  return null;
                },
              ),
              SizedBox(height: 16),
              TextFormField(
                controller: _capacityController,
                decoration: InputDecoration(
                  labelText: '产能 (件/月)',
                  border: OutlineInputBorder(),
                ),
                keyboardType: TextInputType.number,
                validator: (value) {
                  if (value?.isEmpty ?? true) {
                    return '请输入产能';
                  }
                  if (int.tryParse(value!) == null) {
                    return '请输入有效的数字';
                  }
                  return null;
                },
              ),
              SizedBox(height: 16),
              TextFormField(
                controller: _equipmentController,
                decoration: InputDecoration(
                  labelText: '设备',
                  border: OutlineInputBorder(),
                ),
                maxLines: 3,
              ),
              SizedBox(height: 16),
              TextFormField(
                controller: _certificatesController,
                decoration: InputDecoration(
                  labelText: '证书',
                  border: OutlineInputBorder(),
                ),
                maxLines: 3,
              ),
              SizedBox(height: 16),
              TextFormField(
                controller: _employeeCountController,
                decoration: InputDecoration(
                  labelText: '员工数量',
                  border: OutlineInputBorder(),
                ),
                keyboardType: TextInputType.number,
                validator: (value) {
                  if (value?.isEmpty ?? true) {
                    return '请输入员工数量';
                  }
                  if (int.tryParse(value!) == null) {
                    return '请输入有效的数字';
                  }
                  return null;
                },
              ),
              SizedBox(height: 16),
              _buildMediaSection('照片', _photos, (url) {
                setState(() {
                  _photos.add(url);
                });
              }),
              SizedBox(height: 16),
              _buildMediaSection('视频', _videos, (url) {
                setState(() {
                  _videos.add(url);
                });
              }),
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildMediaSection(String title, List<String> media, Function(String) onAdd) {
    return Card(
      child: Padding(
        padding: EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                Text(
                  title,
                  style: TextStyle(
                    fontWeight: FontWeight.bold,
                    fontSize: 16,
                  ),
                ),
                TextButton(
                  onPressed: () => _showAddMediaDialog(title, onAdd),
                  child: Text('添加'),
                ),
              ],
            ),
            if (media.isNotEmpty)
              SizedBox(
                height: 100,
                child: ListView.builder(
                  scrollDirection: Axis.horizontal,
                  itemCount: media.length,
                  itemBuilder: (context, index) {
                    return Container(
                      margin: EdgeInsets.only(right: 8),
                      child: Stack(
                        children: [
                          ClipRRect(
                            borderRadius: BorderRadius.circular(8),
                            child: Container(
                              width: 100,
                              height: 100,
                              color: Colors.grey[300],
                              child: Icon(Icons.image),
                            ),
                          ),
                          Positioned(
                            top: 4,
                            right: 4,
                            child: GestureDetector(
                              onTap: () {
                                setState(() {
                                  media.removeAt(index);
                                });
                              },
                              child: Container(
                                padding: EdgeInsets.all(4),
                                decoration: BoxDecoration(
                                  color: Colors.red,
                                  shape: BoxShape.circle,
                                ),
                                child: Icon(
                                  Icons.close,
                                  color: Colors.white,
                                  size: 16,
                                ),
                              ),
                            ),
                          ),
                        ],
                      ),
                    );
                  },
                ),
              ),
          ],
        ),
      ),
    );
  }

  void _showAddMediaDialog(String title, Function(String) onAdd) {
    final controller = TextEditingController();
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: Text('添加$title'),
        content: TextField(
          controller: controller,
          decoration: InputDecoration(
            labelText: 'URL',
            hintText: '请输入$titleURL',
          ),
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: Text('取消'),
          ),
          TextButton(
            onPressed: () {
              if (controller.text.isNotEmpty) {
                onAdd(controller.text);
                Navigator.pop(context);
              }
            },
            child: Text('确定'),
          ),
        ],
      ),
    );
  }
}
```

## 使用步骤

1. **添加依赖**
   在 `pubspec.yaml` 中添加：
   ```yaml
   dependencies:
     http: ^1.1.0
   ```

2. **设置认证token**
   ```dart
   // 在登录成功后设置token
   FactoryProfileService.setAuthToken('your-jwt-token');
   ```

3. **使用API**
   ```dart
   // 获取工厂信息
   final profile = await FactoryProfileService.getFactoryProfile();
   
   // 更新工厂信息
   final request = UpdateFactoryProfileRequest(
     companyName: '新工厂名称',
     address: '新地址',
     capacity: 1000,
     employeeCount: 50,
   );
   final updatedProfile = await FactoryProfileService.updateFactoryProfile(request);
   ```

## 注意事项

1. **认证**: 确保在调用API前已设置正确的认证token
2. **错误处理**: 所有API调用都应该包含适当的错误处理
3. **数据验证**: 在提交数据前验证用户输入
4. **网络状态**: 考虑网络连接状态和超时处理
5. **图片/视频**: 照片和视频URL应该是可访问的完整URL

## 完整示例项目

完整的Flutter示例项目包含：
- 工厂信息展示页面
- 工厂信息编辑页面
- API服务类
- 数据模型
- 错误处理
- 加载状态管理

这个实现提供了完整的工厂信息管理功能，包括所有新增的字段（照片、视频、员工数量等）。 