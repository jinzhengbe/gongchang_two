#!/bin/bash

# 设置颜色输出
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}开始发布流程...${NC}"

# 1. 更新开发文档
echo -e "${GREEN}正在更新开发文档...${NC}"
current_date=$(date "+%Y-%m-%d")
echo -e "\n## 更新记录 ($current_date)\n- 更新了开发文档\n- 更新了开发日志\n- 提交了代码更新" >> 开发文档.md

# 2. 更新开发日志
echo -e "${GREEN}正在更新开发日志...${NC}"
echo -e "\n## $current_date\n- 完成了代码更新\n- 更新了相关文档" >> 开发日志.md

# 3. Git 操作
echo -e "${GREEN}正在执行 Git 操作...${NC}"
git add .
git commit -m "docs: 更新开发文档和日志 ($current_date)"
git push

echo -e "${GREEN}发布流程完成！${NC}" 