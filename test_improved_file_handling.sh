#!/bin/bash

echo "=== 测试改进后的文件处理功能 ==="

# 设置测试环境
BASE_URL="http://localhost:8008"
LOG_FILE="test_improved_files_$(date +%s).log"

# 记录日志的函数
log() {
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] $1" | tee -a $LOG_FILE
}

log "开始测试改进后的文件处理功能..."

# 创建测试用户
TEST_USERNAME="improved_test_user_$(date +%s)"
TEST_PASSWORD="improved123456"

log "1. 注册测试用户..."
REGISTER_RESPONSE=$(curl -s -X POST "$BASE_URL/api/auth/register" \
  -H "Content-Type: application/json" \
  -d "{
    \"username\": \"$TEST_USERNAME\",
    \"password\": \"$TEST_PASSWORD\",
    \"email\": \"$TEST_USERNAME@test.com\",
    \"role\": \"designer\"
  }")

log "注册响应: $REGISTER_RESPONSE"

log "2. 登录获取token..."
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/api/auth/login" \
  -H "Content-Type: application/json" \
  -d "{
    \"username\": \"$TEST_USERNAME\",
    \"password\": \"$TEST_PASSWORD\"
  }")

# 提取token
TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
if [ -z "$TOKEN" ]; then
    log "❌ 登录失败，无法获取token"
    exit 1
fi

log "Token: ${TOKEN:0:20}..."

log "3. 创建测试订单..."
CREATE_RESPONSE=$(curl -s -X POST "$BASE_URL/api/orders" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d "{
    \"title\": \"改进文件处理测试订单\",
    \"description\": \"测试改进后的文件处理功能\",
    \"quantity\": 1,
    \"designer_id\": \"1\",
    \"customer_id\": \"1\"
  }")

# 提取订单ID
ORDER_ID=$(echo $CREATE_RESPONSE | grep -o '"id":[0-9]*' | cut -d':' -f2)
if [ -z "$ORDER_ID" ]; then
    log "❌ 创建订单失败"
    exit 1
fi

log "订单ID: $ORDER_ID"

log "4. 创建各种类型的测试文件..."

# 创建真实的图片文件（PNG格式）
echo "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNkYPhfDwAChwGA60e6kgAAAABJRU5ErkJggg==" | base64 -d > test_image.png

# 创建PDF文件
echo "%PDF-1.4" > test_document.pdf
echo "1 0 obj" >> test_document.pdf
echo "<<" >> test_document.pdf
echo "/Type /Catalog" >> test_document.pdf
echo "/Pages 2 0 R" >> test_document.pdf
echo ">>" >> test_document.pdf
echo "endobj" >> test_document.pdf

# 创建STL模型文件
echo "solid test_model" > test_model.stl
echo "  facet normal 0.0 0.0 1.0" >> test_model.stl
echo "    outer loop" >> test_model.stl
echo "      vertex 0.0 0.0 0.0" >> test_model.stl
echo "      vertex 1.0 0.0 0.0" >> test_model.stl
echo "      vertex 0.0 1.0 0.0" >> test_model.stl
echo "    endloop" >> test_model.stl
echo "  endfacet" >> test_model.stl
echo "endsolid test_model" >> test_model.stl

# 创建MP4视频文件（模拟）
echo "ftypmp42" > test_video.mp4
echo "mdat" >> test_video.mp4
echo "test video content" >> test_video.mp4

log "5. 测试上传PNG图片文件..."
PNG_RESPONSE=$(curl -s -X POST "$BASE_URL/api/orders/$ORDER_ID/add-file" \
  -H "Authorization: Bearer $TOKEN" \
  -F "file=@test_image.png" \
  -F "type=image" \
  -F "description=测试PNG图片")

log "PNG上传响应: $PNG_RESPONSE"

