import 'dart:io';
import 'package:shelf/shelf.dart';
import 'package:shelf/shelf_io.dart' as io;
import 'package:shelf_proxy/shelf_proxy.dart';
import 'package:http/http.dart' as http;

void main() async {
  final handler = proxyHandler('https://localhost/', client: http.Client());

  final pipeline =
      const Pipeline().addMiddleware(logRequests()).addHandler(handler);

  await io.serve(pipeline, 'localhost', 8080);
  print('Proxy server listening on port 8080');
}
