import 'dart:io';
import 'dart:convert';

/// 修复国际化文件中的键名，将点号(.)替换为下划线(_)
void main() async {
  final translationsDir = Directory('assets/assets/translations');
  
  if (!await translationsDir.exists()) {
    print('错误: 目录不存在 ${translationsDir.path}');
    return;
  }
  
  final files = await translationsDir.list().where((f) => f.path.endsWith('.json')).toList();
  
  for (final file in files) {
    await fixI18nKeys(file.path);
  }
  
  print('\n所有国际化文件修复完成！');
}

Future<void> fixI18nKeys(String filePath) async {
  print('正在处理文件: $filePath');
  
  // 读取文件
  final file = File(filePath);
  final content = await file.readAsString();
  
  // 使用正则表达式找到所有带点号的键名
  // 匹配模式: "key.with.dots": "value"
  final pattern = RegExp(r'"([^"]*\.[^"]*)"\s*:\s*"([^"]*)"');
  
  String newContent = content.replaceAllMapped(pattern, (match) {
    final oldKey = match.group(1)!;
    final value = match.group(2)!;
    final newKey = oldKey.replaceAll('.', '_');
    print('  修复键名: \'$oldKey\' -> \'$newKey\'');
    return '"$newKey": "$value"';
  });
  
  // 写回文件
  await file.writeAsString(newContent);
  
  print('  完成修复: $filePath');
} 