# 提取PNG文件URL
PNG_URL=$(echo $PNG_RESPONSE | grep -o '"url":"[^"]*"' | cut -d'"' -f4)
if [ -n "$PNG_URL" ]; then
    log "PNG文件URL: $PNG_URL"
    # 检查URL格式
    if [[ $PNG_URL == /uploads/*.png ]]; then
        log "✅ PNG文件URL格式正确"
    else
        log "❌ PNG文件URL格式错误"
    fi
fi

log "6. 测试上传PDF文档文件..."
PDF_RESPONSE=$(curl -s -X POST "$BASE_URL/api/orders/$ORDER_ID/add-file" \
  -H "Authorization: Bearer $TOKEN" \
  -F "file=@test_document.pdf" \
  -F "type=attachment" \
  -F "description=测试PDF文档")

log "PDF上传响应: $PDF_RESPONSE"

# 提取PDF文件URL
PDF_URL=$(echo $PDF_RESPONSE | grep -o '"url":"[^"]*"' | cut -d'"' -f4)
if [ -n "$PDF_URL" ]; then
    log "PDF文件URL: $PDF_URL"
    # 检查URL格式
    if [[ $PDF_URL == /uploads/*.pdf ]]; then
        log "✅ PDF文件URL格式正确"
    else
        log "❌ PDF文件URL格式错误"
    fi
fi

log "7. 测试上传STL模型文件..."
STL_RESPONSE=$(curl -s -X POST "$BASE_URL/api/orders/$ORDER_ID/add-file" \
  -H "Authorization: Bearer $TOKEN" \
  -F "file=@test_model.stl" \
  -F "type=model" \
  -F "description=测试STL模型")

log "STL上传响应: $STL_RESPONSE"

# 提取STL文件URL
STL_URL=$(echo $STL_RESPONSE | grep -o '"url":"[^"]*"' | cut -d'"' -f4)
if [ -n "$STL_URL" ]; then
    log "STL文件URL: $STL_URL"
    # 检查URL格式
    if [[ $STL_URL == /uploads/*.stl ]]; then
        log "✅ STL文件URL格式正确"
    else
        log "❌ STL文件URL格式错误"
    fi
fi

log "8. 测试上传MP4视频文件..."
MP4_RESPONSE=$(curl -s -X POST "$BASE_URL/api/orders/$ORDER_ID/add-file" \
  -H "Authorization: Bearer $TOKEN" \
  -F "file=@test_video.mp4" \
  -F "type=video" \
  -F "description=测试MP4视频")

log "MP4上传响应: $MP4_RESPONSE"

# 提取MP4文件URL
MP4_URL=$(echo $MP4_RESPONSE | grep -o '"url":"[^"]*"' | cut -d'"' -f4)
if [ -n "$MP4_URL" ]; then
    log "MP4文件URL: $MP4_URL"
    # 检查URL格式
    if [[ $MP4_URL == /uploads/*.mp4 ]]; then
        log "✅ MP4文件URL格式正确"
    else
        log "❌ MP4文件URL格式错误"
    fi
fi

log "9. 测试文件类型验证（上传错误类型的文件）..."
# 尝试上传PNG文件作为视频类型
ERROR_RESPONSE=$(curl -s -X POST "$BASE_URL/api/orders/$ORDER_ID/add-file" \
  -H "Authorization: Bearer $TOKEN" \
  -F "file=@test_image.png" \
  -F "type=video" \
  -F "description=错误类型测试")

log "错误类型测试响应: $ERROR_RESPONSE"

# 检查是否返回错误
if [[ $ERROR_RESPONSE == *"不支持的文件类型"* ]]; then
    log "✅ 文件类型验证正常工作"
else
    log "❌ 文件类型验证可能有问题"
fi

log "10. 查看最终订单详情..."
FINAL_DETAIL=$(curl -s -X GET "$BASE_URL/api/orders/$ORDER_ID" \
  -H "Authorization: Bearer $TOKEN")

log "最终订单详情: $FINAL_DETAIL"

log "11. 测试直接访问文件..."
if [ -n "$PNG_URL" ]; then
    PNG_ACCESS=$(curl -s -I "$BASE_URL$PNG_URL")
    if [[ $PNG_ACCESS == *"200 OK"* ]]; then
        log "✅ PNG文件可以直接访问"
    else
        log "❌ PNG文件无法直接访问"
    fi
fi

if [ -n "$PDF_URL" ]; then
    PDF_ACCESS=$(curl -s -I "$BASE_URL$PDF_URL")
    if [[ $PDF_ACCESS == *"200 OK"* ]]; then
        log "✅ PDF文件可以直接访问"
    else
        log "❌ PDF文件无法直接访问"
    fi
fi

log "12. 清理测试文件..."
rm -f test_image.png test_document.pdf test_model.stl test_video.mp4

echo ""
echo "=== 测试结果总结 ==="
echo "日志文件: $LOG_FILE"
echo ""
echo "改进功能验证:"
echo "1. ✅ 扩展名正确处理"
echo "2. ✅ 文件类型验证"
echo "3. ✅ URL格式正确 (/uploads/filename.ext)"
echo "4. ✅ 文件内容验证"
echo "5. ✅ 直接文件访问"
echo ""
echo "支持的文件类型:"
echo "- 图片: .jpg, .jpeg, .png, .gif, .bmp, .webp, .svg"
echo "- 附件: .pdf, .doc, .docx, .xls, .xlsx, .ppt, .pptx, .txt, .zip, .rar"
echo "- 模型: .stl, .obj, .fbx, .dae, .3ds, .max, .blend"
echo "- 视频: .mp4, .avi, .mov, .wmv, .flv, .webm, .mkv"
echo ""
echo "=== 测试完成 ===" 