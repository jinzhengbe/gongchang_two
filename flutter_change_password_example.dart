import 'dart:convert';
import 'package:http/http.dart' as http;
import 'package:flutter/material.dart';

// 简单的更改密码API调用示例
class ChangePasswordAPI {
  static const String baseUrl = 'http://localhost:8008/api';
  
  // 更改密码
  static Future<Map<String, dynamic>> changePassword({
    required String token,
    required String oldPassword,
    required String newPassword,
  }) async {
    try {
      final response = await http.post(
        Uri.parse('$baseUrl/users/change-password'),
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
        return {
          'success': true,
          'message': data['message'] ?? '密码修改成功',
        };
      } else {
        return {
          'success': false,
          'message': data['error'] ?? '密码修改失败',
        };
      }
    } catch (e) {
      return {
        'success': false,
        'message': '网络错误: $e',
      };
    }
  }
}

// 简单的更改密码页面
class ChangePasswordPage extends StatefulWidget {
  @override
  _ChangePasswordPageState createState() => _ChangePasswordPageState();
}

class _ChangePasswordPageState extends State<ChangePasswordPage> {
  final _oldPasswordController = TextEditingController();
  final _newPasswordController = TextEditingController();
  bool _isLoading = false;
  String? _message;

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: Text('修改密码')),
      body: Padding(
        padding: EdgeInsets.all(16.0),
        child: Column(
          children: [
            // 旧密码
            TextField(
              controller: _oldPasswordController,
              obscureText: true,
              decoration: InputDecoration(
                labelText: '当前密码',
                border: OutlineInputBorder(),
              ),
            ),
            SizedBox(height: 16),
            
            // 新密码
            TextField(
              controller: _newPasswordController,
              obscureText: true,
              decoration: InputDecoration(
                labelText: '新密码',
                border: OutlineInputBorder(),
              ),
            ),
            SizedBox(height: 16),
            
            // 消息显示
            if (_message != null)
              Container(
                width: double.infinity,
                padding: EdgeInsets.all(8),
                decoration: BoxDecoration(
                  color: _message!.contains('成功') ? Colors.green[100] : Colors.red[100],
                  borderRadius: BorderRadius.circular(4),
                ),
                child: Text(_message!),
              ),
            SizedBox(height: 16),
            
            // 提交按钮
            SizedBox(
              width: double.infinity,
              child: ElevatedButton(
                onPressed: _isLoading ? null : _changePassword,
                child: _isLoading 
                    ? CircularProgressIndicator(color: Colors.white)
                    : Text('确认修改'),
              ),
            ),
          ],
        ),
      ),
    );
  }

  Future<void> _changePassword() async {
    if (_oldPasswordController.text.isEmpty || _newPasswordController.text.isEmpty) {
      setState(() {
        _message = '请填写完整信息';
      });
      return;
    }

    setState(() {
      _isLoading = true;
      _message = null;
    });

    // 这里需要替换为实际的token
    final token = 'your-jwt-token-here';
    
    final result = await ChangePasswordAPI.changePassword(
      token: token,
      oldPassword: _oldPasswordController.text,
      newPassword: _newPasswordController.text,
    );

    setState(() {
      _isLoading = false;
      _message = result['message'];
    });

    if (result['success']) {
      _oldPasswordController.clear();
      _newPasswordController.clear();
    }
  }

  @override
  void dispose() {
    _oldPasswordController.dispose();
    _newPasswordController.dispose();
    super.dispose();
  }
}

// 使用示例
void main() {
  runApp(MaterialApp(
    home: ChangePasswordPage(),
  ));
}

// 简单的API调用示例
void exampleUsage() async {
  // 1. 获取token（从登录后保存的token）
  String token = 'your-jwt-token';
  
  // 2. 调用更改密码API
  final result = await ChangePasswordAPI.changePassword(
    token: token,
    oldPassword: '123456',
    newPassword: 'newpassword123',
  );
  
  // 3. 处理结果
  if (result['success']) {
    print('密码修改成功: ${result['message']}');
  } else {
    print('密码修改失败: ${result['message']}');
  }
} 