# Flutter 订单搜索API使用指南

## 概述

本文档介绍如何在Flutter应用中集成订单搜索功能，包括高级搜索、搜索建议和统计功能。

## 1. 订单搜索服务类

### 创建搜索服务类

```dart
class OrderSearchService {
  final String baseUrl;
  final String token;
  
  OrderSearchService({required this.baseUrl, required this.token});
  
  // 基础搜索
  Future<OrderSearchResponse> searchOrders({
    String? query,
    String? status,
    String? startDate,
    String? endDate,
    int page = 1,
    int pageSize = 20,
    String sortBy = 'created_at',
    String sortOrder = 'desc',
  }) async {
    // 实现搜索逻辑
  }
  
  // 获取搜索建议
  Future<SearchSuggestionResponse> getSearchSuggestions({
    required String query,
    int limit = 10,
  }) async {
    // 实现建议逻辑
  }
  
  // 获取搜索统计
  Future<SearchStatisticsResponse> getSearchStatistics() async {
    // 实现统计逻辑
  }
}
```

## 2. 数据模型

### 搜索请求模型

```dart
class OrderSearchRequest {
  final String? query;
  final String? status;
  final String? startDate;
  final String? endDate;
  final int page;
  final int pageSize;
  final String sortBy;
  final String sortOrder;
  
  OrderSearchRequest({
    this.query,
    this.status,
    this.startDate,
    this.endDate,
    this.page = 1,
    this.pageSize = 20,
    this.sortBy = 'created_at',
    this.sortOrder = 'desc',
  });
}
```

### 搜索响应模型

```dart
class OrderSearchResponse {
  final bool success;
  final OrderSearchData data;
  
  OrderSearchResponse({
    required this.success,
    required this.data,
  });
}

class OrderSearchData {
  final List<OrderSearchItem> orders;
  final int total;
  final int page;
  final int pageSize;
  
  OrderSearchData({
    required this.orders,
    required this.total,
    required this.page,
    required this.pageSize,
  });
}

class OrderSearchItem {
  final int id;
  final String title;
  final String orderNo;
  final String status;
  final List<String> fabrics;
  final FactoryInfo factory;
  final DateTime createdAt;
  final DateTime updatedAt;
  
  OrderSearchItem({
    required this.id,
    required this.title,
    required this.orderNo,
    required this.status,
    required this.fabrics,
    required this.factory,
    required this.createdAt,
    required this.updatedAt,
  });
}

class FactoryInfo {
  final String id;
  final String name;
  
  FactoryInfo({
    required this.id,
    required this.name,
  });
}
```

### 搜索建议模型

```dart
class SearchSuggestionResponse {
  final bool success;
  final SearchSuggestionData data;
  
  SearchSuggestionResponse({
    required this.success,
    required this.data,
  });
}

class SearchSuggestionData {
  final List<SearchSuggestion> suggestions;
  
  SearchSuggestionData({
    required this.suggestions,
  });
}

class SearchSuggestion {
  final String type;
  final String text;
  final String highlight;
  
  SearchSuggestion({
    required this.type,
    required this.text,
    required this.highlight,
  });
}
```

## 3. 使用方法

### 基础搜索

```dart
// 创建搜索服务实例
final searchService = OrderSearchService(
  baseUrl: 'http://your-api-domain:8008',
  token: 'your-auth-token',
);

// 执行搜索
try {
  final response = await searchService.searchOrders(
    query: '连衣裙',
    status: 'published',
    page: 1,
    pageSize: 20,
  );
  
  if (response.success) {
    final orders = response.data.orders;
    final total = response.data.total;
    // 处理搜索结果
  }
} catch (e) {
  // 处理错误
}
```

### 高级搜索

```dart
// 组合搜索条件
final response = await searchService.searchOrders(
  query: '真丝面料',
  status: 'published',
  startDate: '2024-01-01',
  endDate: '2024-12-31',
  sortBy: 'created_at',
  sortOrder: 'desc',
  page: 1,
  pageSize: 10,
);
```

### 搜索建议

```dart
// 获取搜索建议
final suggestions = await searchService.getSearchSuggestions(
  query: '连衣裙',
  limit: 5,
);

if (suggestions.success) {
  final suggestionList = suggestions.data.suggestions;
  // 显示建议列表
}
```

### 搜索统计

```dart
// 获取搜索统计
final statistics = await searchService.getSearchStatistics();

if (statistics.success) {
  final hotKeywords = statistics.data.hotKeywords;
  final totalOrders = statistics.data.totalOrders;
  // 显示统计信息
}
```

## 4. 状态管理

