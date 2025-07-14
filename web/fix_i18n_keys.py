#!/usr/bin/env python3
"""
修复国际化文件中的键名，将点号(.)替换为下划线(_)
"""

import json
import os
import re

def fix_i18n_keys(file_path):
    """修复单个国际化文件中的键名"""
    print(f"正在处理文件: {file_path}")
    
    # 读取文件
    with open(file_path, 'r', encoding='utf-8') as f:
        content = f.read()
    
    # 使用正则表达式找到所有带点号的键名
    # 匹配模式: "key.with.dots": "value"
    pattern = r'"([^"]*\.[^"]*)"\s*:\s*"([^"]*)"'
    
    def replace_key(match):
        old_key = match.group(1)
        value = match.group(2)
        new_key = old_key.replace('.', '_')
        print(f"  修复键名: '{old_key}' -> '{new_key}'")
        return f'"{new_key}": "{value}"'
    
    # 执行替换
    new_content = re.sub(pattern, replace_key, content)
    
    # 写回文件
    with open(file_path, 'w', encoding='utf-8') as f:
        f.write(new_content)
    
    print(f"  完成修复: {file_path}")

def main():
    """主函数"""
    translations_dir = "assets/assets/translations"
    
    # 检查目录是否存在
    if not os.path.exists(translations_dir):
        print(f"错误: 目录不存在 {translations_dir}")
        return
    
    # 处理所有翻译文件
    for filename in os.listdir(translations_dir):
        if filename.endswith('.json'):
            file_path = os.path.join(translations_dir, filename)
            fix_i18n_keys(file_path)
    
    print("\n所有国际化文件修复完成！")

if __name__ == "__main__":
    main() 