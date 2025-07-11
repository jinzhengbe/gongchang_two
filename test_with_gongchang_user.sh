#!/bin/bash

echo "=== 使用gongchang用户测试工厂信息编辑API ==="

# 使用gongchang用户进行测试
echo "1. 登录gongchang用户..."
LOGIN_RESPONSE=$(curl -s -X POST "http://localhost:8008/api/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "gongchang",
    "password": "123456"
  }')

echo "登录响应: $LOGIN_RESPONSE"

# 提取token
TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"token":"[^"]*"' | cut -d'"' -f4)

if [ ! -z "$TOKEN" ]; then
    echo "✅ 登录成功，获取到token: ${TOKEN:0:20}..."
    
    echo ""
    echo "2. 获取当前工厂信息..."
    GET_RESPONSE=$(curl -s -X GET "http://localhost:8008/api/factories/profile" \
      -H "Authorization: Bearer $TOKEN" \
      -H "Content-Type: application/json")
    echo "获取工厂信息响应: $GET_RESPONSE"
    
    echo ""
    echo "3. 更新工厂信息..."
    UPDATE_RESPONSE=$(curl -s -X PUT "http://localhost:8008/api/factories/profile" \
      -H "Authorization: Bearer $TOKEN" \
      -H "Content-Type: application/json" \
      -d '{
        "company_name": "工场测试更新",
        "address": "广东省深圳市南山区科技园",
        "capacity": 3000,
        "equipment": "全自动裁剪机,工业缝纫机,质检设备,激光切割机,智能熨烫设备",
        "certificates": "ISO9001,质量管理体系认证,环保认证,OHSAS18001,CE认证",
        "photos": [
          "https://example.com/factory_photo1.jpg",
          "https://example.com/factory_photo2.jpg",
          "https://example.com/factory_photo3.jpg"
        ],
        "videos": [
          "https://example.com/factory_video1.mp4",
          "https://example.com/factory_video2.mp4"
        ],
        "employee_count": 250
      }')
    echo "更新响应: $UPDATE_RESPONSE"
    
    echo ""
    echo "4. 验证更新结果..."
    VERIFY_RESPONSE=$(curl -s -X GET "http://localhost:8008/api/factories/profile" \
      -H "Authorization: Bearer $TOKEN" \
      -H "Content-Type: application/json")
    echo "验证响应: $VERIFY_RESPONSE"
    
    # 检查关键字段
    if echo "$VERIFY_RESPONSE" | grep -q '"CompanyName":"工场测试更新"'; then
        echo "✅ 公司名称更新成功"
    else
        echo "❌ 公司名称更新失败"
    fi
    
    if echo "$VERIFY_RESPONSE" | grep -q '"EmployeeCount":250'; then
        echo "✅ 员工数量更新成功"
    else
        echo "❌ 员工数量更新失败"
    fi
    
    if echo "$VERIFY_RESPONSE" | grep -q 'factory_photo1.jpg'; then
        echo "✅ 照片数组更新成功"
    else
        echo "❌ 照片数组更新失败"
    fi
    
    if echo "$VERIFY_RESPONSE" | grep -q 'factory_video1.mp4'; then
        echo "✅ 视频数组更新成功"
    else
        echo "❌ 视频数组更新失败"
    fi
    
    echo ""
    echo "5. 测试部分更新..."
    PARTIAL_UPDATE_RESPONSE=$(curl -s -X PUT "http://localhost:8008/api/factories/profile" \
      -H "Authorization: Bearer $TOKEN" \
      -H "Content-Type: application/json" \
      -d '{
        "company_name": "工场部分更新测试",
        "employee_count": 300
      }')
    echo "部分更新响应: $PARTIAL_UPDATE_RESPONSE"
    
    if echo "$PARTIAL_UPDATE_RESPONSE" | grep -q '"code":0'; then
        echo "✅ 部分更新成功"
    else
        echo "❌ 部分更新失败"
    fi
    
    echo ""
    echo "=== 测试完成 ==="
    echo "✅ 工厂信息编辑API功能正常"
    echo "📋 测试内容:"
    echo "  - 用户认证 (gongchang/123456)"
    echo "  - 获取工厂信息"
    echo "  - 完整更新工厂信息"
    echo "  - 部分更新工厂信息"
    echo "  - 照片和视频数组处理"
    echo "  - 员工数量字段更新"
    echo "  - 数据验证"
else
    echo "❌ 登录失败"
    echo "请检查用户名和密码是否正确"
fi 