### 使用Provider/Bloc管理搜索状态

```dart
class OrderSearchProvider extends ChangeNotifier {
  List<OrderSearchItem> _orders = [];
  bool _isLoading = false;
  String _error = '';
  int _total = 0;
  int _currentPage = 1;
  
  // 执行搜索
  Future<void> searchOrders(OrderSearchRequest request) async {
    _isLoading = true;
    _error = '';
    notifyListeners();
    
    try {
      final response = await searchService.searchOrders(
        query: request.query,
        status: request.status,
        startDate: request.startDate,
        endDate: request.endDate,
        page: request.page,
        pageSize: request.pageSize,
        sortBy: request.sortBy,
        sortOrder: request.sortOrder,
      );
      
      if (response.success) {
        _orders = response.data.orders;
        _total = response.data.total;
        _currentPage = response.data.page;
      } else {
        _error = '搜索失败';
      }
    } catch (e) {
      _error = e.toString();
    } finally {
      _isLoading = false;
      notifyListeners();
    }
  }
  
  // 加载更多
  Future<void> loadMore(OrderSearchRequest request) async {
    if (_isLoading) return;
    
    _isLoading = true;
    notifyListeners();
    
    try {
      final response = await searchService.searchOrders(
        query: request.query,
        status: request.status,
        startDate: request.startDate,
        endDate: request.endDate,
        page: _currentPage + 1,
        pageSize: request.pageSize,
        sortBy: request.sortBy,
        sortOrder: request.sortOrder,
      );
      
      if (response.success) {
        _orders.addAll(response.data.orders);
        _total = response.data.total;
        _currentPage = response.data.page;
      }
    } catch (e) {
      _error = e.toString();
    } finally {
      _isLoading = false;
      notifyListeners();
    }
  }
}
```

## 5. UI组件

### 搜索页面

```dart
class OrderSearchPage extends StatefulWidget {
  @override
  _OrderSearchPageState createState() => _OrderSearchPageState();
}

class _OrderSearchPageState extends State<OrderSearchPage> {
  final TextEditingController _searchController = TextEditingController();
  final OrderSearchProvider _searchProvider = OrderSearchProvider();
  
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: Text('订单搜索')),
      body: Column(
        children: [
          // 搜索栏
          _buildSearchBar(),
          // 筛选条件
          _buildFilterSection(),
          // 搜索结果
          Expanded(child: _buildSearchResults()),
        ],
      ),
    );
  }
  
  Widget _buildSearchBar() {
    return Padding(
      padding: EdgeInsets.all(16),
      child: TextField(
        controller: _searchController,
        decoration: InputDecoration(
          hintText: '搜索订单...',
          suffixIcon: IconButton(
            icon: Icon(Icons.search),
            onPressed: _performSearch,
          ),
        ),
        onSubmitted: (_) => _performSearch(),
      ),
    );
  }
  
  Widget _buildFilterSection() {
    return Container(
      padding: EdgeInsets.symmetric(horizontal: 16),
      child: Row(
        children: [
          // 状态筛选
          DropdownButton<String>(
            value: _selectedStatus,
            items: _statusOptions.map((status) {
              return DropdownMenuItem(
                value: status.value,
                child: Text(status.label),
              );
            }).toList(),
            onChanged: (value) {
              setState(() {
                _selectedStatus = value;
              });
              _performSearch();
            },
          ),
          // 排序选项
          DropdownButton<String>(
            value: _selectedSortBy,
            items: _sortOptions.map((sort) {
              return DropdownMenuItem(
                value: sort.value,
                child: Text(sort.label),
              );
            }).toList(),
            onChanged: (value) {
              setState(() {
                _selectedSortBy = value;
              });
              _performSearch();
            },
          ),
        ],
      ),
    );
  }
  
  Widget _buildSearchResults() {
    return Consumer<OrderSearchProvider>(
      builder: (context, provider, child) {
        if (provider.isLoading && provider.orders.isEmpty) {
          return Center(child: CircularProgressIndicator());
        }
        
        if (provider.error.isNotEmpty) {
          return Center(child: Text(provider.error));
        }
        
        if (provider.orders.isEmpty) {
          return Center(child: Text('暂无搜索结果'));
        }
        
        return ListView.builder(
          itemCount: provider.orders.length + 1,
          itemBuilder: (context, index) {
            if (index == provider.orders.length) {
              // 加载更多按钮
              if (provider.hasMore) {
                return _buildLoadMoreButton();
              }
              return SizedBox.shrink();
            }
            
            final order = provider.orders[index];
            return _buildOrderItem(order);
          },
        );
      },
    );
  }
  
  Widget _buildOrderItem(OrderSearchItem order) {
    return Card(
      margin: EdgeInsets.symmetric(horizontal: 16, vertical: 8),
      child: ListTile(
        title: Text(order.title),
        subtitle: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text('订单号: ${order.orderNo}'),
            Text('状态: ${order.status}'),
            Text('面料: ${order.fabrics.join(", ")}'),
            if (order.factory.name.isNotEmpty)
              Text('工厂: ${order.factory.name}'),
          ],
        ),
        trailing: Text(order.createdAt.toString().substring(0, 10)),
        onTap: () {
          // 跳转到订单详情页
          Navigator.pushNamed(
            context,
            '/order-detail',
            arguments: order.id,
          );
        },
      ),
    );
  }
  
  void _performSearch() {
    final request = OrderSearchRequest(
      query: _searchController.text,
      status: _selectedStatus,
      sortBy: _selectedSortBy,
      sortOrder: 'desc',
      page: 1,
      pageSize: 20,
    );
    
    _searchProvider.searchOrders(request);
  }
}
```

## 6. 错误处理

### 网络错误处理

```dart
class NetworkException implements Exception {
  final String message;
  final int? statusCode;
  
  NetworkException(this.message, [this.statusCode]);
  
  @override
  String toString() => 'NetworkException: $message';
}

// 在服务类中处理错误
Future<OrderSearchResponse> searchOrders(OrderSearchRequest request) async {
  try {
    final response = await http.get(
      Uri.parse('$baseUrl/api/order-search'),
      headers: {
        'Authorization': 'Bearer $token',
        'Content-Type': 'application/json',
      },
    );
    
    if (response.statusCode == 200) {
      final data = json.decode(response.body);
      return OrderSearchResponse.fromJson(data);
    } else if (response.statusCode == 401) {
      throw NetworkException('未授权访问', 401);
    } else if (response.statusCode == 400) {
      throw NetworkException('请求参数错误', 400);
    } else {
      throw NetworkException('服务器错误', response.statusCode);
    }
  } catch (e) {
    if (e is NetworkException) {
      rethrow;
    }
    throw NetworkException('网络连接失败: ${e.toString()}');
  }
}
```

## 7. 性能优化

### 搜索防抖

```dart
class Debouncer {
  final Duration delay;
  Timer? _timer;
  
  Debouncer({this.delay = const Duration(milliseconds: 500)});
  
  void run(VoidCallback action) {
    _timer?.cancel();
    _timer = Timer(delay, action);
  }
  
  void dispose() {
    _timer?.cancel();
  }
}

// 在搜索页面中使用
class _OrderSearchPageState extends State<OrderSearchPage> {
  final Debouncer _debouncer = Debouncer();
  
  void _onSearchChanged(String query) {
    _debouncer.run(() {
      _performSearch();
    });
  }
  
  @override
  void dispose() {
    _debouncer.dispose();
    super.dispose();
  }
}
```

### 缓存机制

```dart
class SearchCache {
  static final Map<String, OrderSearchResponse> _cache = {};
  static const Duration _cacheExpiry = Duration(minutes: 5);
  
  static String _generateKey(OrderSearchRequest request) {
    return '${request.query}_${request.status}_${request.page}_${request.pageSize}';
  }
  
  static OrderSearchResponse? get(OrderSearchRequest request) {
    final key = _generateKey(request);
    final cached = _cache[key];
    
    if (cached != null) {
      // 检查缓存是否过期
      final now = DateTime.now();
      if (now.difference(cached.timestamp) < _cacheExpiry) {
        return cached;
      } else {
        _cache.remove(key);
      }
    }
    
    return null;
  }
  
  static void set(OrderSearchRequest request, OrderSearchResponse response) {
    final key = _generateKey(request);
    _cache[key] = response;
  }
  
  static void clear() {
    _cache.clear();
  }
}
```

## 8. 测试

### 单元测试

```dart
void main() {
  group('OrderSearchService Tests', () {
    test('should search orders successfully', () async {
      // 测试搜索功能
    });
    
    test('should handle network errors', () async {
      // 测试错误处理
    });
    
    test('should get search suggestions', () async {
      // 测试搜索建议
    });
  });
}
```

## 9. 注意事项

1. **权限控制**：确保用户已登录并具有相应权限
2. **网络状态**：处理网络连接失败的情况
3. **数据验证**：验证API返回的数据格式
4. **用户体验**：提供加载状态和错误提示
5. **性能优化**：使用分页加载和缓存机制
6. **错误处理**：妥善处理各种异常情况 